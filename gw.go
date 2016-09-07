package main

import (
	"os"
	"log"
	"flag"
	"io"
	"strings"
	"html/template"
	"github.com/qasico/gw/print"
	"github.com/qasico/gw/cmd"
)

var (
	// commands, the list of available command.
	commands = []*cmd.Command{
		cmd.Run,
	}

	// helpTemplate, template of help instructions.
	helpTemplate = `{{if .Runnable}}usage: gw {{.Name}}
{{.UsageText}}{{end}}

{{.Description | trim}}
`

	// usageTemplate, template of usage instructions.
	usageTemplate = `Usage: gw [command] [arguments]

Commands: {{range .}}{{if .Runnable}}
       {{.Name | printf "%-11s"}} {{.UsageText}}{{end}}{{end}}

Use "gw help [command]" for more information about a command.
`
)

func main() {
	flag.Usage = showUsage

	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		showUsage()
	}

	if args[0] == "help" {
		showHelp(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Action != nil {
			cmd.Flag.Usage = func() {
				cmd.ShowUsage()
			}

			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			os.Exit(cmd.Action(cmd, args))

			return
		}
	}

	print.Err("gw: unknown subcommand %q\nRun 'gw help' for usage.\n", args[0])
}

// showUsage, show usage instruction on command line.
func showUsage() {
	compileTemplate(os.Stdout, usageTemplate, commands)
	os.Exit(2)
}

// showHelp, show help template on command line.
func showHelp(args []string) {
	if len(args) == 0 {
		showUsage()
		return
	}

	if len(args) != 1 {
		print.Die("usage: gw help command\n\nToo many arguments given.\n")
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			compileTemplate(os.Stdout, helpTemplate, cmd)
			return
		}
	}

	print.Die("Unknown help topic %#q.  Run 'gw help'.\n", arg)
}

// compileTemplate, showing template and combine template with real data.
func compileTemplate(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": func(s template.HTML) template.HTML {
		return template.HTML(strings.TrimSpace(string(s)))
	}})

	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}