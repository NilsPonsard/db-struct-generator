package commands

import (
	cli "github.com/jawher/mow.cli"
	"github.com/nilsponsard/db-struct-generator/internal/commands/generate"
)

// configure subcommands
func SetupCommands(app *cli.Cli) {
	app.Command("generate", "generate", generate.Generate)
}
