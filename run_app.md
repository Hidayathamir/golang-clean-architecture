## How To Run Application

Follow these steps to set up and run the entire ecosystem (Web, Worker, and Observability tools).

### 1. Start Infrastructure

Run the docker-compose to start Postgres, Kafka, SigNoz (Tracing), and AKHQ (Kafka UI).

```bash
make docker-compose
```

Wait until all containers are healthy. You can check the status using:
```bash
docker ps --format "{{.Names}}\t{{.Status}}"
```

### 2. Run Database Migrations

From a new terminal, run the migrations to set up your database schema.

```bash
make migrate
```

### 3. Run Application Servers

You usually need to run both the Web server (for APIs) and the Worker (for background jobs/Kafka consumers).

**Terminal A: Run Web Server**
```bash
make run
```
*   **Swagger UI**: [http://localhost:3000/swagger](http://localhost:3000/swagger)
*   **API Base URL**: `http://localhost:3000`

**Terminal B: Run Worker**
```bash
make run-worker
```
*   This handles async tasks like sending notifications or processing image uploads from Kafka topics.

### 4. Observability & Management Tools

Once everything is running, you can monitor the system using these tools:

| Tool | URL | Description |
| :--- | :--- | :--- |
| **SigNoz** | [http://localhost:3301](http://localhost:3301) | View distributed traces, logs, and application metrics. |
| **AKHQ** | [http://localhost:8080](http://localhost:8080) | Manage Kafka topics, consumers, and view messages in real-time. |
| **Swagger** | [http://localhost:3000/swagger](http://localhost:3000/swagger) | Interactive API documentation and testing. |

---

### Troubleshooting

*   **Port Conflicts**: Ensure ports `5432` (Postgres), `9093` (Kafka), `3301` (SigNoz), `8080` (AKHQ), and `3000` (Web App) are not in use by other services.
*   **Kafka Readiness**: If the worker fails to start, wait a few more seconds for Kafka to be fully ready even after the container shows "Up".
*   **Resetting Environment**: If you want to start fresh, run:
    ```bash
    make docker-compose # This will down -v (remove volumes) and up again
    ```