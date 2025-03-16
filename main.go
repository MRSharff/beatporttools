package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

var (
	commands = map[string]func(args []string) error{
		"organize": organize,
	}
)

func main() {
	var (
		// log levels
		info  bool
		debug bool
	)

	flags := flag.NewFlagSet("beatporttools", flag.ExitOnError)
	// it would be cool to set up the -v flag to where I could instead get the amount of 'v's and then map that to the
	// log level, since that's a common way to handle log level flags, eg: -v, -vv, -vvv, etc. For now, just have info
	// and debug available as two separate flags.
	flags.BoolVar(&info, "v", false, "show info logs")
	flags.BoolVar(&debug, "vv", false, "show debug logs")
	flags.Usage = func() {
		w := flags.Output()
		fmt.Fprintln(w, "A tool for working with music files downloaded from Beatport")
		fmt.Fprintln(w, "Usage:")
		fmt.Fprintln(w, "\tbeatporttools <command> [arguments]")
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Global Flags:\n")
		flags.PrintDefaults()
		fmt.Fprintf(w, "Commands:\n")
		fmt.Fprintf(w, "\torganize\tReorganizes music downloaded from beatport")
	}
	flags.Parse(os.Args[1:])

	if len(os.Args) < 2 {
		flags.Usage()
		os.Exit(1)
	}

	logLevel := slog.LevelWarn
	if info {
		logLevel = slog.LevelInfo
	}
	if debug {
		logLevel = slog.LevelDebug
	}
	slog.SetLogLoggerLevel(logLevel)

	command, ok := commands[os.Args[1]]
	if !ok {
		flags.Usage()
		os.Exit(1)
	}

	command(os.Args[2:])

}

func organize(args []string) error {
	var (
		source   string
		dest     string
		noPrompt bool
	)

	flags := flag.NewFlagSet("organize", flag.ExitOnError)
	flags.StringVar(&source, "source", ".", "source directory, where your Beatport downloads are located")
	flags.StringVar(&dest, "dest", ".", "destination directory, where you want the release folders to be created")
	flags.BoolVar(&noPrompt, "y", false, "do not prompt for input, accept all prompts")

	flags.Usage = func() {
		w := flags.Output()
		fmt.Fprintln(w, "usage:")
		fmt.Fprintln(w, "\tbeatporttools organize [-source source] [-dest dest] [-y]")
		fmt.Fprintln(w, "flags:")
		flags.PrintDefaults()
		fmt.Fprintln(w, "example:")
		fmt.Fprintln(w, "\tbeatporttools organize -y -source ~/Downloads/beatport_tracks_2025_03 -dest ~/Downloads/beatport_tracks_2025_03_organized")
	}

	err := flags.Parse(args)
	if err != nil {
		return fmt.Errorf("error parsing flags for the organize command: %w", err)
	}
	organizeIntoReleaseFolders(source, dest, noPrompt)
	// todo: probably return an error from organizeIntoReleaseFolders
	return nil
}
