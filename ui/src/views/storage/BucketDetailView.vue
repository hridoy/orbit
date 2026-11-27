<template>
  <div class="cf-bucket-page">
    <!-- 1. BREADCRUMBS & TOP BAR -->
    <div class="cf-bucket-breadcrumb-bar">
      <router-link to="/r2" class="cf-bucket-breadcrumb-link">
        <CfIcon name="r2" size="18" />
        <span>R2 Object Storage</span>
      </router-link>
      <span style="color: #d4d4d4;">›</span>
      <div class="cf-bucket-breadcrumb-name">
        <span>{{ bucket.name }}</span>
        <button class="cf-mini-icon-btn" title="Copy bucket name" @click="copyText(bucket.name, 'name')">
          <CfIcon v-if="copiedKey !== 'name'" name="cis-clipboard" size="14" />
          <CfIcon v-else name="cis-checkmark-1-medium" size="14" style="color: #10b981;" />
        </button>
      </div>
    </div>

    <!-- 2. SEGMENTED PILL TABLIST (Objects, Metrics, Settings) -->
    <div class="cf-bucket-tabs-row">
      <div class="cf-bucket-pill-tabs" role="tablist">
        <button
          role="tab"
          class="cf-bucket-tab"
          :class="{ active: currentTab === 'objects' }"
          @click="currentTab = 'objects'"
        >
          Objects
        </button>
        <button
          role="tab"
          class="cf-bucket-tab"
          :class="{ active: currentTab === 'metrics' }"
          @click="currentTab = 'metrics'"
        >
          Metrics
        </button>
        <button
          role="tab"
          class="cf-bucket-tab"
          :class="{ active: currentTab === 'settings' }"
          @click="currentTab = 'settings'"
        >
          Settings
        </button>
      </div>
    </div>

    <!-- TAB 1: OBJECTS -->
    <div v-if="currentTab === 'objects'" class="cf-bucket-content-body">
      <!-- Metadata Stat Cards Grid (4 Cards) -->
      <div class="cf-bucket-meta-grid">
        <div class="cf-bucket-meta-card">
          <span class="cf-bucket-meta-label">Location</span>
          <span class="cf-bucket-meta-val">Asia-Pacific (APAC)</span>
        </div>
        <div class="cf-bucket-meta-card">
          <span class="cf-bucket-meta-label">Storage class</span>
          <span class="cf-bucket-meta-val">Standard</span>
        </div>
        <div class="cf-bucket-meta-card">
          <span class="cf-bucket-meta-label">Created</span>
          <span class="cf-bucket-meta-val">Nov 3, 2026</span>
        </div>
        <div class="cf-bucket-meta-card">
          <span class="cf-bucket-meta-label">S3 API</span>
          <span class="cf-bucket-meta-val">
            <code :title="s3Endpoint">{{ s3Endpoint }}</code>
            <button class="cf-mini-icon-btn" title="Copy S3 URL" @click="copyText(s3Endpoint, 's3')">
              <svg v-if="copiedKey !== 's3'" width="13" height="13" viewBox="0 0 256 256" fill="currentColor">
                <path d="M216,40H88a8,8,0,0,0-8,8V80H40a8,8,0,0,0-8,8V216a8,8,0,0,0,8,8H168a8,8,0,0,0,8-8V176h40a8,8,0,0,0,8-8V48A8,8,0,0,0,216,40ZM160,208H48V96H160Zm48-48H176V88a8,8,0,0,0-8-8H96V56H208Z"></path>
              </svg>
              <svg v-else width="13" height="13" viewBox="0 0 256 256" fill="#10b981">
                <path d="M229.66,77.66l-128,128a8,8,0,0,1-11.32,0l-56-56a8,8,0,0,1,11.32-11.32L96,188.69,218.34,66.34a8,8,0,0,1,11.32,11.32Z"></path>
              </svg>
            </button>
          </span>
        </div>
      </div>

      <!-- Action Toolbar -->
      <div class="cf-bucket-toolbar-row">
        <div class="cf-bucket-search-box">
          <CfIcon name="cis-magnifying-glass" size="14" style="color: #71717a;" />
          <input
            v-model="searchPrefix"
            type="text"
            class="cf-bucket-search-input"
            placeholder="Search by prefix..."
          />
        </div>

        <div class="cf-bucket-actions-group">
          <!-- View prefixes as folders toggle -->
          <label class="cf-view-prefixes-toggle">
            <input type="checkbox" v-model="viewPrefixesAsFolders" />
            <span>View prefixes as folders</span>
          </label>

          <button v-if="viewPrefixesAsFolders" class="cf-btn-white" @click="showAddFolderModal = true">
            <CfIcon name="cis-folder-1" size="14" />
            <span>+ Add folder</span>
          </button>
          <button class="cf-btn-blue-primary" @click="showUploadDropzone = !showUploadDropzone">
            <CfIcon name="cis-arrow-wall-down" size="14" style="transform: rotate(180deg);" />
            <span>Upload</span>
          </button>
        </div>
      </div>

      <!-- In-Page Nested Folder Breadcrumb Trail -->
      <div v-if="currentPrefixPath.length > 0" class="cf-inpage-folder-breadcrumb">
        <span class="cf-folder-crumb-link" @click="navigateToPrefix('')">{{ bucket.name }}</span>
        <template v-for="(segment, idx) in currentPrefixPath" :key="idx">
          <span class="cf-folder-crumb-sep">/</span>
          <span 
            v-if="idx < currentPrefixPath.length - 1" 
            class="cf-folder-crumb-link" 
            @click="navigateToPrefix(currentPrefixPath.slice(0, idx + 1).join('/') + '/')"
          >
            {{ segment }}
          </span>
          <span v-else class="cf-folder-crumb-current">
            {{ segment }}
          </span>
        </template>
        <button class="cf-mini-icon-btn" :title="copiedKey === 'prefix' ? 'Copied!' : 'Copy prefix path'" @click="copyText(activePrefix, 'prefix')">
          <CfIcon v-if="copiedKey !== 'prefix'" name="cis-clipboard" size="12" />
          <CfIcon v-else name="cis-checkmark-1-medium" size="12" style="color: #10b981;" />
        </button>
      </div>

      <!-- Multi-select Bar (When >= 1 checked) -->
      <div v-if="selectedKeys.length > 0" class="cf-multiselect-active-bar">
        <span>{{ selectedKeys.length }} object(s) selected</span>
        <button class="cf-btn-danger-solid" @click="deleteSelected">
          Delete {{ selectedKeys.length }} object(s)
        </button>
      </div>

      <!-- Drag & Drop Zone Box (Toggleable) -->
      <div v-if="showUploadDropzone" class="cf-bucket-dropzone" @click="triggerFileInput">
        <svg width="36" height="36" viewBox="0 0 256 256" fill="currentColor" style="color: #0051c3;">
          <path d="M224,144v64a8,8,0,0,1-8,8H40a8,8,0,0,1-8-8V144a8,8,0,0,1,16,0v56H208V144a8,8,0,0,1,16,0ZM93.66,77.66,120,51.31V144a8,8,0,0,0,16,0V51.31l26.34,26.35a8,8,0,0,0,11.32-11.32l-40-40a8,8,0,0,0-11.32,0l-40,40A8,8,0,0,0,93.66,77.66Z"></path>
        </svg>
        <div style="font-weight: 500; font-size: 14px;">Drag and drop files or folders here, or click to browse</div>
        <div style="font-size: 12px; color: #71717a;">Supports individual objects up to 5 TB in size via multipart upload.</div>
        <input ref="fileInput" type="file" style="display: none;" multiple @change="handleFileUpload" />
      </div>

      <!-- Objects Table Card -->
      <div class="cf-bucket-table-card">
        <table class="cf-bucket-table">
          <thead>
            <tr>
              <th style="width: 40px; text-align: center;">
                <input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll" />
              </th>
              <th>Name</th>
              <th style="text-align: right;">Size</th>
              <th>Storage class</th>
              <th>Last modified</th>
              <th style="width: 48px;"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="obj in filteredObjects" :key="obj.key">
              <td style="text-align: center;">
                <input
                  type="checkbox"
                  :checked="selectedKeys.includes(obj.key)"
                  @change="toggleSelect(obj.key)"
                />
              </td>
              <td>
                <div class="cf-obj-name-cell">
                  <!-- Folder Icon vs File Icon -->
                  <CfIcon v-if="obj.isFolder" name="cis-folder-1" size="16" style="color: #f59e0b;" />
                  <CfIcon v-else name="cis-page-text" size="16" style="color: #71717a;" />
                  <span class="cf-obj-link" @click="handleItemClick(obj)">
                    {{ obj.key }}
                  </span>
                </div>
              </td>
              <td style="text-align: right; font-variant-numeric: tabular-nums;">
                {{ obj.size }}
              </td>
              <td>{{ obj.storageClass }}</td>
              <td style="color: #71717a; font-variant-numeric: tabular-nums;">{{ obj.modified }}</td>
              <td style="position: relative;">
                <button class="cf-mini-icon-btn" @click="toggleRowMenu(obj.key)">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" fill="currentColor" viewBox="0 0 256 256">
                    <path d="M144,128a16,16,0,1,1-16-16A16,16,0,0,1,144,128ZM60,112a16,16,0,1,0,16,16A16,16,0,0,0,60,112Zm136,0a16,16,0,1,0,16,16A16,16,0,0,0,196,112Z"/>
                  </svg>
                </button>
                <div v-if="activeRowMenuKey === obj.key" class="cf-tokens-menu-popover">
                  <button v-if="!obj.isFolder" class="cf-tokens-menu-item" @click="downloadObj(obj)">Download</button>
                  <button class="cf-tokens-menu-item danger" @click="deleteSingleObj(obj)">Delete</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- TAB 2: METRICS -->
    <div v-else-if="currentTab === 'metrics'" class="cf-bucket-content-body">
      <!-- Top Time Range Filter -->
      <div style="display: flex; justify-content: flex-end;">
        <select class="cf-input-control" style="width: 180px; height: 34px;">
          <option>Last 24 hours</option>
          <option>Last 6 hours</option>
          <option>Last 7 days</option>
          <option>Last 30 days</option>
        </select>
      </div>

      <!-- 5 Summary KPI Cards -->
      <div class="cf-metrics-kpi-grid">
        <div class="cf-bucket-meta-card">
          <span class="cf-bucket-meta-label">Average Storage</span>
          <span class="cf-bucket-meta-val">49.98 MB</span>
        </div>
        <div class="cf-bucket-meta-card">
          <span class="cf-bucket-meta-label">Data Retrieved</span>
          <span class="cf-bucket-meta-val">0 B</span>
        </div>
        <div class="cf-bucket-meta-card">
          <span class="cf-bucket-meta-label">Class A Operations</span>
          <span class="cf-bucket-meta-val">11</span>
        </div>
        <div class="cf-bucket-meta-card">
          <span class="cf-bucket-meta-label">Class B Operations</span>
          <span class="cf-bucket-meta-val">18</span>
        </div>
        <div class="cf-bucket-meta-card">
          <span class="cf-bucket-meta-label">Request Distribution</span>
          <span class="cf-bucket-meta-val">0</span>
        </div>
      </div>

      <!-- 6 Charts Panels (2 Columns) -->
      <div class="cf-metrics-charts-grid">
        <!-- 1. Average Storage -->
        <div class="cf-chart-panel-card">
          <div class="cf-chart-panel-header">
            <span class="cf-chart-title">Average Storage ⓘ</span>
            <span class="cf-chart-legend">Total · 49.98 MB</span>
          </div>
          <div class="cf-chart-canvas-mock">
            <div v-for="b in 18" :key="b" class="cf-chart-bar-col" :style="{ height: '65%' }"></div>
          </div>
        </div>

        <!-- 2. Data Retrieved -->
        <div class="cf-chart-panel-card">
          <div class="cf-chart-panel-header">
            <span class="cf-chart-title">Data Retrieved ⓘ</span>
            <span class="cf-chart-legend">0 B</span>
          </div>
          <div style="height: 160px; display: flex; align-items: center; justify-content: center; color: #71717a; font-size: 13px;">
            No data is available for this time range
          </div>
        </div>

        <!-- 3. Class A Operations -->
        <div class="cf-chart-panel-card">
          <div class="cf-chart-panel-header">
            <span class="cf-chart-title">Class A Operations ⓘ</span>
            <span class="cf-chart-legend">Total · 11</span>
          </div>
          <div class="cf-chart-canvas-mock">
            <div v-for="b in 18" :key="b" class="cf-chart-bar-col" :style="{ height: `${((b * 5 + 10) % 70) + 15}%` }"></div>
          </div>
        </div>

        <!-- 4. Class B Operations -->
        <div class="cf-chart-panel-card">
          <div class="cf-chart-panel-header">
            <span class="cf-chart-title">Class B Operations ⓘ</span>
            <span class="cf-chart-legend">Total · 18</span>
          </div>
          <div class="cf-chart-canvas-mock">
            <div v-for="b in 18" :key="b" class="cf-chart-bar-col" :style="{ height: `${((b * 7 + 15) % 80) + 10}%` }"></div>
          </div>
        </div>

        <!-- 5. Errors by Status Code -->
        <div class="cf-chart-panel-card">
          <div class="cf-chart-panel-header">
            <span class="cf-chart-title">Errors by Status Code ⓘ</span>
            <span class="cf-chart-legend">Total · 6 (404 Not Found)</span>
          </div>
          <div class="cf-chart-canvas-mock">
            <div v-for="b in 18" :key="b" class="cf-chart-bar-col" :style="{ height: `${b === 8 || b === 14 ? 40 : 0}%`, background: '#ef4444' }"></div>
          </div>
        </div>

        <!-- 6. Request Distribution by Region -->
        <div class="cf-chart-panel-card">
          <div class="cf-chart-panel-header">
            <span class="cf-chart-title">Request Distribution by Region</span>
            <span class="cf-chart-legend">APAC · 100%</span>
          </div>
          <div style="height: 160px; display: flex; align-items: center; justify-content: center; color: #71717a; font-size: 13px;">
            All operations originating from APAC Region
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 3: SETTINGS -->
    <div v-else-if="currentTab === 'settings'" class="cf-bucket-content-body">
      <div class="cf-settings-split-grid">
        <!-- Main Area: 13 Stacked Settings Cards -->
        <div class="cf-settings-cards-stack">
          <!-- 1. General -->
          <div id="general" class="cf-setting-card">
            <div class="cf-setting-card-head">
              <h3 class="cf-setting-title">General</h3>
            </div>
            <div class="cf-bucket-meta-grid">
              <div class="cf-bucket-meta-card">
                <span class="cf-bucket-meta-label">Bucket Name</span>
                <span class="cf-bucket-meta-val">{{ bucket.name }}</span>
              </div>
              <div class="cf-bucket-meta-card">
                <span class="cf-bucket-meta-label">Location Hint</span>
                <span class="cf-bucket-meta-val">Asia-Pacific (APAC)</span>
              </div>
              <div class="cf-bucket-meta-card">
                <span class="cf-bucket-meta-label">Created</span>
                <span class="cf-bucket-meta-val">Nov 3, 2026</span>
              </div>
            </div>
          </div>

          <!-- 2. Custom Domains -->
          <div id="domains" class="cf-setting-card">
            <div class="cf-setting-card-head">
              <h3 class="cf-setting-title">Custom Domains ⓘ</h3>
              <button class="cf-btn-blue-primary" style="height: 30px; font-size: 12px;" @click="alert('Add Domain dialog.')">+ Add</button>
            </div>
            <p class="cf-setting-desc">There is no custom domain assigned to this bucket.</p>
          </div>

          <!-- 3. Public Development URL -->
          <div id="public-url" class="cf-setting-card">
            <div class="cf-setting-card-head">
              <h3 class="cf-setting-title">Public Development URL ⓘ</h3>
              <button class="cf-btn-blue-primary" style="height: 30px; font-size: 12px;" @click="alert('Public dev URL enabled.')">Enable</button>
            </div>
            <p class="cf-setting-desc">The public development URL is disabled for this bucket.</p>
          </div>

          <!-- 4. Basin Catalog -->
          <div id="basin" class="cf-setting-card">
            <div class="cf-setting-card-head">
              <h3 class="cf-setting-title">Basin Catalog ⓘ</h3>
              <button class="cf-btn-blue-primary" style="height: 30px; font-size: 12px;" @click="alert('Basin catalog enabled.')">Enable</button>
            </div>
            <p class="cf-setting-desc">Basin Catalog is disabled for this bucket.</p>
          </div>

          <!-- 5. CORS Policy -->
          <div id="cors" class="cf-setting-card">
            <div class="cf-setting-card-head">
              <h3 class="cf-setting-title">CORS Policy ⓘ</h3>
              <button class="cf-btn-blue-primary" style="height: 30px; font-size: 12px;" @click="alert('Add CORS modal.')">+ Add</button>
            </div>
            <p class="cf-setting-desc">There is no CORS Policy defined for this bucket.</p>
          </div>

          <!-- 6. Object Lifecycle Rules -->
          <div id="lifecycle" class="cf-setting-card">
            <div class="cf-setting-card-head">
              <h3 class="cf-setting-title">Object Lifecycle Rules ⓘ</h3>
              <button class="cf-btn-blue-primary" style="height: 30px; font-size: 12px;" @click="alert('Add Lifecycle Rule.')">+ Add</button>
            </div>
            <div class="cf-bucket-table-card">
              <table class="cf-bucket-table">
                <thead>
                  <tr>
                    <th>Rule Name</th>
                    <th>Prefix</th>
                    <th>Actions</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>Default Multipart Abort Rule</td>
                    <td>--</td>
                    <td>Abort uploads after 7 day(s)</td>
                    <td><span class="cf-tokens-badge-active">Enabled</span></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- 7. Bucket Lock Rules -->
          <div id="bucket-lock" class="cf-setting-card">
            <div class="cf-setting-card-head">
              <h3 class="cf-setting-title">Bucket Lock Rules ⓘ</h3>
              <button class="cf-btn-blue-primary" style="height: 30px; font-size: 12px;" @click="alert('Add Lock Rule.')">+ Add</button>
            </div>
            <p class="cf-setting-desc">There are no bucket lock rules defined for this bucket.</p>
          </div>

          <!-- 8. Event Notifications -->
          <div id="events" class="cf-setting-card">
            <div class="cf-setting-card-head">
              <h3 class="cf-setting-title">Event Notifications ⓘ</h3>
            </div>
            <p class="cf-setting-desc">Using event notifications requires you to purchase a Workers Paid plan.</p>
            <div style="background: #fafafa; border: 1px solid #e2e8f0; border-radius: 6px; padding: 1rem; display: flex; align-items: center; justify-content: space-between;">
              <div>
                <div style="font-weight: 600; font-size: 14px;">Workers Paid Plan required</div>
                <div style="font-size: 12px; color: #71717a;">Pay-as-you-go event queues and notification pipelines.</div>
              </div>
              <button class="cf-btn-blue-primary" style="height: 32px; font-size: 13px;" @click="alert('Purchase Workers Paid')">
                Purchase Workers Paid
              </button>
            </div>
          </div>

          <!-- 9. Data Access Logs -->
          <div id="logs" class="cf-setting-card">
            <div class="cf-setting-card-head">
              <h3 class="cf-setting-title">Data Access Logs</h3>
              <select class="cf-input-control" style="width: 140px; height: 32px;">
                <option>● Disabled</option>
                <option>● Enabled</option>
              </select>
            </div>
            <p class="cf-setting-desc">Record object-level read, write, list, multipart upload, and delete operations for this bucket.</p>
          </div>

          <!-- 10. Local Uploads (Beta) -->
          <div id="local-uploads" class="cf-setting-card">
            <div class="cf-setting-card-head">
              <h3 class="cf-setting-title">
                Local Uploads
                <span class="cf-tokens-badge-active" style="background: rgba(0,81,195,0.1); color: #0051c3;">Beta</span>
              </h3>
              <select class="cf-input-control" style="width: 140px; height: 32px;">
                <option>● Disabled</option>
                <option>● Enabled</option>
              </select>
            </div>
            <p class="cf-setting-desc">Improve performance for cross-region uploads. Check your Request Distribution metrics to see where requests originate.</p>
          </div>

          <!-- 11. On Demand Migration -->
          <div id="migration" class="cf-setting-card">
            <div class="cf-setting-card-head">
              <h3 class="cf-setting-title">On Demand Migration ⓘ</h3>
              <button class="cf-btn-blue-primary" style="height: 30px; font-size: 12px;" @click="alert('On demand migration enabled.')">Enable</button>
            </div>
            <p class="cf-setting-desc">Continuously migrate data from an Amazon S3 and Google Cloud bucket to this bucket only as read requests are made.</p>
          </div>

          <!-- 12. Default Storage Class -->
          <div id="storage-class" class="cf-setting-card">
            <h3 class="cf-setting-title">Default Storage Class</h3>
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
              <div
                class="cf-perm-option"
                :class="{ active: defaultStorageClass === 'Standard' }"
                @click="defaultStorageClass = 'Standard'"
              >
                <input type="radio" name="storageClass" :checked="defaultStorageClass === 'Standard'" />
                <div class="cf-perm-text-group">
                  <span class="cf-perm-name">Standard</span>
                  <span class="cf-perm-desc">Recommended for objects that will be accessed at least once a month.</span>
                </div>
              </div>
              <div
                class="cf-perm-option"
                :class="{ active: defaultStorageClass === 'Infrequent Access' }"
                @click="defaultStorageClass = 'Infrequent Access'"
              >
                <input type="radio" name="storageClass" :checked="defaultStorageClass === 'Infrequent Access'" />
                <div class="cf-perm-text-group">
                  <span class="cf-perm-name">Infrequent Access</span>
                  <span class="cf-perm-desc">Recommended for objects that will be accessed less than once a month.</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 13. Danger Zone -->
          <div id="danger" class="cf-setting-card cf-danger-card">
            <h3 class="cf-setting-title" style="color: #dc2626;">Danger Zone</h3>
            <div class="cf-danger-action-row">
              <div>
                <div style="font-weight: 600; font-size: 14px; color: #18181b;">Empty Bucket</div>
                <div class="cf-setting-desc">Delete all objects and folders in this bucket. The bucket and configurations will remain.</div>
              </div>
              <button class="cf-btn-danger-solid" @click="emptyBucket">Empty</button>
            </div>
            <div class="cf-danger-action-row">
              <div>
                <div style="font-weight: 600; font-size: 14px; color: #18181b;">Delete bucket</div>
                <div class="cf-setting-desc">Permanently delete this bucket and all configurations associated with it.</div>
              </div>
              <button class="cf-btn-danger-solid" @click="deleteBucket">Delete</button>
            </div>
          </div>
        </div>

        <!-- Right Rail Table of Contents ("On this page") -->
        <nav class="cf-settings-toc-rail" aria-label="On this page">
          <span class="cf-toc-header">On this page</span>
          <a href="#general" class="cf-toc-link">General</a>
          <a href="#domains" class="cf-toc-link">Custom Domains</a>
          <a href="#public-url" class="cf-toc-link">Public Development URL</a>
          <a href="#basin" class="cf-toc-link">Basin Catalog</a>
          <a href="#cors" class="cf-toc-link">CORS Policy</a>
          <a href="#lifecycle" class="cf-toc-link">Object Lifecycle Rules</a>
          <a href="#bucket-lock" class="cf-toc-link">Bucket Lock Rules</a>
          <a href="#events" class="cf-toc-link">Event Notifications</a>
          <a href="#logs" class="cf-toc-link">Data Access Logs</a>
          <a href="#local-uploads" class="cf-toc-link">Local Uploads <span style="font-size: 10px; color: #0051c3;">Beta</span></a>
          <a href="#migration" class="cf-toc-link">On Demand Migration</a>
          <a href="#storage-class" class="cf-toc-link">Default Storage Class</a>
          <a href="#danger" class="cf-toc-link danger">Danger Zone</a>
        </nav>
      </div>
    </div>

    <!-- MODAL: ADD FOLDER -->
    <div v-if="showAddFolderModal" class="cf-r2-modal-overlay" @click.self="showAddFolderModal = false">
      <div class="cf-r2-modal-card">
        <div class="cf-r2-modal-header">
          <h2 class="cf-r2-modal-title">Create folder</h2>
          <button class="cf-mini-icon-btn" @click="showAddFolderModal = false">✕</button>
        </div>
        <form @submit.prevent="createFolder">
          <div class="cf-r2-modal-body">
            <p class="cf-setting-desc">Folders are simulated prefixes in R2.</p>
            <div class="cf-form-section">
              <label class="cf-form-title" for="modal-folder-name">Folder name</label>
              <div style="display: flex; align-items: center;">
                <input
                  id="modal-folder-name"
                  v-model="newFolderName"
                  type="text"
                  class="cf-input-control"
                  placeholder="e.g. thumbnails"
                  required
                />
                <span style="padding: 0 8px; font-weight: 500; color: #71717a;">/</span>
              </div>
            </div>
          </div>
          <div class="cf-r2-modal-footer">
            <button type="button" class="cf-btn-outline" @click="showAddFolderModal = false">
              Cancel
            </button>
            <button type="submit" class="cf-btn-blue-primary">
              Create folder
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const currentTab = ref('objects')
const searchPrefix = ref('')
const copiedKey = ref('')
const showAddFolderModal = ref(false)
const showUploadDropzone = ref(false)
const newFolderName = ref('')
const activeRowMenuKey = ref(null)
const selectedKeys = ref([])
const defaultStorageClass = ref('Standard')
const fileInput = ref(null)
const viewPrefixesAsFolders = ref(true)

