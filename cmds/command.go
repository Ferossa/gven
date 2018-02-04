package cmds

import (
	"github.com/ferossa/gven/libs"
	"github.com/ferossa/gven/structs"
)

// CommandRepository map where all commands stored
var CommandsRepository map[string]ICommand

func registerCommand(name string, cmd ICommand) {
	if CommandsRepository == nil {
		CommandsRepository = make(map[string]ICommand)
	}
	CommandsRepository[name] = cmd
}

// ICommand interface for all commands
type ICommand interface {
	SetConsole(con libs.IConsole) // set console integration
	RequireConfig() bool          // indicates that command runs only on initialized project
	ShortInfo() string            // gives short one line info for usage output
	Info() string                 // gives full info for help command
	Run(*structs.Context)         // execute command
}

type Command struct {
	console libs.IConsole
}

func (c *Command) SetConsole(con libs.IConsole) {
	c.console = con
}

func (c *Command) RequireConfig() bool {
	return true
}

func (c *Command) ShortInfo() string {
	return "not implemented"
}

func (c *Command) Info() string {
	return "not implemented"
}
