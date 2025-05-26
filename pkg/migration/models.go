package migration

// MigrationRecord represents a migration entry in the database
type MigrationRecord struct {
	Version   string `db:"version"`
	CreatedAt int64  `db:"created_at"`
}
