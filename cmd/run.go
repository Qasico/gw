package cmd

import (
	"io/ioutil"
	"strings"

	"github.com/qasico/gw/print"
	"os"
	"runtime"
	"path"
)

var (
	Run = &Command{
		Usage: "run [-ext=.go]",
		UsageText: "start watching any changes in directory and rebuild it.",
		Description: `
Run command will watching any changes in directory of go project,
it will recompile and restart the application binary.
`,
	}

	mainFiles ListOpts
)

func init() {
	Run.Action = actionRun
	Run.Flag.Var(&mainFiles, "ext", "specify file extension to watch")
}

// actionRun, perform to scan directory and get ready to watching.
func actionRun(cmd *Command, arguments []string) int {
	gps := getGoPath()
	if len(gps) == 0 {
		print.Die("Failed to start:  %s \n", "$GOPATH is not set or empty")
	}

	exit := make(chan bool)
	cwd, _ := os.Getwd()
	appName := path.Base(cwd)

	print.Info("Run on: %s\n", cwd)

	var paths []string
	readDirectory(cwd, &paths)

	var files []string
	for _, arg := range mainFiles {
		if len(arg) > 0 {
			files = append(files, arg)
		}
	}

	Watch(appName, paths, files)
	Build()

	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}

	return 0
}

// readDirectory binds paths with list of existing directory.
func readDirectory(directory string, paths *[]string) {
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}

	useDirectory := false
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), "docs") {
			continue
		}

		if fileInfo.IsDir() == true && fileInfo.Name()[0] != '.' {
			readDirectory(directory + "/" + fileInfo.Name(), paths)
			continue
		}

		if useDirectory == true {
			continue
		}

		if path.Ext(fileInfo.Name()) == ".go" {
			*paths = append(*paths, directory)
			useDirectory = true
		}
	}

	return
}

// getGoPath returns list of go path on system.
func getGoPath() (p []string) {
	gopath := os.Getenv("GOPATH")
	p = strings.Split(gopath, ":")

	return
}