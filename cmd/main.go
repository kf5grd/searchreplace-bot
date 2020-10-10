package cmd

import (
	"searchreplacebot/pkg/logr"
	"io"
	"os"
	"os/signal"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"samhofi.us/x/keybase/v2"
)

var version string

func Run(args []string, stdout io.Writer) error {
	flags := []cli.Flag{
		&cli.PathFlag{
			Name:    "config",
			Aliases: []string{"c"},
			EnvVars: []string{"BOT_CONFIG"},
			Usage:   "Load config from `FILE`",
			Value:   "config.toml",
		},
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			EnvVars: []string{"BOT_DEBUG"},
			Usage:   "Enable debug mode",
		}),
		altsrc.NewPathFlag(&cli.PathFlag{
			Name:    "home",
			Aliases: []string{"H"},
			EnvVars: []string{"BOT_HOME"},
			Usage:   "Set an alternate home directory for the Keybase client",
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"j"},
			EnvVars: []string{"BOT_JSON"},
			Usage:   "Output logs in JSON format",
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:    "replacer-basic",
			Aliases: []string{"r"},
			Usage:   "Replacer string in the format '|find this|replace with this'. The first character is the separator string. This flag can be specified more than once",
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:    "replacer-regex",
			Aliases: []string{"R"},
			Usage:   "Regular expression replacer string in the format '|find this|replace with this'. The first character is the separator string. This flag can be specified more than once",
		}),
	}

	app := cli.App{
		Name:                 "searchreplacebot",
		Version:              version,
		HideVersion:          false,
		Usage:                "A Keybase bot that does a search/replace on messages and sends a reply with the result",
		EnableBashCompletion: true,
		Writer:               stdout,
		Action:               run,
		Before:               altsrc.InitInputSourceWithContext(flags, altsrc.NewTomlSourceFromFlagFunc("config")),
		Flags:                flags,
	}

	if err := app.Run(args); err != nil {
		return err
	}

	return nil
}

func run(c *cli.Context) error {
	var b = bot{
		k:              keybase.New(keybase.SetHomePath(c.Path("home"))),
		log:            logr.New(c.App.Writer, c.Bool("debug"), c.Bool("json")),
		replacersBasic: c.StringSlice("replacer-basic"),
		replacersRegex: c.StringSlice("replacer-regex"),
	}

	for _, r := range b.replacersBasic {
		b.log.Debug("Basic replacer: '%s'", r)
	}

	for _, r := range b.replacersRegex {
		b.log.Debug("Regex replacer: '%s'", r)
	}

	b.advertiseCommands()

	// catch ctrl + c
	var trap = make(chan os.Signal, 1)
	signal.Notify(trap, os.Interrupt)
	go func() {
		for _ = range trap {
			b.log.Debug("Received interrupt signal")
			b.clearCommands()
			os.Exit(0)
		}
	}()

	b.registerHandlers()
	b.log.Info("Running as user %s", b.k.Username)
	b.k.Run(b.handlers, &b.opts)
	return nil
}
