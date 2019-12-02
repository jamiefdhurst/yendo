package yendo

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const defaultMigrateTable = "migration"

// Migrator Handle migrations
type Migrator struct {
	db     Database
	Folder string
	Table  string
}

// NewMigrator Return a new Migrator instance
func NewMigrator(db Database, folder string) *Migrator {
	folder = strings.TrimLeft(folder, "/")
	return &Migrator{db, os.Getenv("GOPATH") + "/" + folder, defaultMigrateTable}
}

// Migrate the database to the latest version
func (m Migrator) Migrate() error {
	var err error

	// Ensure migration table is present
	if err = m.createTable(); err != nil {
		return err
	}

	// Collect migration scripts from folder
	allMigrations, err := m.available()
	if err != nil {
		return err
	}

	// Remove migrations that have been run
	currentMigrations, err := m.previous()
	if err != nil {
		return err
	}
	migrations := m.diff(allMigrations, currentMigrations)

	// Iterate and run all migrations
	for _, migration := range migrations {
		sqlBuf, _ := ioutil.ReadFile(path.Join(m.Folder, migration))
		sql := strings.TrimSpace(string(sqlBuf))
		if len(sql) == 0 {
			return errors.New("Could not read file: " + migration)
		}

		_, err = m.db.Exec(sql)
		if err != nil {
			return err
		}

		// Insert migration entry
		_, err = m.db.Exec("INSERT INTO `" + m.Table + "` VALUES ('" + migration + "', CURRENT_TIMESTAMP())")
		if err != nil {
			return err
		}
	}

	return nil
}

func (m Migrator) createTable() error {
	_, err := m.db.Exec("CREATE TABLE IF NOT EXISTS `" + m.Table + "` (" +
		"`migration` VARCHAR(200) NOT NULL PRIMARY KEY, " +
		"`timestamp` DATETIME NOT NULL " +
		")")

	return err
}

func (m Migrator) available() (allMigrations []string, err error) {
	err = filepath.Walk(m.Folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".sql" {
			allMigrations = append(allMigrations, info.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return allMigrations, nil
}

func (m Migrator) previous() (currentMigrations []string, err error) {
	rows, err := m.db.Query("SELECT `migration` FROM `" + m.Table + "` ORDER BY `timestamp`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var migration string
	for rows.Next() {
		rows.Scan(&migration)
		currentMigrations = append(currentMigrations, migration)
	}

	return currentMigrations, nil
}

func (m Migrator) diff(all []string, current []string) (diff []string) {
	d := make(map[string]bool)

	for _, m := range current {
		d[m] = true
	}

	for _, m := range all {
		if _, ok := d[m]; !ok {
			diff = append(diff, m)
		}
	}

	return
}
