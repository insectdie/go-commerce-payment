# Migration Tool

This tool is used to create and run migrations in the database.

## Run Migration

```go
go run migration.go ./sql "host=localhost port=5432 user=postgres dbname=ao_product_service_batch_3 password=postgres sslmode=disable" up
```

## Down Migration

```go
go run migration.go ./sql "host=localhost port=5432 user=postgres dbname=ao_product_service_batch_3 password=postgres sslmode=disable" down
```

## Create new SQL

```go
go run migration.go ./sql "host=localhost port=5432 user=postgres dbname=ao_product_service_batch_3 sslmode=disable" create create_shops_table sql
```
