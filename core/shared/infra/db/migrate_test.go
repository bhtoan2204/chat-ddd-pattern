package db

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestSplitSQLStatements_PostgresFunctionBlocks(t *testing.T) {
	input := `
CREATE TABLE ledger_transactions (
    transaction_id text PRIMARY KEY
);

CREATE OR REPLACE FUNCTION touch_updated_at()
RETURNS trigger AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE INDEX idx_ledger_entries_transaction_id ON ledger_transactions(transaction_id);
`

	statements := splitPostgresSQLStatements(input)
	if len(statements) != 3 {
		t.Fatalf("expected 3 statements, got %d: %#v", len(statements), statements)
	}

	if !strings.Contains(statements[1], "RETURNS trigger AS $$") {
		t.Fatalf("expected function block to stay intact, got %q", statements[1])
	}
}

func TestPrepareStatementsForPostgres_AllUpMigrations(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "..", "..", "migration", "*.up.sql"))
	if err != nil {
		t.Fatalf("glob migration files failed: %v", err)
	}
	sort.Strings(files)
	if len(files) == 0 {
		t.Fatal("expected migration files")
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("read migration file failed: %v", err)
			}

			statements := splitPostgresSQLStatements(string(content))
			if len(statements) == 0 && strings.TrimSpace(string(content)) != "" {
				t.Fatalf("expected statements for %s", filepath.Base(file))
			}
		})
	}
}
