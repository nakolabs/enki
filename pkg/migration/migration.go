package migration

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
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
		log.Error().Err(err).Msg("failed to read migration directory")
		return
	}

	migrationNumber := len(entries)/2 + 1
	up := fmt.Sprintf("migration/%d_%s.up.sql", migrationNumber, args[0])
	log.Info().Str("file", up).Msg("creating migration file")

	err = os.WriteFile(up, nil, 0755)
	if err != nil {
		log.Error().Err(err).Str("file", up).Msg("failed to create up migration file")
		return
	}

	down := fmt.Sprintf("migration/%d_%s.down.sql", migrationNumber, args[0])
	log.Info().Str("file", down).Msg("creating migration file")

	err = os.WriteFile(down, nil, 0755)
	if err != nil {
		log.Error().Err(err).Str("file", down).Msg("failed to create down migration file")
		return
	}

	log.Info().Str("migration_name", args[0]).Int("number", migrationNumber).Msg("migration created successfully")
}

func (m *Migration) Up() error {
	_, err := m.db.Exec(`CREATE TABLE IF NOT EXISTS migration (
		version VARCHAR(255) NOT NULL,
		created_at bigint NOT NULL
	)`)
	if err != nil {
		log.Error().Err(err).Msg("failed to create migration table")
		return err
	}

	migration := &MigrationRecord{}

	err = m.db.Get(migration, "SELECT * FROM migration ORDER BY created_at DESC")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error().Err(err).Msg("failed to get latest migration")
		return err
	}

	gotLatest := false
	if migration.CreatedAt == 0 {
		gotLatest = true
		log.Info().Msg("no previous migrations found, running all migrations")
	} else {
		log.Info().Str("latest_version", migration.Version).Msg("found latest migration")
	}

	files, err := m.migrateFS.ReadDir("migration")
	if err != nil {
		log.Error().Err(err).Msg("failed to read migration directory")
		return err
	}

	sort.SliceStable(files, func(i, j int) bool {
		numI := extractNumber(files[i].Name())
		numJ := extractNumber(files[j].Name())
		return numI < numJ
	})

	tx, err := m.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return err
	}
	defer tx.Rollback()

	migrationsRun := 0
	for i := 0; i < len(files); i++ {
		fName := files[i].Name()
		if !strings.Contains(fName, ".up.sql") {
			continue
		}

		fileName := strings.ReplaceAll(fName, ".up.sql", "")
		if gotLatest {
			log.Info().Str("file", fName).Str("version", fileName).Msg("executing migration")

			byteFile, err := m.migrateFS.ReadFile(fmt.Sprintf("migration/%s", fName))
			if err != nil {
				log.Error().Err(err).Str("file", fName).Msg("failed to read migration file")
				return err
			}

			_, err = tx.Exec(string(byteFile))
			if err != nil {
				log.Error().Err(err).Str("file", fName).Msg("failed to execute migration")
				return err
			}

			_, err = tx.Exec("INSERT INTO migration (version,created_at) VALUES ($1, $2)", fileName, time.Now().UnixMilli())
			if err != nil {
				log.Error().Err(err).Str("version", fileName).Msg("failed to insert migration record")
				return err
			}

			migrationsRun++
		}

		if fileName == migration.Version {
			gotLatest = true
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("failed to commit migration transaction")
		return err
	}

	log.Info().Int("migrations_executed", migrationsRun).Msg("migration completed successfully")
	return nil
}

func (m *Migration) Down() error {
	_, err := m.db.Exec(`CREATE TABLE IF NOT EXISTS migration (
		version VARCHAR(255) NOT NULL,
		created_at bigint NOT NULL
	)`)
	if err != nil {
		log.Error().Err(err).Msg("failed to create migration table")
		return err
	}

	migration := &MigrationRecord{}

	err = m.db.Get(migration, "SELECT * FROM migration ORDER BY created_at DESC")
	if err == sql.ErrNoRows {
		log.Info().Msg("no migration version to rollback")
		return nil
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get latest migration")
		return err
	}

	log.Info().Str("target_version", migration.Version).Msg("rolling back migration")

	files, err := m.migrateFS.ReadDir("migration")
	if err != nil {
		log.Error().Err(err).Msg("failed to read migration directory")
		return err
	}

	sort.SliceStable(files, func(i, j int) bool {
		numI := extractNumber(files[i].Name())
		numJ := extractNumber(files[j].Name())
		return numI < numJ
	})

	tx, err := m.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return err
	}
	defer tx.Rollback()

	for _, v := range files {
		if !strings.Contains(v.Name(), ".down.sql") {
			continue
		}

		fileName := strings.ReplaceAll(v.Name(), ".down.sql", "")
		if fileName == migration.Version {
			log.Info().Str("file", v.Name()).Str("version", fileName).Msg("executing migration rollback")

			byteFile, err := m.migrateFS.ReadFile(fmt.Sprintf("migration/%s", v.Name()))
			if err != nil {
				log.Error().Err(err).Str("file", v.Name()).Msg("failed to read migration file")
				return err
			}

			_, err = tx.Exec(string(byteFile))
			if err != nil {
				log.Error().Err(err).Str("file", v.Name()).Msg("failed to execute migration rollback")
				return err
			}

			_, err = tx.Exec("DELETE FROM migration WHERE version = $1", fileName)
			if err != nil {
				log.Error().Err(err).Str("version", fileName).Msg("failed to delete migration record")
				return err
			}

			break
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("failed to commit migration rollback")
		return err
	}

	log.Info().Str("rolled_back_version", migration.Version).Msg("migration rollback completed successfully")
	return nil
}

