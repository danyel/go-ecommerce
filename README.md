![Build](https://github.com/danyel/go-ecommerce/actions/workflows/go.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/danyel/go-ecommerce/badge.svg?branch=main)](https://coveralls.io/github/danyel/go-ecommerce?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/danyel/go-ecommerce)](https://goreportcard.com/report/github.com/danyel/go-ecommerce)
![GitHub release](https://img.shields.io/github/v/release/danyel/go-ecommerce)
![License](https://img.shields.io/github/license/danyel/go-ecommerce)

# Go-Commerce

React application with a golang backend.\
With a rabbitmq broker for event sourcing and postgres as a database.\
Golang (v1.25.4) for the backend development.\
Typescript for the frontend.\
The project is set up for dev container.

## Tech-stack

| Functionality | Language/Framework/Tool/Technology |
|---------------|------------------------------------|
| Backend       | Golang                             |
| Frontend      | Typescript                         |
| CSS           | Tailwindcss                        |
| Router        | Chi                                |
| ORM           | Gorm *(maybe migrating to bun)*    |
| Database      | Postgres (18)                      |
| Migration     | Goose                              |
| Container     | Docker                             |
| Broker        | Rabbitmq                           |

## How to start

### Tools

To install all the tools at once use following command or install individually.

#### Setting environment

```shell 
git config --global --add safe.directory /workspace
export GOROOT=/usr/local/go
```

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

###### All services up

```shell
make env_up
```

###### All services down

```shell
make env_down
```

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
| port     | 5432      |

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

### Authentication and frontend navigation

Public product, category, translation-read, and basket-read endpoints remain available anonymously. Creating or changing a basket and the checkout route require an `Authorization: Bearer <app-session-token>` header. Google login is exposed at `POST /api/auth/v1/google` with `{ "id_token": "<Google ID token>" }`; configure `GOOGLE_CLIENT_ID` with the web client ID and never commit credentials. The backend verifies the Google signature using Google's certificate endpoint, issuer, audience, expiry, and verified-email claims before issuing the application session token.

The React app uses the same bearer token for API mutations, provides category links and filtering, and protects checkout. `gocommerce/src/state/event-bus.ts` is an extensible in-process event bus for updating interested components; no SSE/WebSocket transport is enabled because the current backend has no compatible event stream. The development account linker is intentionally in-memory and must be replaced with a persistent `UserLinker` implementation before using multiple application instances or relying on account links across restarts.

### Public API anti-scraping controls

Public catalog responses are paginated (`page` and `page_size`) and the server caps
`page_size` at `API_MAX_PAGE_SIZE` (default `50`). The React catalog requests 24
products at a time. Requests are rate-limited per source IP and bearer-token
fingerprint using `API_RATE_LIMIT_REQUESTS` (default `60`) over
`API_RATE_LIMIT_WINDOW` (default `1m`). Configure browser origins explicitly with
`CORS_ALLOWED_ORIGINS` as a comma-separated list; no wildcard origin is enabled.
`robots.txt` asks compliant crawlers not to crawl `/api/`, but this is advisory.

These controls reduce bulk automated collection and backend load; they cannot make
data rendered to a browser impossible to copy. Public content should therefore
never be treated as secret, and stronger protection requires authenticated access,
business controls, and operational monitoring.

### Run the test

### Demo catalog

The Goose baseline installs an original deterministic demonstration catalog with
a computer-store hierarchy and three products. It uses placeholder images and
does not copy product content from a third-party retailer. The equivalent
idempotent SQL is also available in `data/catalog_seed.sql`:

```shell
make migration
psql "$DATABASE_URL" -f data/catalog_seed.sql
```

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
| update shopping basket item       | PUT /api/shopping-basket/v1/shopping-baskets/{id}  |
| get shopping basket               | GET  /api/shopping-basket/v1/shopping-baskets/{id} |
| product management get products   | GET /api/product-management/v1/products            |
| product management create product | POST /api/product-management/v1/products           |
| product management get product    | GET /api/product-management/v1/products/{id}       |
| product management delete product | DELETE /api/product-management/v1/products/{id}    |
| product management update product | PUT /api/product-management/v1/products/{id}       |
| get categories                    | GET /api/management/v1/categories                 |
| create category                   | POST /api/category/v1/categories                   |
| get translations                  | GET  /api/cms/v1/translations                      |
| get translation                   | GET /api/cms/v1/translations/{language}/{code}     |
| Google login                      | POST /api/auth/v1/google                           |
| management add translation        | POST /api/management/v1/translations               |
| get products                      | GET /api/product/v1/products                       |
| get product                       | GET /api/product/v1/products/{id}                  |

### Schema architecture

The clean-install baseline creates normalized users, external identities and
sessions, hierarchical categories, decimal products with JSON metadata, owned
shopping baskets and unique basket lines, keyed stock reservations, and an
outbox for post-commit events. Product/category/CMS reads remain public;
basket mutations and all catalog/content writes require authentication, with
catalog writes restricted to `ADMIN` or `CATALOG_MANAGER` roles.

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

### Dev container

```shell
git config --global --add safe.directory /workspace
git config --global user.name <name>
git config --global user.email <email>
```