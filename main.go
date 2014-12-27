package main

// git2sqlite.
//
// Takes a git repository and turns it into an sqlite database.


import (
	"log"
	"bytes"
	"strings"
	"io/ioutil"
	"path/filepath"
	"compress/zlib"
)

type GitObject struct {
	path string
	hash string
	kind string
	content string
}

func getObjectContents(objectlocation string) bytes.Buffer {

	// given a path to a git object, returns the contents of the object
	// as a string.

	// read up the file into a buffer

	rawbytes, readerr := ioutil.ReadFile(objectlocation)

	if readerr != nil {
		log.Fatal(readerr)
	}

	buff := bytes.NewReader(rawbytes)

	// ok, now we can use zlib to decompress it.

	r, compresserr := zlib.NewReader(buff)

	if compresserr != nil {
		log.Fatal(compresserr)
	}

	// create & return an output buffer of the decompressed data

	outbuffer := new(bytes.Buffer)
	_, err := outbuffer.ReadFrom(r)

	if err != nil {
		log.Fatal(err)
	}

	return *outbuffer

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

				content := getObjectContents(fullpath)

				// the contents of a git object look like:
				// <type> <length><\x00 byte><content...>
				parts := strings.Split(content.String(), "\x00")
				headerparts := strings.Split(parts[0], " ")

				object := GitObject{path:fullpath, hash:fullhash, kind: headerparts[0], content:parts[1] }
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
		log.Printf("    %s", s.kind)
	}

}
