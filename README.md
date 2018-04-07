<p align="center"><a href="#readme"><img src="https://gh.kaos.st/go-confluence.svg"/></a></p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage-example">Usage example</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<p align="center">
  <a href="https://godoc.org/pkg.re/essentialkaos/go-confluence.v1"><img src="https://godoc.org/pkg.re/essentialkaos/go-confluence.v1?status.svg"></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/go-confluence"><img src="https://goreportcard.com/badge/github.com/essentialkaos/go-confluence"></a>
  <a href="https://travis-ci.org/essentialkaos/go-confluence"><img src="https://travis-ci.org/essentialkaos/go-confluence.svg"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-go-confluence-master"><img alt="codebeat badge" src="https://codebeat.co/badges/c367cff1-4b71-43de-9a47-9fb34e8c34df" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg"></a>
</p>

`go-confluence` is a Go package for wroking with [Confluence REST API](https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/).

Currently, this package support only getting data from API (i.e., you cannot create or modify data using this package).

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.7+ workspace ([instructions](https://golang.org/doc/install)), then:

````
go get pkg.re/essentialkaos/go-confluence.v1
````

For update to latest stable release, do:

```
go get -u pkg.re/essentialkaos/go-confluence.v1
```

### Usage example

```go
package main

import (
  "fmt"
  cf "pkg.re/essentialkaos/go-confluence.v1"
)

func main() {
  api, err := cf.NewAPI("https://confluence.domain.com", "john", "MySuppaPAssWOrd")

  if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
  }

  content, err := cf.GetContentByID(
    "18173522", cf.ContentIDParameters{
      Version: 4,
      Expand:  []string{"space", "body.view", "version"},
    },
  )

  if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
  }

  fmt.Println("ID: %s\n", content.ID)
}

```

### Build Status

| Branch     | Status |
|------------|--------|
| `master` (_Stable_) | [![Build Status](https://travis-ci.org/essentialkaos/go-confluence.svg?branch=master)](https://travis-ci.org/essentialkaos/go-confluence) |
| `develop` (_Unstable_) | [![Build Status](https://travis-ci.org/essentialkaos/go-confluence.svg?branch=develop)](https://travis-ci.org/essentialkaos/go-confluence) |

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
