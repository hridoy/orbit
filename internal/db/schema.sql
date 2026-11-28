-- ============================================================================
-- ORBIT CONTROL PLANE SCHEMA
-- Engine: Embedded SQLite 3 (Pure Go modernc.org/sqlite, CGO_ENABLED=0)
-- Target: Zero-Docker Sovereign PaaS
-- ============================================================================

CREATE TABLE IF NOT EXISTS nodes (
    node_id             TEXT PRIMARY KEY,
    hostname            TEXT NOT NULL,
    public_ip           TEXT NOT NULL,
    wireguard_ip        TEXT,
    wireguard_pubkey    TEXT,
    role                TEXT NOT NULL DEFAULT 'worker',
    cpu_cores           INTEGER NOT NULL DEFAULT 1,
    ram_total_bytes     INTEGER NOT NULL DEFAULT 0,
    disk_total_bytes    INTEGER NOT NULL DEFAULT 0,
    agent_version       TEXT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'online',
    last_heartbeat_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_nodes_status ON nodes(status);

CREATE TABLE IF NOT EXISTS apps (
    app_id              TEXT PRIMARY KEY,
    name                TEXT NOT NULL,
    domain              TEXT UNIQUE,
    node_id             TEXT NOT NULL,
    linux_user          TEXT NOT NULL,
    active_port         INTEGER NOT NULL,
    idle_port           INTEGER NOT NULL,
    binary_entrypoint   TEXT NOT NULL DEFAULT './server',
    exec_args           TEXT NOT NULL DEFAULT '[]',
    ram_limit_mb        INTEGER NOT NULL DEFAULT 256,
    cpu_quota_pct       INTEGER NOT NULL DEFAULT 100,
    pids_limit          INTEGER NOT NULL DEFAULT 128,
    health_check_path   TEXT NOT NULL DEFAULT '/healthz',
    health_timeout_sec  INTEGER NOT NULL DEFAULT 10,
    auto_restart        INTEGER NOT NULL DEFAULT 1,
    status              TEXT NOT NULL DEFAULT 'stopped',
    current_version     TEXT,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(node_id) REFERENCES nodes(node_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_apps_domain ON apps(domain);
CREATE INDEX IF NOT EXISTS idx_apps_node_id ON apps(node_id);

CREATE TABLE IF NOT EXISTS app_deployments (
    deployment_id       TEXT PRIMARY KEY,
    app_id              TEXT NOT NULL,
    version_hash        TEXT NOT NULL,
    port_used           INTEGER NOT NULL,
    status              TEXT NOT NULL,
    deployed_by         TEXT NOT NULL,
    duration_ms         INTEGER NOT NULL DEFAULT 0,
    logs                TEXT NOT NULL DEFAULT '',
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(app_id) REFERENCES apps(app_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_deployments_app ON app_deployments(app_id, created_at DESC);

CREATE TABLE IF NOT EXISTS databases (
    db_id               TEXT PRIMARY KEY,
    name                TEXT NOT NULL,
    engine              TEXT NOT NULL,
    version             TEXT NOT NULL,
    node_id             TEXT NOT NULL,
    port                INTEGER NOT NULL UNIQUE,
    socket_path         TEXT NOT NULL,
    ram_profile         TEXT NOT NULL DEFAULT 'micro',
    storage_limit_mb    INTEGER NOT NULL DEFAULT 10240,
    data_directory      TEXT NOT NULL,
    admin_user          TEXT NOT NULL,
    encrypted_password  TEXT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'stopped',
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(node_id) REFERENCES nodes(node_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_databases_node ON databases(node_id);

CREATE TABLE IF NOT EXISTS storage_buckets (
    bucket_id           TEXT PRIMARY KEY,
    owner_id            TEXT NOT NULL,
    is_public           INTEGER NOT NULL DEFAULT 0,
    quota_bytes         INTEGER NOT NULL DEFAULT 10737418240,
    used_bytes          INTEGER NOT NULL DEFAULT 0,
    versioning_enabled  INTEGER NOT NULL DEFAULT 0,
    object_lock_enabled INTEGER NOT NULL DEFAULT 0,
    retention_days      INTEGER NOT NULL DEFAULT 0,
    lifecycle_days      INTEGER NOT NULL DEFAULT 0,
    webhook_url         TEXT NOT NULL DEFAULT '',
    erasure_enabled     INTEGER NOT NULL DEFAULT 0,
    sse_enabled         INTEGER NOT NULL DEFAULT 0,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS storage_iam_keys (
    key_id              TEXT PRIMARY KEY,
    access_key          TEXT NOT NULL UNIQUE,
    secret_key_enc      TEXT NOT NULL,
    label               TEXT NOT NULL,
    bucket_scope        TEXT NOT NULL DEFAULT '*',
    prefix_scope        TEXT NOT NULL DEFAULT '',
    permissions         TEXT NOT NULL DEFAULT 'read',
    policy_json         TEXT NOT NULL DEFAULT '',
    active              INTEGER NOT NULL DEFAULT 1,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS storage_peers (
    peer_id             TEXT PRIMARY KEY,
    name                TEXT NOT NULL,
    endpoint            TEXT NOT NULL UNIQUE,
    auth_token_enc      TEXT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'unknown',
    last_sync_at        DATETIME,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS storage_objects (
    object_id           TEXT PRIMARY KEY,
    bucket_id           TEXT NOT NULL,
    object_key          TEXT NOT NULL,
    size_bytes          INTEGER NOT NULL,
    content_type        TEXT NOT NULL DEFAULT 'application/octet-stream',
    etag                TEXT NOT NULL,
    disk_rel_path       TEXT NOT NULL,
    retain_until        DATETIME,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(bucket_id) REFERENCES storage_buckets(bucket_id) ON DELETE CASCADE,
    UNIQUE(bucket_id, object_key)
);

CREATE INDEX IF NOT EXISTS idx_storage_objects_lookup ON storage_objects(bucket_id, object_key);

CREATE TABLE IF NOT EXISTS storage_object_versions (
    version_id          TEXT PRIMARY KEY,
    bucket_id           TEXT NOT NULL,
    object_key          TEXT NOT NULL,
    size_bytes          INTEGER NOT NULL,
    content_type        TEXT NOT NULL DEFAULT 'application/octet-stream',
    etag                TEXT NOT NULL,
    disk_rel_path       TEXT NOT NULL,
    is_delete_marker    INTEGER NOT NULL DEFAULT 0,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(bucket_id) REFERENCES storage_buckets(bucket_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_storage_versions_lookup ON storage_object_versions(bucket_id, object_key, created_at DESC);

CREATE TABLE IF NOT EXISTS storage_multipart_uploads (
    upload_id           TEXT PRIMARY KEY,
    bucket_id           TEXT NOT NULL,
    object_key          TEXT NOT NULL,
    initiated_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(bucket_id) REFERENCES storage_buckets(bucket_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS storage_multipart_parts (
    upload_id           TEXT NOT NULL,
    part_number         INTEGER NOT NULL,
    etag                TEXT NOT NULL,
    size_bytes          INTEGER NOT NULL,
    part_path           TEXT NOT NULL,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(upload_id, part_number),
    FOREIGN KEY(upload_id) REFERENCES storage_multipart_uploads(upload_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS secrets (
    secret_id           TEXT PRIMARY KEY,
    scope_type          TEXT NOT NULL,
    scope_id            TEXT NOT NULL,
    secret_key          TEXT NOT NULL,
    ciphertext_base64   TEXT NOT NULL,
    nonce_base64        TEXT NOT NULL,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(scope_type, scope_id, secret_key)
);

CREATE INDEX IF NOT EXISTS idx_secrets_scope ON secrets(scope_type, scope_id);

CREATE TABLE IF NOT EXISTS api_tokens (
    token_id            TEXT PRIMARY KEY,
    name                TEXT NOT NULL,
    token_hash          TEXT NOT NULL UNIQUE,
    role                TEXT NOT NULL DEFAULT 'operator',
    expires_at          DATETIME,
    last_used_at        DATETIME,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tls_certificates (
    domain              TEXT PRIMARY KEY,
    cert_pem            TEXT NOT NULL,
    key_pem             TEXT NOT NULL,
    expires_at          DATETIME NOT NULL,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Key-check row: encrypted with the current master key so startup can tell which key is active.
CREATE TABLE IF NOT EXISTS vault_meta (
    name                TEXT PRIMARY KEY,
    ciphertext_base64   TEXT NOT NULL,
    nonce_base64        TEXT NOT NULL
);

-- Where an app came from (GitHub) and which commit is live, for update checks.
CREATE TABLE IF NOT EXISTS app_sources (
    app_id              TEXT PRIMARY KEY,
    repo_url            TEXT NOT NULL,
    git_ref             TEXT NOT NULL DEFAULT 'main',
    installed_commit    TEXT NOT NULL DEFAULT '',
    latest_commit       TEXT NOT NULL DEFAULT '',
    checked_at          DATETIME,
    FOREIGN KEY(app_id) REFERENCES apps(app_id) ON DELETE CASCADE
);

-- Dashboard login. One row (id = 1) holds the operator account created on first run.
CREATE TABLE IF NOT EXISTS admin_users (
    id                  INTEGER PRIMARY KEY CHECK (id = 1),
    username            TEXT NOT NULL,
    password_hash       TEXT NOT NULL,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- One row per public domain. Private keys and EAB secrets are stored encrypted by internal/tlsengine.
CREATE TABLE IF NOT EXISTS tls_domains (
    domain              TEXT PRIMARY KEY,
    provider            TEXT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'pending',
    last_error          TEXT NOT NULL DEFAULT '',
    cert_pem            TEXT NOT NULL DEFAULT '',
    key_enc             TEXT NOT NULL DEFAULT '',
    eab_kid             TEXT NOT NULL DEFAULT '',
    eab_hmac_enc        TEXT NOT NULL DEFAULT '',
    expires_at          DATETIME,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Key/value server settings, such as the panel (dashboard) domain.
CREATE TABLE IF NOT EXISTS server_settings (
    key                 TEXT PRIMARY KEY,
    value               TEXT NOT NULL,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Saved third-party credentials (Cloudflare API token, ZeroSSL EAB). Secrets are encrypted.
CREATE TABLE IF NOT EXISTS integrations (
    provider            TEXT PRIMARY KEY,
    secret_enc          TEXT NOT NULL,
    extra               TEXT NOT NULL DEFAULT '',
    hint                TEXT NOT NULL DEFAULT '',
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Dashboard members. role is admin, operator, or viewer. The password hash uses the same
-- PBKDF2 format as admin_users.
CREATE TABLE IF NOT EXISTS dashboard_users (
    username            TEXT PRIMARY KEY,
    password_hash       TEXT NOT NULL,
    role                TEXT NOT NULL,
    disabled            INTEGER NOT NULL DEFAULT 0,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Public ports opened for a database. The agent forwards the public port to the database's
-- loopback port.
CREATE TABLE IF NOT EXISTS database_exposures (
    db_id               TEXT PRIMARY KEY,
    public_port         INTEGER NOT NULL UNIQUE,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Two-factor (TOTP) secrets for dashboard members. The secret is encrypted with the vault.
-- last_step is the last accepted time step, so one code cannot be used twice.
CREATE TABLE IF NOT EXISTS dashboard_totp (
    username            TEXT PRIMARY KEY,
    secret_ciphertext   TEXT NOT NULL,
    secret_nonce        TEXT NOT NULL,
    enabled             INTEGER NOT NULL DEFAULT 0,
    last_step           INTEGER NOT NULL DEFAULT 0,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
