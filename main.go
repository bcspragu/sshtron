package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/ssh"
)

func handler(conn net.Conn, gm *GameManager, config *ssh.ServerConfig) {
	// Before use, a handshake must be performed on the incoming
	// net.Conn.
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		fmt.Println("Failed to handshake with new client")
		return
	}
	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple
		// terminal interface.
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			fmt.Println("could not accept channel.")
			return
		}

		// TODO: Remove this -- only temporary while we launch on HN
		//
		// To see how many concurrent users are online
		fmt.Printf("Player joined. Current stats: %d users, %d games\n",
			gm.SessionCount(), gm.GameCount())

		// Reject all out of band requests accept for the unix defaults, pty-req and
		// shell.
		go func(in <-chan *ssh.Request) {
			for req := range in {
				switch req.Type {
				case "pty-req":
					req.Reply(true, nil)
					continue
				case "shell":
					req.Reply(true, nil)
					continue
				}
				req.Reply(false, nil)
			}
		}(requests)

		gm.HandleNewChannel(channel, sshConn.User())
	}
}

func main() {
	sshPort := flag.String("ssh_port", ":2022", "The ssh port to run the server on")
	passwdFile := flag.String("password_file", "", "If provided, require a password to play on our server")
	flag.Parse()

	var (
		pwCallback   func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error)
		noClientAuth bool
	)

	if *passwdFile == "" {
		noClientAuth = true
	} else {
		pwCallback = callbackFromFile(*passwdFile)
	}

	// Everyone can login!
	config := &ssh.ServerConfig{
		PasswordCallback: pwCallback,
		NoClientAuth:     noClientAuth,
	}

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		panic("Failed to load private key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Failed to parse private key")
	}

	config.AddHostKey(private)

	// Create the GameManager
	gm := NewGameManager()

	fmt.Printf("Listening on SSH port %s...\n", *sshPort)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", *sshPort)
	if err != nil {
		panic("failed to listen for connection")
	}

	for {
		nConn, err := listener.Accept()
		if err != nil {
			panic("failed to accept incoming connection")
		}

		go handler(nConn, gm, config)
	}
}

func callbackFromFile(file string) func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			return nil, errors.New("failed to read password file")
		}
	}

	return func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		if !bytes.Equal(dat, password) {
			return nil, errors.New("invalid password")
		}
		return nil, nil
	}
}
