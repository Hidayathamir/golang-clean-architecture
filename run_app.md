## How To Run Application

I recommend run with docker for testing.

1. Run docker compose.

```bash
make docker-compose
```

Wait until this container status is up.

```shell
$ docker ps --format "{{.Names}}\t{{.Status}}"
signoz-otel-collector	Up 33 minutes
signoz	Up 33 minutes (healthy)
akhq	Up 34 minutes (healthy)
signoz-clickhouse	Up 34 minutes (healthy)
kafka-clean-arch	Up 34 minutes
zookeeper-clean-arch	Up 34 minutes
postgres-clean-arch	Up 34 minutes
signoz-zookeeper-1	Up 34 minutes (healthy)
```

2. Keep docker container running. From another terminal, run migration and run application.

```bash
make migrate
```

This will run migration.

```bash
make run
```

This will run application. Go to swagger http://localhost:3000/swagger.
