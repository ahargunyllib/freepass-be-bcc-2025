# BCC Conference API

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)

## Description

This is a backend application for BCC Conference API, part of Freepass challenge by BCC FILKOM UB. It is designed following the best practices and recommendations to ensure a clean, maintainable, and scalable codebase. It is also built with security in mind to protect against common security threats.

The application is built on [Go v1.23.4](https://tip.golang.org/doc/go1.22) and [PostgreSQL](https://www.postgresql.org/). It uses [Fiber](https://docs.gofiber.io/) as the HTTP framework and [pgx](https://github.com/jackc/pgx) as the driver and [sqlx](github.com/jmoiron/sqlx) as the query builder.

## Documentation

For database schema documentation, see [here](https://dbdocs.io/ahargunyllib/freepass-be-bcc-2025), powered by [dbdocs.io](https://dbdocs.io/).

For API documentation, see [here](https://bcc-conference-api-docs.ahargunyllib.tech), powered by [Apidog](https://apidog.com/).

## Architecture

The project is structured following the Clean Architecture, Layered Architecture, Domain-Driven Design principles, Hexagonal Architecture, and SOLID principles.

## Features

- **Migration**: database schema migration using [golang-migrate](https://github.com/golang-migrate/migrate)
- **Validation**: request data validation utilizing [Package validator](https://github.com/go-playground/validator)
- **Logging**: implemented with [zerolog](https://github.com/rs/zerolog)
- **Testing**: unit and integration tests powered by [Testify](https://github.com/stretchr/testify) with formatted output using [gotestsum](https://github.com/gotestyourself/gotestsum)
- **Error handling**: centralized error management system
- **Email functionality**: implemented using [Gomail](https://github.com/go-gomail/gomail)
- **Environment variables**: managed with [Viper](https://github.com/spf13/viper)
- **Security**: HTTP headers secured by [Fiber-Helmet](https://docs.gofiber.io/api/middleware/helmet)
- **CORS**: Cross-Origin Resource-Sharing enabled through [Fiber-CORS](https://docs.gofiber.io/api/middleware/cors)
- **Compression**: gzip compression provided by [Fiber-Compress](https://docs.gofiber.io/api/middleware/compress)
- **Linting**: code quality ensured with [golangci-lint](https://golangci-lint.run)
- **Docker support**
