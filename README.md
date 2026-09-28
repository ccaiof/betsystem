# BetSystem

Serviço em Go para processamento de operações de apostas e movimentação de
carteiras. O ambiente local usa Docker Compose com PostgreSQL, Keycloak e
MiniStack.

A arquitetura e as decisões técnicas estão descritas em
[ARCHITECTURE.md](ARCHITECTURE.md).

## Pré-requisitos

- Docker com Docker Compose v2;
- GNU Make;
- Go 1.27.1, somente se for executar comandos Go fora do container.

## Subir o ambiente

Na raiz do projeto, execute:

```sh
make up
```

Esse comando inicia:

| Serviço   | Endereço              | Finalidade            |
| ---------- | ---------------------- | --------------------- |
| API        | https://localhost:8081 | Aplicação BetSystem |
| Keycloak   | http://localhost:8080  | IdP OIDC              |
| MiniStack  | http://localhost:4566  | Serviços AWS locaiss |
| PostgreSQL | localhost:5432         | Persistência         |

A aplicação espera que PostgreSQL, Keycloak e MiniStack estejam saudáveis antes
de iniciar. Para acompanhar os logs:

```sh
docker compose -f docker/docker-compose.dev.yml logs -f app
```

## Migrations

Depois de subir os containers, aplique as migrations:

```sh
make migrate-up
```

Consultar o estado:

```sh
make migrate-status
```

Reverter a última migration:

```sh
make migrate-down
```

A reversão deve ser usada apenas em desenvolvimento. Para apagar containers,
volumes e todos os dados locais:

```sh
make delete
```

## Executar a aplicação

O container de desenvolvimento inicia a aplicação com Air, que recompila o
binário quando os arquivos são alterados. Depois de executar `make up` e
`make migrate-up`, a API fica disponível em:

```text
https://localhost:8081
```

O endpoint público de health check disponível atualmente é:

```sh
curl https://localhost:8081/health/live
curl https://localhost:8081/health/ready
```

A rota autenticada de teste é:

```sh
curl https://localhost:8081/auth/me
```

Ela exige um bearer token válido emitido pelo realm `betsystem` do Keycloak.
As credenciais e clientes de teste são provisionados pelo arquivo
`keycloak/betsystem-realm.json`. Não use credenciais de desenvolvimento em
ambientes compartilhados.

## Testes

Os comandos abaixo executam os testes dentro do container da aplicação e
carregam as variáveis de `.env.test`.

### Todos os testes

```sh
make test
```

### Todos os testes com saída detalhada

```sh
make test-v
```

### Testes com detector de race condition

```sh
make test-race
```

Esse comando executa `go test -race ./...`. Ele é especialmente importante para
validar código que coordena operações concorrentes.

### Testes de integração

```sh
make test-integration
```

Esse comando:

1. inicia o ambiente Docker Compose;
2. aplica as migrations;
3. executa os testes em `tests/integration` com `-p 1`.

Os testes de integração dependem de PostgreSQL, Keycloak e MiniStack. Se o
ambiente já estiver em execução, o comando pode reutilizar os containers e os
volumes existentes.

### Testes diretamente com Go

Com as dependências disponíveis localmente, também é possível executar:

```sh
go test ./...
go test -race ./...
go vet ./...
```

A forma recomendada para integração continua sendo `make test-integration`,
porque ela garante que o PostgreSQL e as demais dependências sejam executados
no ambiente definido pelo Compose.

## Variáveis de ambiente

A aplicação usa, no Compose, as seguintes configurações principais:

```text
APP_ENV=development
KEYCLOAK_ISSUER=http://keycloak:8080/realms/betsystem
KEYCLOAK_CLIENT_ID=betsystem-api
DB_URL=postgres://postgres:postgres@db:5432/betsystem?sslmode=disable
```

O arquivo `.env.test` contém os identificadores dos clientes usados nos testes
OIDC. Segredos reais não devem ser adicionados ao repositório.

## Parar o ambiente

Para parar os containers preservando os volumes:

```sh
make down
```

Para reiniciar os serviços:

```sh
make restart
```

Para remover também os volumes e dados locais:

```sh
make delete
```

## Estado atual

O projeto já contém o domínio de `Money`, carteira, transação e ledger, a
integração básica com PostgreSQL, a autenticação OIDC com Keycloak, migrations
e composição inicial com Fx.

Ainda estão pendentes ou incompletos no fluxo executável o consumidor SQS, a
publicação da outbox, o processamento completo de `WIN`, `LOSS`, `REFUND` e
`ROLLBACK`, os endpoints de negócio e os cenários distribuídos de recuperação.
Essas limitações estão detalhadas em [ARCHITECTURE.md](ARCHITECTURE.md).