const bucket = ref({
  id: 'code-collection',
  name: 'code-collection',
  region: 'APAC'
})

const accountId = '330b08634cd698ef2ff8587d66ac4b2c'
const s3Endpoint = `https://${accountId}.r2.cloudflarestorage.com/${bucket.value.name}`

// Active nested prefix
const activePrefix = ref(route.query.prefix || '')

// Synchronize with URL query ?prefix=
watch(
  () => route.query.prefix,
  (newP) => {
    activePrefix.value = newP || ''
  },
  { immediate: true }
)

const currentPrefixPath = computed(() => {
  if (!activePrefix.value) return []
  return activePrefix.value.split('/').filter(Boolean)
})

const navigateToPrefix = (newPrefix) => {
  activePrefix.value = newPrefix
  router.push({
    path: `/r2/${bucket.value.id}`,
    query: newPrefix ? { prefix: newPrefix } : {}
  })
}

// Hierarchical object database
const rawDatabase = ref([
  {
    key: 'thumbnails/',
    isFolder: true,
    size: '—',
    storageClass: '—',
    modified: '—'
  },
  {
    key: 'thumbnails/6208/',
    isFolder: true,
    size: '—',
    storageClass: '—',
    modified: '—'
  },
  {
    key: 'thumbnails/6208/thumb.jpg',
    isFolder: false,
    size: '106.31 kB',
    storageClass: 'Standard',
    modified: 'Nov 04, 2026, 03:42 PM'
  },
  {
    key: 'thumbnails/6208/banner.png',
    isFolder: false,
    size: '412.50 kB',
    storageClass: 'Standard',
    modified: 'Nov 04, 2026, 03:45 PM'
  },
  {
    key: '6208/',
    isFolder: true,
    size: '—',
    storageClass: '—',
    modified: '—'
  },
  {
    key: 'frontend-bundle.tar.gz',
    isFolder: false,
    size: '14.24 MB',
    storageClass: 'Standard',
    modified: 'Nov 22, 2026, 03:22 PM'
  },
  {
    key: 'thumb.jpg',
    isFolder: false,
    size: '11.53 kB',
    storageClass: 'Standard',
    modified: 'Nov 20, 2026, 11:15 AM'
  },
  {
    key: 'dataset-2026-q4.parquet',
    isFolder: false,
    size: '35.70 MB',
    storageClass: 'Standard',
    modified: 'Nov 18, 2026, 08:40 PM'
  }
])

