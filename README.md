# Golang Clean Architecture Template

## Description

This is golang clean architecture template.

## Architecture

![Clean Architecture](architecture.png)

1. External system perform request (HTTP, gRPC, Messaging, etc)
2. The Delivery creates various Model from request data
3. The Delivery calls Use Case, and execute it using Model data
4. The Use Case create Entity data for the business logic
5. The Use Case calls Repository, and execute it using Entity data
6. The Repository use Entity data to perform database operation
7. The Repository perform database operation to the database
8. The Use Case create various Model for Gateway or from Entity data
9. The Use Case calls Gateway, and execute it using Model data
10. The Gateway using Model data to construct request to external system 
11. The Gateway perform request to external system (HTTP, gRPC, Messaging, etc)

## How To Use This Project

1. Check required tool

```shell
make check-tools
```

2. Rename go module name

```shell
make rename-go-mod
```

Now the project is yours.

Next, see [How To Run Application](run_app.md).

Check also other command in [Makefile](Makefile).

Check also [My Note](README_my_note.md).

For contributor instructions, see [Repository Guidelines](AGENTS.md).
