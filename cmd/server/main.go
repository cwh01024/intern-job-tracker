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

	// Initialize database
	database, err := db.New(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()
	log.Println("✅ Database initialized")

	// Initialize components
	repo := repository.NewJobRepository(database)
	jobNotifier := notifier.NewDefaultIMessageNotifier()
	jobScraper := scraper.NewScraper(nil)
	jobScheduler := scheduler.New(repo, jobScraper, jobNotifier, *recipient)

	// Run once mode
	if *runOnce {
		log.Println("Running one-time job check...")
		if err := jobScheduler.RunNow(); err != nil {
			log.Fatalf("Job check failed: %v", err)
		}
		log.Println("✅ Job check complete")
		return
	}

	// Start scheduler
	if *recipient != "" {
		if err := jobScheduler.StartWithSchedule(*schedule); err != nil {
			log.Fatalf("Failed to start scheduler: %v", err)
		}
		log.Printf("✅ Scheduler started (recipient: %s, schedule: %s)", *recipient, *schedule)
	} else {
		log.Println("⚠️  No recipient configured - scheduler disabled")
		log.Println("   Run with -recipient='+1234567890' to enable notifications")
	}

	// Initialize API
	handler := api.NewHandler(repo, jobScheduler)
	router := handler.Router()

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("\nShutting down...")
		jobScheduler.Stop()
		os.Exit(0)
	}()

	// Start server
	addr := ":" + *port
	log.Printf("🚀 Server starting on http://localhost%s", addr)
	log.Println("   Dashboard: http://localhost" + addr)
	log.Println("   API: http://localhost" + addr + "/api/jobs")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
