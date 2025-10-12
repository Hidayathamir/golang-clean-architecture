## My Note

This project is forked from [khannedy/golang-clean-architecture](https://github.com/khannedy/golang-clean-architecture) with the following features:

---

### 🚀 Features

1. **General**
   - **Better delivery return handling** — see [`return response.Data(ctx, http.StatusOK, res)`](internal/delivery/http/address_controller.go).
   - **Swagger auto generation** — see [example](internal/delivery/http/address_controller.go). Generate with `make swag`. Access: [http://localhost:3000/swagger](http://localhost:3000/swagger)
   - **Command shortcuts via Makefile** — see [Makefile](Makefile).
   - **Gateway REST API client** — example: [Slack client](internal/gateway/rest/slack_client.go).
   - **Simplified repository (no generics)** — see [UserRepository](internal/repository/user_repository.go).
   - **Simple Kafka producer call** — usage: [`u.AddressProducer.Send`](internal/usecase/address/create.go).
   - **Split usecase by domain** — example: [Address usecase](internal/usecase/address).
   - **Run with Docker** — see [Run Application](#run-application) & [docker-compose.yml](docker-compose.yml).

2. **Error Handling**
   - **Consistent error wrapping & mapping** — e.g. [`errkit.BadRequest(err)`](internal/usecase/address/create.go) handled by [`response.Error`](internal/config/fiber.go).
   - **Auto function-name enrichment** — [`errkit.AddFuncName`](internal/usecase/address/create.go); names generated via `make errkitcli`.
   - **Example response:**
     ```json
     {
       "data": null,
       "error_message": "conflict",
       "error_detail": [
         "http.(*UserController).Register",
         "user.(*UserUsecaseImpl).Create",
         "[409] conflict",
         "user already exists"
       ]
     }
     ```

3. **Logging**
   - **Middleware-driven structured logging** — business logic stays clean; see [usecase](internal/usecase/address/create.go) and [usecase middleware logger](internal/usecase/address/address_usecase_mw_logger.go).
   - **Trace ID per request** — see [trace ID middleware](internal/delivery/http/middleware/trace_id_middleware.go).
   - **Trace-friendly logs** — sample:
     ```json
     {
       "err": "user.(*UserUsecaseImpl).Create:: [409] conflict:: user already exists",
       "fields": {
         "req": { "id": "joko", "password": "joko", "name": "Joko" },
         "res": null
       },
       "level": "error",
       "msg": "user.(*UserUsecaseMwLogger).Create",
       "source": "/home/hidayat/data-d/myrepo/golang-clean-architecture/internal/usecase/user/user_usecase_mw_logger.go:30",
       "time": "2025-10-11T08:06:12+07:00",
       "trace_id": "080eb9d7-a5db-44e8-8493-0ae302f612d9"
     }
     ```

4. **Testing**
   - **Interface-first design for easy mocking** — see [address usecase interface](internal/usecase/address/address_usecase.go). Generate mock with `make generate`.
   - **Unit test examples** — see [usecase/address tests](internal/usecase/address).
   - **Integration tests with containerized DB** — see setup in [test/init_test.go](test/init_test.go).

---