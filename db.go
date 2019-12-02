package yendo

import (
	"database/sql"
	"errors"
	"fmt"

	// MySQL driver required
	_ "github.com/go-sql-driver/mysql"
	"github.com/jamiefdhurst/yendo/result"
)

const dsnFormat = "%s:%s@tcp(%s:%d)/%s?timeout=1s"
const defaultPort = 3306

// Database Define a common interface for all database drivers
type Database interface {
	Close() error
	Connect() error
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Query(sql string, args ...interface{}) (result.Rows, error)
	QueryRow(sql string, args ...interface{}) result.Row
}

type dsnFormatter interface {
	Format() string
}

// Dsn Connection details wrapper
type Dsn struct {
	dsnFormatter
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func (d *Dsn) format() string {
	if d.Port == 0 {
		d.Port = defaultPort
	}

	return fmt.Sprintf(dsnFormat, d.User, d.Password, d.Host, d.Port, d.Name)
}

// MySQL Database connection object
type MySQL struct {
	Database
	dsn        Dsn
	connection *sql.DB
}

// NewMySQL Create a new MySQL entry with default requirements
func NewMySQL(dsn Dsn) *MySQL {
	return &MySQL{dsn: dsn}
}

// Connect Start the connection
func (m *MySQL) Connect() error {
	var err error

	m.connection, _ = sql.Open("mysql", m.dsn.format())
	if err = m.connection.Ping(); err != nil {
		return errors.New("Unable to connect to database with provided credentials")
	}

	return nil
}

// Exec Execute a query on the database, returning a simple result
func (m *MySQL) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return m.connection.Exec(sql, args...)
}

// Query Query the database
func (m *MySQL) Query(sql string, args ...interface{}) (result.Rows, error) {
	return m.connection.Query(sql, args...)
}

// QueryRow Query the database for a single row only
func (m *MySQL) QueryRow(sql string, args ...interface{}) result.Row {
	return m.connection.QueryRow(sql, args...)
}

// Close Terminate the current connection (should be used with defer)
func (m *MySQL) Close() error {
	return m.connection.Close()
}
