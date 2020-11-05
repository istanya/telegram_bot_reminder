# Telegram Bot Reminder

Bot help to create scheduled reminders. You may set date and write reminder message then bot send this message at the right time.

You may use the bot in this [link](https://t.me/ReReMind_bot).(@ReReMind_bot)

## Setup local development

### Install tools

- [Golang](https://golang.org/)

- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

### Configuring

Copy the `app.env.example` file to `app.env`:

```bash
cp app.env.example app.env
```

Edit the `app.env` file.

### Setup infrastructure

- Install all dependency:

    ```bash
    dep ensure
    ```

- Run db migration:

    ```bash
    make migrateup
    ```

### How to run

- Run server:

    ```bash
    go run *.go
    ```
