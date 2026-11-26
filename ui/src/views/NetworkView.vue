<template>
  <div>
    <div class="breadcrumbs">
      <span>Networking</span>
      <span class="divider">/</span>
      <span class="current">Network Protocols</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Edge Network Protocols</h1>
        <p class="page-subtitle">Manage modern transport protocols including HTTP/3 (QUIC), 0-RTT, gRPC, and WebSockets across the global Anycast network.</p>
      </div>
    </div>

    <!-- Quick Telemetry Grid -->
    <div class="stats-overview-grid" style="margin-bottom: 1.5rem;">
      <CfStatCard
        label="HTTP/3 (with QUIC)"
        value="Enabled"
        delta-type="up"
        delta="Active"
        tooltip="UDP Port 443 Anycast"
      />
      <CfStatCard
        label="0-RTT Resumption"
        value="Active"
        delta-type="up"
        delta="TLS 1.3"
        tooltip="TLS 1.3 Zero Round-Trip"
      />
      <CfStatCard
        label="IPv6 Compatibility"
        value="Global IPv6"
        delta-type="up"
        delta="Dual-Stack"
        tooltip="Dual-stack IPv4/IPv6 gateway"
      />
      <CfStatCard
        label="WebSockets"
        value="Unlimited"
        delta-type="up"
        delta="TCP Stream"
        tooltip="Full duplex TCP streaming"
      />
    </div>

    <!-- Protocol Setting Rows -->
    <div style="display: flex; flex-direction: column; gap: 1.25rem;">
      <!-- HTTP/3 Card -->
      <div class="ui-kit-section">
        <div style="display: flex; justify-content: space-between; align-items: flex-start;">
          <div>
            <h3 class="ui-kit-title">HTTP/3 (with QUIC)</h3>
            <p class="ui-kit-desc">Accelerates web page loading by using QUIC, an encrypted UDP-based transport protocol that eliminates head-of-line blocking and improves performance on mobile networks.</p>
          </div>
          <button 
            class="btn-cf btn-cf-sm" 
            :class="network.http3 ? 'btn-cf-primary' : 'btn-cf-secondary'"
            @click="network.http3 = !network.http3"
          >
            {{ network.http3 ? 'Enabled' : 'Disabled' }}
          </button>
        </div>
      </div>

      <!-- 0-RTT Card -->
      <div class="ui-kit-section">
        <div style="display: flex; justify-content: space-between; align-items: flex-start;">
          <div>
            <h3 class="ui-kit-title">0-RTT Connection Resumption</h3>
            <p class="ui-kit-desc">Improves performance for clients that have previously connected to your site by resuming sessions without the TLS handshake round-trip latency overhead.</p>
          </div>
          <button 
            class="btn-cf btn-cf-sm" 
            :class="network['0rtt'] ? 'btn-cf-primary' : 'btn-cf-secondary'"
            @click="network['0rtt'] = !network['0rtt']"
          >
            {{ network['0rtt'] ? 'Enabled' : 'Disabled' }}
          </button>
        </div>
      </div>

      <!-- WebSockets Card -->
      <div class="ui-kit-section">
        <div style="display: flex; justify-content: space-between; align-items: flex-start;">
          <div>
            <h3 class="ui-kit-title">WebSockets Support</h3>
            <p class="ui-kit-desc">Allow real-time, bidirectional communication between your clients and origin servers using standard WebSocket protocol handshakes.</p>
          </div>
          <button 
            class="btn-cf btn-cf-sm" 
            :class="network.websockets ? 'btn-cf-primary' : 'btn-cf-secondary'"
            @click="network.websockets = !network.websockets"
          >
            {{ network.websockets ? 'Enabled' : 'Disabled' }}
          </button>
        </div>
      </div>

      <!-- gRPC Card -->
      <div class="ui-kit-section">
        <div style="display: flex; justify-content: space-between; align-items: flex-start;">
          <div>
            <h3 class="ui-kit-title">gRPC Proxying</h3>
            <p class="ui-kit-desc">Route bidirectional streaming gRPC connections to your backend services through Cloudflare's HTTP/2 proxy layer with zero reconfiguration.</p>
          </div>
          <button 
            class="btn-cf btn-cf-sm" 
            :class="network.grpc ? 'btn-cf-primary' : 'btn-cf-secondary'"
            @click="network.grpc = !network.grpc"
          >
            {{ network.grpc ? 'Enabled' : 'Disabled' }}
          </button>
        </div>
      </div>

      <!-- Onion Routing -->
      <div class="ui-kit-section">
        <div style="display: flex; justify-content: space-between; align-items: flex-start;">
          <div>
            <h3 class="ui-kit-title">Onion Routing (Tor Network)</h3>
            <p class="ui-kit-desc">Allow privacy-focused visitors using Tor to connect via Cloudflare's cryptographic .onion services, reducing CAPTCHAs while preserving origin security.</p>
          </div>
          <button 
            class="btn-cf btn-cf-sm" 
            :class="onionRouting ? 'btn-cf-primary' : 'btn-cf-secondary'"
            @click="onionRouting = !onionRouting"
          >
            {{ onionRouting ? 'Enabled' : 'Disabled' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { mockData } from '@/data/mockData.js'

const network = mockData.activeZone.network
const onionRouting = ref(true)
</script>
