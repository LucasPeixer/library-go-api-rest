package initializers

import (
	"log"
	"os"
)

// DbDSN é a DSN (data source name) para conexão com banco de dados.
var DbDSN string

// JwtKey é a chave secreta para geração de chaves JWT.
var JwtKey []byte

// LoadEnv carrega as variáveis de ambiente necessárias.
func LoadEnv() {
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
