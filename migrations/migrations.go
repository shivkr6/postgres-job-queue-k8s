package migrations

import _ "embed"

//go:embed 001_create_jobs.sql
var CreateJobsSQL string
