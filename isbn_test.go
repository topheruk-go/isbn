package isbn

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/jackc/pgx/v4"
	_ "github.com/mattn/go-sqlite3"
)

func TestIsbnValidation(t *testing.T) {
	tt := []struct {
		desc string
		data string
		err  string
	}{
		{desc: "invalid: format", data: "978071670344a", err: "invalid ISBN format"},
		{desc: "invalid: length 12 != 13", data: "978071670344", err: "invalid ISBN length 12"},
		{desc: "valid: isbn 13", data: "9780716703440"},
		{desc: "invalid: isbn 13", data: "9780716703410", err: "invalid ISBN value"},
		{desc: "invalid: length 14 != 13", data: "97807167034403", err: "invalid ISBN length 14"},
		{desc: "valid: isbn 10", data: "0716703440"},
		{desc: "valid: isbn 10 w/ dashes", data: "0-7167-0344-0"},
		{desc: "valid: isbn 13 w/ dashes", data: "978-0-7167-0344-0"},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			isbn, err := Parse(tc.data)
			if err != nil {
				assert.Error(t, err, tc.err)
				return
			}

			assert.Equal(t, len(isbn.String()), 13)
		})
	}
}

// func TestIsbnJSON(t *testing.T) {
// 	data := `["9780716703440","0716703440"]`

// 	var isbns []ISBN
// 	err := json.NewDecoder(strings.NewReader(data)).Decode(&isbns)
// 	assert.NilError(t, err)

// 	var buf bytes.Buffer
// 	assert.NilError(t, json.NewEncoder(&buf).Encode(&isbns[1]))

// 	// cant think of a better way to test this
// 	assert.Equal(t, len(buf.String()), len("9780716703440")+1)
// }

// func TestIsbnSql(t *testing.T) {
// 	isbn, _ := Parse("9780716703440")

// 	db, _ := sql.Open("sqlite3", ":memory:")

// 	type testcase struct {
// 		ID   int
// 		Isbn *ISBN // nullable
// 	}

// 	tt := []struct {
// 		desc   string
// 		schema string
// 	}{
// 		{"isbn as BLOB", "CREATE TABLE foo (id INTEGER PRIMARY KEY, isbn BLOB)"},
// 		{"isbn as TEXT", "CREATE TABLE foo (id INTEGER PRIMARY KEY, isbn TEXT)"},
// 	}
// 	for _, tc := range tt {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			_, err := db.Exec(tc.schema)
// 			assert.NilError(t, err)

// 			_, err = db.Exec(`INSERT INTO foo (isbn) VALUES ($1)`, isbn)
// 			assert.NilError(t, err)

// 			var tc testcase
// 			assert.NilError(t, db.QueryRow(`SELECT id, isbn FROM foo WHERE id = 1`).Scan(&tc.ID, &tc.Isbn))
// 			assert.Equal(t, tc.Isbn.String(), isbn.String())

// 			db.Exec(`DROP TABLE foo`)
// 		})
// 	}
// }

// func TestIsbnPsql(t *testing.T) {
// 	db, err := pgx.Connect(context.Background(), os.ExpandEnv("host=${POSTGRES_HOSTNAME} port=${DB_PORT} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable"))
// 	assert.NilError(t, err)

// 	isbn, _ := Parse("9780716703440")
// 	type testcase struct {
// 		ID    int
// 		Isbn  *ISBN     // nullable
// 		Title *[13]byte // nullable
// 	}

// 	tt := []struct {
// 		desc   string
// 		schema string
// 	}{
// 		{"isbn as VARCHAR(13)", "CREATE TEMP TABLE foo (id SERIAL PRIMARY KEY, isbn VARCHAR(13), title VARCHAR(255))"},
// 	}
// 	for _, tc := range tt {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			_, err := db.Exec(context.Background(), tc.schema)
// 			assert.NilError(t, err)

// 			_, err = db.Exec(context.Background(), `INSERT INTO foo (isbn) VALUES ($1)`, isbn)
// 			assert.NilError(t, err)

// 			var tc testcase
// 			assert.NilError(t, db.QueryRow(context.Background(), `SELECT id, isbn, title FROM foo WHERE id = 1`).Scan(&tc.ID, &tc.Isbn, &tc.Title))
// 			assert.Equal(t, tc.Isbn.String(), isbn.String())

// 			// db.Exec(context.Background(),`DROP TABLE test`)
// 		})
// 	}
// }
