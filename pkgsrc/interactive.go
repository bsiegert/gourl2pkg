package pkgsrc

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/AlecAivazis/survey.v1"
)

var categories []string

func Categories(pkgsrcdir string) ([]string, error) {
	if len(categories) > 0 {
		return categories, nil
	}

	f, err := os.Open(pkgsrcdir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dir, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		if !fi.IsDir() {
			continue
		}
		switch fi.Name() {
		// These directories do not contain packages:
		case "CVS":
		case "distfiles", "packages":
		case "doc":
		case "mk":
		case "regress":

		default:
			categories = append(categories, fi.Name())
		}
	}
	sort.Strings(categories)

	return categories, nil
}

func (p *PkgMeta) InteractiveSetup() {
	fmt.Printf("\nPackaging repository %q.\n", p.GoSrcpath)
	survey.AskOne(&survey.Input{
		Message: "Enter a distname, including version: ",
		Default: p.Distname,
	}, &p.Distname, survey.Required)
}
