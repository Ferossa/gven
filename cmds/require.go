package cmds

import (
	"../structs"
	"log"
	"path"
	"strings"
)

func init() {
	registerCommand("require", new(RequireCommand))
}

type RequireCommand struct {
	Command
}

func (c *RequireCommand) RequireConfig() bool {
	return false
}

func (c *RequireCommand) Run(ctx *structs.Context) {
	args := c.console.Args()
	log.Printf("Require %#v", args)

	// initializing config if necessary
	if ctx.Config == nil {
		log.Println("Initializing project")
		ctx.Config = new(structs.Config)
		err := ctx.Config.Init(ctx.ProjectPath)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Project initialized")
	}

	// parsing command line flags
	targetsFlag := c.console.Flag("t", "")
	isDev := c.console.FlagBool("d", false)

	// checking arguments
	var targets []string
	if targetsFlag == "" {
		targets = ctx.GetTargets()
	} else {
		targets = strings.Split(targetsFlag, ",")
	}

	if len(targets) == 0 {
		log.Printf("No targets to require")
	}

	args = args[1:]
	if len(args) == 0 {
		c.ShortInfo()
		return
	}

	dep := strings.SplitN(args[0], ":", 2)
	if len(dep) == 1 {
		dep = append(dep, "*")
	}

	var repoUrl, repoType string
	if len(args) > 1 {
		repoUrl = args[1]
		if len(args) > 2 {
			repoType = args[2]
		}
	}

	for _, target := range targets {
		info, ok := ctx.Config.Targets[target]
		if !ok {
			binary := target
			if target == "." {
				binary = path.Base(ctx.ProjectPath)
			}
			info = structs.Target{
				Output:       path.Join("bin", binary),
				Dependencies: make(map[string]string),
				Development:  make(map[string]string),
			}
		}

		depName := dep[0]
		depVersion := dep[1]
		if isDev {
			info.Development[depName] = depVersion
			// remove dep if exists
			delete(info.Dependencies, depName)
		} else {
			info.Dependencies[depName] = depVersion
			// remove dev if exists
			delete(info.Development, depName)
		}

		if repoUrl != "" {
			repo := structs.Repository{
				Url:  repoUrl,
				Type: repoType,
			}
			ctx.Config.Repositories[depName] = repo
		}

		ctx.Config.Targets[target] = info
	}

	if err := ctx.Config.Save(ctx.ProjectPath); err != nil {
		log.Printf("error saving config: %v", err)
	}

	log.Println("Dependency added")
}

func (c *RequireCommand) ShortInfo() string {
	return "require [-t targets] package:version [url] [type] -- add dependency to targets"
}
