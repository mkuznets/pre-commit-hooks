package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Define the top-level command and its options
	flag.Usage = func() {
		fmt.Printf("Usage: %s COMMAND [OPTIONS] [ARGS...]\n\n", filepath.Base(os.Args[0]))
		fmt.Printf("Commands:\n")
		fmt.Printf("  format-go")
		fmt.Println()
		flag.PrintDefaults()
	}
	flag.Parse()

	// Ensure that a subcommand is provided
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Determine the subcommand
	switch flag.Arg(0) {
	case "format-go":
		if err := formatGoCommand(flag.Args()[1:]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func formatGoCommand(args []string) error {
	// Define the options for the 'format-go' subcommand
	formatFlags := flag.NewFlagSet("format-go", flag.ExitOnError)
	localPrefix := formatFlags.String("l", "", "local import prefix")
	formatFlags.Usage = func() {
		fmt.Printf("Usage: %s format-go [OPTIONS] FILES\n", filepath.Base(os.Args[0]))
		formatFlags.PrintDefaults()
	}
	if err := formatFlags.Parse(args); err != nil {
		return err
	}

	// Ensure that files are provided
	if formatFlags.NArg() == 0 {
		formatFlags.Usage()
		os.Exit(1)
	}

	// Apply the gofumpt format
	gofumptArgs := []string{"-l", "-w"}
	gofumptArgs = append(gofumptArgs, formatFlags.Args()...)

	gofumptCmd := exec.Command("gofumpt", gofumptArgs...)
	gofumptCmd.Stdout = os.Stdout
	gofumptCmd.Stderr = os.Stderr

	fmt.Printf("Running: %v\n", gofumptCmd.Args)
	if err := gofumptCmd.Run(); err != nil {
		return err
	}

	// Apply the gci write
	gciArgs := []string{"write"}
	if *localPrefix != "" {
		gciArgs = append(gciArgs, "-s", fmt.Sprintf("prefix(%s)", *localPrefix))
	}
	gciArgs = append(gciArgs,
		"-s", "standard",
		"-s", "default",
		"-s", "blank",
		"-s", "dot",
		"--skip-generated",
	)
	gciArgs = append(gciArgs, formatFlags.Args()...)

	gciCmd := exec.Command("gci", gciArgs...)
	gciCmd.Stdout = os.Stdout
	gciCmd.Stderr = os.Stderr

	fmt.Printf("Running: %v\n", gciCmd.Args)
	if err := gciCmd.Run(); err != nil {
		return err
	}

	return nil
}
