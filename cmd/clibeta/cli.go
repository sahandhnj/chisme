package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sahand.dev/chisme/internal/commandrunner"
	"strconv"
	"strings"
)

func main() {
	loadEnv()
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run cli.go <COMMAND>\n")
		os.Exit(1)
	}

	command := strings.Join(os.Args[1:], " ")
	privateKey, err := os.ReadFile(getEnv("SSH_PRIVATE_KEY_PATH", ""))
	if err != nil {
		log.Fatalf("Error reading private key: %s", err)
	}
	sshConfig := commandrunner.SSHConfig{
		Host:               getEnv("SSH_HOST", "host"),
		Port:               getEnvAsInt("SSH_PORT", 22),
		User:               getEnv("SSH_USER", "u"),
		PrivateKey:         privateKey,
		PrivateKeyPassword: getEnv("SSH_PRIVATE_KEY_PASSWORD", "p"),
	}
	commandRunner, err := commandrunner.NewSSHCommandRunner(sshConfig)
	if err != nil {
		log.Fatalf("Error creating SSH command runner: %s", err)
	}

	//if getEnv("SSH_ASKPASS_PATH", "") != "" {
	//	commandRunner.AskPassPath = getEnv("SSH_ASKPASS_PATH", "")
	//}
	ec := commandrunner.ExecCommand{Command: command, Elevated: true}

	ec.Input = strings.NewReader("PASS")

	output, errorsChan, err := commandRunner.RunCommandAsync(ec)
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

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}
	return value
}

func loadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}
