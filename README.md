
- PotsgreSQL
- Biblioteca padrão para routing e manipulações no banco
- Templ
- Tailwind CSS
- Docker
- [Air](https://github.com/air-verse/air) Live reloading após alterações
- [HTMX](https://htmx.org/) + [_Hyperscript](https://hyperscript.org/) para conteúdo dinâmico no front end sem a necessidade de javascript

## Quickstart
### Usando Docker

Favor primeiro verificar caso haja alguma outra instância do Postgres rodando na mesma porta (5432).

```sh
docker compose up
```
Servidor estará disponível em `http://localhost:8888` por padrão.


### Desenvolvimento
```sh
docker compose up db adminer -d
air
```
