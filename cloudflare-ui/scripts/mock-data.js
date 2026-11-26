/* ==========================================================================
   Cloudflare Dashboard (Kumo Design System) — Comprehensive Master Mock Data
   Extracted & Mirrored 1:1 from Production Dashboard (Account: 330b08634cd698ef2ff8587d66ac4b2c)
   Envato Commercial Standard — ThemeForest / CodeCanyon
   ========================================================================== */

window.CLOUDFLARE_MOCK_DATA = {
  account: {
    id: "330b08634cd698ef2ff8587d66ac4b2c",
    name: "Hridoy@duck.com's Account",
    email: "hridoy@duck.com",
    plan: "Cloudflare Enterprise & Pro Tier",
    hostname: "iad01.cloudflare.com",
    status: "All Systems Normal",
    region: "Global Anycast Network (330+ Cities)",
    billingCard: "•••• •••• •••• 3550 (Exp: 10/2027)"
  },

  domains: [
    {
      id: "northlab.dev",
      name: "northlab.dev",
      status: "Active",
      ssl: "Full (Strict)",
      requests_24h: "9.76k",
      requests_delta: "↘ 59.2%",
      bandwidth_24h: "482 MB",
      threats_blocked: 12,
      nameservers: ["ns1.cloudflare.com", "ns2.cloudflare.com"],
      last_modified: "Just now",
      plan: "Free Plan"
    },
    {
      id: "sovereign-app.dev",
      name: "sovereign-app.dev",
      status: "Active",
      ssl: "Full (Strict)",
      requests_24h: "142,890",
      requests_delta: "↗ 14.8%",
      bandwidth_24h: "892 MB",
      threats_blocked: 0,
      nameservers: ["ns1.cloudflare.com", "ns2.cloudflare.com"],
      last_modified: "2 mins ago",
      plan: "Pro Plan"
    },
    {
      id: "api-gateway.io",
      name: "api-gateway.io",
      status: "Active",
      ssl: "Full (Strict)",
      requests_24h: "412,500",
      requests_delta: "↗ 8.2%",
      bandwidth_24h: "2.4 GB",
      threats_blocked: 14,
      nameservers: ["ns1.cloudflare.com", "ns2.cloudflare.com"],
      last_modified: "Yesterday",
      plan: "Business Plan"
    },
    {
      id: "orbit-cluster.internal",
      name: "orbit-cluster.internal",
      status: "Active",
      ssl: "mTLS Mutual TLS",
      requests_24h: "98,420",
      requests_delta: "→ 0.0%",
      bandwidth_24h: "410 MB",
      threats_blocked: 0,
      nameservers: ["ns1.cloudflare.com", "ns2.cloudflare.com"],
      last_modified: "1 hour ago",
      plan: "Enterprise"
    }
  ],

  // Specific Zone Scope Data for 'northlab.dev'
  activeZone: {
    name: "northlab.dev",
    plan: "Free Website",
    status: "Active",
    dnsRecords: [
      { id: "dns_1", type: "A", name: "northlab.dev", content: "104.21.54.212", proxied: true, ttl: "Auto", modified: "1 min ago" },
      { id: "dns_2", type: "A", name: "codecollection", content: "172.67.188.45", proxied: true, ttl: "Auto", modified: "10 mins ago" },
      { id: "dns_3", type: "CNAME", name: "www", content: "northlab.dev", proxied: true, ttl: "Auto", modified: "1 hour ago" },
      { id: "dns_4", type: "CNAME", name: "api", content: "code-collection.workers.dev", proxied: true, ttl: "Auto", modified: "2 days ago" },
      { id: "dns_5", type: "TXT", name: "@", content: "v=spf1 include:_spf.mx.cloudflare.net ~all", proxied: false, ttl: "Auto", modified: "1 week ago" },
      { id: "dns_6", type: "TXT", name: "_dmarc", content: "v=DMARC1; p=reject; rua=mailto:dmarc@northlab.dev", proxied: false, ttl: "3600", modified: "1 week ago" },
      { id: "dns_7", type: "MX", name: "@", content: "route1.mx.cloudflare.net (Priority: 10)", proxied: false, ttl: "Auto", modified: "2 weeks ago" },
      { id: "dns_8", type: "AAAA", name: "northlab.dev", content: "2606:4700:3030::ac43:bc2d", proxied: true, ttl: "Auto", modified: "3 weeks ago" }
    ],
    ssl: {
      mode: "Full (Strict)",
      universalSsl: "Active",
      tls13: true,
      alwaysHttps: true,
      hsts: true,
      minTls: "TLS 1.2"
    },
    caching: {
      cachingLevel: "Standard",
      browserCacheTtl: "4 hours",
      tieredCache: "Enabled (Argo Tiered Caching)",
      purgeStatus: "Ready"
    },
    network: {
      http3: true,
      "0rtt": true,
      ipv6: true,
      websockets: true,
      grpc: true
    }
  },

  workers: [
    {
      app_id: "code-collection",
      name: "code-collection",
      subdomain: "code-collection.hridoy.workers.dev",
      custom_domain: "codecollection.northlab.dev",
      route: "codecollection.northlab.dev/*",
      repo: "github:hridoy/code-collection",
      branch: "main",
      deploy_cmd: "npx wrangler deploy",
      current_version: "v2.8.4 (commit 9e4f21)",
      status: "Active",
      invocations_24h: "13",
      invocations_delta: "↗ 8.3%",
      errors_24h: "0",
      median_cpu_time: "11 ms",
      median_cpu_delta: "↘ 7.9%",
      last_deployed: "2 hours ago via Git Integration",
      compatibility_date: "2026-11-03",
      compatibility_flags: ["nodejs_compat"],
      secrets: [
        { key: "APP_ORIGIN", value: "https://codecollection.northlab.dev", type: "Environment Variable" },
        { key: "DATABASE_URL", value: "libsql://code_collection-db.turso.io", type: "Environment Variable" },
        { key: "ENVATO_TOKEN", value: "env_tok_••••••••••••9a1c", type: "Secret (Encrypted)" },
        { key: "GOOGLE_CLIENT_ID", value: "98127394812-apps.googleusercontent.com", type: "Environment Variable" },
        { key: "GOOGLE_CLIENT_SECRET", value: "GOCSPX-••••••••••••38fd", type: "Secret (Encrypted)" },
        { key: "TOKEN_ENCRYPTION_KEY", value: "aes_256_••••••••••••77a0", type: "Secret (Encrypted)" }
      ],
      bindings: [
        { name: "DB", type: "D1 Database", target: "code_collection (UUID: 670002a1-53c5-4f67-b8c9-6aff880f8052)" },
        { name: "FILES", type: "R2 Bucket", target: "code-collection" },
        { name: "HYPERDRIVE", type: "Hyperdrive", target: "code-collection-hyperdrive" },
        { name: "IMPORTS", type: "Queue Producer", target: "code-collection-imports" }
      ],
      deployments: [
        { version: "v2.8.4", commit: "9e4f21", author: "Hridoy Duck", date: "Nov 22, 2026 06:40 UTC", status: "Active (Production)" },
        { version: "v2.8.3", commit: "8a1c02", author: "Hridoy Duck", date: "Nov 21, 2026 18:22 UTC", status: "Superseded" },
        { version: "v2.8.2", commit: "f4201b", author: "Hridoy Duck", date: "Nov 20, 2026 12:15 UTC", status: "Superseded" }
      ]
    },
    {
      app_id: "api-gateway-service",
      name: "api-gateway-service",
      subdomain: "api-gateway.hridoy.workers.dev",
      custom_domain: "api.sovereign-app.dev",
      route: "api.sovereign-app.dev/*",
      repo: "github:sovereign/api-gateway",
      branch: "main",
      deploy_cmd: "npx wrangler deploy",
      current_version: "v3.1.0 (commit c78210)",
      status: "Active",
      invocations_24h: "248,190",
      invocations_delta: "↗ 12.4%",
      errors_24h: "0.00%",
      median_cpu_time: "1.4 ms",
      median_cpu_delta: "↘ 0.2ms",
      last_deployed: "Yesterday",
      compatibility_date: "2026-11-01",
      compatibility_flags: ["nodejs_compat"],
      secrets: [
        { key: "JWT_SECRET", value: "jwt_••••••••••••55b1", type: "Secret (Encrypted)" }
      ],
      bindings: [
        { name: "USERS_DB", type: "D1 Database", target: "production_master_d1" },
        { name: "CACHE_KV", type: "KV Namespace", target: "edge_api_cache" }
      ],
      deployments: [
        { version: "v3.1.0", commit: "c78210", author: "Hridoy Duck", date: "Nov 21, 2026 09:12 UTC", status: "Active (Production)" }
      ]
    }
  ],

  databases: [
    {
      id: "670002a1-53c5-4f67-b8c9-6aff880f8052",
      name: "code_collection",
      engine: "Cloudflare D1 (SQLite Engine)",
      location: "Asia-Pacific (APAC / WNAM)",
      size: "42.8 MB",
      tables_count: 8,
      read_units_24h: "4.8M reads",
      write_units_24h: "890k writes",
      queries_count: "5.69M",
      created_at: "Nov 02, 2026",
      tables: [
        { name: "users", rows: 1420, columns: 8, size: "1.4 MB" },
        { name: "code_snippets", rows: 8940, columns: 12, size: "28.2 MB" },
        { name: "tags", rows: 142, columns: 4, size: "48 KB" },
        { name: "collections", rows: 620, columns: 6, size: "310 KB" },
        { name: "api_keys", rows: 48, columns: 7, size: "64 KB" },
        { name: "downloads", rows: 48120, columns: 5, size: "12.4 MB" },
        { name: "sessions", rows: 312, columns: 5, size: "180 KB" },
        { name: "schema_migrations", rows: 14, columns: 3, size: "16 KB" }
      ],
      sampleQueryResults: [
        { id: 1, title: "Auth0 JWT Middleware with Next.js", lang: "typescript", author: "hridoy", downloads: 4120, status: "published" },
        { id: 2, title: "Tailwind CSS v4 Dark Mode Customizer", lang: "css", author: "alex_chen", downloads: 2890, status: "published" },
        { id: 3, title: "Cloudflare Workers AI Llama-3 Streamer", lang: "javascript", author: "hridoy", downloads: 5410, status: "published" },
        { id: 4, title: "PostgreSQL Connection Pooler in Go", lang: "go", author: "sarah_j", downloads: 1840, status: "published" },
        { id: 5, title: "D1 Database SQLite Schema Migrator", lang: "sql", author: "hridoy", downloads: 3200, status: "published" }
      ]
    },
    {
      id: "9042a1bc-77d0-4821-a129-9bb71012fa01",
      name: "production_master_d1",
      engine: "Cloudflare D1 (SQLite Engine)",
      location: "North America (WNAM)",
      size: "128.4 MB",
      tables_count: 14,
      read_units_24h: "14.2M reads",
      write_units_24h: "2.4M writes",
      queries_count: "16.6M",
      created_at: "Oct 18, 2026",
      tables: [
        { name: "accounts", rows: 4200, columns: 10, size: "8.2 MB" },
        { name: "billing_records", rows: 84000, columns: 9, size: "64.0 MB" },
        { name: "workspaces", rows: 6800, columns: 6, size: "12.4 MB" }
      ],
      sampleQueryResults: []
    }
  ],

  buckets: [
    {
      bucket_id: "code-collection",
      name: "code-collection",
      region: "Asia-Pacific (APAC)",
      s3_endpoint: "https://330b08634cd698ef2ff8587d66ac4b2c.r2.cloudflarestorage.com/code-collection",
      custom_domain: "files.codecollection.northlab.dev",
      storage_class: "Standard (Zero Egress)",
      created_at: "Nov 02, 2026",
      size: "49.98 MB",
      objects_count: 564,
      class_a_ops: 564,
      class_b_ops: 437,
      public_dev_url: false,
      objects: [
        { key: "thumbnails/app-hero-light.png", type: "image/png", size: "1.4 MB", modified: "Nov 22, 2026 04:12" },
        { key: "thumbnails/dark-mode-card.svg", type: "image/svg+xml", size: "84 KB", modified: "Nov 21, 2026 19:40" },
        { key: "archives/v2.8.4-release.tar.gz", type: "application/gzip", size: "14.8 MB", modified: "Nov 22, 2026 06:40" },
        { key: "assets/combined-icons.svg", type: "image/svg+xml", size: "199 KB", modified: "Nov 20, 2026 11:20" },
        { key: "fonts/inter-variable.woff2", type: "font/woff2", size: "320 KB", modified: "Nov 18, 2026 08:00" },
        { key: "backups/d1-snapshot-latest.sqlite", type: "application/vnd.sqlite3", size: "32.1 MB", modified: "Nov 22, 2026 00:00" }
      ]
    },
    {
      bucket_id: "production-assets",
      name: "production-assets",
      region: "North America (ENAM)",
      s3_endpoint: "https://330b08634cd698ef2ff8587d66ac4b2c.r2.cloudflarestorage.com/production-assets",
      custom_domain: "static.sovereign-app.dev",
      storage_class: "Standard (Zero Egress)",
      created_at: "Oct 12, 2026",
      size: "1.2 GB",
      objects_count: 14200,
      class_a_ops: 1420,
      class_b_ops: 18920,
      public_dev_url: true,
      objects: []
    }
  ],

  r2Usage: {
    classA: "1,984 ops",
    classB: "19,357 ops",
    totalStorage: "1.25 GB",
    egressFee: "$0.00 / GB (Free Worldwide)",
    endpoint: "https://330b08634cd698ef2ff8587d66ac4b2c.r2.cloudflarestorage.com",
    auth: "AWS Signature V4 (HMAC-SHA256)"
  },

  analytics: {
    timeRange: "24h",
    totalRequests: "9.76k",
    totalRequestsDelta: "↘ 59.2%",
    workerInvocations: "13",
    workerInvocationsDelta: "↗ 8.3%",
    workerErrors: "0",
    cacheHitRate: "1.51%",
    cacheHitRateDelta: "↗ 1286.9%",
    cpuTimeP90: "11 ms",
    cpuTimeDelta: "↘ 7.9%",
    buildMinutes: "0 mins",
    statusBreakdown: [
      { code: "2xx Success", percent: 98.4, count: "9,604 reqs", color: "fill-blue" },
      { code: "3xx Redirect", percent: 1.1, count: "107 reqs", color: "fill-orange" },
      { code: "4xx Client Errors", percent: 0.5, count: "49 reqs", color: "fill-orange" },
      { code: "5xx Server Errors", percent: 0.0, count: "0 reqs", color: "fill-red" }
    ],
    resourceQuotas: [
      { name: "Workers RAM Quotas", percent: 18.5, label: "18.5% Allocated (V8 Isolate)", color: "fill-green" },
      { name: "CPU Time Budget", percent: 11.0, label: "11.0 ms P90 (Limit: 50 ms)", color: "fill-blue" },
      { name: "D1 Storage Quotas", percent: 4.2, label: "4.2% of 10 GB Included", color: "fill-purple" }
    ]
  },

  teamMembers: [
    { id: "usr_1", name: "Hridoy Duck", email: "hridoy@duck.com", role: "Super Administrator", groups: "All Permissions", twoFactor: true, lastActive: "Just now", initials: "HD", color: "var(--color-kumo-brand)" },
    { id: "usr_2", name: "Alex Chen", email: "alex.chen@northlab.dev", role: "Administrator", groups: "Developers, DNS Managers", twoFactor: true, lastActive: "14 mins ago", initials: "AC", color: "#059669" },
    { id: "usr_3", name: "Sarah Jenkins", email: "s.jenkins@enterprise.io", role: "Developer", groups: "Workers Deployers", twoFactor: true, lastActive: "2 hours ago", initials: "SJ", color: "#7c3aed" },
    { id: "usr_4", name: "Devin Murphy", email: "devin@finance-ops.net", role: "Billing Only", groups: "Finance Viewers", twoFactor: false, lastActive: "Yesterday", initials: "DM", color: "#dc2626" },
    { id: "usr_5", name: "Marcus Brody", email: "m.brody@security-audit.com", role: "Read Only", groups: "Auditors", twoFactor: true, lastActive: "3 days ago", initials: "MB", color: "#64748b" }
  ],

  billing: {
    currentPlan: "Cloudflare Pro + Paid Addons",
    monthlyPrice: "$20.00 / month",
    billingCycle: "Nov 1, 2026 – Nov 30, 2026",
    nextPayment: "Dec 1, 2026",
    currentEstimatedTotal: "$24.50",
    subscriptions: [
      { service: "Workers & Pages", tier: "Workers Free Base", status: "Subscribed", cost: "$0.00 / mo" },
      { service: "Storage & Databases (R2)", tier: "R2 Paid Tier (Standard)", status: "Subscribed", cost: "$4.50 / mo" },
      { service: "Zero Trust", tier: "Teams Free Base (up to 50 users)", status: "Subscribed", cost: "$0.00 / mo" },
      { service: "Zone: northlab.dev", tier: "Free Plan", status: "Subscribed", cost: "$0.00 / mo" },
      { service: "Zone: sovereign-app.dev", tier: "Pro Plan ($20/mo)", status: "Subscribed", cost: "$20.00 / mo" }
    ],
    paymentMethods: [
      { id: "pm_1", brand: "Mastercard", last4: "3550", expiry: "10/2027", isDefault: true, holder: "Hridoy Duck", billingEmail: "hridoy@duck.com" }
    ],
    invoices: [
      { id: "INV-CF-2026-11-01", date: "Nov 01, 2026", amount: "$24.50", status: "Paid", description: "Pro Plan + R2 Storage Addon" },
      { id: "INV-CF-2026-10-01", date: "Oct 01, 2026", amount: "$20.00", status: "Paid", description: "Cloudflare Pro Monthly Subscription" },
      { id: "INV-CF-2026-09-01", date: "Sep 01, 2026", amount: "$21.15", status: "Paid", description: "Pro Plan + D1 Storage Overage" },
      { id: "INV-CF-2026-08-01", date: "Aug 01, 2026", amount: "$20.00", status: "Paid", description: "Cloudflare Pro Monthly Subscription" }
    ]
  },

  apiTokens: [
    { id: "tok_1", name: "Orbit CI/CD Production Token", status: "Active", permissions: ["Workers: Edit", "D1: Edit", "R2: Write"], created: "Nov 02, 2026", lastUsed: "4 mins ago", masked: "cf_pat_9a4f••••••••••••38f1", full: "cf_pat_9a4f891b2c40de83a71829f0019238f1" },
    { id: "tok_2", name: "Grafana & Prometheus Observability", status: "Active", permissions: ["Analytics: Read", "Logs: Read"], created: "Oct 14, 2026", lastUsed: "12 secs ago", masked: "cf_pat_7b2c••••••••••••89e4", full: "cf_pat_7b2c0199e414fa99b1093847562089e4" },
    { id: "tok_3", name: "Cert-Manager ACME DNS01 Automated", status: "Active", permissions: ["DNS: Edit", "Zone: Read"], created: "Oct 01, 2026", lastUsed: "2 days ago", masked: "cf_pat_1c9d••••••••••••55a0", full: "cf_pat_1c9d8819aa43bb76cc091283940155a0" },
    { id: "tok_4", name: "Legacy Cold Storage Exporter", status: "Revoked", permissions: ["R2: Read"], created: "Aug 20, 2026", lastUsed: "1 month ago", masked: "cf_pat_00df••••••••••••99b2", full: "cf_pat_00df8192ab4400cd81726354891099b2" }
  ],

  auditLogs: [
    { timestamp: "2026-11-22 08:40:12", event: "ZONE_DNS_RECORD_UPDATE", subject: "A:codecollection.northlab.dev", actor: "hridoy@duck.com", ip: "198.51.100.24", zone: "northlab.dev", status: "Success" },
    { timestamp: "2026-11-22 06:40:02", event: "WORKER_DEPLOY_PRODUCTION", subject: "worker:code-collection:v2.8.4", actor: "Git Integration (commit 9e4f21)", ip: "cloudflare:git-hook", zone: "northlab.dev", status: "Success" },
    { timestamp: "2026-11-22 02:15:30", event: "D1_QUERY_EXECUTION", subject: "database:code_collection", actor: "hridoy@duck.com", ip: "198.51.100.24", zone: "Global", status: "Success" },
    { timestamp: "2026-11-21 19:42:15", event: "R2_OBJECT_UPLOAD", subject: "bucket:code-collection/thumbnails/", actor: "s3_api_token", ip: "198.51.100.24", zone: "Global", status: "Success" },
    { timestamp: "2026-11-21 14:10:00", event: "SSL_MODE_CHANGE", subject: "ssl_setting:Full Strict", actor: "hridoy@duck.com", ip: "198.51.100.24", zone: "northlab.dev", status: "Success" },
    { timestamp: "2026-11-20 11:20:44", event: "FIREWALL_RULE_CREATE", subject: "rule:Block Known Malicious UAs", actor: "hridoy@duck.com", ip: "198.51.100.24", zone: "northlab.dev", status: "Success" }
  ],

  settings: {
    profile: {
      firstName: "Hridoy",
      lastName: "Duck",
      email: "hridoy@duck.com",
      timezone: "(UTC-05:00) Eastern Time (US & Canada)",
      language: "English (United States)",
      avatarInitials: "HD"
    },
    security: {
      twoFactorEnabled: true,
      sessionTimeoutMinutes: 60,
      ipPinning: true,
      activeSessions: [
        { device: "MacBook Pro (macOS 15.1)", browser: "Chrome 131.0", ip: "198.51.100.24", location: "New York, USA", current: true },
        { device: "iPhone 16 Pro (iOS 18.1)", browser: "Safari Mobile", ip: "198.51.100.89", location: "New York, USA", current: false }
      ]
    },
    notifications: {
      emailAlerts: true,
      ddosAlerts: true,
      billingAlerts: true,
      weeklyDigests: false
    }
  },

  workersAi: {
    gateway: {
      name: "enterprise-ai-gateway",
      totalRequests: "1,420,890",
      cacheHitRate: "38.4%",
      tokensGenerated: "89.4M tokens",
      costSaved: "$1,840.50 (0% Markup)",
      avgLatency: "114 ms",
      providersConnected: ["Workers AI (Serverless GPUs)", "OpenAI", "Anthropic", "Cohere"]
    },
    models: [
      { id: "llama-3.3-70b", name: "Meta Llama 3.3 70B Instruct", category: "Text Generation", task: "Chat & Reasoning", speed: "142 tps", context: "128k tokens", status: "Active (FP8 Accelerated)", author: "Meta" },
      { id: "mistral-7b-v0.3", name: "Mistral 7B Instruct v0.3", category: "Text Generation", task: "Instruction Following", speed: "210 tps", context: "32k tokens", status: "Active", author: "Mistral AI" },
      { id: "whisper-large-v3", name: "OpenAI Whisper Large v3", category: "Speech Recognition", task: "Audio to Text", speed: "Realtime (0.2x RTF)", context: "30s chunks", status: "Active", author: "OpenAI" },
      { id: "bge-large-en", name: "BAAI BGE-Large-EN v1.5", category: "Vector Embeddings", task: "Text to 1024-d Vectors", speed: "840 tps", context: "512 tokens", status: "Active", author: "BAAI" },
      { id: "sdxl-turbo", name: "Stability AI SDXL Turbo", category: "Image Generation", task: "Text-to-Image (1 Step)", speed: "0.8s / image", context: "1024x1024", status: "Active", author: "Stability AI" }
    ],
    samplePrompt: "Explain how Cloudflare Workers run on V8 Isolates instead of virtual machines.",
    sampleResponse: "Cloudflare Workers utilize Google Chrome's V8 engine isolates rather than traditional containers or virtual machines. An isolate is a lightweight execution context that can spin up in under 5 milliseconds with near-zero memory footprint (around 3MB per isolate vs hundreds of megabytes for Docker containers). This architecture enables instant serverless execution across all 330+ Cloudflare edge locations with completely eliminated cold-start penalties."
  },

  turnstile: {
    widgets: [
      { id: "0x4AAAAAAAX12345678", name: "Production Login Shield", domain: "sovereign-app.dev", mode: "Managed Challenge", solves_24h: "48,920", passRate: "99.8%", status: "Active" },
      { id: "0x4AAAAAAAY87654321", name: "Sign Up Bot Blocker", domain: "northlab.dev", mode: "Non-Interactive", solves_24h: "12,410", passRate: "99.4%", status: "Active" },
      { id: "0x4AAAAAAAZ11223344", name: "API Checkout Invisible", domain: "api-gateway.io", mode: "Invisible", solves_24h: "84,100", passRate: "99.9%", status: "Active" }
    ]
  },

  workersKv: {
    namespaces: [
      { id: "a1b2c3d4e5f6-kv-01", name: "EDGE_CONFIG_STORE", keys_count: 1420, size: "14.2 MB", read_ops_24h: "1.8M", write_ops_24h: "42k" },
      { id: "b2c3d4e5f6a1-kv-02", name: "USER_AUTH_TOKENS", keys_count: 8940, size: "28.9 MB", read_ops_24h: "4.2M", write_ops_24h: "189k" },
      { id: "c3d4e5f6a1b2-kv-03", name: "RATE_LIMIT_BUCKETS", keys_count: 34100, size: "84.1 MB", read_ops_24h: "12.4M", write_ops_24h: "2.1M" }
    ],
    sampleKeys: [
      { key: "config:feature_flags:production", value: '{"beta_dashboard": true, "ai_streaming": true, "max_upload_mb": 500}', ttl: "Never Expire", metadata: "version: 2.4.1" },
      { key: "session:user_hridoy_duck:active", value: '{"ip": "198.51.100.24", "role": "superadmin", "exp": 1763812000}', ttl: "48 mins remaining", metadata: "ua: Chrome 131" },
      { key: "rate_limit:ip:198.51.100.24", value: '{"requests": 14, "bucket_max": 100, "reset_epoch": 1763805600}', ttl: "42 secs remaining", metadata: "policy: api_burst" }
    ]
  }
};