func (m *Migration) Seed() {
	files, err := m.seedFS.ReadDir("seeder")
	if err != nil {
		log.Error().Err(err).Msg("failed to read seeder directory")
		return
	}

	const defaultPassword = "genesis"
	pass, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash default password")
		return
	}

	fmt.Println("Using default password for seeding:", string(pass))

	log.Info().Int("seed_files", len(files)).Msg("starting database seeding")

	for _, file := range files {
		log.Info().Str("file", file.Name()).Msg("executing seed file")

		b, err := m.seedFS.ReadFile(fmt.Sprintf("seeder/%s", file.Name()))
		if err != nil {
			log.Error().Err(err).Str("file", file.Name()).Msg("failed to read seed file")
			return
		}

		_, err = m.db.Exec(string(b))
		if err != nil {
			log.Error().Err(err).Str("file", file.Name()).Msg("failed to execute seed file")
			return
		}
	}

	// m.db.Exec()

	log.Info().Int("seed_files_executed", len(files)).Msg("database seeding completed successfully")
}

func (m *Migration) Fresh() error {
	log.Info().Msg("starting fresh migration - dropping all tables and types")

	// Get all table names
	rows, err := m.db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
	`)
	if err != nil {
		log.Error().Err(err).Msg("failed to get table names")
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Error().Err(err).Msg("failed to scan table name")
			return err
		}
		tables = append(tables, tableName)
	}

	// Get all custom types
	typeRows, err := m.db.Query(`
		SELECT typname 
		FROM pg_type 
		WHERE typnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'public')
		AND typtype = 'e'
	`)
	if err != nil {
		log.Error().Err(err).Msg("failed to get type names")
		return err
	}
	defer typeRows.Close()

	var types []string
	for typeRows.Next() {
		var typeName string
		if err := typeRows.Scan(&typeName); err != nil {
			log.Error().Err(err).Msg("failed to scan type name")
			return err
		}
		types = append(types, typeName)
	}

	tx, err := m.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return err
	}
	defer tx.Rollback()

	// Drop all tables
	if len(tables) > 0 {
		for _, table := range tables {
			log.Info().Str("table", table).Msg("dropping table")
			_, err = tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
			if err != nil {
				log.Error().Err(err).Str("table", table).Msg("failed to drop table")
				return err
			}
		}
		log.Info().Int("tables_dropped", len(tables)).Msg("all tables dropped successfully")
	} else {
		log.Info().Msg("no tables found to drop")
	}

	// Drop all custom types
	if len(types) > 0 {
		for _, typeName := range types {
			log.Info().Str("type", typeName).Msg("dropping type")
			_, err = tx.Exec(fmt.Sprintf("DROP TYPE IF EXISTS %s CASCADE", typeName))
			if err != nil {
				log.Error().Err(err).Str("type", typeName).Msg("failed to drop type")
				return err
			}
		}
		log.Info().Int("types_dropped", len(types)).Msg("all types dropped successfully")
	} else {
		log.Info().Msg("no custom types found to drop")
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit drop operation transaction")
		return err
	}

	// Re-run migrations
	return m.Up()
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
			err := m.Up()
			if err != nil {
				log.Error().Err(err).Msg("migration up failed")
			}
		},
	}

	down := &cobra.Command{
		Use:   "down",
		Short: "down migration",
		Run: func(cmd *cobra.Command, args []string) {
			err := m.Down()
			if err != nil {
				log.Error().Err(err).Msg("migration down failed")
			}
		},
	}

	seed := &cobra.Command{
		Use:   "seed",
		Short: "seed migration",
		Run: func(cmd *cobra.Command, args []string) {
			m.Seed()
		},
	}

	fresh := &cobra.Command{
		Use:   "fresh",
		Short: "fresh migration",
		Run: func(cmd *cobra.Command, args []string) {
			m.Fresh()
		},
	}

	cmd.AddCommand(create, up, down, seed, fresh)
	return cmd
}
