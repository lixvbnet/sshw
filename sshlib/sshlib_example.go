package sshlib

import (
	"log"
)

func main() {
	client := &SSHClient{
		Host:    "localhost",
		Port:    22,
	}
	err := client.Login("root", "")
	if err != nil {
		log.Fatal(err)
	}

	fd, state, err := client.RequestTerminal()
	defer client.RestoreTerminal(fd, state)
	if err != nil {
		log.Fatal(err)
	}

	client.Shell()

	// Do some I/O conversation
	const PROMPT = "~] "
	client.ReadUntil(PROMPT)
	client.WriteLine("echo hello")
	client.ReadUntil(PROMPT)

	client.Wait()
}
