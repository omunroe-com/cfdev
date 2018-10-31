package version

import (
	"code.cloudfoundry.org/cfdev/config"
	"code.cloudfoundry.org/cfdev/iso"
	"code.cloudfoundry.org/cfdev/semver"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

type UI interface {
	Say(message string, args ...interface{})
}

//go:generate mockgen -package mocks -destination mocks/isoreader.go code.cloudfoundry.org/cfdev/cmd/start MetaDataReader
type IsoReader interface {
	Read(isoPath string) (iso.Metadata, error)
}

type Args struct {
	DepsIsoPath string
}

type Version struct {
	UI        UI
	Version   *semver.Version
	Config    config.Config
	IsoReader IsoReader
}

func (v *Version) Execute(args Args) {
	message := []string{fmt.Sprintf("CLI: %s", v.Version.Original)}

	if !exists(args.DepsIsoPath) {
		v.UI.Say(strings.Join(message, "\n"))
		return
	}

	metadata, err := v.IsoReader.Read(args.DepsIsoPath)
	if err != nil {
		v.UI.Say(strings.Join(message, "\n"))
		return
	}

	for _, version := range metadata.Versions {
		message = append(message, fmt.Sprintf("%s: %s", version.Name, version.Value))
	}

	v.UI.Say(strings.Join(message, "\n"))
}

func (v *Version) Cmd() *cobra.Command {
	args := Args{}
	cmd := &cobra.Command{
		Use: "version",
		Run: func(_ *cobra.Command, _ []string) {
			v.Execute(args)
		},
	}

	pf := cmd.PersistentFlags()
	pf.StringVarP(
		&args.DepsIsoPath,
		"file",
		"f",
		filepath.Join(v.Config.CacheDir, "cf-deps.iso"),
		"path to .dev file containing bosh & cf bits",
	)
	return cmd
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		return false
	}

	return true
}
