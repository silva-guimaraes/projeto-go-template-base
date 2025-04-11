
# projeto-go-template-base
## Stack
- PotsgreSQL
- Biblioteca padrão para routing e manipulações no banco
- Templ
- Tailwind CSS
- Docker
- [Air](https://github.com/air-verse/air) Live reloading após alterações
- [HTMX](https://htmx.org/) + [_Hyperscript](https://hyperscript.org/) para conteúdo dinâmico no front end sem a necessidade de javascript

Todo o projeto serve de ambos um ponta pé inicial para um projeto novo ou playground para testar novas bibliotecas e tecnologias.
Não espero que outros usem isso.

## Raciocínio pela escolha da stack
- Eliminar lógica do front-end.
- Minimizar dependência em javascript.
- Minimizar etapas de compilação.
    - Atualmente temos apenas: o servidor, templates `.templ`, e classes em *tailwind*.
- Tecnologias com histórico/promessas de estabilidade.
- Minimizar bugs em ambiente de produção.
- Desenvolvimento local.
- Simplicidade.

## Quickstart
### Usando Docker
```sh
make docker-build
```
Cria o contêiner e instala as imagens e ferramentas necessárias por de baixo dos panos.
O servidor estará disponível em `http://localhost:8888` por padrão.
Favor verificar caso haja algum outro serviço rodando nas seguintes portas: `5432`, `8888`, `9090`, `8090`.

### Desenvolvimento
```sh
make dev
```
Necessita que seu sistema tenha as dependências listadas acima.

## TODO:

- [x] Prometheus
- [x] Organizar dados para a coleta
- [ ] Grafana
- [ ] Deploy para uma VPS
- [ ] Explorar alternativas para `database/sql`
- [ ] Explorar alternativas para `testing`
- [ ] CI/CD