// Filter items according to active prefix and view mode
const filteredObjects = computed(() => {
  const currentP = activePrefix.value
  let pool = []

  if (!viewPrefixesAsFolders.value) {
    // Flat view: all files shown with full keys
    pool = rawDatabase.value.filter(o => !o.isFolder)
  } else if (!currentP) {
    // Root level view
    pool = [
      { key: 'thumbnails/', isFolder: true, size: '—', storageClass: '—', modified: '—' },
      { key: '6208/', isFolder: true, size: '—', storageClass: '—', modified: '—' },
      { key: 'frontend-bundle.tar.gz', isFolder: false, size: '14.24 MB', storageClass: 'Standard', modified: 'Nov 22, 2026, 03:22 PM' },
      { key: 'thumb.jpg', isFolder: false, size: '11.53 kB', storageClass: 'Standard', modified: 'Nov 20, 2026, 11:15 AM' },
      { key: 'dataset-2026-q4.parquet', isFolder: false, size: '35.70 MB', storageClass: 'Standard', modified: 'Nov 18, 2026, 08:40 PM' }
    ]
  } else if (currentP === 'thumbnails/') {
    // Inside thumbnails/
    pool = [
      { key: '6208/', isFolder: true, fullKey: 'thumbnails/6208/', size: '—', storageClass: '—', modified: '—' }
    ]
  } else if (currentP === 'thumbnails/6208/') {
    // Inside thumbnails/6208/
    pool = [
      { key: 'thumb.jpg', isFolder: false, fullKey: 'thumbnails/6208/thumb.jpg', size: '106.31 kB', storageClass: 'Standard', modified: 'Nov 04, 2026, 03:42 PM' },
      { key: 'banner.png', isFolder: false, fullKey: 'thumbnails/6208/banner.png', size: '412.50 kB', storageClass: 'Standard', modified: 'Nov 04, 2026, 03:45 PM' }
    ]
  } else {
    pool = []
  }

  const q = searchPrefix.value.toLowerCase().trim()
  if (!q) return pool
  return pool.filter(o => o.key.toLowerCase().includes(q))
})

