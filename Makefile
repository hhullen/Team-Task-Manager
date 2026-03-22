PROJECT_NAME=team-task-manager

PWD=$(pwd)
RM=rm -rf
EXTENSION=

MOCKGEN_BIN=$(shell go env GOPATH)/bin/mockgen$(EXTENSION)
SQLC_BIN=$(shell go env GOPATH)/bin/sqlc$(EXTENSION)
SWAG_BIN=$(shell go env GOPATH)/bin/swag$(EXTENSION)
SWAG_DOCS_DIR=internal/docs
SHADOW_BIN=$(shell go env GOPATH)/bin/shadow$(EXTENSION)
ERRCHECK_BIN=$(shell go env GOPATH)/bin/errcheck$(EXTENSION)

COVERAGE_FILE=coverage.out
NOT_FILTERED_SUFF=_not_filtered
FILTER_COVERAGE_FROM_MOCK=grep -v "mock" $(COVERAGE_FILE)$(NOT_FILTERED_SUFF) | grep -v "sqlc" > $(COVERAGE_FILE)

SERVICE_DATASTRUCT_DIR=internal/datastruct
API_DIR=internal/api/v1

MIGRATOR_DIR=./cmd/migrator
MIGRATOR_BIN=$(MIGRATOR_DIR)/migrator$(EXTENSION)

TEAM_TASK_MANAGER_DIR=./cmd/$(PROJECT_NAME)
TEAM_TASK_MANAGER_BIN=$(TEAM_TASK_MANAGER_DIR)/$(PROJECT_NAME)$(EXTENSION)

LOCAL_DB_NAME=$(PROJECT_NAME)-local-database
LOCAL_DB_DATA_NAME=local_$(PROJECT_NAME)_mysql_data
LOCAL_REDIS_NAME=$(PROJECT_NAME)-local-redis
LOCAL_REDIS_DATA_NAME=local_$(PROJECT_NAME)_redis_data

ifeq ($(OS),Windows_NT)
	SHELL=powershell.exe
	EXTENSION=.exe
	PWD=$(shell powershell -Command "(Get-Location).Path")
	RM=echo
	RM_POSTFIX=| Remove-Item -Force -ErrorAction SilentlyContinue; exit 0
	FILTER_COVERAGE_FROM_MOCK=(Get-Content $(COVERAGE_FILE)$(NOT_FILTERED_SUFF)) | Where-Object { $$_ -notmatch "mock" } | Where-Object { $$_ -notmatch "sqlc" } | Set-Content $(COVERAGE_FILE)
endif


.PHONY: deps errcheck linter generate-sqlc generate-mocks generage-swag migrations-up migrations-down migrations-status start-local-database stop-local-database clean-local-database start-local-redis stop-local-redis clean-local-redis run-local run-local-fast service service-stop service-rebuild coverage-info coverage-html clean-go-cache clean

deps:
	go mod download

$(ERRCHECK_BIN):
	go install -a github.com/kisielk/errcheck@latest

errcheck: $(ERRCHECK_BIN)
	errcheck.exe -verbose -ignoregenerated ./...

$(SHADOW_BIN):
	go install -a golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest

linter: $(SHADOW_BIN)
	go vet -vettool=$(SHADOW_BIN) ./...

$(SQLC_BIN):
	go install -a github.com/sqlc-dev/sqlc/cmd/sqlc@latest

generate-sqlc: $(SQLC_BIN)
	sqlc generate

$(MOCKGEN_BIN):
	go install -a github.com/golang/mock/mockgen@latest

generate-mocks: $(MOCKGEN_BIN)
	go generate -run "mockgen" ./...

$(SWAG_BIN):
	go install -a github.com/swaggo/swag/cmd/swag@latest

generate-swag: $(SWAG_BIN)
	swag init -d $(TEAM_TASK_MANAGER_DIR),$(SERVICE_DATASTRUCT_DIR),$(API_DIR) -o $(SWAG_DOCS_DIR)

$(MIGRATOR_BIN):
	go build -o $(MIGRATOR_BIN) $(MIGRATOR_DIR)

