package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"postgres-job-queue/internal/queue"
	"postgres-job-queue/migrations"
)

func main() {
	cmd, err := parseCommand(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n", err)
		printUsage()
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

	switch cmd.name {
	case "migrate":
		if err := db.Migrate(ctx, migrations.CreateJobsSQL); err != nil {
			log.Fatalf("migrate: %v", err)
		}

		fmt.Println("migration complete")
		return
	case "enqueue":
		if err := db.Enqueue(ctx, cmd.payload); err != nil {
			log.Fatalf("enqueue: %v", err)
		}

		fmt.Println("job enqueued")
		return
	case "seed":
		if err := db.Seed(ctx, cmd.count); err != nil {
			log.Fatalf("seed: %v", err)
		}

		fmt.Printf("seeded %d jobs\n", cmd.count)
		return
	}

	fmt.Println("queue CLI is ready")
	fmt.Println("postgres connection OK")
}

type command struct {
	name    string
	payload json.RawMessage
	count   int
}

func parseCommand(args []string) (command, error) {
	if len(args) == 0 {
		return command{}, nil
	}

	switch args[0] {
	case "migrate":
		if len(args) != 1 {
			return command{}, fmt.Errorf("migrate does not accept extra arguments")
		}

		return command{name: "migrate"}, nil
	case "enqueue":
		if len(args) != 2 {
			return command{}, fmt.Errorf("enqueue requires exactly one JSON payload argument")
		}

		payload := json.RawMessage(args[1])
		if !json.Valid(payload) {
			return command{}, fmt.Errorf("enqueue payload must be valid JSON")
		}

		return command{name: "enqueue", payload: payload}, nil
	case "seed":
		if len(args) != 2 {
			return command{}, fmt.Errorf("seed requires --count=N")
		}

		const countPrefix = "--count="
		if !strings.HasPrefix(args[1], countPrefix) {
			return command{}, fmt.Errorf("seed requires --count=N")
		}

		count, err := strconv.Atoi(strings.TrimPrefix(args[1], countPrefix))
		if err != nil {
			return command{}, fmt.Errorf("seed count must be an integer")
		}
		if count <= 0 {
			return command{}, fmt.Errorf("seed count must be greater than 0")
		}

		return command{name: "seed", count: count}, nil
	default:
		return command{}, fmt.Errorf("unknown command %q", args[0])
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  queue")
	fmt.Fprintln(os.Stderr, "  queue migrate")
	fmt.Fprintln(os.Stderr, "  queue enqueue '<json-payload>'")
	fmt.Fprintln(os.Stderr, "  queue seed --count=N")
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
