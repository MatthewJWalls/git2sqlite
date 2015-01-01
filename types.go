
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

type GitReference struct {
	path string
	hash string
}

