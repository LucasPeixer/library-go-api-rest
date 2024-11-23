package main

import (
	"bufio"
	"fmt"
	"go-api/db"
	"go-api/repository"
	"go-api/usecase"
	"go-api/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	DbDSN := os.Getenv("DB_DSN")
	if DbDSN == "" {
		log.Fatal("DB_DSN environment variable not set")
	}

	// Cria uma conexão com o banco de dados utilizando o DSN
	dbConn, err := db.CreateDB(DbDSN)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer dbConn.Close()

	userRepo := repository.NewUserRepository(dbConn)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Cria o scanner para input de dados
	scanner := bufio.NewScanner(os.Stdin)

	var i struct {
		Name     string
		Cpf      string
		Phone    string
		Email    string
		Password string
		RoleId   int
	}

	fmt.Print("Enter name: ")
	scanner.Scan()
	i.Name = scanner.Text()

	fmt.Print("Enter CPF (len=11): ")
	scanner.Scan()
	i.Cpf = scanner.Text()

	if !utils.IsValidCPF(i.Cpf) {
		log.Fatalf("Invalid CPF '%s'", i.Cpf)
	}

	fmt.Print("Enter phone: ")
	scanner.Scan()
	i.Phone = scanner.Text()

	fmt.Print("Enter email: ")
	scanner.Scan()
	i.Email = scanner.Text()

	fmt.Print("Enter password: ")
	scanner.Scan()
	i.Password = scanner.Text()

	fmt.Print("Enter role ID: ")
	scanner.Scan()
	_, err = fmt.Sscanf(scanner.Text(), "%d", &i.RoleId)
	if err != nil {
		log.Fatal("Invalid role ID")
	}

	// Faz o hashing da senha
	hashedPassword, err := utils.HashPassword(i.Password)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	// Cria um novo usuário com o useCase
	err = userUseCase.Register(i.Name, i.Cpf, i.Phone, i.Email, hashedPassword, i.RoleId)
	if err != nil {
		log.Fatalf("Error registering user: %v", err)
	}

	fmt.Println("User created successfully!")
}
