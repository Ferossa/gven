package structs

// Context project context settings
type Context struct {
	IsDevMode   bool
	GoPath      string
	ProjectPath string
	Config      *Config
}

// GetTargets returns all build targets
func (c *Context) GetTargets() (r []string) {
	for n, _ := range c.Config.Targets {
		r = append(r, n)
	}

	return
}

// GetDependencies returns all dependencies for current environment
func (c *Context) GetDependencies() (r []string) {
	for _, target := range c.Config.Targets {
		for n, _ := range target.Dependencies {
			r = append(r, n)
		}

		if !c.IsDevMode {
			continue
		}

		for n, _ := range target.Development {
			r = append(r, n)
		}
	}

	return
}
