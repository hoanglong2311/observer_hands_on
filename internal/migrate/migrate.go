// Package migrate applies embedded SQL migration files to the database.
package migrate

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

// Run applies all embedded SQL files in lexicographic order.
// Files are idempotent (use IF NOT EXISTS), so re-running is safe.
func Run(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := sqlFiles.ReadDir("sql")
	if err != nil {
		return err
	}
	for _, e := range entries {
		sql, err := sqlFiles.ReadFile("sql/" + e.Name())
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("apply %s: %w", e.Name(), err)
		}
		slog.Info("migration applied", "file", e.Name())
	}
	return nil
}
