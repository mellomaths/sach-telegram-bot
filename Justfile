
run:
    go run cmd/*.go

build:
    go build -o sach-telegram-bot ./cmd

exec:
    ./sach-telegram-bot

sqlc:
    sqlc generate

pg-migration-up:
    cd internal/adapters/postgresql/migrations && goose postgres postgres://postgres:postgres@192.168.1.100:5432/sach-bot up

pg-migration-down:
    cd internal/adapters/postgresql/migrations && goose postgres postgres://postgres:postgres@192.168.1.100:5432/sach-bot down

create-migration NAME:
    goose -s create {{NAME}} sql
