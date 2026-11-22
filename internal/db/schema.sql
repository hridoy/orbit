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
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS storage_objects (
    object_id           TEXT PRIMARY KEY,
    bucket_id           TEXT NOT NULL,
    object_key          TEXT NOT NULL,
    size_bytes          INTEGER NOT NULL,
    content_type        TEXT NOT NULL DEFAULT 'application/octet-stream',
    etag                TEXT NOT NULL,
    disk_rel_path       TEXT NOT NULL,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(bucket_id) REFERENCES storage_buckets(bucket_id) ON DELETE CASCADE,
    UNIQUE(bucket_id, object_key)
);

CREATE INDEX IF NOT EXISTS idx_storage_objects_lookup ON storage_objects(bucket_id, object_key);

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
