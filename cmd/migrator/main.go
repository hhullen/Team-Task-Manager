package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"team-task-manager/internal/clients/mysql"
	secretprovider "team-task-manager/internal/secret_provider"
	"team-task-manager/internal/supports"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	migrationsDir = "./migrations/mysql"
	cmdUp         = "up"
	cmdDown       = "down"
	cmdStatus     = "status"
)

var executors = map[string]func(db *sql.DB, dir string, opts ...goose.OptionsFunc) error{
	cmdUp:     goose.Up,
	cmdDown:   goose.Down,
	cmdStatus: goose.Status,
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			log.Fatalf("%v", p)
		}
	}()

	if len(os.Args) < 2 {
		log.Fatalf("migration action is required: %s/%s/%s", cmdUp, cmdDown, cmdStatus)
	}
	command := os.Args[1]

	goose.SetBaseFS(nil)

	Migrate(command)
}

func Migrate(command string) {
	ctx := context.Background()
	defer ctx.Done()

	secrets := secretprovider.NewSecretProvider()

	host, err := secrets.ReadSecret("db_host")
	if err != nil {
		panic(err)
	}
	port, err := secrets.ReadSecret("db_port")
	if err != nil {
		panic(err)
	}
	user, err := secrets.ReadSecret("db_migrator_user")
	if err != nil {
		panic(err)
	}
	password, err := secrets.ReadSecret("db_migrator_password")
	if err != nil {
		panic(err)
	}
	name, err := secrets.ReadSecret("db_name")
	if err != nil {
		panic(err)
	}

	if !supports.IsInContainer() {
		host = "localhost"
	}

	db, err := mysql.NewMySQLConn(ctx, host, port, user, password, name)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("mysql"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	ExecMigration(db, command, migrationsDir)

}

func ExecMigration(db *sql.DB, command, migrationsDir string) {
	executor, exists := executors[command]
	if !exists {
		log.Fatalf("Wrong comand send: %s. Required: %s/%s/%s", command, cmdUp, cmdDown, cmdStatus)
	}

	if err := executor(db, migrationsDir); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Printf("%s successfully migrated\n", migrationsDir)
}
