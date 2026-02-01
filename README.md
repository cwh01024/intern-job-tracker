# 🚀 Intern Job Tracker

Automatically monitor summer intern SWE positions at top tech companies and receive iMessage notifications when new positions open!

## Features

- 🔍 **Automated Scraping**: Checks Google, Amazon, Uber, DoorDash career pages daily
- 📱 **iMessage Notifications**: Sends alerts via macOS Messages app when new jobs found
- 📊 **Dashboard**: Modern web interface to view all tracked positions
- ⏰ **Configurable Schedule**: Default daily at 9 AM, fully customizable
- 🗃️ **SQLite Storage**: Persistent job tracking with no external dependencies

## Quick Start

```bash
# Clone the repository
git clone https://github.com/cwh01024/intern-job-tracker.git
cd intern-job-tracker

# Install dependencies
go mod download

# Run the server (without notifications)
go run ./cmd/server

# Run with iMessage notifications enabled
go run ./cmd/server -recipient="+1234567890"
```

Open http://localhost:8080 to view the dashboard.

## Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `8080` | Server port |
| `-db` | `jobs.db` | Database file path |
| `-recipient` | `""` | iMessage recipient (phone or Apple ID) |
| `-schedule` | `0 9 * * *` | Cron schedule (default: 9 AM daily) |
| `-run-once` | `false` | Run job check once and exit |

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/jobs` | List all discovered jobs |
| GET | `/api/jobs/:id` | Get specific job details |
| GET | `/api/stats` | Get job statistics |
| POST | `/api/refresh` | Trigger manual job check |

## Project Structure

```
├── cmd/server/         # Main application
├── internal/
│   ├── api/           # HTTP handlers
│   ├── db/            # Database connection
│   ├── model/         # Data models
│   ├── notifier/      # iMessage integration
│   ├── repository/    # Data access layer
│   ├── scheduler/     # Cron job scheduling
│   └── scraper/       # Career page scraping
├── web/               # Frontend dashboard
└── migrations/        # SQL schema
```

## Running Tests

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/repository/... -v
```

## Companies Tracked

- **Google** - google.com/careers
- **Amazon** - amazon.jobs
- **Uber** - uber.com/careers
- **DoorDash** - careers.doordash.com

## Requirements

- **macOS** (for iMessage notifications)
- **Go 1.21+**
- Messages app configured with iMessage

## License

MIT
