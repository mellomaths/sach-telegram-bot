
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

docker-build:
    docker build -t sach-telegram-bot .

docker-run:
    docker run -d --name sach-telegram-bot -p 8080:8080 sach-telegram-bot

docker-stop:
    docker stop sach-telegram-bot

docker-remove:
    docker rm sach-telegram-bot
