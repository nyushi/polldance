package polldance

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

//type EventHandler func(*EventData) error
type FilterCommand struct {
	Command string
	Debug   bool
}

func (f *FilterCommand) Handler(ev *EventData) error {
	println(1234)
	commands := strings.Split(f.Command, " ")
	cmd := exec.Command(commands[0], commands[1:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to pipe stdin: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to pipe stdout: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to pipe stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to run process: %w", err)
	}

	stdin.Write([]byte(ev.Data))
	stdin.Close()

	b, _ := ioutil.ReadAll(stdout)
	ev.Data = string(b)

	b, _ = ioutil.ReadAll(stderr)
	if len(b) > 0 {
		return fmt.Errorf("stderr: %s", string(b))
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait process: %w", err)
	}

	return nil
}
