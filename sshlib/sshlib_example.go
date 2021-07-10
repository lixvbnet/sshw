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
	fd, state, err := client.Shell()
	defer client.RestoreTerminal(fd, state)
	if err != nil {
		log.Fatal(err)
	}

	client.Wait()
}
