package libs

import (
	"bytes"
	"os"
	"os/exec"
	"log"
)

type Command struct {
	Env map[string]string
}

func (c *Command) Exec(name string, args []string) (out string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	command := exec.Command(name, args...)
	command.Env = append(command.Env, "PATH="+os.Getenv("PATH"))
	for k, v := range c.Env {
		command.Env = append(command.Env, k + "=" + v)
	}
	if _, ok := c.Env["GOPATH"]; !ok {
		command.Env = append(command.Env, "GOPATH="+cwd)
	}

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	command.Stdout = &stdOut
	command.Stderr = &stdErr

	err = command.Run()
	if err != nil {
		log.Println(stdErr.String())
		return
	}

	if stdOut.Len() > 0 {
		out = stdOut.String()
	}
	return
}