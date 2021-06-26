<p align="center"><a href="#readme"><img src="https://gh.kaos.st/go-confluence.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/g/go-confluence.v5"><img src="https://gh.kaos.st/godoc.svg" alt="PkgGoDev" /></a>
  <a href="https://kaos.sh/r/go-confluence"><img src="https://kaos.sh/r/go-confluence.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/w/go-confluence/ci"><img src="https://kaos.sh/w/go-confluence/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/go-confluence/codeql"><img src="https://kaos.sh/w/go-confluence/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="https://kaos.sh/b/go-confluence"><img src="https://kaos.sh/b/c367cff1-4b71-43de-9a47-9fb34e8c34df.svg" alt="Codebeat badge" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage-example">Usage example</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

`go-confluence` is a Go package for working with [Confluence REST API](https://docs.atlassian.com/ConfluenceServer/rest/7.3.3/).

Currently, this package support only getting data from API (_i.e., you cannot create or modify data using this package_).

### Installation

Make sure you have a working Go 1.15+ workspace (_[instructions](https://golang.org/doc/install)_), then:

````
go get -d pkg.re/essentialkaos/go-confluence.v5
````

For update to latest stable release, do:

```
go get -d -u pkg.re/essentialkaos/go-confluence.v5
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
| `master` (_Stable_) | [![CI](https://kaos.sh/w/go-confluence/ci.svg?branch=master)](https://kaos.sh/w/go-confluence/ci?query=branch:master) |
| `develop` (_Unstable_) | [![CI](https://kaos.sh/w/go-confluence/ci.svg?branch=develop)](https://kaos.sh/w/go-confluence/ci?query=branch:develop) |

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
