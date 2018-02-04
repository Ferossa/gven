package vcs

import (
	"../libs"
)

var VCSRepository map[string]IVersionControlSystem

func registerVCS(name string, vcs IVersionControlSystem) {
	if VCSRepository == nil {
		VCSRepository = make(map[string]IVersionControlSystem)
	}
	VCSRepository[name] = vcs
}

// IVersionControlSystem interface for version control systems
type IVersionControlSystem interface {
	SetConsole(con libs.IConsole) // set console integration
	Update(url string, dir string, version string) error
}

type VersionControlSystem struct {
	console libs.IConsole
}

func (v *VersionControlSystem) SetConsole(con libs.IConsole) {
	v.console = con
}
