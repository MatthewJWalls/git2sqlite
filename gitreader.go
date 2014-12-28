package main

import (
	"log"
	"bytes"
	"strings"
	"io/ioutil"
	"path/filepath"
	"compress/zlib"
)

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

