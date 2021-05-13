# Sator API

Sator API, monolith based on go-kit

## Prerequisites

To use make commands you must have installed the following things:
- **[Golang](https://golang.org)** - to compile the application from source
- **[Docker](https://www.docker.com/get-started)** and **[docker-compose](https://docs.docker.com/compose/install/)** - to run database on local machine
- **[sql-migration](https://github.com/rubenv/sql-migrate)** - to run migrations
- **[sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)** - to generate SQL repository boilerplate using your SQL-queries and migrations
- **[Make](https://www.gnu.org/software/make/)** to use make-helpers


## First run on localhost

```shell
make up && sleep 10 && make migrate run-local
```
then the API will be available on `localhost:8080`

## Development

comming soon

## Useful tools

- [TablePlus](https://tableplus.com) - DB manager
- [Insomnia](https://insomnia.rest) - open source API client
- [VS Code](https://code.visualstudio.com) - IDE (in addition, you can find recommended plugins in the `.vscode` folder)