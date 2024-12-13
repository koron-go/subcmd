/*
Package subcmd provides sub commander.
*/
package subcmd

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Runner defines a base interface for Command and Set.
// Runner interface is defined for use only with DefineSet function.
type Runner interface {
	// Name returns name of runner.
	Name() string

	// Desc returns description of runner.
	Desc() string

	// Run runs runner with context and arguments.
	Run(ctx context.Context, args []string) error
}

// CommandFunc is handler of sub-command, and an entry point.
type CommandFunc func(ctx context.Context, args []string) error

// Command repreesents a sub-command, and implements Runner interface.
type Command struct {
	name  string
	desc  string
	runFn CommandFunc
}

var _ Runner = Command{}

// DefineCommand defines a Command with name, desc, and function.
func DefineCommand(name, desc string, fn CommandFunc) Command {
	return Command{
		name:  name,
		desc:  desc,
		runFn: fn,
	}
}

// Name returns name of the command.
func (c Command) Name() string {
	return c.name
}

// Desc returns description of the command.
func (c Command) Desc() string {
	return c.desc
}

// Run executes sub-command. It will invoke CommandFunc which passed to DefineCommand.
func (c Command) Run(ctx context.Context, args []string) error {
	ctx = withName(ctx, c)
	if c.runFn == nil {
		names := strings.Join(Names(ctx), " ")
		return fmt.Errorf("no function declared for command: %s", names)
	}
	return c.runFn(ctx, args)
}

// Set provides set of Commands or nested Sets.
type Set struct {
	name    string
	desc    string
	Runners []Runner
}

var _ Runner = Set{}

// DefineSet defines a set of Runners with name, and desc.
func DefineSet(name, desc string, runners ...Runner) Set {
	return Set{
		name:    name,
		desc:    desc,
		Runners: runners,
	}
}

// DefineRootSet defines a set of Runners which used as root of Set (maybe
// passed to Run).
func DefineRootSet(runners ...Runner) Set {
	return Set{name: rootName(), Runners: runners}
}

// Name returns name of Set.
func (s Set) Name() string {
	return s.name
}

// Desc returns description of Set.
func (s Set) Desc() string {
	return s.desc
}

// childRunner retrieves a child Runner with name
func (s Set) childRunner(name string) Runner {
	for _, r := range s.Runners {
		if r.Name() == name {
			return r
		}
	}
	return nil
}

type errorSetRun struct {
	src Set
	msg string
}

func (err *errorSetRun) Error() string {
	// align width of name columns
	w := 12
	for _, r := range err.src.Runners {
		if n := len(r.Name()) + 1; n > w {
			w = (n + 3) / 4 * 4
		}
	}
	// format error message
	bb := &bytes.Buffer{}
	fmt.Fprintf(bb, "%s.\n\nAvailable sub-commands are:\n", err.msg)
	for _, r := range err.src.Runners {
		fmt.Fprintf(bb, "\n\t%-*s%s", w, r.Name(), r.Desc())
	}
	return bb.String()
}

func (s Set) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return &errorSetRun{src: s, msg: "no commands selected"}
	}
	name := args[0]
	child := s.childRunner(name)
	if child == nil {
		return &errorSetRun{src: s, msg: "command not found"}
	}
	return child.Run(withName(ctx, s), args[1:])
}

// Run runs a Runner with ctx and args.
func Run(r Runner, args ...string) error {
	return r.Run(context.Background(), args)
}

type keyNames struct{}

// Names retrives names layer of current sub command.
func Names(ctx context.Context) []string {
	if names, ok := ctx.Value(keyNames{}).([]string); ok {
		return names
	}
	return nil
}

func withName(ctx context.Context, r Runner) context.Context {
	return context.WithValue(ctx, keyNames{}, append(Names(ctx), r.Name()))
}

func stripExeExt(in string) string {
	_, out := filepath.Split(in)
	ext := filepath.Ext(out)
	if ext == ".exe" {
		return out[:len(out)-len(ext)]
	}
	return out
}

func rootName() string {
	exe, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("failed to obtain executable name: %s", err))
	}
	return stripExeExt(exe)
}

func FlagSet(ctx context.Context) *flag.FlagSet {
	name := strings.Join(Names(ctx), " ")
	return flag.NewFlagSet(name, flag.ExitOnError)
}
