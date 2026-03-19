package main

import (
	"log"
	"os"
	"team-task-manager/internal/migrator"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	defaultSecretsDir          = "./secrets/"
	defaultContainerSecretsDir = "/run/secrets/"
)

func main() {
	defer func() {
		if p := recover(); p != nil {
			log.Fatalf("%v", p)
		}
	}()

	if len(os.Args) < 2 {
		log.Fatalf("migration action is required")
	}
	command := os.Args[1]

	goose.SetBaseFS(nil)

	migrator.Migrate(command)
}
