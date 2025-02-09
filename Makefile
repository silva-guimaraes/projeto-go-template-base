.PHONY: goose build run goose tw templ
include .env

run: templ tw
	go run ./server

build: tw
	go build -v -o base ./server

tw:
	tailwindcss -i views/static/css/input.css -o views/static/css/output.css

templ:
	templ generate


db_teste: goose
	psql -h localhost -U ${PG_USER} ${PG_DB} < database/testdata/teste.sql

goose:
	goose postgres "postgres://${PG_USER}:${PG_PASS}@${PG_ADDR}/${PG_DB}" up -dir database/migrations
