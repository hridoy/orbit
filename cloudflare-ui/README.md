# Cloudflare Dashboard (Kumo Design System) — Commercial Admin Template

[![Envato ThemeForest](https://img.shields.io/badge/Envato-ThemeForest%20Ready-green.svg)](https://themeforest.net/)
[![Pure Vanilla](https://img.shields.io/badge/Framework-100%25%20Vanilla%20HTML5%20%2F%20CSS3%20%2F%20JS-blue.svg)]()
[![Zero Dependencies](https://img.shields.io/badge/Dependencies-Zero%20External%20NPM-orange.svg)]()

A pixel-perfect, 1:1 authentic commercial admin template mirroring the complete **Cloudflare Dashboard (Kumo Design System)**. Extracted directly from production and crafted for commercial sale on Envato (ThemeForest / CodeCanyon).

---

## 🚀 Live Demo & Key Highlights

- **29+ Comprehensive Screens**:
  - **Account Home**: 4-stat metric strip, active zones, quick action tiles.
  - **Zone Overview (`northlab.dev`)**: Zone-scoped edge dashboard.
  - **DNS Records**: A, AAAA, CNAME, TXT, MX records table, Proxied (Orange cloud vs Grey cloud) toggle switch, Add Record modal, BIND export.
  - **SSL / TLS**: 4-card mode selector (Off, Flexible, Full, Full Strict), Universal SSL, TLS 1.3, HSTS.
  - **Security & WAF**: Managed OWASP rulesets, rate limiting, bot mitigation, blocked threats graphs.
  - **Caching & CDN**: Caching levels, Tiered Cache, 1-click Purge Everything modal.
  - **Rules**: Edge Transform rules, URL rewrites, Page rules.
  - **Network Protocols**: HTTP/3 (QUIC), 0-RTT Connection Resumption, IPv6, WebSockets, gRPC.
  - **Workers & Pages (Compute)**: 2-tier application cards, V8 isolate performance.
  - **Worker Detail (`code-collection`)**: Metrics, deployment release history, real-time live log stream, encrypted runtime variables and secrets.
  - **D1 SQLite Databases**: Databases roster, size and read/write unit tracking.
  - **D1 Console (`code_collection`)**: Interactive SQL query runner with tabular results and latency timing, tables & schema browser.
  - **R2 Object Storage**: S3-compatible layout with zero egress fee badges.
  - **R2 File Manager (`code-collection`)**: Folder tree, drag-and-drop object list, prefix search.
  - **Observability & Analytics**: HTTP status breakdown proportion bar, resource quotas.
  - **Zero Trust**: Sovereign firewall policy, identity-aware tunnels.
  - **Members & RBAC**: Team roster, avatar initials, role badges, Invite Member modal.
  - **Billing & Subscriptions**: Plan tier cards, monthly itemized usage, saved cards, invoices table with PDF triggers.
  - **API Tokens**: Granular scoped tokens, 1-click copy box, Global API key.
  - **Audit Logs**: Immutable security activity stream with filtering.
  - **Settings**: Profile information, active browser sessions with 1-click revocation.
  - **UI Components Kit**: Full showcase of buttons, pills, inputs, alerts, and tables.
  - **Authentication Screens**: Split-screen Login, Register with password strength meter, 2FA OTP pin entry, and Forgot Password.
  - **Error Screens**: Cloudflare 404 with Ray ID, Error 520 server connection diagram, and Scheduled Maintenance countdown.
- **Authentic 68-Token Design System**: Extracted directly from Cloudflare production CSS variables and OKLCH color spaces.
- **244-Symbol SVG Icon Spritesheet**: Includes official Cloudflare vector symbols (`assets/combined-icons.svg`).
- **Live Theme Customizer**: Floating gear trigger with instant Dark/Light mode toggle, 5 Accent Brand Colors (Orange `#f6821f`, Blue `#0051c3`, Emerald `#059669`, Purple `#7c3aed`, Red `#dc2626`), and Mini-Sidebar toggle (68px).
- **⌘K Spotlight Command Palette**: Instant modal search with live filtering across all 22 resources.
- **Zero Runtime Dependencies**: No Node.js build step required. Works out of the box on any web server or static host.

---

## 📁 Directory Structure

```text
cloudflare-ui/
├── index.html                  # Master Interactive SPA Shell
├── documentation.html          # Commercial Buyer Documentation
├── README.md                   # Project Overview & Envato Guide
├── assets/
│   ├── cloudflare-logo.svg     # Official Cloudflare Brandmark
│   └── combined-icons.svg      # 244-symbol SVG Sprite Sheet
├── styles/
│   ├── kumo-tokens.css         # Authentic Root Tokens & Palettes
│   ├── kumo-base.css           # Global Reset, Inter Typography, Shell Flexbox
│   ├── kumo-sidebar.css        # Sticky 256px Tree-Guideline Left Sidebar
│   ├── kumo-header.css         # Sticky 56px Top Header with Search & User Menu
│   ├── kumo-cards.css          # Metric Stat Cards, 2-Tier Item Cards, S3 Layout
│   ├── kumo-tables.css         # Enterprise Tables, Badges, Hover Actions
│   ├── kumo-components.css     # Form Controls, Buttons, Modals, Segmented Pickers
│   ├── kumo-modals.css         # ⌘K Command Palette, Action Modals
│   ├── kumo-pages.css          # Sub-Navigation Tabs, DNS Toggles, SQL Console
│   ├── kumo-customizer.css     # Theme Customizer Drawer & 5 Accent Color Engines
│   └── kumo-responsive.css     # Mobile Drawer & Collapsed Mini-Sidebar Mode
├── scripts/
│   ├── mock-data.js            # Comprehensive Production-Grade JSON Mock Data
│   ├── icons.js                # SVG Sprite Helper Utilities
│   └── app.js                  # Master SPA Controller & Router
└── pages/                      # 28 Standalone HTML Multi-Page Templates
```

---

## 🎨 Theme Customizer & Brand Colors

Change the accent color dynamically via JavaScript or by setting the `data-accent` attribute on the `<html>` tag:

```html
<html lang="en" data-accent="orange">   <!-- Cloudflare Orange (Default) -->
<html lang="en" data-accent="blue">     <!-- Deep Enterprise Blue -->
<html lang="en" data-accent="emerald">  <!-- Emerald Green -->
<html lang="en" data-accent="purple">   <!-- Royal Purple -->
<html lang="en" data-accent="red">      <!-- Ruby Red -->
```

---

## 🛠️ Quick Start

Simply serve the `cloudflare-ui` directory with any static file server:

```bash
# Python
python3 -m http.server 8080 --directory cloudflare-ui

# Node.js
npx serve cloudflare-ui

# Nginx / Apache / Caddy / Cloudflare Pages
# Drop the folder directly into your public document root!
```

---

## 📜 Envato Commercial License

Formatted and compliant with standard and extended commercial licensing terms for Envato Market (ThemeForest & CodeCanyon).
