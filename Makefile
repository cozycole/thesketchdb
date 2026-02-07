include .env

ENV ?= dev

ifeq ($(ENV),prod)
	DB := $(DB_URL)
else
	DB := $(DEV_DB_URL)
endif

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]
	
# ==================================================================================== #
# DB MIGRATIONS
# ==================================================================================== #

## db/migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## migrations/up: apply all up database migrations to ENV={prod,dev}
.PHONY: migrations/up
migrations/up: confirm
	@echo 'Running up migrations to $(ENV) db...'
	migrate -path=./migrations -database $(DB) up

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: dev/app
dev/app:
	air -c ./.air.toml 

.PHONY: dev/tailwind
dev/tailwind:
	npx tailwind \
	  -c ./tailwind.config.js \
	  -i ./ui/static/css/styles.css \
	  -o ./ui/static/css/dist/styles.css \
	  --watch 

.PHONY: dev/esbuild
dev/esbuild:
	npx esbuild ./ui/static/js/main.js \
	  --bundle \
	  --outfile=./ui/static/js/dist/main.js \
	  --sourcemap \
	  --watch=forever 

.PHONY: dev/browser-sync
dev/browser-sync:
	npx browser-sync start \
	  --watch \
	  --files 'cmd/**/*, internal/**/*, ui/**/*' \
	  --reload-delay 1800 \
	  --port 4001 \
	  --proxy 'localhost:8080'

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DEV_DB_URL}

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

## prod/build: build thesketchdb image and deploy it
.PHONY: prod/build
prod/build:
	docker compose build thesketchdb
	docker compose up thesketchdb -d

## prod/psql: connect to the database container using psql
.PHONY: prod/psql
prod/psql:
	docker exec -it thesketchdb-db-1 psql -U colet -d thesketchdb

