
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

type SQLiteDatabase struct {
	fileName string
}

// recreates the database from scratch
func (this SQLiteDatabase) Create() {

	os.Remove(this.fileName)

	db, openerr := sql.Open("sqlite3", this.fileName)

	if openerr != nil {
		log.Fatal(openerr)
	}

	defer db.Close()

	stmt := `
    create table commits (
      commit_hash text not null primary key,
      parent_a text,
      parent_b text,
      author text not null,
      committer text not null,
      date number not null,
      tree_hash not null,
      message text
    );

    create table trees (
      tree_hash text not null, 
      blob_hash text not null,
      file_name text not null,
      file_mode text not null
    );

    create table blobs (
      blob_hash text not null primary key, 
      content text
    );

    create table refs (path text not null primary key, hash text);

    `
	_, stmterr := db.Exec(stmt)

	if stmterr != nil {
		log.Fatal(stmterr)
	}

}

// writes contents of the git repository to the database
func (this SQLiteDatabase) WriteRepository(repo GitRepository) {

	blobs, trees, commits := repo.Objects()

	this.writeBlobs(blobs)
	this.writeTrees(trees)
	this.writeCommits(commits)
	this.writeRefsToSQLite(repo.References())

}

// utility function: inserts git reference data
func (this SQLiteDatabase) writeRefsToSQLite(refs []GitReference ) {

	db, openerr := sql.Open("sqlite3", this.fileName)
	
	if openerr != nil {
		log.Fatal(openerr)
	}

	defer db.Close()

	tx, _ := db.Begin()
	stmt, stmerr := tx.Prepare("insert into refs values (?, ?)")

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

// utility function: insert commit object data
func (this SQLiteDatabase) writeCommits(commits []GitCommit) {

	db, openerr := sql.Open("sqlite3", this.fileName)

	if openerr != nil {
		log.Fatal(openerr)
	}
	
	defer db.Close()
	tx, _ := db.Begin()

	for _, o := range(commits) {

		stmt, stmterr := tx.Prepare("insert into commits values (?, ?, ?, ?, ?, ?, ?, ?)")

		if stmterr != nil {
			log.Fatal(stmterr)
		}

		var parentA string
		var parentB string

		if len(o.parents) == 2 {
			parentA = o.parents[0]
			parentB = o.parents[1]
		} else if len(o.parents) == 1 {
			parentA = o.parents[0]
		}

        _, execerr := stmt.Exec(
			o.hash, 
			parentA, 
			parentB, 
			o.author, 
			o.committer, 
			o.date, 
			o.tree, 
			o.message,
		)

		if execerr != nil {
			log.Fatal(execerr)
		}

	}

	tx.Commit()

}

// utility function: insert blob object data
func (this SQLiteDatabase) writeBlobs(blobs []GitBlob) {

	db, openerr := sql.Open("sqlite3", this.fileName)

	if openerr != nil {
		log.Fatal(openerr)
	}
	
	defer db.Close()
	tx, _ := db.Begin()

	for _, o := range(blobs) {

		stmt, stmterr := tx.Prepare("insert into blobs values (?, ?)")

		if stmterr != nil {
			log.Fatal(stmterr)
		}

		_, execerr := stmt.Exec(o.hash, o.content)

		if execerr != nil {
			log.Fatal(execerr)
		}

	}

	tx.Commit()

}

// utility function: insert tree object data
func (this SQLiteDatabase) writeTrees(trees []GitTree) {

	db, openerr := sql.Open("sqlite3", this.fileName)

	if openerr != nil {
		log.Fatal(openerr)
	}
	
	defer db.Close()
	tx, _ := db.Begin()

	for _, o := range(trees) {

		for _, f := range(o.files){

			stmt, stmterr := tx.Prepare("insert into trees values (?, ?, ?, ?)")

			if stmterr != nil {
				log.Fatal(stmterr)
			}

			_, execerr := stmt.Exec(o.hash, f.hash, f.name, f.mode)

			if execerr != nil {
				log.Fatal(execerr)
			}

		}

	}

	tx.Commit()

}
