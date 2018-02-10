package cmds

import (
	"github.com/ferossa/gven/libs"
	"github.com/ferossa/gven/structs"
	"log"
	"os"
	"path"
	"strings"
)

func init() {
	registerCommand("build", new(BuildCommand))
}

type BuildCommand struct {
	Command
}

func (c *BuildCommand) Run(ctx *structs.Context) {
	args := c.console.Args()
	targets := args[1:]
	log.Printf("Build %v", strings.Join(targets, " "))

	if len(targets) == 0 {
		targets = ctx.GetTargets()
	}

	for _, t := range targets {
		log.Println("Building target", t)

		tc := ctx.Config.Targets[t]

		buildPath := ""
		switch {
		case path.IsAbs(tc.Output):
			buildPath = tc.Output
		case strings.HasPrefix(tc.Output, "${GOPATH}"):
			buildPath = ctx.GoPath + tc.Output[9:]
		default:
			buildPath = path.Join(ctx.ProjectPath, tc.Output)
		}

		command := &libs.Command{}
		if !ctx.Config.OverrideGopath {
			command.Env = map[string]string{"GOPATH": os.Getenv("GOPATH")}
		}
		_, err := command.Exec("go", []string{"build", "-o", buildPath, t})
		if err != nil {
			log.Println(err)
		} else {
			log.Println(t, "built")
		}
	}

	log.Println("All targets built")
}

func (c *BuildCommand) ShortInfo() string {
	return "build [targets] -- builds specified targets"
}
