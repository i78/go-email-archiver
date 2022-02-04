package main

import (
	"email-archiver-cli/internal/cli"
	"fmt"
	"github.com/alecthomas/kong"
)

var CLI struct {
	Init         cli.InitCommand        `cmd help:"Initializes a Mail Repository in the current working directory" json:"init,omitempty"`
	GenerateKeys cli.GenerateKeyCommand `cmd help:"Generate or rotate encryption keys" json:"generate_keys"`
	Verbose      verboseFlag            `help:"Enable debug logging." json:"verbose,omitempty"`
}

type verboseFlag bool

func main() {
	fmt.Println("Go")
	ctx := kong.Parse(&CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
