package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

var instance *sql.DB

//go:embed migrations/*.sql
var embedMigrations embed.FS

func New() *sql.DB {

	// Evita criar mais de um banco
	if instance != nil {
		return instance
	}

	// Carrega as variáveis de ambiente definadas no arquivo .env na raiz do projeto.
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	// Para rodar os testes desse pacote, É necessário que haja um symlink de um .env nesta
	// pasta (database/) apontando para o .env no topo do projeto, já que 'go test' usa o
	// diretório raiz do pacote (./database) para executar os testes enquanto 'go run' usa o diretório
	// raiz do projeto.

	// variáveis ambientais de .env são carregadas para o ambiente do nosso processo, basta
	// usar a forma padrão para pegar os valores dessas variáveis.
	var (
		user     = os.Getenv("PG_USER")
		pass     = os.Getenv("PG_PASS")
		dbName   = os.Getenv("PG_DB")
		hostAddr = os.Getenv("PG_ADDR")
	)

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", user, pass, hostAddr, dbName)
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		panic(err)
	}

	// Verifica se a conexão foi um sucesso
	if err := db.Ping(); err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	// Aplica todas as migrações disponíveis
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	instance = db
	log.Println("Nova instância do banco de dados foi criada.")
	return instance
}

func MustBeginTx() *sql.Tx {
	if instance == nil {
		panic(fmt.Errorf("Banco de dados não foi inicializado. Não foi possível criar nova transação."))
	}
	tx, err := instance.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return tx
}
