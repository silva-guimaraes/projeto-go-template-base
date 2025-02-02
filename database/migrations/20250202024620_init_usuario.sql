-- +goose Up
-- +goose StatementBegin
-- +goose StatementEnd

-- Diff code generated with pgModeler (PostgreSQL Database Modeler)
-- pgModeler version: 1.2.0-alpha1
-- Diff date: 2025-02-01 23:56:43
-- Source model: dev1
-- Database: dev1
-- PostgreSQL version: 16.0

-- [ Diff summary ]
-- Dropped objects: 0
-- Created objects: 1
-- Changed objects: 0

SET search_path=public,pg_catalog;
-- ddl-end --


-- [ Created objects ] --
-- object: public.usuario | type: TABLE --
-- DROP TABLE IF EXISTS public.usuario CASCADE;
CREATE TABLE public.usuario (
	id integer NOT NULL GENERATED ALWAYS AS IDENTITY ,
	email text NOT NULL,
	senha_hash text NOT NULL,
	nome text NOT NULL,
	ctime timestamp NOT NULL DEFAULT current_timestamp,
	CONSTRAINT usuario_pk PRIMARY KEY (id)
);
-- ddl-end --
ALTER TABLE public.usuario OWNER TO dev;
-- ddl-end --


