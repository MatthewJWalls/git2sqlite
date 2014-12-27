package main

// git2sqlite.
//
// Takes a git repository and turns it into an sqlite database.


import (
	"log"
	"io/ioutil"
	"path/filepath"
)

type GitObject struct {
	path string
	hash string
	kind string
	content string
}

func getObjects(gitlocation string) []GitObject {

	// given the location of a .git directory, reads up the objects 
	// inside of that directory.

	dirs, err := ioutil.ReadDir(filepath.Join(gitlocation, "objects"))

	if err != nil {
		log.Fatal(err)
	}

	objects := make([]GitObject, 0, 50)

	for _, dir := range dirs {

		if len(dir.Name()) == 2 {

			prefix := dir.Name()
			objectPath := filepath.Join(gitlocation, "objects", prefix)

			objectFiles, err2 := ioutil.ReadDir(objectPath)

			if err2 != nil {
				log.Fatal(err)
			}

			for _, o := range objectFiles {
				fullhash := prefix + o.Name()
				fullpath := filepath.Join(objectPath, o.Name())
				object := GitObject{path:fullpath, hash:fullhash }
				objects = append(objects, object)
			}
			
		}

	}

	return objects

}

func main() {

	// Just hacking out a prototype for now as a feasibility study.

	log.Println("Starting.")

	for i, s := range(getObjects(".git")) {
		log.Printf("%d) %s", i, s.path)
		log.Printf("    %s", s.hash)
	}

}
