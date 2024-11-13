package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DbDSN é a DSN (data source name) para conexão com banco de dados.
var DbDSN string

// JwtKey é a chave secreta para geração de chaves JWT.
var JwtKey []byte

// LoadEnv carrega as variáveis de ambiente necessárias.
func LoadEnv() {
	// Carrega as variáveis do arquivo .env se existir
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	DbDSN = os.Getenv("DB_DSN")
	if DbDSN == "" {
		log.Fatal("DB_DSN environment variable not set")
	}

	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		log.Fatal("JWT_KEY environment variable not set")
	}
	JwtKey = []byte(jwtKey)
}
