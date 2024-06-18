package commandrunner

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
	"time"
)

// SSHCommandRunner implements CommandRunner for running bash commands over an SSH connection
type SSHCommandRunner struct {
	config      SSHConfig
	AskPassPath string
}

// SSHConfig holds the configuration for the sshrunner connection
type SSHConfig struct {
	Host               string
	Port               int
	User               string
	PrivateKey         []byte
	PrivateKeyPassword string
}

// NewSSHCommandRunner creates a new SSHCommandRunner
func NewSSHCommandRunner(config SSHConfig) (*SSHCommandRunner, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	return &SSHCommandRunner{
		config: config,
	}, nil
}

// validateConfig checks that all required fields are present in the config
func validateConfig(config SSHConfig) error {
	if config.Host == "" {
		return fmt.Errorf("host is required")
	}
	if config.Port == 0 {
		return fmt.Errorf("port is required")
	}
	if config.User == "" {
		return fmt.Errorf("user is required")
	}
	if len(config.PrivateKey) == 0 {
		return fmt.Errorf("private key is required")
	}
	return nil
}

// RunCommand runs a command over a sshrunner connection and returns the output as a scanner
func (s *SSHCommandRunner) RunCommand(ec ExecCommand) (*bufio.Scanner, error) {
	client, err := connectToSSH(&s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sshrunner: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	applyCommandSettings(&ec, s.AskPassPath, session)

	output, err := session.CombinedOutput(ec.Command)
	if err != nil {
		return nil, err
	}

	return bufio.NewScanner(bytes.NewReader(output)), nil
}

// RunCommandAsync runs a command over a sshrunner connection asynchronously and returns a channel with the output lines
func (s *SSHCommandRunner) RunCommandAsync(ec ExecCommand) (<-chan string, <-chan error, error) {
	output := make(chan string, 10)
	errorsChan := make(chan error)

	client, err := connectToSSH(&s.config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to sshrunner: %w", err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	stdOut, stdErr, err := setupSshPipes(session)
	if err != nil {
		closeResources(session, client)
		return nil, nil, err
	}

	applyCommandSettings(&ec, s.AskPassPath, session)

	var wg sync.WaitGroup
	wg.Add(1)

	go handleSshOutput(stdOut, stdErr, output, errorsChan, &wg)
	go runSshCommand(session, ec.Command, output, errorsChan, client, &wg)

	return output, errorsChan, nil
}

// connectToSSH connects to an ssh server using the provided configuration
func connectToSSH(config *SSHConfig) (*ssh.Client, error) {
	clientConfig, err := setupSSHConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup ssh config: %w", err)
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ssh: %w", err)
	}

	return client, nil
}

// setupSSHConfig sets up the ssh client config by reading the private key and parsing it
func setupSSHConfig(config *SSHConfig) (*ssh.ClientConfig, error) {
	signer, err := getPublicKeySignerFromPrivateKey(config.PrivateKey, config.PrivateKeyPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key signer from private key: %w", err)
	}

	return &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}, nil
}

// applyCommandSettings applies the root elevation and Input settings to the command
func applyCommandSettings(ec *ExecCommand, askPassPath string, session *ssh.Session) {
	if ec.Elevated {
		applyCommandRootElevation(&ec.Command, askPassPath)
	}

	if ec.Input != nil {
		session.Stdin = ec.Input
	}
}

// getPublicKeySignerFromPrivateKey returns a signer from the provided private key
func getPublicKeySignerFromPrivateKey(privateKey []byte, password string) (ssh.Signer, error) {
	var signer ssh.Signer
	var err error
	switch password {
	case "":
		signer, err = ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
	default:
		signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(password))
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key with passphrase: %w", err)
		}
	}
	return signer, nil
}

// setupSshPipes sets up the stdout and stderr pipes for the session
func setupSshPipes(session *ssh.Session) (io.Reader, io.Reader, error) {
	stdOut, err := session.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stdErr, err := session.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}
	return stdOut, stdErr, nil
}

// runSshCommand runs the provided command on the session and sends the output to the output channel
func runSshCommand(session *ssh.Session, command string, output chan<- string, errorsChan chan<- error, client *ssh.Client, wg *sync.WaitGroup) {
	defer closeResources(session, client)
	defer close(output)
	defer close(errorsChan)
	if err := session.Run(command); err != nil {
		errorsChan <- fmt.Errorf("failed to run command: %w", err)
	}
	wg.Wait()
}

// handleSshOutput reads from the stdout and stderr pipes and sends the output to the output channel
func handleSshOutput(stdOut, stdErr io.Reader, output chan<- string, errorsChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(io.MultiReader(stdErr, stdOut))
	for scanner.Scan() {
		output <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		errorsChan <- fmt.Errorf("error reading stderr: %w", err)
	}
}

// closeResources closes the session and client
func closeResources(session *ssh.Session, client *ssh.Client) {
	if session != nil {
		session.Close()
	}
	if client != nil {
		client.Close()
	}
}