migrations-up: $(MIGRATOR_BIN)
	$(MIGRATOR_BIN) up

migrations-down: $(MIGRATOR_BIN)
	$(MIGRATOR_BIN) down

migrations-down-all: $(MIGRATOR_BIN)
	$(MIGRATOR_BIN) down-all

migrations-status: $(MIGRATOR_BIN)
	$(MIGRATOR_BIN) status

DB_ENV+= -e MYSQL_ROOT_PASSWORD_FILE=/run/secrets/db_root_password 
DB_ENV+= -e MYSQL_DATABASE_FILE=/run/secrets/db_name

DB_VOL+= -v $(PWD)/secrets/db_root_password:/run/secrets/db_root_password:ro
DB_VOL+= -v $(PWD)/secrets/db_app_password:/run/secrets/db_app_password:ro
DB_VOL+= -v $(PWD)/secrets/db_app_user:/run/secrets/db_app_user:ro
DB_VOL+= -v $(PWD)/secrets/db_migrator_password:/run/secrets/db_migrator_password:ro
DB_VOL+= -v $(PWD)/secrets/db_migrator_user:/run/secrets/db_migrator_user:ro
DB_VOL+= -v $(PWD)/secrets/db_name:/run/secrets/db_name:ro
DB_VOL+= -v $(LOCAL_DB_DATA_NAME):/var/lib/mysql
DB_VOL+= -v $(PWD)/init/init.sh:/docker-entrypoint-initdb.d/init.sh:ro
DB_VOL+= -v $(PWD)/init/init_roles.sql:/sql_init/init_roles.sql:ro

start-local-database:
	docker run -d --rm -p 3306:3306 $(DB_ENV) $(DB_VOL) --name $(LOCAL_DB_NAME) mysql:8.0.45-debian

stop-local-database:
	docker container stop $(LOCAL_DB_NAME)

clean-local-database:
	docker volume rm $(LOCAL_DB_DATA_NAME)

REDIS_VOL+= -v $(LOCAL_REDIS_DATA_NAME):/data
REDIS_VOL+= -v  $(PWD)/secrets/redis_password:/run/secrets/redis_password:ro
REDIS_CMD=sh -c redis-server --requirepass $(cat /run/secrets/redis_password) --appendonly yes

start-local-redis:
	docker run -d --rm -p 6379:6379 $(REDIS_VOL) --name $(LOCAL_REDIS_NAME) redis:8.2.4-alpine3.22 $(REDIS_CMD)

stop-local-redis:
	docker container stop $(LOCAL_REDIS_NAME)

clean-local-redis:
	docker volume rm $(LOCAL_REDIS_DATA_NAME)

run-local: $(TEAM_TASK_MANAGER_BIN)
	$(TEAM_TASK_MANAGER_BIN)

run-local-fast:
	go run $(TEAM_TASK_MANAGER_DIR)/main.go

$(TEAM_TASK_MANAGER_BIN):
	go build -o $(TEAM_TASK_MANAGER_BIN) $(TEAM_TASK_MANAGER_DIR)

service:
	docker compose up

service-stop:
	docker compose down

service-rebuild:
	docker compose down -v
	docker compose up --build --renew-anon-volumes --force-recreate

$(COVERAGE_FILE):
	go test "-coverpkg=./..." "-coverprofile=$(COVERAGE_FILE)$(NOT_FILTERED_SUFF)" -v -short ./internal/...
	$(FILTER_COVERAGE_FROM_MOCK)

coverage-info: $(COVERAGE_FILE)
	go tool cover "-func=coverage.out"

coverage-html: $(COVERAGE_FILE)
	go tool cover "-html=coverage.out"

unit-test:
	go test -v -short ./internal/...

integration-test:
	go test -v ./tests/integration/...

clean-go-cache:
	go clean -cache -modcache
	go env -w GOPROXY=https://proxy.golang.org,direct

clean:
	$(RM) $(MIGRATOR_BIN) $(TEAM_TASK_MANAGER_BIN) $(COVERAGE_FILE) $(COVERAGE_FILE)$(NOT_FILTERED_SUFF) $(RM_POSTFIX)
