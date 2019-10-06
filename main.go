package main

import (
	"encoding/json"
	"github.com/ferossa/gven/cmds"
	"github.com/ferossa/gven/libs"
	"github.com/ferossa/gven/structs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const version = "v0.1.4"

func main() {
	// intro log
	log.Printf("GoLang dependency manager gven %s", version)

	con := new(libs.Console)
	con.Parse()

	// command name
	var cn string
	args := con.Args()
	if len(args) > 0 {
		cn = args[0]
	}

	// get command by name
	if c, ok := cmds.CommandsRepository[cn]; ok {
		c.SetConsole(con)

		// init context
		ctx, err := initContext()
		if err != nil {
			log.Fatalf("Initialization error: %v", err)
			return
		}

		// read config
		var cfg *structs.Config
		cfg, err = loadConfig(ctx)
		if c.RequireConfig() && err != nil {
			log.Fatalf("Config loading error: %v", err)
			return
		}

		// prepare context
		ctx.Config = cfg

		// execute command
		c.Run(ctx)
		return
	}

	// TODO: show usage info
	log.Println("usage: gven command [args]")
}

func initContext() (ctx *structs.Context, err error) {
	ctx = new(structs.Context)

	projectPath, err := os.Getwd()
	if err != nil {
		return
	}

	ctx.ProjectPath = projectPath
	ctx.GoPath = os.Getenv("GOPATH")
	ctx.IsDevMode = os.Getenv("GODEV") == "1"

	return
}

func loadConfig(ctx *structs.Context) (cfg *structs.Config, err error) {
	configFile := filepath.Join(ctx.ProjectPath, "gven.json")
	if _, err = os.Stat(configFile); !os.IsNotExist(err) {
		var configData []byte
		configData, err = ioutil.ReadFile(configFile)
		if err != nil {
			return
		}

		cfg = new(structs.Config)
		err = json.Unmarshal(configData, cfg)
		if err != nil {
			cfg = nil
			return
		}
	}

	return
}
