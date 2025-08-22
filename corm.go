package corm

import (
	"os/exec"
	"strings"
)

type Command interface {
	Name() string
	String() string
	Run() error
	Output() ([]byte, error)
}

type Operator string

const (
	OpPipe      Operator = "|"
	OpAnd       Operator = "&"
	OpSemicolon Operator = ";"
	OpAndAnd    Operator = "&&"
	OpOrOr      Operator = "||"
)

type Pipeline struct {
	commands  []Command
	operators []Operator
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		commands:  make([]Command, 0),
		operators: make([]Operator, 0),
	}
}

func (p *Pipeline) Add(cmd Command) *Pipeline {
	p.commands = append(p.commands, cmd)
	return p
}

func (p *Pipeline) Pipe(cmd Command) *Pipeline {
	if len(p.commands) > 0 {
		p.operators = append(p.operators, OpPipe)
	}
	p.commands = append(p.commands, cmd)
	return p
}

func (p *Pipeline) And(cmd Command) *Pipeline {
	if len(p.commands) > 0 {
		p.operators = append(p.operators, OpAnd)
	}
	p.commands = append(p.commands, cmd)
	return p
}

func (p *Pipeline) AndAnd(cmd Command) *Pipeline {
	if len(p.commands) > 0 {
		p.operators = append(p.operators, OpAndAnd)
	}
	p.commands = append(p.commands, cmd)
	return p
}

func (p *Pipeline) OrOr(cmd Command) *Pipeline {
	if len(p.commands) > 0 {
		p.operators = append(p.operators, OpOrOr)
	}
	p.commands = append(p.commands, cmd)
	return p
}

func (p *Pipeline) Semicolon(cmd Command) *Pipeline {
	if len(p.commands) > 0 {
		p.operators = append(p.operators, OpSemicolon)
	}
	p.commands = append(p.commands, cmd)
	return p
}

func (p *Pipeline) String() string {
	if len(p.commands) == 0 {
		return ""
	}
	if len(p.commands) == 1 {
		return p.commands[0].String()
	}

	var parts []string
	parts = append(parts, p.commands[0].String())

	for i, op := range p.operators {
		parts = append(parts, string(op), p.commands[i+1].String())
	}

	return strings.Join(parts, " ")
}

func (p *Pipeline) Run() error {
	return exec.Command(p.String()).Run()
}

func Pipe(cmds ...Command) *Pipeline {
	p := NewPipeline()
	for _, cmd := range cmds {
		p.Pipe(cmd)
	}
	return p
}

func And(cmds ...Command) *Pipeline {
	p := NewPipeline()
	for _, cmd := range cmds {
		p.And(cmd)
	}
	return p
}

func AndAnd(cmds ...Command) *Pipeline {
	p := NewPipeline()
	for _, cmd := range cmds {
		p.AndAnd(cmd)
	}
	return p
}

func OrOr(cmds ...Command) *Pipeline {
	p := NewPipeline()
	for _, cmd := range cmds {
		p.OrOr(cmd)
	}
	return p
}

func Semicolon(cmds ...Command) *Pipeline {
	p := NewPipeline()
	for _, cmd := range cmds {
		p.Semicolon(cmd)
	}
	return p
}

type BaseCommand struct{}

func (c *BaseCommand) Name() string {
	panic("method not implemented")
}

func (c *BaseCommand) Run() error {
	return exec.Command(c.String()).Run()
}

func (c *BaseCommand) Output() ([]byte, error) {
	return exec.Command(c.String()).CombinedOutput()
}

func (c *BaseCommand) String() string {
	panic("method not implemented")
}
