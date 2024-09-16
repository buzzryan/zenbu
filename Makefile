# Initialize .env file
init-env:
	cp .env-example .env

# Run local server (no container)
local-db:
	docker run --name zenbu-mysql --platform linux/amd64 -d -p 3306:3306 --env-file .env-rds mysql:8.0.27

local-ddb:
	docker run --name zenbu-ddb -d -p 20019:8000 amazon/dynamodb-local

local-db-migrate:
	set -a; source .env; set +a; export MYSQL_ENDPOINT=localhost:3306; go run cmd/migration/rdb/main.go

ddb:
	set -a; source .env; set +a; go run cmd/migration/ddb/main.go

local:
	set -a; source .env; set +a; export MYSQL_ENDPOINT=localhost:3306; go run cmd/server/main.go

local-air:
	air