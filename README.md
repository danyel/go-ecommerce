![Build](https://github.com/danyel/go-ecommerce/actions/workflows/go.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/danyel/go-ecommerce/badge.svg?branch=main)](https://coveralls.io/github/danyel/go-ecommerce?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/danyel/go-ecommerce)](https://goreportcard.com/report/github.com/danyel/go-ecommerce)
![GitHub release](https://img.shields.io/github/v/release/danyel/go-ecommerce)
![License](https://img.shields.io/github/license/danyel/go-ecommerce)

# Go-Commerce

A Rest API for an e-commerce application, written in golang.\
I wanted to know how golang works and how to create a rest api.

## Tech-stack

| Functionality | Framework       |
|---------------|-----------------|
| Router        | Chi             |
| Database      | Gorm (postgres) |
| Migration     | goose           |
| DBMS          | docker          |

## How to start

### Goose

```shell
brew install goose
```

### To fetch the dependencies

```shell
go mod tidy
```

### Start the database

```shell
docker compose up -d
```

| Key      | Value     |
|----------|-----------|
| username | ecommerce |
| password | ecommerce |
| database | ecommerce |
| port     | 5401      |

```sql
CREATE SCHEMA ecommerce;
```

### Migrate the database

```shell
goose up
```

### Application

```shell
go run cmd/main.go
```

### Run the test

#### Integration tests

```shell
go test github.com/danyel/ecommerce/test/integration
```

#### Mock tests

```shell
go test github.com/danyel/ecommerce/test/mock
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
