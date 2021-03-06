package commands

import (
	"errors"
	"flag"
	"fmt"

	"github.com/paketo-buildpacks/packit/cargo"
)

//go:generate faux --interface BuildpackInspector --output fakes/buildpack_inspector.go
type BuildpackInspector interface {
	Dependencies(path string) ([]cargo.Config, error)
}

//go:generate faux --interface Formatter --output fakes/formatter.go
type Formatter interface {
	Markdown(cargo.Config)
}

type Summarize struct {
	buildpackInspector BuildpackInspector
	formatter          Formatter
}

func NewSummarize(buildpackInspector BuildpackInspector, formatter Formatter) Summarize {
	return Summarize{
		buildpackInspector: buildpackInspector,
		formatter:          formatter,
	}
}

func (s Summarize) Execute(args []string) error {
	var (
		buildpackTarballPath string
		format               string
	)

	fset := flag.NewFlagSet("summarize", flag.ContinueOnError)
	fset.StringVar(&buildpackTarballPath, "buildpack", "", "path to a buildpackage tarball (required)")
	fset.StringVar(&format, "format", "markdown", "format of output options are (markdown)")
	err := fset.Parse(args)
	if err != nil {
		return err
	}

	if buildpackTarballPath == "" {
		return errors.New("missing required flag --buildpack")
	}

	configs, err := s.buildpackInspector.Dependencies(buildpackTarballPath)
	if err != nil {
		return fmt.Errorf("failed to inspect buildpack dependencies: %w", err)
	}

	switch format {
	case "markdown":
		for _, config := range configs {
			s.formatter.Markdown(config)
		}
	default:
		return fmt.Errorf("unknown format %q, please choose from the following formats (\"markdown\")", format)
	}

	return nil
}
