package main

import (
	"log"
	"os"
	"strconv"

	"github.com/jamiefdhurst/yendo"
)

func main() {

	// Set up the database connection
	connDetails := yendo.Dsn{
		Host:     os.Getenv("DB_HOST"),
		Port:     3306,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
	conn := yendo.NewMySQL(connDetails)
	conn.Connect()
	defer conn.Close()
	log.Println("Connected to database...")

	// Run the migrations
	migration := yendo.NewMigrator(conn, "/src/github.com/jamiefdhurst/yendo/example/sql/")
	migration.Migrate()

	log.Println("Migrations complete...")

	// Interact with the database to prove the migrations
	log.Println("Inserting data...")
	conn.Exec("INSERT INTO `foobar` VALUES (NULL, '1234', '567890', 'abcdef')")
	conn.Exec("INSERT INTO `foobar` VALUES (NULL, '5678', '901234', 'ghijkl')")
	conn.Exec("INSERT INTO `foobar` VALUES (NULL, '9012', '345678', 'mnopqr')")

	log.Println("Querying the database...")

	var count int
	conn.QueryRow("SELECT COUNT(*) FROM `foobar`").Scan(&count)
	log.Println("Result returned from the query: " + strconv.Itoa(count))

	log.Println("Example application run complete.")
}
