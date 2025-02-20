---
title: MySQL
keywords: [mysql]
description: Connecting to a MySQL database.
---

# MySQL Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/mysql) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/mysql)

This project demonstrates how to connect to a MySQL database in a Go application using the Velocity framework.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package
- MySQL
- [Go MySQL Driver](https://github.com/go-sql-driver/mysql)

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/mysql
    ```

2. Install dependencies:
    ```sh
    go get
    ```

3. Set up your MySQL database and update the connection string in the code.

## Running the Application

1. Start the application:
    ```sh
    go run main.go
    ```

## Example

Here is an example of how to connect to a MySQL database in a Velocity application:

```go
package main

import (
    "database/sql"
    "log"

    "go.khulnasoft.com/velocity"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // Database connection
    dsn := "username:password@tcp(127.0.0.1:3306)/dbname"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Velocity instance
    app := velocity.New()

    // Routes
    app.Get("/", func(c *velocity.Ctx) error {
        var greeting string
        err := db.QueryRow("SELECT 'Hello, World!'").Scan(&greeting)
        if err != nil {
            return err
        }
        return c.SendString(greeting)
    })

    // Start server
    log.Fatal(app.Listen(":3000"))
}
```

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [MySQL Documentation](https://dev.mysql.com/doc/)
- [Go MySQL Driver Documentation](https://pkg.go.dev/github.com/go-sql-driver/mysql)
