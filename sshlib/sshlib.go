package sshlib

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/term"
)

var (
	Ciphers = []string{
		"aes128-ctr",
		"aes192-ctr",
		"aes256-ctr",
		"aes128-gcm@openssh.com",
		"chacha20-poly1305@openssh.com",
		"arcfour256",
		"arcfour128",
		"arcfour",
		"aes128-cbc",
		"3des-cbc",
		"blowfish-cbc",
		"cast128-cbc",
		"aes192-cbc",
		"aes256-cbc",
	}

	KeyExchanges = []string{
		"diffie-hellman-group-exchange-sha256",
		"diffie-hellman-group18-sha512",
		"diffie-hellman-group16-sha512",
		"diffie-hellman-group14-sha256",
	}

	Outputs = []io.Writer{os.Stdout}
)

type SSHClient struct {
	Host    string
	Port    int
	client  *ssh.Client
	session *ssh.Session
	*BufferedReader
	*WriteCloser
}

func (c *SSHClient) CurrentClient() *ssh.Client {
	return c.client
}

func (c *SSHClient) CurrentSession() *ssh.Session {
	return c.session
}

// Login to remote server with user/password
func (c *SSHClient) Login(user, password string) (err error) {
	authMethods := []ssh.AuthMethod{ssh.Password(password)}
	return c.LoginAuth(user, authMethods)
}

// Login to remote server with multiple auth methods
func (c *SSHClient) LoginAuth(user string, authMethods []ssh.AuthMethod) (err error) {
	clientConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}
	clientConfig.SetDefaults()
	clientConfig.Ciphers = append(clientConfig.Ciphers, Ciphers...)
	clientConfig.KeyExchanges = append(clientConfig.KeyExchanges, KeyExchanges...)

	// get client
	c.client, err = ssh.Dial("tcp", net.JoinHostPort(c.Host, strconv.Itoa(c.Port)), clientConfig)
	if err != nil {
		return
	}

	// open a session
	c.session, err = c.client.NewSession()
	if err != nil {
		log.Fatal("create session failed: ", err)
		return
	}

	// *********** I/O *********** //
	stdoutPipe, err := c.session.StdoutPipe()
	if err != nil {
		return
	}
	c.BufferedReader = NewBufferedReader(stdoutPipe, Outputs...)
	c.session.Stderr = os.Stderr
	stdinPipe, err := c.session.StdinPipe()
	if err != nil {
		return
	}
	c.WriteCloser = &WriteCloser{stdinPipe}
	return nil
}

func (c *SSHClient) RequestTerminal() (fd int, state *term.State, err error) {
	// *********** TERMINAL related ************ //
	fd = int(os.Stdin.Fd())
	state, err = terminal.MakeRaw(fd)
	if err != nil {
		return
	}
	// changed fd to int(os.Stdout.Fd()) becaused terminal.GetSize(fd) doesn't work in Windows (https://github.com/golang/go/issues/20388)
	w, h, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = c.session.RequestPty("xterm", h, w, modes)
	if err != nil {
		return
	}

	// interval get terminal size
	// fix resize issue
	go func() {
		var (
			ow = w
			oh = h
		)
		for {
			// changed fd to int(os.Stdout.Fd()) becaused terminal.GetSize(fd) doesn't work in Windows (https://github.com/golang/go/issues/20388)
			cw, ch, err := terminal.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				break
			}

			if cw != ow || ch != oh {
				err = c.session.WindowChange(ch, cw)
				if err != nil {
					break
				}
				ow = cw
				oh = ch
			}
			time.Sleep(time.Second)
		}
	}()
	return
}

func (c *SSHClient) RestoreTerminal(fd int, state *term.State) error {
	return terminal.Restore(fd, state)
}

// Interactive Shell (Need terminal)
func (c *SSHClient) Shell() error {
	return c.session.Shell()
}

// Start cmd (Need terminal)
func (c *SSHClient) Start(cmd string) error {
	err := c.session.Start(cmd)
	return err
}

// Run cmd and get its output (Does not need terminal)
func (c *SSHClient) Run(cmd string) (output string, err error) {
	err = c.session.Run(cmd)
	output = string(c.ReadUntilEOF())
	return
}

func (c *SSHClient) Wait() error {
	// send keepalive
	go func() {
		for {
			time.Sleep(time.Second * 10)
			c.client.SendRequest("keepalive@openssh.com", false, nil)
		}
	}()

	// keep reading remote output
	go func() {
		buf := make([]byte, 1024)
		for {
			_, err := c.Read(buf)
			if err != nil {
				c.session.Close()
				c.client.Close()
				break
			}
		}
	}()

	// change stdin to user
	go func() {
		// loop io.Copy as it might also return with err == nil in some rare cases
		// (e.g. when using Windows Terminal, press PAUSE key Fn + B)
		for {
			_, err := io.Copy(c, os.Stdin)
			if err != nil {
				c.session.Close()
				c.client.Close()
				break
			}
		}
	}()

	return c.session.Wait()
}
