
- PotsgreSQL
- Biblioteca padrão para routing e manipulações no banco
- Templ
- Tailwind CSS
- Docker
- [Air](https://github.com/air-verse/air) Live reloading após alterações
- [HTMX](https://htmx.org/) + [_Hyperscript](https://hyperscript.org/) para conteúdo dinâmico no front end sem a necessidade de javascript

Todo o projeto serve de ambos um ponta pé inicial para um projeto novo ou playground para testar novas bibliotecas e tecnologias.
Não espero que outros usem isso.

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

<hr>

> "Perfection is achieved, not when there is nothing more to add, but when there is nothing left to take away"
> —Antoine de Saint-Exupéry
