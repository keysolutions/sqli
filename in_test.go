package sqli

import "testing"

func TestIn(t *testing.T) {
	tests := []struct {
		query    string
		args     []interface{}
		expQuery string
		expArgs  []interface{}
	}{
		{
			"SELECT id FROM table WHERE id IN ($1)", []interface{}{[]int{1, 2, 3}},
			"SELECT id FROM table WHERE id IN ($1, $2, $3)", []interface{}{1, 2, 3},
		},
		{
			"SELECT id FROM table WHERE id IN (?)", []interface{}{[]int{1, 2, 3}},
			"SELECT id FROM table WHERE id IN (?, ?, ?)", []interface{}{1, 2, 3},
		},
		{
			"SELECT id FROM table WHERE id IN ($1) AND id = $2", []interface{}{[]int{1, 2, 3}, 4},
			"SELECT id FROM table WHERE id IN ($1, $2, $3) AND id = $4", []interface{}{1, 2, 3, 4},
		},
		{
			"SELECT id FROM table WHERE id IN ($1)", []interface{}{},
			"SELECT id FROM table WHERE id IN ()", []interface{}{},
		},
	}
	for _, test := range tests {
		q, a := In(test.query, test.args...)
		if q != test.expQuery {
			t.Errorf("expected '%s' to equal '%s'", q, test.expQuery)
		}
		if !equal(a, test.expArgs) {
			t.Errorf("expected '%v' to equal '%v'", a, test.expArgs)
		}
	}
}

// equal compares to interface slices for equality.
func equal(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
