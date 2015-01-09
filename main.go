
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
	"os"
	"path/filepath"
)

func main() {

	// check args

	var location string

	if len(os.Args) == 2 {
		location, _ = filepath.Abs(os.Args[1])
	} else {
		location, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	REPOSITORY_NAME := filepath.Join(location, ".git")
	DATABASE_NAME := filepath.Base(location)+".db"

	repo := GitRepository{REPOSITORY_NAME}
	db := SQLiteDatabase{DATABASE_NAME}

	db.Create()
	db.WriteRepository(repo)

	log.Println("Finished.")

}

