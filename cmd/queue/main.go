package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"postgres-job-queue/internal/queue"
	"postgres-job-queue/migrations"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 || (len(args) == 1 && args[0] != "migrate") {
		fmt.Fprintln(os.Stderr, "usage: queue [migrate]")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://queue:queue@localhost:5432/queue?sslmode=disable"
	}

	db, err := queue.Open(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connect to postgres: %v", err)
	}
	defer db.Close()

	if err := waitForPostgres(ctx, db); err != nil {
		log.Fatalf("ping postgres: %v", err)
	}

	if len(args) == 1 && args[0] == "migrate" {
		if err := db.Migrate(ctx, migrations.CreateJobsSQL); err != nil {
			log.Fatalf("migrate: %v", err)
		}

		fmt.Println("migration complete")
		return
	}

	fmt.Println("queue CLI is ready")
	fmt.Println("postgres connection OK")
}

func waitForPostgres(ctx context.Context, db *queue.DB) error {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	var lastErr error
	for {
		if err := db.Ping(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}

		select {
		case <-ctx.Done():
			return lastErr
		case <-ticker.C:
		}
	}
}
