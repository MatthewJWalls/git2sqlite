package main

/*
    Functions for reading the contents of the Git object database.

    Git has three kinds of object: Commits, Blobs (files) and Trees (directories).
    They are stored one-per-file under .git/objects, and each one is compressed
    with zlib.

    All Git objects have a header and a content section. The header just states 
    what kind of object it is (tree/blob/commit), and the content is different
    for each kind of object.

    We also support the reading of references (which are things like branches). A
    reference just points to a commit hash.

    usage of this class:

        repo = GitRepository{"some/repo/.git"}

        blobs, trees, commits := repo.Objects()
        refs := repo.References()

*/

import (
	"log"
	"fmt"
	"bytes"
	"strings"
	"io/ioutil"
	"path/filepath"
	"compress/zlib"
)

type GitRepository struct {
	location string
}

// returns a list of all object hashes in the repo
func (this GitRepository) Hashes() []string {

	hashes := make([]string, 0, 20)

	dirs, err := ioutil.ReadDir(filepath.Join(this.location, "objects"))

	if err != nil {
		log.Fatal(err)
	}

	// for all directories in objects directory

	for _, dir := range dirs {

		// ignore if they're not 2 characters

		if len(dir.Name()) != 2 {
			continue
		}

		// list all files in that directory

		minorpath := filepath.Join(this.location, "objects", dir.Name())
		objectFiles, err2 := ioutil.ReadDir(minorpath)
		
		if err2 != nil {
			log.Fatal(err)
		}
		
		// add the hash

		for _, o := range objectFiles {
			hash := dir.Name() + o.Name()
			hashes = append(hashes, hash)
		}
			
	}

	return hashes

}

// returns all the blobs, trees and commit objects in a git repo
func (this GitRepository) Objects() ([]GitBlob, []GitTree, []GitCommit) {

	blobs := make([]GitBlob, 0, 10)
	trees := make([]GitTree, 0, 10)
	commits := make([]GitCommit, 0, 10)

	objectToBlob := func(o GitObject) GitBlob {
		return GitBlob{o.path, o.hash, o.content}
	}

	objectToTree := func(o GitObject) GitTree {

		// tree format is <mode> <filename><null byte><20 byte sha1>

		files := make([]string, 0, 0)
		peg := 0

		for i := 0; i < len(o.content); i+=1 {
			if o.content[i] == '\x00' {
				filename := o.content[peg+1:i]
				files = append(files, strings.Split(fmt.Sprintf("%s", filename), " ")[1])
				i += 20;
				peg = i;
			}
		}

		return GitTree{o.path, o.hash, files}

	}

	objectToCommit := func(o GitObject) GitCommit {

		top := strings.SplitN(o.content, "\n\n", 2)[0]
		msg := strings.SplitN(o.content, "\n\n", 2)[1]

		commit := GitCommit{path: o.path, hash: o.hash, message:msg, parents:make([]string, 0, 2)}

		for _, line := range(strings.Split(top, "\n")) {

			parts := strings.Split(line, " ")

			switch parts[0] {
			case "tree":
				commit.tree = parts[1]
			case "parent":
				commit.parents = append(commit.parents, parts[1])
			case "author":
				commit.author = parts[1]
			case "committer":
				commit.committer = parts[1]
				commit.date = parts[len(parts)-2]
			}			

		}

		return commit

	}

	for _, h := range this.Hashes() {
		
		path := filepath.Join(this.location, "objects", h[:2], h[2:])
		content := this.getObjectContents(path)
		
		parts := strings.SplitN(content.String(), "\x00", 2)
		kind := strings.Split(parts[0], " ")[0]
		
		object := GitObject{path:path, hash:h, kind:kind, content:parts[1] }

		switch kind {
		case "blob":
			blobs = append(blobs, objectToBlob(object))
		case "tree":
			trees = append(trees, objectToTree(object))
		case "commit":
			commits = append(commits, objectToCommit(object))
		}

	}
	
	return blobs, trees, commits

}

func (this GitRepository) getObjectContents(objectlocation string) bytes.Buffer {

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

func (this GitRepository) References() []GitReference {

	references := make([]GitReference, 0, 10)

	heads, err := ioutil.ReadDir(filepath.Join(this.location, "refs", "heads"))

	if err != nil {
		log.Fatal(err)
	}

	for _, head := range heads {

		path := filepath.Join(this.location, "refs", "heads", head.Name())
		name := filepath.Join("heads", head.Name())

		rawbytes, readerr := ioutil.ReadFile(path)

		if readerr != nil {
			log.Fatal(readerr)
		}

		reference := GitReference{name, strings.TrimSpace(string(rawbytes))}
		references = append(references, reference)
		
	}

	return references

}

