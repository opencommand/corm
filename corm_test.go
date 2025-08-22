package corm_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/opencommand/corm"
)

type SudoCommand struct {
	opt_i bool
	opt_s bool
	cmd   corm.Command
	corm.BaseCommand
}

func Sudo(cmd corm.Command) *SudoCommand {
	return &SudoCommand{cmd: cmd}
}

func (c *SudoCommand) Name() string {
	return "sudo"
}

func (c *SudoCommand) Shell() *SudoCommand {
	c.opt_s = true
	return c
}

func (c *SudoCommand) Interactive() *SudoCommand {
	c.opt_i = true
	return c
}

func (c *SudoCommand) String() string {
	args := []string{}
	if c.opt_i {
		args = append(args, "-i")
	}
	if c.opt_s {
		args = append(args, "-s")
	}
	return fmt.Sprintf("%s %s %s", c.Name(), strings.Join(args, " "), c.cmd)
}

type GoCommand struct {
	corm.BaseCommand
}

func Go() *GoCommand {
	return &GoCommand{}
}

func (c *GoCommand) Name() string {
	return "go"
}

func (c *GoCommand) String() string {
	return c.Name()
}

func (c *GoCommand) Test(pkg string) *GoTestCommand {
	return &GoTestCommand{pkg: pkg}
}

type GoTestCommand struct {
	pkg          string
	coverprofile string
	covermode    string
	coverpkg     string
	corm.BaseCommand
}

func (c *GoTestCommand) Name() string {
	return "go test"
}

func (c *GoTestCommand) String() string {
	args := []string{}
	if c.coverprofile != "" {
		args = append(args, fmt.Sprintf("-coverprofile=%s", c.coverprofile))
	}
	if c.covermode != "" {
		args = append(args, fmt.Sprintf("-covermode=%s", c.covermode))
	}
	if c.coverpkg != "" {
		args = append(args, fmt.Sprintf("-coverpkg=%s", c.coverpkg))
	}
	return fmt.Sprintf("%s %s %s", c.Name(), strings.Join(args, " "), c.pkg)
}

func (c *GoTestCommand) Coverprofile(coverprofile string) *GoTestCommand {
	c.coverprofile = coverprofile
	return c
}

func (c *GoTestCommand) Covermode(covermode string) *GoTestCommand {
	c.covermode = covermode
	return c
}

func (c *GoTestCommand) Coverpkg(coverpkg string) *GoTestCommand {
	c.coverpkg = coverpkg
	return c
}

func TestGoCommand(t *testing.T) {
	cmd1 := Sudo(Go().Test("./...").Covermode("atomic").Coverprofile("coverage.out")).Shell().Interactive()
	cmd2 := Sudo(Go().Test("./...").Covermode("atomic").Coverprofile("coverage.out"))
	cmd3 := Go().Test("./...")

	result1 := corm.NewPipeline().
		Add(cmd1).
		Pipe(cmd2).
		And(cmd3)
	fmt.Println("构建器模式:", result1.String())

	result2 := corm.Pipe(cmd1, cmd2, cmd3)
	fmt.Println("函数式:", result2.String())
}
