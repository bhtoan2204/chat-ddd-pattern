package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/sijms/go-ora/v2"
)

var onceIns sync.Once
var singleton MigrateTool
var mutex = &sync.Mutex{}

type MigrateTool interface {
	Migrate(source, connStr string) error
}

type migrateTool struct{}

func NewMigrateTool() MigrateTool {
	onceIns.Do(func() {
		singleton = &migrateTool{}
	})

	return singleton
}

func (mgt *migrateTool) Migrate(source, connStr string) error {
	mutex.Lock()
	defer mutex.Unlock()

	path, err := normalizeFileSource(source)
	if err != nil {
		return err
	}

	db, err := sql.Open("oracle", connStr)
	if err != nil {
		return fmt.Errorf("open oracle failed: %w", err)
	}
	defer db.Close()

	if err := ensureSchemaMigrations(db); err != nil {
		return err
	}

	applied, err := getAppliedVersions(db)
	if err != nil {
		return err
	}

	files, err := listMigrationFiles(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if applied[file.Version] {
			continue
		}
		if err := applyMigrationFile(db, file.Path, file.Version); err != nil {
			return err
		}
	}

	return nil
}

type migrationFile struct {
	Path    string
	Version int
}

func normalizeFileSource(source string) (string, error) {
	path := strings.TrimPrefix(source, "file://")
	if path == "" {
		return "", fmt.Errorf("migration source path is empty")
	}
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get cwd failed: %w", err)
		}
		path = filepath.Join(cwd, path)
	}
	return filepath.Clean(path), nil
}

func listMigrationFiles(dir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migration dir failed: %w", err)
	}
	var files []migrationFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		version, err := parseVersion(name)
		if err != nil {
			return nil, err
		}
		files = append(files, migrationFile{
			Path:    filepath.Join(dir, name),
			Version: version,
		})
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Version < files[j].Version
	})
	return files, nil
}

func parseVersion(name string) (int, error) {
	parts := strings.SplitN(name, "_", 2)
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid migration filename: %s", name)
	}
	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid migration version in %s: %w", name, err)
	}
	return version, nil
}

func ensureSchemaMigrations(db *sql.DB) error {
	stmt := `
CREATE TABLE schema_migrations (
	version NUMBER(10) PRIMARY KEY,
	applied_at TIMESTAMP DEFAULT SYSTIMESTAMP NOT NULL
)`
	if _, err := db.Exec(stmt); err != nil && !isObjectExistsError(err) {
		return fmt.Errorf("create schema_migrations failed: %w", err)
	}
	return nil
}

func getAppliedVersions(db *sql.DB) (map[int]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("read schema_migrations failed: %w", err)
	}
	defer rows.Close()
	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan schema_migrations failed: %w", err)
		}
		applied[version] = true
	}
	return applied, nil
}

func applyMigrationFile(db *sql.DB, path string, version int) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read migration file failed: %w", err)
	}
	statements := splitSQLStatements(string(content))
	if len(statements) == 0 {
		return fmt.Errorf("migration file has no statements: %s", path)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	for _, stmt := range statements {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			if isObjectExistsError(err) {
				continue
			}
			_ = tx.Rollback()
			return fmt.Errorf("exec migration failed: %w", err)
		}
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations(version) VALUES (:1)", version); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update schema_migrations failed: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration failed: %w", err)
	}
	return nil
}

func splitSQLStatements(input string) []string {
	lines := strings.Split(input, "\n")
	var cleaned strings.Builder
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "--") || trim == "" {
			continue
		}
		cleaned.WriteString(line)
		cleaned.WriteString("\n")
	}
	raw := strings.Split(cleaned.String(), ";")
	var out []string
	for _, stmt := range raw {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		out = append(out, stmt)
	}
	return out
}

func isObjectExistsError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "ORA-00955")
}
