package yendo

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/jamiefdhurst/yendo/result"
)

const defaultFolder = "/src/github.com/jamiefdhurst/yendo/"
const testData = "/test/data/"

type MockErrorRow struct{}

func (m MockErrorRow) Scan(...interface{}) error {
	return errors.New("Simulating row error")
}

type MockDbFailsAllQueries struct {
	Queries int
}

func (m *MockDbFailsAllQueries) Connect() error { return nil }
func (m *MockDbFailsAllQueries) Exec(sql string, args ...interface{}) (sql.Result, error) {
	m.Queries++
	return nil, errors.New("Simulating error")
}
func (m *MockDbFailsAllQueries) Query(sql string, args ...interface{}) (result.Rows, error) {
	m.Queries++
	return nil, errors.New("Simulating error")
}
func (m *MockDbFailsAllQueries) QueryRow(sql string, args ...interface{}) result.Row {
	m.Queries++
	return new(MockErrorRow)
}
func (m *MockDbFailsAllQueries) Close() error { return nil }

func TestMigrate_ErrorWhenCreateTableFails(t *testing.T) {
	mockDb := &MockDbFailsAllQueries{}

	subject := NewMigrator(mockDb, "")

	err := subject.Migrate()
	if err == nil {
		t.Error("Expected error from creating the new table")
	}
	if mockDb.Queries != 1 {
		t.Errorf("Expected 1 query to have run, received: %d", mockDb.Queries)
	}
}

type MockDbQueryError struct{}

func (m *MockDbQueryError) Connect() error { return nil }
func (m *MockDbQueryError) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}
func (m *MockDbQueryError) Query(sql string, args ...interface{}) (result.Rows, error) {
	return nil, errors.New("Simulating error")
}
func (m *MockDbQueryError) QueryRow(sql string, args ...interface{}) result.Row {
	return nil
}
func (m *MockDbQueryError) Close() error { return nil }

type MockEmptyRows struct{}

func (m MockEmptyRows) Close() error                   { return nil }
func (m MockEmptyRows) Columns() ([]string, error)     { return nil, nil }
func (m MockEmptyRows) Next() bool                     { return false }
func (m MockEmptyRows) Scan(dest ...interface{}) error { return nil }

func TestMigrate_ErrorWhenAvailableMigrationsFails(t *testing.T) {
	mockDb := &MockDbQueryError{}

	subject := Migrator{mockDb, "non-existant-folder", "migration"}

	err := subject.Migrate()
	if err == nil {
		t.Error("Expected error from walking the available migrations folder")
	}
}

func TestMigrate_ErrorWhenPreviousMigrationsFails(t *testing.T) {
	mockDb := &MockDbQueryError{}

	subject := NewMigrator(mockDb, "")

	err := subject.Migrate()
	if err == nil {
		t.Error("Expected error from calling previous migrations")
	}
}

type MockDbErrorOnExecOnly struct{}

func (m *MockDbErrorOnExecOnly) Connect() error { return nil }
func (m *MockDbErrorOnExecOnly) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return nil, errors.New("Simulated error")
}
func (m *MockDbErrorOnExecOnly) Query(sql string, args ...interface{}) (result.Rows, error) {
	return MockEmptyRows{}, nil
}
func (m *MockDbErrorOnExecOnly) QueryRow(sql string, args ...interface{}) result.Row {
	return nil
}
func (m *MockDbErrorOnExecOnly) Close() error { return nil }

type MockDbErrorOnSecondExec struct {
	execs int
}

func (m *MockDbErrorOnSecondExec) Connect() error { return nil }
func (m *MockDbErrorOnSecondExec) Exec(sql string, args ...interface{}) (sql.Result, error) {
	m.execs++
	if m.execs == 2 {
		return nil, errors.New("Simulated error")
	}
	return nil, nil
}
func (m *MockDbErrorOnSecondExec) Query(sql string, args ...interface{}) (result.Rows, error) {
	return MockEmptyRows{}, nil
}
func (m *MockDbErrorOnSecondExec) QueryRow(sql string, args ...interface{}) result.Row {
	return nil
}
func (m *MockDbErrorOnSecondExec) Close() error { return nil }

func testMigrator(dbConnection Database, file string) Migrator {
	folder := "/" + strings.TrimLeft(os.Getenv("DB_MIGRATION_FOLDER"), "/")
	if folder == "/" {
		folder = os.Getenv("GOPATH") + defaultFolder
	}
	return Migrator{dbConnection, folder + testData + file, "migration"}
}

func TestMigrate_ErrorWhenReadingFile(t *testing.T) {
	mockDb := &MockDbErrorOnSecondExec{}

	subject := testMigrator(mockDb, "sql-file-error")

	err := subject.Migrate()
	if err == nil {
		t.Error("Expected error from reading files")
	}
}

func TestMigrate_ErrorWhenExecutingQueryFromFile(t *testing.T) {
	mockDb := &MockDbErrorOnSecondExec{}

	subject := testMigrator(mockDb, "sql-single")

	err := subject.Migrate()
	if err == nil {
		t.Error("Expected error from executing the first file found")
	}
}

type MockDbErrorOnThirdExec struct {
	execs int
}

