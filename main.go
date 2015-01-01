
/*
    git2sqlite - converts git repositories to sqlite databases.

    When ran against a git repository, will output an sqlite database
    with the following tables:

        refs :: <path, hash>
        blobs :: <hash, content>
        trees :: <hash, content>
        commits :: <hash, content>

    This project is a work in progress.
*/

package main

import (
	"log"
)

func main() {

	// Just hacking out a prototype for now as a feasibility study.

	log.Println("Starting.")

	DATABASE_NAME := "testing.db"
	REPOSITORY_NAME := ".git"

	initDatabase(DATABASE_NAME)

	// objects

	objects := getObjects(REPOSITORY_NAME)

	for i, o := range(objects) {
		log.Printf("%4d| %s (%s)", i, o.hash, o.kind)
	}

	writeObjectsToSQLite(objects, DATABASE_NAME)

	// references

	log.Println()
	log.Println("References")
	log.Println()

	references := getReferences(REPOSITORY_NAME)

	for i, o := range(references) {
		log.Printf("%4d| %s (%s)", i, o.hash, o.path)
	}

	writeRefsToSQLite(references, DATABASE_NAME)

}

