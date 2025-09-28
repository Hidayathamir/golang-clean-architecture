## How To Run Application

I recommend run with docker for testing.

1. Run docker compose.

```bash
make docker-compose
```

This will run docker compose that will run mysql & kafka. Go to http://localhost:8080/ to see kafka.

2. Keep docker container running. From another terminal, run migration and run application.

```bash
make migrate
```

This will run migration.

```bash
make run
```

This will run application. Go to swagger http://localhost:3000/swagger.