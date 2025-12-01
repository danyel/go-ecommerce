![Build](https://github.com/danyel/go-ecommerce/actions/workflows/go.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/danyel/go-ecommerce/badge.svg?branch=main)](https://coveralls.io/github/danyel/go-ecommerce?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/danyel/go-ecommerce)](https://goreportcard.com/report/github.com/danyel/go-ecommerce)
![GitHub release](https://img.shields.io/github/v/release/danyel/go-ecommerce)
![License](https://img.shields.io/github/license/danyel/go-ecommerce)

# Go-Commerce

React application with a golang backend.\
With a rabbitmq broker for event sourcing and postgres as a database.\
Golang (v1.25.4) for the backend development.\
Typescript for the frontend.

## Tech-stack

| Functionality | Language/Framework/Tool/Technology |
|---------------|------------------------------------|
| Backend       | Golang                             |
| Frontend      | Typescript                         |
| CSS           | Tailwindcss                        |
| Router        | Chi                                |
| ORM           | Gorm *(maybe migrating to bun)*    |
| Database      | Postgres                           |
| Migration     | Goose                              |
| Container     | Docker                             |
| Broker        | Rabbitmq                           |

## How to start

### Tools

To install all the tools at once use following command or install individually.

```shell
make tools
```

#### Air hot reload

https://github.com/air-verse/air

```shell
go install github.com/air-verse/air@latest
```

#### Goose

https://github.com/pressly/goose

```shell
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### To fetch the dependencies

```shell
go mod tidy
```

### Start backend

#### Prerequisites

- docker

##### Start the database

```shell
docker compose up -d ecommerce-database
```

###### Configuration details

| Key      | Value     |
|----------|-----------|
| username | ecommerce |
| password | ecommerce |
| database | ecommerce |
| port     | 5401      |

```sql
CREATE SCHEMA ecommerce;
```

###### Migrate the database

```shell
make migration
```

#### Broker (rabbitmq)

###### Start the broker

```shell
docker compose up -d rabbitmq
```

###### Configuration details

| Key      | Value     |
|----------|-----------|
| username | developer |
| password | developer |
| url      | localhost |
| port     | 5672      |
| protocol | amqp      |

```shell
make run
```

### Start frontend

#### Prerequisites

- Install npm and node

##### Install dependencies

```shell
npm i
```

#### Run frontend

```shell
make ui
```

Access the application on http://localhost:5173/

### Run the test

###### Integration tests

```shell
make integration_tests
```

###### Mock tests

```shell
make mock_tests
```

## API

| Functionality                     | Endpoint                                           |
|-----------------------------------|----------------------------------------------------|
| create shopping basket            | POST /api/shopping-basket/v1/shopping-baskets      |
| add item to shopping basket       | POST /api/shopping-basket/v1/shopping-baskets/{id} |
| get shopping basket               | GET  /api/shopping-basket/v1/shopping-baskets/{id} |
| product management get products   | GET /api/product-management/v1/products            |
| product management create product | POST /api/product-management/v1/products           |
| product management get product    | GET /api/product-management/v1/products/{id}       |
| product management delete product | DELETE /api/product-management/v1/products/{id}    |
| product management update product | PUT /api/product-management/v1/products/{id}       |
| get categories                    | GET /api/category/v1/categories                    |
| get translations                  | GET  /api/cms/v1/translations                      |
| get translation                   | GET /api/cms/v1/translations/{language}/{code}     |
| management add translation        | POST /api/management/v1/translations               |
| get products                      | GET /api/product/v1/products                       |
| get product                       | GET /api/product/v1/products/{id}                  |

### Makefile commands

###### Build the projects

```shell
make build
```

###### Hotswap the application

```shell
make air
```

###### Run the react app

```shell
make ui 
```

###### Start the database migration

```shell
make migration
```

###### Build and test the application

```shell
make full
```

###### Install tools

```shell
make tools
```

###### Run integration tests

```shell
make integration_tests
```

###### Run mock tests

```shell
make mock_tests
```