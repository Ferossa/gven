package vcs

import (
	"os"
	"strings"
)

func init() {
	registerVCS("git", new(Git))
}

// Git git version control system
type Git struct {
	VersionControlSystem
}

func (v *Git) Update(url string, dir string, version string) (err error) {
	// preparations
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		// if dir is empty - clone repository
		_, err = v.console.Exec("", "git", []string{"clone", url, dir})
	} else {
		// or fetch all updates
		_, err = v.console.Exec(dir, "git", []string{"fetch"})
	}

	// check tag
	var res string
	if version != "*" {
		res, err = v.console.Pipe(dir, "git", []string{"ls-remote", "-q", "-t", "-h"}, "grep", []string{version})
		if err != nil {
			return
		}

		lines := strings.Split(strings.TrimSpace(res), "\n")
		if len(lines) == 1 {
			// found exactly what we need
			parts := strings.SplitN(strings.TrimSpace(lines[0]), "\t", 2)
			if len(parts) == 1 {
				parts = append(parts, parts[0])
			}

			if parts[0] == version {
				// checking to commit
				res, err = v.console.Exec(dir, "git", []string{"checkout", parts[0]})
			} else {
				// checking to version
				checkoutVersion := strings.TrimPrefix(parts[1], "refs/heads/")
				if checkoutVersion == parts[1] {
					checkoutVersion = strings.TrimPrefix(checkoutVersion, "refs/tags/")
				}
				res, err = v.console.Exec(dir, "git", []string{"checkout", checkoutVersion})
			}
			return
		}
	}

	// if we didn't find exact match - trying to fetch by mask
	if version == "*" {
		res, err = v.console.Exec(dir, "git", []string{"tag", "--format=%(refname:strip=2)", "--sort=-v:refname", "-l"})
	} else {
		res, err = v.console.Exec(dir, "git", []string{"tag", "--format=%(refname:strip=2)", "--sort=-v:refname", "-l", version})
	}

	if err != nil {
		return
	}

	lines := strings.Split(strings.TrimSpace(res), "\n")
	found := strings.TrimSpace(lines[0])
	if found != "" {
		_, err = v.console.Exec(dir, "git", []string{"checkout", found})
	} else {
		_, err = v.console.Exec(dir, "git", []string{"pull"})
	}

	return
}
