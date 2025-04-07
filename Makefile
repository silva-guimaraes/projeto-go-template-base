.PHONY: dev build build-docker tw templ db_teste goose

include .env

dev:
	sudo docker compose up db adminer prometheus -d
	air

build: tw templ
	go build -v -o base ./server

build-docker:
	sudo docker compose up --force-recreate --build -d

tw:
	tailwindcss -i views/static/css/input.css -o views/static/css/output.css

templ:
	templ generate


db-teste: goose
	psql -h localhost -U ${PG_USER} ${PG_DB} < database/testdata/teste.sql

goose:
	goose postgres "postgres://${PG_USER}:${PG_PASS}@${PG_ADDR}/${PG_DB}" up -dir database/migrations
