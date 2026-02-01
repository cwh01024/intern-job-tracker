// API Base URL
const API_BASE = '/api';

// DOM Elements
const jobsTable = document.getElementById('jobs-tbody');
const totalJobsEl = document.getElementById('total-jobs');
const notifiedJobsEl = document.getElementById('notified-jobs');
const companiesCountEl = document.getElementById('companies-count');
const companyFilter = document.getElementById('company-filter');
const refreshBtn = document.getElementById('refresh-btn');
const toast = document.getElementById('toast');

// State
let jobs = [];
let stats = {};

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    loadData();
    setupEventListeners();
});

// Event Listeners
function setupEventListeners() {
    refreshBtn.addEventListener('click', handleRefresh);
    companyFilter.addEventListener('change', renderJobs);
}

// Load all data
async function loadData() {
    try {
        await Promise.all([loadJobs(), loadStats()]);
    } catch (error) {
        showToast('Failed to load data', 'error');
        console.error(error);
    }
}

// Load jobs
async function loadJobs() {
    const response = await fetch(`${API_BASE}/jobs`);
    if (!response.ok) throw new Error('Failed to fetch jobs');
    
    jobs = await response.json() || [];
    populateCompanyFilter();
    renderJobs();
}

// Load stats
async function loadStats() {
    const response = await fetch(`${API_BASE}/stats`);
    if (!response.ok) throw new Error('Failed to fetch stats');
    
    stats = await response.json();
    renderStats();
}

// Render stats
function renderStats() {
    totalJobsEl.textContent = stats.total_jobs || 0;
    notifiedJobsEl.textContent = stats.notified || 0;
    companiesCountEl.textContent = Object.keys(stats.by_company || {}).length;
}

// Populate company filter
function populateCompanyFilter() {
    const companies = [...new Set(jobs.map(j => j.company))];
    companyFilter.innerHTML = '<option value="">All Companies</option>';
    
    companies.forEach(company => {
        const option = document.createElement('option');
        option.value = company;
        option.textContent = company;
        companyFilter.appendChild(option);
    });
}

// Render jobs table
function renderJobs() {
    const filterValue = companyFilter.value;
    const filteredJobs = filterValue 
        ? jobs.filter(j => j.company === filterValue)
        : jobs;

    if (filteredJobs.length === 0) {
        jobsTable.innerHTML = `
            <tr>
                <td colspan="6">
                    <div class="empty-state">
                        <div class="empty-state-icon">🔍</div>
                        <p>No jobs found yet. Click "Refresh Jobs" to check for new positions.</p>
                    </div>
                </td>
            </tr>
        `;
        return;
    }

    jobsTable.innerHTML = filteredJobs.map(job => `
        <tr>
            <td>
                <span class="company-badge ${job.company.toLowerCase()}">${job.company}</span>
            </td>
            <td>${escapeHtml(job.title)}</td>
            <td>${job.location || 'N/A'}</td>
            <td>${formatDate(job.discovered_at)}</td>
            <td>
                <span class="status-badge ${job.notified ? 'notified' : 'pending'}">
                    ${job.notified ? '✓ Notified' : '⏳ Pending'}
                </span>
            </td>
            <td>
                <a href="${job.url}" target="_blank" rel="noopener" class="btn-apply">
                    Apply →
                </a>
            </td>
        </tr>
    `).join('');
}

// Handle refresh button click
async function handleRefresh() {
    refreshBtn.disabled = true;
    refreshBtn.innerHTML = '<span class="btn-icon">⏳</span> Refreshing...';

    try {
        const response = await fetch(`${API_BASE}/refresh`, { method: 'POST' });
        if (!response.ok) throw new Error('Refresh failed');
        
        showToast('Refresh complete! Reloading data...', 'success');
        await loadData();
    } catch (error) {
        showToast('Failed to refresh: ' + error.message, 'error');
    } finally {
        refreshBtn.disabled = false;
        refreshBtn.innerHTML = '<span class="btn-icon">🔄</span> Refresh Jobs';
    }
}

// Utility functions
function formatDate(dateStr) {
    if (!dateStr) return 'N/A';
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showToast(message, type = 'info') {
    toast.textContent = message;
    toast.className = `toast ${type} show`;
    
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}