func (m *MockDbErrorOnThirdExec) Connect() error { return nil }
func (m *MockDbErrorOnThirdExec) Exec(sql string, args ...interface{}) (sql.Result, error) {
	m.execs++
	if m.execs == 3 {
		return nil, errors.New("Simulated error")
	}
	return nil, nil
}
func (m *MockDbErrorOnThirdExec) Query(sql string, args ...interface{}) (result.Rows, error) {
	return MockEmptyRows{}, nil
}
func (m *MockDbErrorOnThirdExec) QueryRow(sql string, args ...interface{}) result.Row {
	return nil
}
func (m *MockDbErrorOnThirdExec) Close() error { return nil }

func TestMigrate_ErrorWhenExecutingInsertAfterQueryFromFile(t *testing.T) {
	mockDb := &MockDbErrorOnThirdExec{}

	subject := testMigrator(mockDb, "sql-single")

	err := subject.Migrate()
	if err == nil {
		t.Error("Expected error from executing the insert statement once the query has been susccessful")
	}
}

func fixturesEmpty(t *testing.T) {
	// Drop any existing migrations
	dbConnection := testConnection()
	dbConnection.Connect()
	dbConnection.Exec("DROP TABLE IF EXISTS `migration`")
	dbConnection.Exec("DROP TABLE IF EXISTS `test1`")
	dbConnection.Exec("DROP TABLE IF EXISTS `test2`")
	dbConnection.Close()
}

func TestMigrate_SuccessForNewSingleMigration(t *testing.T) {
	fixturesEmpty(t)
	dbConnection := testConnection()
	dbConnection.Connect()
	defer dbConnection.Close()

	subject := testMigrator(dbConnection, "sql-single")

	err := subject.Migrate()
	if err != nil {
		t.Errorf("Expected no error but received: %s", err)
	}
	var migrations int
	dbConnection.QueryRow("SELECT COUNT(`migration`) FROM `migration`").Scan(&migrations)
	if migrations != 1 {
		t.Errorf("Expected 1 migration but received %d", migrations)
	}
	var table string
	dbConnection.QueryRow("SHOW TABLES LIKE 'test1'").Scan(&table)
	if table != "test1" {
		t.Errorf("Expected table test1 to be available but received '%s'", table)
	}
}

func TestMigrate_SuccessForNewMultipleMigrations(t *testing.T) {
	fixturesEmpty(t)
	dbConnection := testConnection()
	dbConnection.Connect()
	defer dbConnection.Close()

	subject := testMigrator(dbConnection, "sql-multiple")

	err := subject.Migrate()
	if err != nil {
		t.Errorf("Expected no error but received: %s", err)
	}
	var migrations int
	dbConnection.QueryRow("SELECT COUNT(`migration`) FROM `migration`").Scan(&migrations)
	if migrations != 2 {
		t.Errorf("Expected 2 migrations but received %d", migrations)
	}
	var table string
	dbConnection.QueryRow("SHOW TABLES LIKE 'test2'").Scan(&table)
	if table != "test2" {
		t.Errorf("Expected table test2 to be available but received '%s'", table)
	}
}

func fixturesExistingMigration(t *testing.T) {
	// Drop any existing migrations
	dbConnection := testConnection()
	dbConnection.Connect()
	dbConnection.Exec("CREATE TABLE IF NOT EXISTS `migration` (`migration` VARCHAR(200) NOT NULL PRIMARY KEY, `timestamp` DATETIME NOT NULL)")
	dbConnection.Exec("TRUNCATE TABLE `migration`")
	dbConnection.Exec("INSERT INTO `migration` VALUES ('002.sql', '2019-01-01 00:00:00')")
	dbConnection.Exec("CREATE TABLE IF NOT EXISTS `test1` (`id` INT NOT NULL PRIMARY KEY)")
	dbConnection.Exec("DROP TABLE IF EXISTS `test2`")
	dbConnection.Close()
}

func TestMigrate_SuccessForExistingAndNewSingleMigrations(t *testing.T) {
	fixturesExistingMigration(t)
	dbConnection := testConnection()
	dbConnection.Connect()
	defer dbConnection.Close()

	subject := testMigrator(dbConnection, "sql-single")

	err := subject.Migrate()
	if err != nil {
		t.Errorf("Expected no error but received: %s", err)
	}
	var migrations int
	dbConnection.QueryRow("SELECT COUNT(`migration`) FROM `migration`").Scan(&migrations)
	if migrations != 1 {
		t.Errorf("Expected 1 migration but received %d", migrations)
	}
}

func TestMigrate_SuccessForExistingAndNewMultipleMigrations(t *testing.T) {
	fixturesExistingMigration(t)
	dbConnection := testConnection()
	dbConnection.Connect()
	defer dbConnection.Close()

	subject := testMigrator(dbConnection, "sql-multiple")

	err := subject.Migrate()
	if err != nil {
		t.Errorf("Expected no error but received: %s", err)
	}
	var migrations int
	dbConnection.QueryRow("SELECT COUNT(`migration`) FROM `migration`").Scan(&migrations)
	if migrations != 2 {
		t.Errorf("Expected 2 migrations but received %d", migrations)
	}
	var table string
	dbConnection.QueryRow("SHOW TABLES LIKE 'test2'").Scan(&table)
	if table != "test2" {
		t.Errorf("Expected table test2 to be available but received '%s'", table)
	}
}
