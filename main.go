
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

	repo := GitRepository{REPOSITORY_NAME}

	// objects

	blobs, trees, commits := repo.Objects()

	for _, o := range(blobs) {
		log.Printf("%s (blob)", o.hash)
	}

	for _, o := range(trees) {
		log.Printf("%s (tree)", o.hash)
	}

	for _, o := range(commits) {
		log.Printf("%s (commit)", o.hash)
	}

	//writeObjectsToSQLite(objects, DATABASE_NAME)

	// references

	log.Println()
	log.Println("References")
	log.Println()

	for _, o := range(repo.References()) {
		log.Printf("%s (%s)", o.hash, o.path)
	}

	//writeRefsToSQLite(references, DATABASE_NAME)

}

