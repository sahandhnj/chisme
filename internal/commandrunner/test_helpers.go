package commandrunner

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"strings"
	"testing"
)

// readFromChannels reads from the output and outputError channels and returns the output as string and error
func readFromChannels(t *testing.T, output <-chan string, outputError <-chan error) (string, error) {
	t.Helper()

	var got strings.Builder

	done := make(chan struct{})
	go func() {
		for line := range output {
			got.WriteString(line)
		}
		close(done)
	}()

	var err error
	select {
	case err = <-outputError:
	case <-done:
	}

	return got.String(), err
}

// startTestSSHServer starts a simple SSH server for testing purposes on 127.0.0.1:9090
func startTestSSHServer(t *testing.T) func() {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	private, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		t.Fatalf("failed to create signer from private key: %v", err)
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}
	config.AddHostKey(private)

	// setting up ssh server on port 9090
	listener, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		t.Fatalf("failed to listen on a port: %v", err)
	}

	go func() {
		for {
			nConn, err := listener.Accept()
			if err != nil {
				//log.Printf("failed to accept incoming connection: %v", err)
				continue
			}

			go func(nConn net.Conn) {
				_, chans, reqs, err := ssh.NewServerConn(nConn, config)
				if err != nil {
					log.Printf("failed to handshake: %v", err)
					return
				}

				// Discard all global out-of-band Requests
				go ssh.DiscardRequests(reqs)

				for newChannel := range chans {
					if newChannel.ChannelType() != "session" {
						newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
						continue
					}

					channel, requests, err := newChannel.Accept()
					if err != nil {
						log.Printf("could not accept channel: %v", err)
						continue
					}

					go func() {
						for req := range requests {
							if req.Type == "exec" {
								cmd := string(req.Payload[4:])
								if strings.Contains(cmd, "echo") {
									req.Reply(true, nil)
									io.WriteString(channel, strings.TrimLeft(cmd, "echo ")+"\n")
									channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
								} else {
									req.Reply(true, nil)
									io.WriteString(channel, "unknown command\n")
									channel.SendRequest("exit-status", false, []byte{0, 0, 0, 1})
								}
								channel.Close()
							}
						}
					}()
				}
			}(nConn)
		}
	}()

	return func() {
		listener.Close()
	}
}

// generateClientPrivateKey generates a private key for the client
func generateClientPrivateKey(t *testing.T) []byte {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate client private key: %v", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return privateKeyPEM
}
