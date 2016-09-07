package cmd

import (
	"flag"
	"html/template"
	"strings"
	"fmt"
	"os"
)

type Command struct {
	// Action, is runnable function
	Action      func(c *Command, arguments []string) int

	// Usage, is example usage command.
	Usage       string

	// UsageText, short description of command.
	UsageText   template.HTML

	// Description, is description of command.
	Description template.HTML

	// Flag, set any flag on this command.
	Flag        flag.FlagSet
}

// Name return string name of command.
func (c *Command) Name() string {
	name := c.Usage
	if i := strings.Index(name, " "); i >= 0 {
		name = name[:i]
	}

	return name
}

// ShowUsage, show usage instruction of this command.
func (c *Command) ShowUsage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.Usage)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(string(c.Description)))
	os.Exit(2)
}

// Runnable, check is command is runnable
func (c *Command) Runnable() bool {
	return c.Action != nil
}

type ListOpts []string

func (opts *ListOpts) String() string {
	return fmt.Sprint(*opts)
}

func (opts *ListOpts) Set(value string) error {
	*opts = append(*opts, value)
	return nil
}