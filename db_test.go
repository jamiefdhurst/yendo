package yendo

import (
	"os"
	"testing"
)

func TestDsn_Format(t *testing.T) {
	// No password
	subject := &Dsn{Host: "test", Port: 1234, User: "user", Name: "db"}

	expected := "user:@tcp(test:1234)/db?timeout=1s"
	formatted := subject.format()
	if expected != formatted {
		t.Errorf("No password test failed, expected '%s' but received '%s'", expected, formatted)
	}

	// Default port
	subject = &Dsn{Host: "test", User: "user", Password: "pass1234", Name: "db"}

	expected = "user:pass1234@tcp(test:3306)/db?timeout=1s"
	formatted = subject.format()
	if expected != formatted {
		t.Errorf("Default port test failed, expected '%s' but received '%s'", expected, formatted)
	}

	// Standard
	subject = &Dsn{Host: "127.0.0.1", Port: 3307, User: "user", Password: "pass1234", Name: "db"}

	expected = "user:pass1234@tcp(127.0.0.1:3307)/db?timeout=1s"
	formatted = subject.format()
	if expected != formatted {
		t.Errorf("Expected '%s' but received '%s'", expected, formatted)
	}
}

func TestNewMySQL(t *testing.T) {
	var subject Database
	subject = NewMySQL(Dsn{})

	_, ok := subject.(Database)
	if !ok {
		t.Error("Expected an instance of Database")
	}
}

func TestConnect_ErrorOnPing(t *testing.T) {
	dsn := Dsn{Host: "not a real host", User: "nope", Name: "123"}
	subject := NewMySQL(dsn)

	err := subject.Connect()

	if err == nil {
		t.Error("Expected error from problematic connection string")
	}
}

func testConnection() *MySQL {
	return NewMySQL(Dsn{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	})
}

func TestConnectAndClose(t *testing.T) {
	subject := testConnection()

	err := subject.Connect()

	if err != nil {
		t.Errorf("Expected no error from connect but received '%s'", err)
	}

	err = subject.Close()

	if err != nil {
		t.Errorf("Expected no error from close but received '%s'", err)
	}

}

func TestExec(t *testing.T) {
	subject := testConnection()
	_ = subject.Connect()
	defer subject.Close()

	result, err := subject.Exec("SELECT 1")
	if err != nil {
		t.Errorf("Expected successful result but received error '%s'", err)
	}

	rows, _ := result.RowsAffected()
	if err != nil || rows > 0 {
		t.Error("Expected query to have been executed and no rows to have been affected")
	}
}

func TestQuery(t *testing.T) {
	subject := testConnection()
	_ = subject.Connect()
	defer subject.Close()

	rows, err := subject.Query("SELECT 1 AS example")
	if err != nil {
		t.Errorf("Expected query to have been executed but received error '%s'", err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	if len(columns) != 1 || columns[0] != "example" {
		t.Error("Expected column of 'example' to have been returned")
	}

	var test int
	for rows.Next() {
		rows.Scan(&test)
		if test != 1 {
			t.Error("Expected row with value of '1' to have been returned")
		}
	}
}

func TestQueryRow(t *testing.T) {
	subject := testConnection()
	_ = subject.Connect()
	defer subject.Close()

	var test int
	err := subject.QueryRow("SELECT 1 AS example").Scan(&test)
	if err != nil {
		t.Errorf("Expected query to have been executed but received error '%s'", err)
	}
	if test != 1 {
		t.Error("Expected row with value of '1' to have been returned")
	}
}
