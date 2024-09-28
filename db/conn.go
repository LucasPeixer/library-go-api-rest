package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	connect := "host=vividly-climbing-mallard.data-1.use1.tembo.io port=5432 user=postgres password=8fjHOEvlIsib4fhH dbname=armazem_DB sslmode=require"

	db, err := sql.Open("postgres", connect)

	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao pingar o banco de dados: ", err)
	}

	log.Println("Conexão estabelecida com sucesso!")

	return db, nil
}