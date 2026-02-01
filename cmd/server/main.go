package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"intern-job-tracker/internal/api"
	"intern-job-tracker/internal/db"
	"intern-job-tracker/internal/notifier"
	"intern-job-tracker/internal/repository"
	"intern-job-tracker/internal/scheduler"
	"intern-job-tracker/internal/scraper"
)

func main() {
	// Command line flags
	port := flag.String("port", "8080", "Server port")
	dbPath := flag.String("db", "jobs.db", "Database file path")
	recipient := flag.String("recipient", "", "iMessage recipient (phone or Apple ID)")
	schedule := flag.String("schedule", "0 9 * * *", "Cron schedule for job checks")
	runOnce := flag.Bool("run-once", false, "Run job check once and exit")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("")

	// Initialize database
	database, err := db.New(*dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.Close()
	log.Println("✅ Database initialized")

	// Initialize repositories
	jobRepo := repository.NewJobRepository(database)
	companyRepo := repository.NewCompanyRepository(database)
	runLogRepo := repository.NewRunLogRepository(database)

	// Initialize components
	jobNotifier := notifier.NewDefaultIMessageNotifier()
	jobScraper := scraper.NewScraper(nil)
	jobScheduler := scheduler.New(jobRepo, companyRepo, runLogRepo, jobScraper, jobNotifier, *recipient)

	// Run once mode
	if *runOnce {
		if *recipient == "" {
			log.Println("⚠️  No recipient specified. Use -recipient=\"+1234567890\"")
		}
		if err := jobScheduler.RunNow(); err != nil {
			log.Fatalf("❌ Job check failed: %v", err)
		}
		return
	}

	// Start scheduler
	if *recipient != "" {
		if err := jobScheduler.StartWithSchedule(*schedule); err != nil {
			log.Fatalf("❌ Failed to start scheduler: %v", err)
		}
		log.Printf("✅ Scheduler started (recipient: %s, schedule: %s)", *recipient, *schedule)
	} else {
		log.Println("⚠️  No recipient configured - scheduler disabled")
		log.Println("   Run with -recipient=\"+1234567890\" to enable notifications")
	}

	// Initialize API
	handler := api.NewHandler(jobRepo, companyRepo, runLogRepo, jobScheduler)
	router := handler.Router()

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("\n⏹️  Shutting down...")
		jobScheduler.Stop()
		os.Exit(0)
	}()

	// Start server
	addr := ":" + *port
	log.Println("═══════════════════════════════════════════")
	log.Printf("🚀 Server starting on http://localhost%s", addr)
	log.Println("───────────────────────────────────────────")
	log.Println("📊 Dashboard: http://localhost" + addr)
	log.Println("📡 API Endpoints:")
	log.Println("   GET  /api/jobs       - List all jobs")
	log.Println("   GET  /api/companies  - List companies")
	log.Println("   POST /api/companies  - Add company")
	log.Println("   GET  /api/metrics    - View metrics")
	log.Println("   GET  /api/logs       - View run history")
	log.Println("   POST /api/refresh    - Trigger job check")
	log.Println("═══════════════════════════════════════════")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
