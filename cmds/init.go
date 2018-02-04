package cmds

import (
	"github.com/ferossa/gven/structs"
	"log"
)

func init() {
	registerCommand("init", new(InitCommand))
}

type InitCommand struct {
	Command
}

func (c *InitCommand) RequireConfig() bool {
	return false
}

func (c *InitCommand) Run(ctx *structs.Context) {
	log.Println("Init")

	if ctx.Config != nil {
		log.Println("Already initialized")
		return
	}

	ctx.Config = new(structs.Config)
	err := ctx.Config.Init(ctx.ProjectPath)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Project initialized")
}

func (c *InitCommand) ShortInfo() string {
	return "init -- initialize project in current directory"
}
