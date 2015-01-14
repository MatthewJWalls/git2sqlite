package main

import (
	"testing"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func TestCreate(t *testing.T) {

	db := SQLiteDatabase{"unittest.db"}
	db.Create()

	if _, err := os.Stat("unittest.db"); os.IsNotExist(err) {
		t.Error("No db created.")
	}

	blob := GitBlob{"", "test", "test"}
	blobslice := make([]GitBlob, 1, 1)
	blobslice[0] = blob

	conn, _ := sql.Open("sqlite3", "unittest.db")
	tx, _ := conn.Begin()

	db.writeBlobs(tx, blobslice)
	tx.Commit()
	conn.Close()
	os.Remove("unittest.db")

}

