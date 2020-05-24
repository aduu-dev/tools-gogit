package install

import (
	"fmt"
	"path/filepath"

	"aduu.dev/utils/helper"
	"k8s.io/klog/v2"
)

// Remove removes the line containing defaultBashComment from the pre-commit
// and post-commit hooks residing under the base path.
func Remove(base string) (err error) {
	hooksFolder := filepath.Join(base, hooksPath())

	// The hooks folder must exist.
	exists, err := helper.DoesPathExistErr(hooksFolder)
	if err != nil {
		return
	}

	if !exists {
		return fmt.Errorf("the hooks folder does not exist")
	}

	if err = EnsureRemoveComment(preCommitFilepath(base), defaultBashComment); err != nil {
		return
	}

	if err = EnsureRemoveComment(postCommitFilepath(base), defaultBashComment); err != nil {
		return
	}

	klog.InfoS("Removed gogit replace lines",
		"from-pre-commit", preCommitFilepath(base),
		"from-post-commit", postCommitFilepath(base),
	)

	return nil
}
