package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml"
)

var (
	rootDir      string
	depName      string
	printGrouped bool
	printMissing bool
)

// go run main.go -dir "/Users/slomek/go/src/github.com/shipwallet" -dep "github.com/GoogleCloudPlatform/cloudsql-proxy"
// "/Users/slomek/go/src/github.com/shipwallet"
// "github.com/shipwallet/proto"

type version struct {
	RepoName string
	FilePath string
	Version  string
}

func main() {
	flag.StringVar(&rootDir, "dir", ".", "root dir to look for dependencies")
	flag.StringVar(&depName, "dep", "github.com/slomek/dep-check-ver", "name of the dependency")
	flag.BoolVar(&printGrouped, "group", false, "print repos grouped by dependency version")
	flag.BoolVar(&printMissing, "missing", false, "print repos without the given dependency")
	flag.Parse()

	var versions []version
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "vendor") {
			return filepath.SkipDir
		}

		if !strings.HasSuffix(path, "Gopkg.lock") {
			return nil
		}

		dVersion := dependencyVersion(path, depName)
		repoName := filepath.Base(filepath.Dir(path))
		v := version{
			RepoName: repoName,
			FilePath: path,
			Version:  dVersion,
		}

		if !printMissing && v.Version == "" {
			return nil
		}

		if printGrouped {
			versions = append(versions, v)
			return nil
		}

		printOne(v)
		return nil
	})

	if printGrouped {
		gg := groupByVersion(versions)
		printGroups(gg)
	}
}

func dependencyVersion(filePath, depName string) string {
	gopkg, _ := toml.LoadFile(filePath)

	projects := gopkg.Get("projects").([]*toml.Tree)
	for _, p := range projects {
		if p.Get("name").(string) == depName {
			if p.Has("version") {
				return fmt.Sprintf("version: %v", p.Get("version"))
			} else if p.Has("branch") {
				return fmt.Sprintf("branch: %v", p.Get("branch"))
			} else if p.Has("revision") {
				return fmt.Sprintf("revision: %v", p.Get("revision"))
			}
		}
	}
	return ""
}

func printOne(v version) {
	fmt.Printf("%20s %s\n", v.RepoName, v.Version)
}

func groupByVersion(vv []version) map[string][]version {
	m := make(map[string][]version)
	for _, v := range vv {
		vrr, ok := m[v.Version]
		if !ok {
			vrr = []version{}
		}
		m[v.Version] = append(vrr, v)
	}

	return m
}

func printGroups(gg map[string][]version) {
	kk := groupsKeys(gg)
	sort.StringSlice(kk).Sort()

	for _, k := range kk {
		printGroup(gg[k])
	}
}

func groupsKeys(gg map[string][]version) []string {
	kk := make([]string, 0, len(gg))
	for k := range gg {
		kk = append(kk, k)
	}
	return kk
}

func printGroup(vv []version) {
	if len(vv) == 0 {
		return
	}

	fmt.Printf("%s (%d repos)\n", vv[0].Version, len(vv))
	for _, v := range vv {
		fmt.Printf(" - %s\n", v.RepoName)
	}
	fmt.Println()
}
