package migrator

import (
	"context"
	"database/sql"
	"log"
	"team-task-manager/internal/clients/mysql"
	secretprovider "team-task-manager/internal/secret_provider"
	"team-task-manager/internal/supports"

	"github.com/pressly/goose/v3"
)

const (
	migrationsDir = "./migrations/mysql"
	cmdUp         = "up"
	cmdDown       = "down"
	cmdDownAll    = "down-all"
	cmdStatus     = "status"

	dbDialect = "mysql"

	defaultSecretsDir          = "./secrets/"
	defaultContainerSecretsDir = "/run/secrets/"
)

var executors = map[string]func(db *sql.DB, dir string, opts ...goose.OptionsFunc) error{
	cmdUp:      goose.Up,
	cmdDown:    goose.Down,
	cmdStatus:  goose.Status,
	cmdDownAll: downAll,
}

func Migrate(command string) {
	ctx := context.Background()
	defer ctx.Done()

	secretDir := defaultSecretsDir
	if supports.IsInContainer() {
		secretDir = defaultContainerSecretsDir
	}

	secrets := secretprovider.NewSecretProvider(secretDir)

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
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	ExecMigration(db, command, migrationsDir)

}

func ExecMigration(db *sql.DB, command, migrationsDir string) {
	executor, exists := executors[command]
	if !exists {
		log.Fatalf("Wrong comand send: %s. Required: %s/%s/%s/%s", command, cmdUp, cmdDown, cmdDownAll, cmdStatus)
	}

	goose.SetBaseFS(nil)
	if err := goose.SetDialect(dbDialect); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if err := executor(db, migrationsDir); err != nil {
		log.Fatalf("failed applying migrations: %v", err)
	}

	log.Printf("%s successfully migrated\n", migrationsDir)
}

func downAll(db *sql.DB, dir string, opts ...goose.OptionsFunc) error {
	return goose.DownTo(db, dir, 0, opts...)
}
