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

If for some reason some dependancy wouldn't start:
```
sudo docker compose -f ./docker/compose.deps.yaml up -d --force-recreate postgres
sudo docker compose -f ./docker/compose.deps.yaml up -d --force-recreate migrations
```

Also you can check the logs:
```
sudo docker logs postgres
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

Creating and checking the account:
```
curl -v --request POST --header "Content-Type: application/json" --data '{"userId":"7254cbd7-201a-4545-a68d-2f78b725702d", "currency":"usd","fullname":"user1","initialDebitAmount":120}' 'localhost:8080/accounts'

curl -v --request GET 'localhost:8080/accounts?user_id=7254cbd7-201a-4545-a68d-2f78b725702d' | jq .
```

Creating and checking the transaction:
```
curl -v --request POST --header "Content-Type: application/json" --data '{"fromAccountId":"a9decc3d-8a21-44f9-8e67-3865199a0088", "toAccountId":"c4b645d2-72bc-4ded-819c-294da62766d9","amount":12.2}' 'localhost:8080/transactions'

curl -v --request GET 'localhost:8080/transactions/b155b9dd-f942-49ef-945f-f5f33f9e0e13' | jq .
```
