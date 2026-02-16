#### Simple app for managing accounts and transactions

REST API for creating accounts and transactions.
You can find oapi spec in the `api` directory.

### Running in docker

If you want to check out the app you can run it with docker compose.

Starting dependencies:
```
sudo docker compose -f ./docker/compose.deps.yaml up -d
```

Starting the app:
```
sudo docker compose -f ./docker/compose.deps.yaml up -d
```

### Building locally

Preparing and building:
```
make mod
make generate
make build_app
```

You can find config examples in the config directory.
Config should provide addresses and creadentials for all dependencies (postgres, rabbitmq).
For starters you could use `config_local.yaml`.

You will also need db migrations, you can find those in the internal/ports/db/migrations.
Run migrations:
```
make migrate_up
```

and

```
make migrate_down
```

Checking the app:
```
curl -v --request POST --header "Content-Type: application/json" --data '{"userId":"7254cbd7-201a-4545-a68d-2f78b725702d", "currency":"usd","fullname":"user1","initialDebitAmount":120}'

curl -v --request GET 'localhost:8080/accounts?user_id=7254cbd7-201a-4545-a68d-2f78b725702d' | jq .
```
