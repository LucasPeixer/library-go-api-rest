package initializers

import (
	"database/sql"
	"go-api/db"
	"log"
)

var DB *sql.DB

// InitDB inicializa a conexão com o banco de dados.
func InitDB() {
	if DbDSN == "" {
		log.Fatal("DB_DSN environment variable not set")
	}
	db_, err := db.CreateDB(DbDSN)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	DB = db_
}
