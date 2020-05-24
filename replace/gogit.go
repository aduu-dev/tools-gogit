package replace

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"aduu.dev/utils/helper"
	"github.com/go-git/go-git/v5"
	copy2 "github.com/otiai10/copy"
	"golang.org/x/mod/modfile"
	"k8s.io/klog/v2"
)

var (
	errPathDoesNotExist   = fmt.Errorf("path does not exist")
	errBackupDoesNotExist = fmt.Errorf("backup file go.mod.b does not exist")
	errBackupExists       = fmt.Errorf("backup file go.mod.b does exist already")
)

func backupFilename() string {
	return "go.mod.b"
}

func turnToNonPointerSlice(reps []*modfile.Replace) (out []modfile.Replace) {
	out = make([]modfile.Replace, 0, len(reps))

	for _, rep := range reps {
		out = append(out, *rep)
	}

	return out
}

func getGoModFilepathAndData(base string) (goModFilepath string, data []byte, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("something is wrong with the go.mod directory: %w", err)
		}
	}()

	path, err := filepath.Abs(base)
	if err != nil {
		return
	}

	exists, err := helper.DoesPathExistErr(path)
	if err != nil {
		return
	}

	if !exists {
		return "", nil, errPathDoesNotExist
	}

	goModFilepath = filepath.Join(path, "go.mod")
	exists, err = helper.DoesPathExistErr(goModFilepath)
	if err != nil {
		return
	}

	if !exists {
		return "", nil, fmt.Errorf("error: no go.mod file in the given path %#v", goModFilepath)
	}

	data, err = ioutil.ReadFile(goModFilepath)
	if err != nil {
		return
	}

	return
}

func removeLocalReplaceDirectives(directives []*modfile.Replace) (localReplaces []*modfile.Replace) {
	for _, rep := range directives {
		if isLocalDirective(rep) {
			localReplaces = append(localReplaces, rep)
		}
	}

	return localReplaces
}

func isLocalDirective(rep *modfile.Replace) bool {
	newPath := rep.New.Path

	return strings.HasPrefix(newPath, "..") || strings.HasPrefix(newPath, "./")
}

// isGomodStaged returns true if go.mod is staged
func isGomodStaged(base string) (staged bool, err error) {
	r, err := git.PlainOpen(base)
	if err != nil {
		return
	}

	w, err := r.Worktree()
	if err != nil {
		return
	}

	status, err := w.Status()
	if err != nil {
		return
	}

	gomodStatus := status.File("go.mod")

	switch gomodStatus.Staging {
	// In these states I assume the intention is to commit the modfied go.mod.
	case git.Added, git.Copied, git.Modified, git.Renamed, git.UpdatedButUnmerged:
		return true, nil
	default:
		return false, nil
	}
}

// RemoveLocalReplacesFromGomod removes go.mod replace directives which are pointing to local folders.
func RemoveLocalReplacesFromGomod(arg string, workOnStagedOnly bool) (err error) {
	backup := filepath.Join(arg, backupFilename())

	exists, err := helper.DoesPathExistErr(backup)
	if err != nil {
		return
	}

	if exists {
		return errBackupExists
	}

	gomodFilepath, data, err := getGoModFilepathAndData(arg)

	if err != nil {
		return
	}

	if len(gomodFilepath) == 0 {
		return fmt.Errorf("goModFilepath is not set")
	}
	// Create backup.
	klog.InfoS("Creating backup", "from", gomodFilepath, "backup", backup)

	if err = copy2.Copy(gomodFilepath, backup); err != nil {
		return fmt.Errorf("failed to create backup at %#v: %w", backup, err)
	}

	file, err := modfile.Parse(gomodFilepath, data, nil)
	if err != nil {
		return fmt.Errorf("failed to parse modfile at %#v: %w", gomodFilepath, err)
	}

	// If we do not require go.mod to be staged we can start removing the local directives immediately.
	if !workOnStagedOnly {
		err = removeLocalDirectivesInFile(file, gomodFilepath)
		if err != nil {
			return err
		}
	} else {
		// Find out the staging status of go.mod.
		isGomodStaged, err := isGomodStaged(arg)
		if err != nil {
			return err
		}

		// If we require go.mod to be staged and go.mod is also staged then remove the replace directives.
		if isGomodStaged {
			err = removeLocalDirectivesInFile(file, gomodFilepath)
			if err != nil {
				return err
			}
		}
	}

	klog.InfoS("Finished removing local replace directives", "go.mod", gomodFilepath, "backup", backup)

	return nil
}

func removeLocalDirectivesInFile(file *modfile.File, gomodFilepath string) (err error) {
	localReplaces := removeLocalReplaceDirectives(file.Replace)

	// Only write out a modified version in case we actually removed a replace directive.
	//
	// Else we create unnecessary noise inside git commits.
	if len(localReplaces) > 0 {
		for _, localReplace := range localReplaces {
			if err = file.DropReplace(localReplace.Old.Path, localReplace.Old.Version); err != nil {
				return err
			}
		}

		var dataOut []byte
		dataOut, err = file.Format()
		if err != nil {
			return err
		}

		// Write out modified "go.mod".
		if err = ioutil.WriteFile(gomodFilepath, dataOut, 0755); err != nil {
			return fmt.Errorf("failed to write modified go.mod file to %#v", gomodFilepath)
		}
	}

	return nil
}

// UndoRemovingLocalReplacesFromGomod replaces the local go.mod with the backup.
func UndoRemovingLocalReplacesFromGomod(arg string, workOnStaged bool) (err error) {
	// Run tests and get go.mod filepath.
	goModFilepath, _, err := getGoModFilepathAndData(arg)
	if err != nil {
		return
	}

	backup := filepath.Join(arg, backupFilename())

	// Check backup exists.
	exists, err := helper.DoesPathExistErr(backup)
	if err != nil {
		return
	}

	if !exists {
		return errBackupDoesNotExist
	}

	// Copy from backup to "go.mod".
	if err = copy2.Copy(backup, goModFilepath); err != nil {
		return
	}

	// Remove backup.
	if err = os.Remove(backup); err != nil {
		return
	}

	klog.InfoS("Undid local go.mod change", "go.mod", goModFilepath, "backup(removed)", backup)

	return nil
}
