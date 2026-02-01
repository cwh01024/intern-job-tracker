-- Schema for Intern Job Tracker

CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    location TEXT,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    notified BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id INTEGER REFERENCES jobs(id),
    sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT
);

CREATE TABLE IF NOT EXISTS config (
    key TEXT PRIMARY KEY,
    value TEXT
);

-- Custom companies to track
CREATE TABLE IF NOT EXISTS companies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    career_url TEXT NOT NULL,
    search_term TEXT DEFAULT 'intern',
    enabled BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Run logs for tracking execution history
CREATE TABLE IF NOT EXISTS run_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    companies_checked INTEGER DEFAULT 0,
    jobs_found INTEGER DEFAULT 0,
    new_jobs INTEGER DEFAULT 0,
    notifications_sent INTEGER DEFAULT 0,
    duration_ms INTEGER DEFAULT 0,
    status TEXT DEFAULT 'success',
    error_message TEXT
);

-- Insert default companies if not exists
INSERT OR IGNORE INTO companies (name, career_url, search_term) VALUES 
    ('Google', 'https://www.google.com/about/careers/applications/jobs/results?q=software+intern&location=United+States', 'intern'),
    ('Amazon', 'https://www.amazon.jobs/en/search?base_query=software+intern&loc_query=United+States', 'intern'),
    ('Uber', 'https://www.uber.com/us/en/careers/list/?query=intern%20software&location=USA', 'intern'),
    ('DoorDash', 'https://careers.doordash.com/jobs/search?query=intern', 'intern');
