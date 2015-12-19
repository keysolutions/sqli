package sqli

import "testing"

func newTestDB(t *testing.T, name string) *DB {
	db, err := Open("test", name)
	if err != nil {
		t.Fatal("Open: %v", err)
	}
	exec(t, db, "WIPE")
	exec(t, db, "CREATE|users|id=int64,name=string")
	exec(t, db, "INSERT|users|id=?,name=John", 1)
	exec(t, db, "INSERT|users|id=?,name=Jane", 2)
	return db
}

func exec(t *testing.T, db *DB, query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
	if err != nil {
		t.Fatalf("Exec %q: %v", query, err)
	}
}

// scanFn implements the Scanner interface for bare functions.
type scanFn func(row Row) error

func (s scanFn) Scan(row Row) error {
	return s(row)
}

// testFunc tests each interface by passing in the function to be tested.
func testFunc(t *testing.T, name string, dbfn func(*DB, Scanner, string, ...interface{}) error) {
	db := newTestDB(t, name)
	defer db.Close()
	var userName string
	fn := func(row Row) error {
		if err := row.Scan(&userName); err != nil {
			return err
		}
		return nil
	}
	if err := dbfn(db, scanFn(fn), "SELECT|users|name|id=?", 1); err != nil {
		t.Fatalf("Query: %v", err)
	}
	if userName != "John" {
		t.Errorf("%v: unexpected name: %v", name, userName)
	}
}

func TestQuery(t *testing.T) {
	testFunc(t, "Query", (*DB).Query)
}

func TestQueryIn(t *testing.T) {
	testFunc(t, "QueryIn", (*DB).QueryIn)
}

func TestQueryRow(t *testing.T) {
	testFunc(t, "QueryRow", (*DB).QueryRow)
}
