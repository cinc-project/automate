include ../../Makefile.common_go

PACKAGE_PATH = github.com/chef/automate/components/infra-proxy-service/server

BINS = ${PACKAGE_PATH}/cmd/infra-proxy-service
MIGRATION_READMES = storage/postgres/migration/sql/README.md

packages:=${PACKAGE_PATH}/...
ifdef CI
  verbose?="-v"
endif

PG_URL ?= "postgresql://postgres@127.0.0.1:5432/infra_proxy_test?sslmode=disable&password=docker"

all: lint build test

static: lint ${MIGRATION_READMES}
unit: build test

.PHONY: ${MIGRATION_READMES}
${MIGRATION_READMES}:
	../../scripts/generate_and_check_migration_files.sh $@

build:
	mkdir -p bin/
	go build -o bin/infra-proxy-service ./cmd/infra-proxy-service

test:
	@PG_URL=$(PG_URL) go test $(verbose) -cover -count=1 -parallel=1 -p 1 $(packages)

setup_docker_pg:
	docker run --name infra-proxy-postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=infra_proxy_test -e POSTGRES_PASSWORD=docker -p 5432:5432 -d postgres:9
	sleep 15 # let docker come up
	# This creates the extension we need to use UUIDs in the migrations.
	# Done in habitat in prod. Not done in code because you must be a superuser.
	docker exec infra-proxy-postgres psql -d $(PG_URL) -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\""

kill_docker_pg:
	docker rm -f infra-proxy-postgres || true

# this command lists all the changes since master, and looks for modifications
# to the migration files -- if there's any of (M)odify, (R)ename, or (D)elete,
# it will trigger an error
migrations-append-only-check:
	@git diff --name-status master... | awk '/[RMD][0-9]*\tcomponents\/infra-proxy\/postgres\/migration\/sql\//{ print "modified/deleted/renamed: " $$2; ec=1 } END{ exit ec }'

.PHONY: all static unit build compile proto test dep-ensure have-dep lint
