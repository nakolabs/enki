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

	log.Info().Int("seed_files_executed", len(files)).Msg("database seeding completed successfully")
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

	cmd.AddCommand(create, up, down, seed)
	return cmd
}
