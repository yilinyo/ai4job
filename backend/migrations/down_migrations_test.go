package migrations_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJobDomainDownMigrationsDropOnlyTheirOwnTables(t *testing.T) {
	tests := map[string]string{
		"000004_create_job_postings.down.sql":        "DROP TABLE IF EXISTS job_postings;",
		"000005_create_resume_applications.down.sql": "DROP TABLE IF EXISTS resume_applications;",
		"000006_create_interview_progress.down.sql":  "DROP TABLE IF EXISTS interview_progress;",
		"000008_create_dashboard_todos.down.sql":     "DROP TABLE IF EXISTS dashboard_todos;",
		"000009_create_activity_logs.down.sql":       "DROP TABLE IF EXISTS activity_logs;",
	}

	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := os.ReadFile(filepath.Join(".", name))
			if err != nil {
				t.Fatalf("read migration: %v", err)
			}
			if strings.TrimSpace(string(got)) != want {
				t.Fatalf("migration content = %q, want %q", strings.TrimSpace(string(got)), want)
			}
		})
	}
}
