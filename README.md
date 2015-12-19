# SQLi – Go SQL Query Interface

SQLi is an extension of the Go `database/sql` package that provides a simplified query interface, removing duplication of common access patterns required by the standard library.

## Query Interface

Queries are performed on types that implement the `sqli.Scanner` interface. `Query` and `QueryRow`, similar to those found in `database/sql`, execute the query and return the rows. But instead of returning the rows object, the Scan interface function is called for each row.

### IN queries

The additional `QueryIn` functions like `Query` but help build variadic parameterized queries by expanding slice values. Example: ```QueryIn(scanner, "SELECT name FROM people WHERE id IN (?)", []int{1, 2, 3})```. Placeholders can be supplied in both `?` or `$n` format depending on your database driver requirements.

## Usage

```go
package main

import (
    "log"

    "github.com/keysolutions/sqli"
)

type Person struct {
    Name string
    Age  int32
}

func (p *Person) Scan(row sql.Row) err {
    return row.Scan(&p.Name, &p.Age)
}

func main() {
    db, err := sqli.Open("driver", "name")
    if err != nil {
        log.Fatalf("Unable to open database: %v", err)
    }
    defer db.Close()

    var person Person
    db.QueryRow(&person, "SELECT name, age FROM people WHERE id = ?", 1)
}
```

The `sqli.Scanner` interface can be used to build more complex types.

```go
type PeopleScanner struct {
    People []*Person
}

func (p *PeopleScanner) Scan(row sqli.Row) error {
    var person Person
    if err := row.Scan(&person.Name, &person.Age); err != nil {
        return err
    }
    p.People = append(p.People, &person)
    return nil
}
```

Or even middleware functions to wrap duplicate logic.

```go
type middleware func(row sqli.Row) error

func (m middleware) Scan(row sqli.Row) error {
	return m(row)
}

func MyMiddleware(scanner sqli.Scanner) middleware {
	return func(row sqli.Row) error {
		// Do something here
		return scanner.Scan(row)
	}
}
```