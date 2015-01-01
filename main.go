package main

// git2sqlite.
//
// Takes a git repository and turns it into an sqlite database.


import (
	"log"
)

type GitObject struct {
	path string
	hash string
	kind string
	content string
}

type GitReference struct {
	path string
	hash string
}

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

}

