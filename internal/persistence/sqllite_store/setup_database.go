package sqllite_store

import (
	"database/sql"
	"fmt"
)

// SetupDatabase initializes the SQLite database and creates the packages table if it doesn't exist
func SetupDatabase(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS packages (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT NOT NULL,
	    installed_version TEXT NOT NULL,
	    version TEXT NOT NULL,
	    installed BOOLEAN NOT NULL,
	    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}
