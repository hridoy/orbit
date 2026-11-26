<template>
  <div>
    <div class="breadcrumbs">
      <span>Manage Account</span>
      <span class="divider">/</span>
      <span class="current">Billing & Subscriptions</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Billing & Subscriptions</h1>
        <p class="page-subtitle">Current Plan: <strong>{{ billing.currentPlan }}</strong> ({{ billing.monthlyPrice }}). Next scheduled invoice on {{ billing.nextPayment }}.</p>
      </div>
      <button class="btn-cf btn-cf-primary" @click="showPlanModal = true">Change Plan</button>
    </div>

    <!-- Tier Plan Cards Grid -->
    <div class="plan-cards-grid" style="margin-bottom: 2rem;">
      <!-- Free Plan -->
      <div class="plan-card">
        <div>
          <h3>Free Plan</h3>
          <div class="plan-price">$0 <span>/ month</span></div>
          <ul class="plan-features-list">
            <li>✓ Global Anycast DNS resolution</li>
            <li>✓ Universal SSL / TLS Certificate</li>
            <li>✓ 100k Workers Invocations / day</li>
            <li>✓ Zero Egress bandwidth fees</li>
            <li>✓ Community & Forum Support</li>
          </ul>
        </div>
        <button class="btn-cf btn-cf-secondary" disabled>Default Free</button>
      </div>

      <!-- Pro Plan (Current Active) -->
      <div class="plan-card current">
        <span class="plan-badge-current">Active Plan</span>
        <div>
          <h3>Pro Plan</h3>
          <div class="plan-price">$20 <span>/ month</span></div>
          <ul class="plan-features-list">
            <li>✓ Everything included in Free</li>
            <li>✓ Web Application Firewall (WAF)</li>
            <li>✓ Lossless Polish Image Optimization</li>
            <li>✓ 20 Page & Transform Rules</li>
            <li>✓ Enhanced HTTP/3 Prioritization</li>
          </ul>
        </div>
        <button class="btn-cf btn-cf-primary" disabled>Current Subscription</button>
      </div>

      <!-- Business Plan -->
      <div class="plan-card">
        <div>
          <h3>Business Plan</h3>
          <div class="plan-price">$200 <span>/ month</span></div>
          <ul class="plan-features-list">
            <li>✓ Everything included in Pro</li>
            <li>✓ 100% Uptime Financial SLA</li>
            <li>✓ Custom Nameservers</li>
            <li>✓ PCI DSS 4.0 Compliance</li>
            <li>✓ Priority 24/7 Email & Chat Support</li>
          </ul>
        </div>
        <button class="btn-cf btn-cf-secondary" @click="showPlanModal = true">Upgrade to Business</button>
      </div>
    </div>

    <!-- Active Subscriptions Table -->
    <div class="kumo-table-container" style="margin-bottom: 2rem;">
      <div class="table-header-title">
        <span>Active Subscriptions & Zone Add-ons</span>
        <span style="font-size: 0.8rem; color: var(--text-color-kumo-subtle);">Estimated month-to-date: {{ billing.currentEstimatedTotal }}</span>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Service / Product</th>
            <th>Tier Description</th>
            <th>Status</th>
            <th>Monthly Cost</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="sub in billing.subscriptions" :key="sub.service">
            <td><strong>{{ sub.service }}</strong></td>
            <td>{{ sub.tier }}</td>
            <td>
              <span class="cf-pill active">
                <span class="status-dot"></span>
                <span>{{ sub.status }}</span>
              </span>
            </td>
            <td><strong>{{ sub.cost }}</strong></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Invoices History Table -->
    <div class="kumo-table-container">
      <div class="table-header-title">
        <span>Invoices & Statements ({{ billing.invoices.length }})</span>
        <span style="font-size: 0.8rem; color: var(--text-color-kumo-subtle);">
          Payment method: {{ billing.paymentMethods[0].brand }} ending in {{ billing.paymentMethods[0].last4 }}
        </span>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Invoice ID</th>
            <th>Date</th>
            <th>Description</th>
            <th>Amount</th>
            <th>Status</th>
            <th>Receipt</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="inv in billing.invoices" :key="inv.id">
            <td><code>{{ inv.id }}</code></td>
            <td>{{ inv.date }}</td>
            <td>{{ inv.description }}</td>
            <td><strong>{{ inv.amount }}</strong></td>
            <td>
              <span class="cf-pill active">
                <span class="status-dot"></span>
                <span>{{ inv.status }}</span>
              </span>
            </td>
            <td>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="downloadInvoice(inv.id)">Download PDF ↗</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Upgrade / Change Plan Modal -->
    <div v-if="showPlanModal" class="cf-modal-backdrop" @click.self="showPlanModal = false">
      <div class="cf-modal-card" style="max-width: 500px;">
        <div class="cf-modal-header">
          <div class="cf-modal-title">Change Account Plan</div>
          <button class="cf-modal-close-btn" @click="showPlanModal = false">✕</button>
        </div>
        <div class="cf-modal-body">
          <p style="font-size: 0.88rem; color: var(--text-color-kumo-subtle); margin-bottom: 1rem;">
            Select a new subscription tier for your account. Changes will take effect immediately and will be prorated on your next billing cycle.
          </p>

          <div class="cf-form-group">
            <label>Target Tier</label>
            <select class="cf-select">
              <option value="pro">Cloudflare Pro ($20/month)</option>
              <option value="business">Cloudflare Business ($200/month)</option>
              <option value="enterprise">Cloudflare Enterprise (Custom SLA)</option>
            </select>
          </div>

          <div class="ui-alert-banner info">
            <div>ℹ️</div>
            <div>Your saved credit card (Mastercard ending in 3550) will be charged for any prorated adjustments.</div>
          </div>
        </div>
        <div class="cf-modal-footer">
          <button class="btn-cf btn-cf-secondary" @click="showPlanModal = false">Cancel</button>
          <button class="btn-cf btn-cf-primary" @click="confirmPlanChange">Confirm Upgrade</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { mockData } from '@/data/mockData.js'

const billing = mockData.billing
const showPlanModal = ref(false)

function downloadInvoice(id) {
  alert(`Generating official PDF statement for ${id}... Download initiated.`)
}

function confirmPlanChange() {
  showPlanModal.value = false
  alert('Subscription tier updated successfully!')
}
</script>
