package main

// git2sqlite.
//
// Takes a git repository and turns it into an sqlite database.


import (
	"log"
	"io/ioutil"
	"path/filepath"
)

func getObjects(gitlocation string) []string {

	// given the location of a .git directory, reads up the objects 
	// inside of that directory.

	dirs, err := ioutil.ReadDir(filepath.Join(gitlocation, "objects"))

	if err != nil {
		log.Fatal(err)
	}

	objects := make([]string, 0, 50)

	for _, dir := range dirs {

		if len(dir.Name()) == 2 {

			prefix := dir.Name()
			objectPath := filepath.Join(gitlocation, "objects", prefix)

			objectFiles, err2 := ioutil.ReadDir(objectPath)

			if err2 != nil {
				log.Fatal(err)
			}

			for _, o := range objectFiles {
				objects = append(objects, prefix+o.Name())
			}
			
		}

	}

	return objects

}

func main() {

	// Just hacking out a prototype for now as a feasibility study.

	log.Println("Starting.")

	for i, s := range(getObjects(".git")) {
		log.Printf("%d) %s", i, s)
	}

}
