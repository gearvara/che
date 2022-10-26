package gearvarabot

import (
	"flag"
	"fmt"
	stdlog "log"
	"os"

	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	"github.com/btwiuse/gearvarabot"
)

const commandName = "gearvarabot"

var (
	config = struct {
		VsCode bool
	}{}

	log = stdlog.New(os.Stderr, commandName+" ", 0)
)

func init() {
	caddycmd.RegisterCommand(caddycmd.Command{
		Name:  commandName,
		Func:  run,
		Usage: "[--vscode]",
		Short: "Gearvara Bot",
		Long: `
Gearvara Bot

Requires the TELEGRAM_BOT_TOKEN environment variable to be set.
`,
		Flags: func() *flag.FlagSet {
			fs := flag.NewFlagSet(commandName, flag.ExitOnError)
			fs.BoolVar(&config.VsCode, "vscode", config.VsCode, "Generate VSCode configuration")
			return fs
		}(),
	})
}

func run(fs caddycmd.Flags) (int, error) {
	w := "vscode"
	if config.VsCode {
		w = "true"
	} else {
		w = "false"
	}

	fmt.Println("hello", w)

	gearvarabot.Main()

	return 0, nil
}
