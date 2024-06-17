package main

import (
	"fmt"
	"log"
	"os"
	"sahand.dev/chisme/internal/commandrunner"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run cli.go <COMMAND>\n")
		os.Exit(1)
	}

	command := strings.Join(os.Args[1:], " ")
	cmdRunner := &commandrunner.BashCommandRunner{}

	output, errorsChan, err := cmdRunner.RunCommandAsync(command)
	if err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for line := range output {
			fmt.Println(line)
		}
	}()

	select {
	case err := <-errorsChan:
		if err != nil {
			log.Fatalf("Command execution error: %v", err)
		}
	case <-done:
	}
}
