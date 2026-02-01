// API Base URL
const API_BASE = '/api';

// State
let jobs = [];
let companies = [];
let logs = [];
let metrics = {};

// DOM Elements
const refreshBtn = document.getElementById('refresh-btn');
const toast = document.getElementById('toast');

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    setupTabs();
    setupEventListeners();
    loadAllData();
});

// Tab Management
function setupTabs() {
    document.querySelectorAll('.tab').forEach(tab => {
        tab.addEventListener('click', () => {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));

            tab.classList.add('active');
            document.getElementById(`${tab.dataset.tab}-tab`).classList.add('active');
        });
    });
}

// Event Listeners
function setupEventListeners() {
    refreshBtn.addEventListener('click', handleRefresh);
    document.getElementById('company-filter').addEventListener('change', renderJobs);
    document.getElementById('add-company-btn').addEventListener('click', () => openCompanyModal());
    document.getElementById('modal-cancel').addEventListener('click', closeCompanyModal);
    document.getElementById('company-form').addEventListener('submit', handleCompanySubmit);
}

// Load all data
async function loadAllData() {
    try {
        await Promise.all([loadJobs(), loadCompanies(), loadLogs(), loadMetrics()]);
    } catch (error) {
        showToast('Failed to load data', 'error');
        console.error(error);
    }
}

// Load jobs
async function loadJobs() {
    const response = await fetch(`${API_BASE}/jobs`);
    jobs = await response.json() || [];
    populateCompanyFilter();
    renderJobs();
    renderJobsByCompany();
}

// Load companies
async function loadCompanies() {
    const response = await fetch(`${API_BASE}/companies`);
    companies = await response.json() || [];
    renderCompanies();
    document.getElementById('companies-count').textContent = companies.filter(c => c.enabled).length;
}

// Load logs
async function loadLogs() {
    const response = await fetch(`${API_BASE}/logs?limit=50`);
    logs = await response.json() || [];
    renderLogs();
    renderRecentLogs();
}

// Load metrics
async function loadMetrics() {
    const response = await fetch(`${API_BASE}/metrics`);
    metrics = await response.json() || {};
    renderMetrics();
}

// Render metrics on dashboard
function renderMetrics() {
    document.getElementById('total-jobs').textContent = metrics.jobs?.total || 0;
    document.getElementById('total-runs').textContent = metrics.runs?.total_runs || 0;
    document.getElementById('new-jobs-found').textContent = metrics.runs?.total_new_jobs_found || 0;
}

// Render recent logs on dashboard
function renderRecentLogs() {
    const container = document.getElementById('recent-logs');
    const recentLogs = logs.slice(0, 5);

    if (recentLogs.length === 0) {
        container.innerHTML = '<p class="empty">No runs yet</p>';
        return;
    }

    container.innerHTML = recentLogs.map(log => `
        <div class="log-item ${log.status}">
            <div class="log-time">${formatDate(log.run_at)}</div>
            <div class="log-details">
                <span>${log.new_jobs} new</span> / 
                <span>${log.jobs_found} found</span>
                <span class="log-duration">${log.duration_ms}ms</span>
            </div>
        </div>
    `).join('');
}

// Render jobs by company chart
function renderJobsByCompany() {
    const container = document.getElementById('jobs-by-company');
    const byCompany = {};
    jobs.forEach(j => byCompany[j.company] = (byCompany[j.company] || 0) + 1);

    const sorted = Object.entries(byCompany).sort((a, b) => b[1] - a[1]);
    const max = sorted[0]?.[1] || 1;

    container.innerHTML = sorted.map(([company, count]) => `
        <div class="bar-row">
            <span class="bar-label">${company}</span>
            <div class="bar-container">
                <div class="bar" style="width: ${(count / max) * 100}%"></div>
                <span class="bar-value">${count}</span>
            </div>
        </div>
    `).join('') || '<p class="empty">No jobs yet</p>';
}

// Populate company filter
function populateCompanyFilter() {
    const filter = document.getElementById('company-filter');
    const uniqueCompanies = [...new Set(jobs.map(j => j.company))];
    filter.innerHTML = '<option value="">All Companies</option>' +
        uniqueCompanies.map(c => `<option value="${c}">${c}</option>`).join('');
}

// Render jobs table
function renderJobs() {
    const filter = document.getElementById('company-filter').value;
    const filtered = filter ? jobs.filter(j => j.company === filter) : jobs;
    const tbody = document.getElementById('jobs-tbody');

    if (filtered.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5"><div class="empty-state"><div class="empty-icon">🔍</div><p>No jobs found</p></div></td></tr>';
        return;
    }

    tbody.innerHTML = filtered.map(job => `
        <tr>
            <td><span class="company-badge ${job.company.toLowerCase()}">${job.company}</span></td>
            <td>${escapeHtml(job.title)}</td>
            <td>${job.location || 'N/A'}</td>
            <td>${formatDate(job.discovered_at)}</td>
            <td><a href="${job.url}" target="_blank" class="btn-apply">Apply →</a></td>
        </tr>
    `).join('');
}

