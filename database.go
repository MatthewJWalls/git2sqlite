
/*
    Functions for writing Git objects to an SQLite database.
*/

package main

import (
	"os"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func initDatabase(dbname string) {

	// executes the required sqlite DML statements

	os.Remove(dbname)

	db, openerr := sql.Open("sqlite3", dbname)

	if openerr != nil {
		log.Fatal(openerr)
	}

	defer db.Close()

	stmt := `
    create table commits (hash text not null primary key, content text);
    create table trees (hash text not null primary key, content text);
    create table blobs (hash text not null primary key, content text);
    create table refs (path text not null primary key, hash text);
    `
	_, stmterr := db.Exec(stmt)

	if stmterr != nil {
		log.Fatal(stmterr)
	}

}

func writeRefsToSQLite(refs []GitReference, dbname string) {

	// inserts git reference data into an sqlite database

	db, openerr := sql.Open("sqlite3", dbname)
	
	if openerr != nil {
		log.Fatal(openerr)
	}

	defer db.Close()

	tx, _ := db.Begin()
	stmt, stmerr := tx.Prepare("insert into refs (path, hash) values (?, ?)")

	if stmerr != nil {
		log.Fatal(stmerr)
	}

	for _, o := range(refs) {
		_, execerr := stmt.Exec(o.path, o.hash)

		if execerr != nil {
			log.Fatal(execerr)
		}

	}

	tx.Commit()

}

func writeObjectsToSQLite(objects []GitObject, dbname string) {

	// inserts git object data into an sqlite database

	db, openerr := sql.Open("sqlite3", dbname)

	if openerr != nil {
		log.Fatal(openerr)
	}
	
	defer db.Close()

	tx, _ := db.Begin()

	for _, o := range(objects) {

		if o.kind == "commit" {

			stmt, stmterr := tx.Prepare("insert into commits (hash, content) values (?, ?)")

			if stmterr != nil {
				log.Fatal(stmterr)
			}

			_, execerr := stmt.Exec(o.hash, o.content)

			if execerr != nil {
				log.Fatal(execerr)
			}

		} else if o.kind == "tree" {


			stmt, stmterr := tx.Prepare("insert into trees (hash, content) values (?, ?)")

			if stmterr != nil {
				log.Fatal(stmterr)
			}

			_, execerr := stmt.Exec(o.hash, o.content)

			if execerr != nil {
				log.Fatal(execerr)
			}

		} else if o.kind == "blob" {

			stmt, stmterr := tx.Prepare("insert into blobs (hash, content) values (?, ?)")

			if stmterr != nil {
				log.Fatal(stmterr)
			}

			_, execerr := stmt.Exec(o.hash, o.content)

			if execerr != nil {
				log.Fatal(execerr)
			}

		}

	}

	tx.Commit()

}
