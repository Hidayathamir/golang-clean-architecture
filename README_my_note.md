## My Note

This project is fork from https://github.com/khannedy/golang-clean-architecture with this feature:

1. Better delivery return handling. See [return response.Data(ctx, http.StatusOK, res)](internal/delivery/http/address_controller.go).
2. Better logging using middleware. See [usecase](internal/usecase/address/create.go) clean, logging in [usecase middleware log](internal/usecase/address/address_usecase_mw_logger.go).
3. Better error handling. See [errkit.BadRequest(err)](internal/usecase/address/create.go) will handled in [response.Error](internal/config/fiber.go).
4. Better error handling 2. See [errkit.AddFuncName](internal/usecase/address/create.go). String func name is auto generate with command `make errkitcli`. Example response json:
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
5. Request has trace id. See [example](internal/delivery/http/middleware/trace_id_middleware.go).
6. Better searching log with trace id. Example log:
```json
{
    "err": "user.(*UserUsecaseImpl).Create:: [409] conflict:: user already exists",
    "fields": {
        "req": {
            "id": "joko",
            "password": "joko",
            "name": "Joko"
        },
        "res": null
    },
    "level": "error",
    "msg": "user.(*UserUsecaseMwLogger).Create",
    "source": "/home/hidayat/data-d/myrepo/golang-clean-architecture/internal/usecase/user/user_usecase_mw_logger.go:30",
    "time": "2025-10-11T08:06:12+07:00",
    "trace_id": "080eb9d7-a5db-44e8-8493-0ae302f612d9"
}
```
7. Swagger auto generated. See [example](internal/delivery/http/address_controller.go). See http://localhost:3000/swagger
8. Using interface make it easier to test. See [example](internal/usecase/address/address_usecase.go).
9. Unit test example. See [usecase address](internal/usecase/address).
10. Command shortcut using [Makefile](Makefile).
11. Gateway rest api client. See [example](internal/gateway/rest/slack_client.go).
12. Simple repository without generic. See [example](internal/repository/user_repository.go).
13. Simple call kafka producer. See [u.AddressProducer.Send](internal/usecase/address/create.go).
14. Splitted usecase. See [example](internal/usecase/address).
15. Test will use container db. See [example](test/init_test.go).
16. Run application with docker container. See [Run Application](#run-application) & [docker-compose.yml](docker-compose.yml).
