package version

import (
	"os"
	"runtime"
	"strconv"
	"text/template"
	"time"
)

var (
	Version   string
	GitCommit string = "HEAD"
	GitState  string = "dirty"
	BuildDate string = "0"
)

var versionTemplate = `Version:     {{.Version}}
Git commit:  {{.GitCommit}}{{if eq .GitState "dirty"}}
Git State:   {{.GitState}}{{end}}
Built:       {{.BuildDate}}
Go version:  {{.GoVersion}}
OS/Arch:     {{.Os}}/{{.Arch}}
`

type VersionInfo struct {
	Version   string
	GoVersion string
	GitCommit string
	GitState  string
	BuildDate string
	Os        string
	Arch      string
}

func New() *VersionInfo {
	i, err := strconv.ParseInt(BuildDate, 10, 64)
	if err != nil {
		panic(err)
	}

	tu := time.Unix(i, 0)

	return &VersionInfo{
		Version:   Version,
		GoVersion: runtime.Version(),
		GitCommit: GitCommit,
		GitState:  GitState,
		BuildDate: tu.String(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

func (i *VersionInfo) Show() {
	tmpl, err := template.New("version").Parse(versionTemplate)
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(os.Stdout, i); err != nil {
		panic(err)
	}
}
