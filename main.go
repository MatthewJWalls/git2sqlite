package main

// git2sqlite.
//
// Takes a git repository and turns it into an sqlite database.


import (
	"os"
	"log"
	"bytes"
	"strings"
	"io/ioutil"
	"path/filepath"
	"compress/zlib"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type GitObject struct {
	path string
	hash string
	kind string
	content string
}

func initDatabase(dbname string) {

	// executes the required sqlite DML statements

	os.Remove(dbname)

	db, openerr := sql.Open("sqlite3", dbname)

	if openerr != nil {
		log.Fatal(openerr)
	}

	defer db.Close()

	stmt := "create table commits (hash text not null primary key, content text)"
	_, stmterr := db.Exec(stmt)

	if stmterr != nil {
		log.Fatal(stmterr)
	}

}

func writeObjectsToSQLite(objects []GitObject, dbname string) {

	// executes the required sqlite DML statements

	db, openerr := sql.Open("sqlite3", dbname)

	if openerr != nil {
		log.Fatal(openerr)
	}

	defer db.Close()

	tx, _ := db.Begin()

	stmt, stmterr := tx.Prepare("insert into commits (hash, content) values (?, ?)")

	if stmterr != nil {
		log.Fatal(stmterr)
	}

	for _, o := range(objects) {
		if o.kind == "commit" {
			_, execerr := stmt.Exec(o.hash, o.content)

			if execerr != nil {
				log.Fatal(execerr)
			}

		}
	}

	tx.Commit()

}

func getObjectContents(objectlocation string) bytes.Buffer {

	// given a path to a git object, returns the contents of the object
	// as a buffer.

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

	// git takes the first 2 characters of the object hash, and puts the
	// file under a directory with the name of the 2 character prefix. The
	// object filename is the hash, minus the first two characters. For
	// example an object with hash "abcdef" would be written to location:
	// .git/objects/ab/cdef

	for _, dir := range dirs {

		// ignore anything other than the 2-character prefix directories

		if len(dir.Name()) != 2 {
			continue
		}

		prefix := dir.Name()
		prefixPath := filepath.Join(gitlocation, "objects", prefix)
		
		objectFiles, err2 := ioutil.ReadDir(prefixPath)
		
		if err2 != nil {
			log.Fatal(err)
		}

		for _, o := range objectFiles {
			fullhash := prefix + o.Name()
			fullpath := filepath.Join(prefixPath, o.Name())

			content := getObjectContents(fullpath)

			// the contents of a git object look like:
			// <type> <length><\x00 byte><content...>
			parts := strings.Split(content.String(), "\x00")
			headerparts := strings.Split(parts[0], " ")

			object := GitObject{path:fullpath, hash:fullhash, kind: headerparts[0], content:parts[1] }
			objects = append(objects, object)
		}
			
	}

	return objects

}

func main() {

	// Just hacking out a prototype for now as a feasibility study.

	log.Println("Starting.")

	DATABASE_NAME := "testing.db"
	REPOSITORY_NAME := ".git"

	initDatabase(DATABASE_NAME)

	objects := getObjects(REPOSITORY_NAME)

	for i, s := range(objects) {
		log.Printf("%d) %s", i, s.path)
		log.Printf("    %s", s.hash)
		log.Printf("    %s", s.kind)
	}

	writeObjectsToSQLite(objects, DATABASE_NAME)

}
