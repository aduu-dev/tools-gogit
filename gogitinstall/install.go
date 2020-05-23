package gogitinstall

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"aduu.dev/utils/helper"
	"k8s.io/klog/v2"
)

var (
	errHooksFolderDoesNotExist    = fmt.Errorf("hooks folder does not exist")
	errPreCommitDoesExistAlready  = fmt.Errorf("pre-commit does exist already")
	errPostCommitDoesExistAlreaxy = fmt.Errorf("post-commit exists already")
)

func hooksPath() string {
	return strings.ReplaceAll(".git/hooks", "/", string(filepath.Separator))
}

func preCommitFilepath(base string) string {
	return filepath.Join(base, hooksPath(), "pre-commit")
}

func postCommitFilepath(base string) string {
	return filepath.Join(base, hooksPath(), "post-commit")
}

func preCommitContent() []byte {
	return []byte(`#!/bin/bash

set -o xtrace

aduu gogit replace .
git add go.mod`)
}

func postCommitContent() []byte {
	return []byte(`#!/bin/bash

set -o xtrace

aduu gogit replace --undo .`)
}

// InstallHooks installs pre-commit hooks which do remove local replace directives temporarily during a commit.
func InstallHooks(base string) (err error) {
	hooksFolder := filepath.Join(base, hooksPath())

	exists, err := helper.DoesPathExistErr(hooksFolder)
	if err != nil {
		return
	}

	if !exists {
		return errHooksFolderDoesNotExist
	}

	exists, err = helper.DoesPathExistErr(preCommitFilepath(base))
	if err != nil {
		return err
	}

	if exists {
		return errPreCommitDoesExistAlready
	}

	exists, err = helper.DoesPathExistErr(postCommitFilepath(base))
	if err != nil {
		return
	}

	if exists {
		return errPostCommitDoesExistAlreaxy
	}

	if err = ioutil.WriteFile(preCommitFilepath(base), preCommitContent(), 0755); err != nil {
		return
	}

	if err = ioutil.WriteFile(postCommitFilepath(base), postCommitContent(), 0755); err != nil {
		return
	}

	klog.InfoS("Successuflly installed commit hooks")

	return nil
}
