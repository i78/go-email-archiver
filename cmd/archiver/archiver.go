package main

import (
	"email-archiver-cli/internal/cli"
	"github.com/alecthomas/kong"
)

var CLI struct {
	Init         cli.InitCommand        `cmd:"" help:"Initializes a Mail Repository in the current working directory"`
	GenerateKeys cli.GenerateKeyCommand `cmd:"" help:"Generate or rotate encryption keys"`
	Config       cli.ConfigCommand      `cmd:"" help:"Generate or update configurations"`
	Archive      cli.ArchiveCommand     `cmd:"" help:"Start archiving emails"`
	Verbose      verboseFlag            `help:"Enable debug logging."`
}

type verboseFlag bool

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
