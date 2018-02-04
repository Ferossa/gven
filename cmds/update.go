package cmds

import (
	"github.com/ferossa/gven/structs"
	"github.com/ferossa/gven/vcs"
	"log"
	"path"
	"strings"
)

func init() {
	registerCommand("update", new(UpdateCommand))
}

type UpdateCommand struct {
	Command
}

func (c *UpdateCommand) Run(ctx *structs.Context) {
	args := c.console.Args()
	pkgs := args[1:]
	log.Printf("Update %s", strings.Join(pkgs, " "))

	all := ctx.GetDependencies()
	if len(pkgs) == 0 {
		pkgs = all
	}

	// make map of dependencies
	allMap := make(map[string]bool)
	for _, d := range all {
		allMap[d] = true
	}

	for _, pkg := range pkgs {
		if _, ok := allMap[pkg]; !ok {
			log.Printf("Unknown dependency %s", pkg)
			continue
		}

		// get package source
		var repoUrl string
		var repoType string
		if repoInfo, ok := ctx.Config.Repositories[pkg]; ok {
			repoUrl = repoInfo.Url
			repoType = repoInfo.Type
		} else {
			repoUrl = "https://" + pkg
		}

		if repoType == "" {
			repoType = "git"
		}

		log.Printf("Updating %s from %s repository %s", pkg, repoType, repoUrl)

		for name, target := range ctx.Config.Targets {
			// check if target has dependency
			ver, ok := target.Dependencies[pkg]
			if !ok {
				ver, ok = target.Development[pkg]
				if !ok {
					continue
				}
			}

			// update package
			sourcePath := ctx.ProjectPath
			if name != "." {
				sourcePath = path.Join(ctx.ProjectPath, "src", name)
			}
			pkgDir := path.Join(sourcePath, "vendor", pkg)

			cs, ok := vcs.VCSRepository[repoType]
			if !ok {
				log.Printf("Unknown repository type %s", repoType)
				return
			}

			cs.SetConsole(c.console)
			err := cs.Update(repoUrl, pkgDir, ver)
			if err != nil {
				log.Printf("Error updating to version %s: %v", ver, err)
				return
			}
		}
	}

	log.Println("All dependencies updated")
}

func (c *UpdateCommand) ShortInfo() string {
	return "update [packages] -- update all or specified dependencies"
}
