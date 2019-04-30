# dep-check-ver

The idea behind this tool is to check what version of a given dependency is used in your repositories. This came to life to help me find if there is a need to update one of the libraries that contained a bug, not allowing us to use newer version of Go.

It is meant to work againts repositories using dep as a dependency manager, as it checks `Gopkg.lock` file in order to see what versions of dependencies are used.

## Installation

```
go get -u github.com/slomek/dep-check-ver
```

## Usage

```
$ dep-check-ver -help   
Usage of dep-check-ver:
  -dep string
        name of the dependency (default "github.com/slomek/dep-check-ver")
  -dir string
        root dir to look for dependencies (default ".")
  -group
        print repos grouped by dependency version
  -missing
        print repos without the given dependency
```

### Simple list

```
$ dep-check-ver -dir path/to/repo-list -dep github.com/spf13/cobra
                repoA version: v0.0.3
                repoB version: v0.0.1
                repoC version: v0.0.3
                repoD version: v0.0.3
                repoE version: v0.0.1
```

### Group by version

```
$ dep-check-ver -dir path/to/repo-list -dep github.com/spf13/cobra
version: v0.0.1 (2 repos)
 - repoB
 - repoE
version: v0.0.3 (3 repos)
 - repoA
 - repoC
 - repoD
```