// Render companies grid
function renderCompanies() {
    const grid = document.getElementById('companies-grid');

    if (companies.length === 0) {
        grid.innerHTML = '<p class="empty">No companies configured</p>';
        return;
    }

    grid.innerHTML = companies.map(c => `
        <div class="company-card ${c.enabled ? '' : 'disabled'}">
            <div class="company-header">
                <h4>${escapeHtml(c.name)}</h4>
                <span class="status-dot ${c.enabled ? 'active' : 'inactive'}"></span>
            </div>
            <p class="company-url">${truncateUrl(c.career_url)}</p>
            <p class="company-search">Search: "${c.search_term}"</p>
            <div class="company-actions">
                <button class="btn-small" onclick="editCompany(${c.id})">Edit</button>
                <button class="btn-small btn-danger" onclick="deleteCompany(${c.id})">Delete</button>
            </div>
        </div>
    `).join('');
}

// Render logs table
function renderLogs() {
    const tbody = document.getElementById('logs-tbody');

    if (logs.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6"><div class="empty-state"><p>No runs yet</p></div></td></tr>';
        return;
    }

    tbody.innerHTML = logs.map(log => `
        <tr class="${log.status}">
            <td>${formatDateTime(log.run_at)}</td>
            <td>${log.companies_checked}</td>
            <td>${log.jobs_found}</td>
            <td><span class="highlight">${log.new_jobs}</span></td>
            <td>${log.duration_ms}ms</td>
            <td><span class="status-badge ${log.status}">${log.status}</span></td>
        </tr>
    `).join('');
}

// Company Modal
function openCompanyModal(company = null) {
    document.getElementById('modal-title').textContent = company ? 'Edit Company' : 'Add Company';
    document.getElementById('company-id').value = company?.id || '';
    document.getElementById('company-name').value = company?.name || '';
    document.getElementById('company-url').value = company?.career_url || '';
    document.getElementById('company-search').value = company?.search_term || 'intern';
    document.getElementById('company-modal').classList.remove('hidden');
}

function closeCompanyModal() {
    document.getElementById('company-modal').classList.add('hidden');
    document.getElementById('company-form').reset();
}

async function handleCompanySubmit(e) {
    e.preventDefault();
    const id = document.getElementById('company-id').value;
    const data = {
        name: document.getElementById('company-name').value,
        career_url: document.getElementById('company-url').value,
        search_term: document.getElementById('company-search').value || 'intern',
        enabled: true
    };

    try {
        if (id) {
            await fetch(`${API_BASE}/companies/${id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            showToast('Company updated!', 'success');
        } else {
            await fetch(`${API_BASE}/companies`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            showToast('Company added!', 'success');
        }
        closeCompanyModal();
        await loadCompanies();
    } catch (error) {
        showToast('Failed to save company', 'error');
    }
}

window.editCompany = function (id) {
    const company = companies.find(c => c.id === id);
    if (company) openCompanyModal(company);
};

window.deleteCompany = async function (id) {
    if (!confirm('Delete this company?')) return;
    try {
        await fetch(`${API_BASE}/companies/${id}`, { method: 'DELETE' });
        showToast('Company deleted', 'success');
        await loadCompanies();
    } catch (error) {
        showToast('Failed to delete', 'error');
    }
};

// Handle refresh
async function handleRefresh() {
    refreshBtn.disabled = true;
    refreshBtn.innerHTML = '<span class="btn-icon">⏳</span> Running...';

    try {
        await fetch(`${API_BASE}/refresh`, { method: 'POST' });
        showToast('Job check complete!', 'success');
        await loadAllData();
    } catch (error) {
        showToast('Failed: ' + error.message, 'error');
    } finally {
        refreshBtn.disabled = false;
        refreshBtn.innerHTML = '<span class="btn-icon">🔄</span> Run Check Now';
    }
}

// Utility functions
function formatDate(dateStr) {
    if (!dateStr) return 'N/A';
    return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
}

function formatDateTime(dateStr) {
    if (!dateStr) return 'N/A';
    return new Date(dateStr).toLocaleString('en-US', {
        month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
    });
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function truncateUrl(url) {
    try {
        const u = new URL(url);
        return u.hostname + (u.pathname.length > 20 ? u.pathname.slice(0, 20) + '...' : u.pathname);
    } catch {
        return url.slice(0, 40) + '...';
    }
}

function showToast(message, type = 'info') {
    toast.textContent = message;
    toast.className = `toast ${type} show`;
    setTimeout(() => toast.classList.remove('show'), 3000);
}
