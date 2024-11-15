package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v2/ffcli"
	"github.com/runetale/yggdrasill/cli/yggdrasill"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	// add sub command
	if len(args) == 1 && (args[0] == "-V" || args[0] == "--version" || args[0] == "-v") {
		args = []string{"version"}
	}

	fs := flag.NewFlagSet("yggdrasill", flag.ExitOnError)
	cmd := &ffcli.Command{
		Name:       "yggdrasill",
		ShortUsage: "yggdrasill <subcommands> [command flags]",
		ShortHelp:  "",
		LongHelp:   "",
		Subcommands: []*ffcli.Command{
			yggdrasill.NotchCmd,
		},
		FlagSet: fs,
		Exec:    func(context.Context, []string) error { return flag.ErrHelp },
	}

	if err := cmd.Parse(args); err != nil {
		return err
	}

	if err := cmd.Run(context.Background()); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	return nil
}
