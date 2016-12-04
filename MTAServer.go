package main

import (
	"bufio"
	"bytes"
	"container/ring"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// MTAServer type represents an MTA Server process
type MTAServer struct {
	Path         string
	Process      *exec.Cmd
	Stdin        io.WriteCloser
	Stdout       io.ReadCloser
	OutputBuffer *ring.Ring
}

// NewMTAServer instantiates a new MTA server instance
func NewMTAServer(path string) *MTAServer {
	server := new(MTAServer)
	server.Path = path
	server.OutputBuffer = ring.New(5000)

	return server
}

func (server *MTAServer) Start() error {
	// Spawn process
	server.Process = exec.Command(server.Path, "-n", "-t", "-u")

	// Get stdin
	var err error
	server.Stdin, err = server.Process.StdinPipe()
	if err != nil {
		return err
	}

	// Get stdout
	server.Stdout, err = server.Process.StdoutPipe()
	if err != nil {
		return err
	}

	// Capture output into ring buffer
	scanner := bufio.NewScanner(server.Stdout)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Println(text)

			server.OutputBuffer.Value = text
			server.OutputBuffer = server.OutputBuffer.Next()
		}
	}()

	server.Process.Stderr = os.Stderr

	return server.Process.Start()
}

func (server *MTAServer) Stop() error {
	if server.Process == nil || server.Process.Process == nil {
		return errors.New("Process not started")
	}

	return server.Process.Process.Signal(os.Interrupt)
}

func (server *MTAServer) Restart() error {
	// Send stop signal
	err := server.Stop()
	if err != nil {
		return err
	}

	// Wait for the server to stop
	err = server.Process.Wait()
	if err != nil {
		return err
	}

	// Start server
	return server.Start()
}

func (server *MTAServer) ExecCommand(command string) error {
	_, err := io.WriteString(server.Stdin, command+"\n")

	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (server *MTAServer) TailBuffer() string {
	// Make string from output buffer
	var buffer bytes.Buffer

	server.OutputBuffer.Do(func(line interface{}) {
		if line != nil {
			buffer.WriteString(line.(string))
			buffer.WriteString("\n")
		}
	})

	return buffer.String()
}
