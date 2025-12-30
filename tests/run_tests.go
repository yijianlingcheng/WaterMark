package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	verbose      = flag.Bool("v", false, "verbose output")
	cover        = flag.Bool("cover", false, "enable coverage")
	coverProfile = flag.String("coverprofile", "", "coverage profile file")
	race         = flag.Bool("race", false, "enable race detector")
	packages     = flag.String("pkgs", "./...", "packages to test")
)

func main() {
	flag.Parse()

	args := []string{"test"}

	if *verbose {
		args = append(args, "-v")
	}

	if *cover {
		args = append(args, "-cover")
		if *coverProfile != "" {
			args = append(args, "-coverprofile="+*coverProfile)
		}
	}

	if *race {
		args = append(args, "-race")
	}

	args = append(args, *packages)

	log.Println("Running tests with command: go", strings.Join(args, " "))

	cmd := exec.CommandContext(context.Background(), "go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Tests failed: %v\n", err)
		os.Exit(1)
	}

	log.Println("All tests passed successfully!")

	if *cover && *coverProfile != "" {
		log.Printf("Coverage report saved to: %s\n", *coverProfile)
		log.Println("To view coverage report, run: go tool cover -html=" + *coverProfile)
	}
}
