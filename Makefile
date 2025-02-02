.PHONY: goose build
include .env

run: goose tw
	go run ./server

tw:
	tailwindcss -i views/static/css/input.css -o views/static/css/output.css

db_teste: goose
	psql -h localhost -U ${PG_USER} ${PG_DB} < database/testdata/teste.sql

goose:
	goose postgres "postgres://${PG_USER}:${PG_PASS}@${PG_ADDR}/${PG_DB}" up -dir database/migrations
