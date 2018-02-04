package libs

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type IConsole interface {
	Args() []string
	Flag(name string, def string) string
	FlagInt64(name string, def int64) int64
	FlagBool(name string, def bool) bool
	Exec(f, name string, args []string) (res string, err error)
	Pipe(f, name string, args []string, pipeTo string, argsTo []string) (res string, err error)
}

type Console struct {
	args  []string
	flags map[string]string
}

// Parse parse command line arguments
func (c *Console) Parse() {
	c.args = make([]string, 0)
	c.flags = make(map[string]string)

	args := os.Args[1:]

	var i int
	for {
		if i >= len(args) {
			break
		}

		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			// flag
			arg = strings.TrimLeft(arg, "-")
			argParts := strings.SplitN(arg, "=", 2)

			var value string
			if len(argParts) > 1 {
				value = argParts[1]
			}

			if strings.HasPrefix(value, `"`) {
				value = value[1 : len(value)-2]
			}
			c.flags[argParts[0]] = value
		} else {
			// argument
			c.args = append(c.args, arg)
		}

		i++
	}
}

func (c *Console) Args() []string {
	return c.args
}

func (c *Console) Flag(name string, def string) string {
	if v, ok := c.flags[name]; ok {
		return v
	}

	return def
}

func (c *Console) FlagInt64(name string, def int64) int64 {
	v, ok := c.flags[name]
	if !ok {
		return def
	}

	pv, _ := strconv.ParseInt(v, 10, 64)
	return pv
}

func (c *Console) FlagBool(name string, def bool) bool {
	v, ok := c.flags[name]
	if !ok {
		return def
	}

	pv, _ := strconv.ParseBool(v)
	return pv
}

func (c *Console) Exec(f, name string, args []string) (res string, err error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = f

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err = cmd.Run()
	if err != nil {
		res = stdErr.String()
		return
	}

	res = stdOut.String()
	return
}

func (c *Console) Pipe(f, name string, args []string, pipeTo string, argsTo []string) (res string, err error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = f

	pcmd := exec.Command(pipeTo, argsTo...)
	pcmd.Stdin, _ = cmd.StdoutPipe()

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stderr = &stdErr
	pcmd.Stdout = &stdOut
	pcmd.Stderr = &stdErr

	err = pcmd.Start()
	if err != nil {
		res = stdErr.String()
		return
	}

	err = cmd.Run()
	if err != nil {
		res = stdErr.String()
		return
	}

	err = pcmd.Wait()
	if err != nil {
		res = stdErr.String()
		return
	}

	res = stdOut.String()
	return
}
