package migration

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Migration struct {
	db        *sqlx.DB
	migrateFS embed.FS
	seedFS    embed.FS
}

func New(db *sqlx.DB, files, seedFS embed.FS) *Migration {
	return &Migration{db, files, seedFS}
}

// Extracts the leading number from the filename
func extractNumber(s string) int {
	parts := strings.Split(s, "_")
	if len(parts) > 0 {
		numStr := parts[0]
		num, err := strconv.Atoi(numStr)
		if err == nil {
			return num
		}
	}
	return 0
}

func (m *Migration) Create(args []string) {

	if len(args) != 1 {
		fmt.Println("\nusage: \n\nmigration create migration_name")
		return
	}

	entries, err := os.ReadDir("migration")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	up := fmt.Sprintf("migration/%d_%s.up.sql", len(entries)/2+1, args[0])
	fmt.Println("create migration file on " + up)
	err = os.WriteFile(up, nil, 0755)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	down := fmt.Sprintf("migration/%d_%s.down.sql", len(entries)/2+1, args[0])
	fmt.Println("create migration file on " + down)
	err = os.WriteFile(down, nil, 0755)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("success")
}

func (m *Migration) Up() {

	_, err := m.db.Exec(`CREATE TABLE IF NOT EXISTS migration (
		version VARCHAR(255) NOT NULL,
		created_at bigint NOT NULL
	)`)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	migration := &struct {
		Version   string `db:"version"`
		CreatedAt int64  `db:"created_at"`
	}{}

	err = m.db.Get(migration, "SELECT * FROM migration ORDER BY created_at DESC")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error(err.Error())
		return
	}

	gotLatest := false
	if migration.CreatedAt == 0 {
		gotLatest = true
	}

	files, err := m.migrateFS.ReadDir("migration")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	sort.SliceStable(files, func(i, j int) bool {
		numI := extractNumber(files[i].Name())
		numJ := extractNumber(files[j].Name())
		return numI < numJ
	})

	tx, err := m.db.Beginx()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer tx.Rollback()

	for i := 0; i < len(files); i++ {
		fName := files[i].Name()
		if !strings.Contains(fName, ".up.sql") {
			continue
		}

		fileName := strings.ReplaceAll(fName, ".up.sql", "")
		if gotLatest {
			fmt.Printf("migration/%s \n", fName)
			byteFile, err := m.migrateFS.ReadFile(fmt.Sprintf("migration/%s", fName))
			if err != nil {
				slog.Error(err.Error())
				return
			}

			_, err = tx.Exec(string(byteFile))
			if err != nil {
				slog.Error(err.Error())
				return
			}

			_, err = tx.Exec("INSERT INTO migration (version,created_at) VALUES ($1, $2)", fileName, time.Now().UnixMilli())
			if err != nil {
				slog.Error(err.Error())
				return
			}

		}

		if fileName == migration.Version {
			gotLatest = true
		}
	}

	err = tx.Commit()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("success")
}

func (m *Migration) Down() {
	_, err := m.db.Exec(`CREATE TABLE IF NOT EXISTS migration (
		version VARCHAR(255) NOT NULL,
		created_at bigint NOT NULL
	)`)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	migration := &struct {
		Version   string `db:"version"`
		CreatedAt int64  `db:"created_at"`
	}{}

	err = m.db.Get(migration, "SELECT * FROM migration ORDER BY created_at DESC")
	if err == sql.ErrNoRows {
		slog.Info("no migration version to down to down")
		return
	}
	if err != nil {
		slog.Error(err.Error())
		return
	}

	files, err := m.migrateFS.ReadDir("migration")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	sort.SliceStable(files, func(i, j int) bool {
		numI := extractNumber(files[i].Name())
		numJ := extractNumber(files[j].Name())
		return numI < numJ
	})

	tx, err := m.db.Beginx()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer tx.Rollback()

	for _, v := range files {
		if !strings.Contains(v.Name(), ".down.sql") {
			continue
		}

		fileName := strings.ReplaceAll(v.Name(), ".down.sql", "")
		if fileName == migration.Version {

			fmt.Printf("migration/%s \n", v.Name())
			byteFile, err := m.migrateFS.ReadFile(fmt.Sprintf("migration/%s", v.Name()))
			if err != nil {
				slog.Error(err.Error())
				return
			}

			_, err = tx.Exec(string(byteFile))
			if err != nil {
				slog.Error(err.Error())
				return
			}

			_, err = tx.Exec("DELETE FROM migration WHERE version = $1", fileName)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			break
		}
	}

	err = tx.Commit()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("success")
}

func (m *Migration) Seed() {
	files, err := m.seedFS.ReadDir("seeder")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	for _, file := range files {
		fmt.Printf("seeder/%s \n", file.Name())
		b, err := m.seedFS.ReadFile(fmt.Sprintf("seeder/%s", file.Name()))
		if err != nil {
			slog.Error(err.Error())
			return
		}
		_, err = m.db.Exec(string(b))
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}

	fmt.Printf("success\n")
}

func (m *Migration) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migration",
		Short: "migration stuff",
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "create new migration",
		Run: func(cmd *cobra.Command, args []string) {
			m.Create(args)
		},
	}

	up := &cobra.Command{
		Use:   "up",
		Short: "up migration",
		Run: func(cmd *cobra.Command, args []string) {
			m.Up()
		},
	}

	down := &cobra.Command{
		Use:   "down",
		Short: "down migration",
		Run: func(cmd *cobra.Command, args []string) {
			m.Down()
		},
	}

	seed := &cobra.Command{
		Use:   "seed",
		Short: "seed migration",
		Run: func(cmd *cobra.Command, args []string) {
			m.Seed()
		},
	}

	cmd.AddCommand(create, up, down, seed)
	return cmd
}
