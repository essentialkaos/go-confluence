<p align="center"><a href="#readme"><img src="https://gh.kaos.st/go-confluence.svg"/></a></p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/essentialkaos/go-confluence"><img src="https://pkg.go.dev/badge/github.com/essentialkaos/go-confluence" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/go-confluence"><img src="https://goreportcard.com/badge/github.com/essentialkaos/go-confluence"></a>
  <a href="https://travis-ci.com/essentialkaos/go-confluence"><img src="https://travis-ci.com/essentialkaos/go-confluence.svg"></a>
  <a href="https://github.com/essentialkaos/go-confluence/actions?query=workflow%3ACodeQL"><img src="https://github.com/essentialkaos/go-confluence/workflows/CodeQL/badge.svg" /></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-go-confluence-master"><img alt="codebeat badge" src="https://codebeat.co/badges/c367cff1-4b71-43de-9a47-9fb34e8c34df" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage-example">Usage example</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

`go-confluence` is a Go package for working with [Confluence REST API](https://docs.atlassian.com/ConfluenceServer/rest/7.3.3/).

Currently, this package support only getting data from API (_i.e., you cannot create or modify data using this package_).

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.13+ workspace (_[instructions](https://golang.org/doc/install)_), then:

````
go get pkg.re/essentialkaos/go-confluence.v5
````

For update to latest stable release, do:

```
go get -u pkg.re/essentialkaos/go-confluence.v5
```

### Usage example

```go
package main

import (
  "fmt"
  cf "pkg.re/essentialkaos/go-confluence.v5"
)

func main() {
  api, err := cf.NewAPI("https://confluence.domain.com", "john", "MySuppaPAssWOrd")
  api.SetUserAgent("MyApp", "1.2.3")

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
| `master` (_Stable_) | [![Build Status](https://travis-ci.com/essentialkaos/go-confluence.svg?branch=master)](https://travis-ci.com/essentialkaos/go-confluence) |
| `develop` (_Unstable_) | [![Build Status](https://travis-ci.com/essentialkaos/go-confluence.svg?branch=develop)](https://travis-ci.com/essentialkaos/go-confluence) |

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
