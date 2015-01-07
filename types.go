
/*
    Data structures to store the contents of the Git object database.
*/

package main

type GitObject struct {
	path string
	hash string
	kind string
	content string
}

type GitCommit struct {
	path string
	hash string
	parents []string
	author string
	committer string
	date string
	tree string
	message string
}

type GitTree struct {
	path string
	hash string
	files []GitTreeFile
}

type GitBlob struct {
	path string
	hash string
	content string
}

type GitTreeFile struct {
	mode string
    name string
	hash string
}

type GitReference struct {
	path string
	hash string
}
