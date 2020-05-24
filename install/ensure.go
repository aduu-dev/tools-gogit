package install

import (
	"io/ioutil"
)

// EnsureRemoveComment removes the line containing the given comment.
func EnsureRemoveComment(file string, comment string) (err error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	newContent, err := removeFromShellFile(string(content), comment)
	if err != nil {
		return
	}

	if err = ioutil.WriteFile(file, []byte(newContent), 0755); err != nil {
		return
	}

	return nil
}

// EnsureAddLinesWithComment ensures that the given line with the given comment is added to the file.
func EnsureAddLinesWithComment(file string, line string, comment string) (err error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	newContent, err := addToShellFile(string(content), line, comment)
	if err != nil {
		return
	}

	if err = ioutil.WriteFile(file, []byte(newContent), 0755); err != nil {
		return
	}

	return nil
}