const isAllSelected = computed(() => {
  return filteredObjects.value.length > 0 && selectedKeys.value.length === filteredObjects.value.length
})

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedKeys.value = []
  } else {
    selectedKeys.value = filteredObjects.value.map(o => o.key)
  }
}

const toggleSelect = (key) => {
  const idx = selectedKeys.value.indexOf(key)
  if (idx > -1) {
    selectedKeys.value.splice(idx, 1)
  } else {
    selectedKeys.value.push(key)
  }
}

const toggleRowMenu = (key) => {
  activeRowMenuKey.value = activeRowMenuKey.value === key ? null : key
}

const copyText = async (text, key) => {
  try {
    await navigator.clipboard.writeText(text)
    copiedKey.value = key
    setTimeout(() => {
      if (copiedKey.value === key) copiedKey.value = ''
    }, 2000)
  } catch (err) {
    console.warn(err)
  }
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileUpload = (e) => {
  const files = e.target.files
  if (!files || files.length === 0) return
  for (let i = 0; i < files.length; i++) {
    const f = files[i]
    rawDatabase.value.unshift({
      key: (activePrefix.value ? activePrefix.value : '') + f.name,
      isFolder: false,
      size: `${(f.size / 1024).toFixed(2)} kB`,
      storageClass: 'Standard',
      modified: 'Just now'
    })
  }
  showUploadDropzone.value = false
  alert(`Uploaded ${files.length} file(s) successfully.`)
}

const createFolder = () => {
  if (!newFolderName.value.trim()) return
  const clean = newFolderName.value.trim().replace(/\/+$/, '') + '/'
  const targetKey = (activePrefix.value ? activePrefix.value : '') + clean
  rawDatabase.value.unshift({
    key: targetKey,
    isFolder: true,
    size: '—',
    storageClass: '—',
    modified: '—'
  })
  showAddFolderModal.value = false
  newFolderName.value = ''
}

const deleteSelected = () => {
  if (confirm(`Delete ${selectedKeys.value.length} selected object(s)?`)) {
    selectedKeys.value = []
  }
}

const deleteSingleObj = (obj) => {
  const full = obj.fullKey || obj.key
  if (confirm(`Permanently delete "${full}"?`)) {
    activeRowMenuKey.value = null
  }
}

const downloadObj = (obj) => {
  const full = obj.fullKey || obj.key
  alert(`Downloading "${full}" from Cloudflare R2 edge...`)
  activeRowMenuKey.value = null
}

const handleItemClick = (obj) => {
  if (obj.isFolder) {
    const nextPrefix = obj.fullKey || (activePrefix.value + obj.key)
    navigateToPrefix(nextPrefix)
  } else {
    // Navigate to Cloudflare Deep Object Details Page
    const targetKey = obj.fullKey || (activePrefix.value ? activePrefix.value + obj.key : obj.key)
    router.push({
      path: `/r2/${bucket.value.id}/objects/${encodeURIComponent(targetKey)}/details`,
      query: activePrefix.value ? { prefix: activePrefix.value } : {}
    })
  }
}

const emptyBucket = () => {
  if (confirm('Are you sure you want to empty this bucket? All objects and folders will be deleted permanently.')) {
    rawDatabase.value = []
    alert('Bucket emptied successfully.')
  }
}

const deleteBucket = () => {
  if (confirm(`Are you sure you want to delete bucket "${bucket.value.name}"? This action cannot be undone.`)) {
    router.push('/r2')
  }
}
</script>
