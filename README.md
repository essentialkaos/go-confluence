<p align="center"><a href="#readme"><img src="https://gh.kaos.st/go-confluence.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/g/go-confluence.v6"><img src="https://gh.kaos.st/godoc.svg" alt="PkgGoDev" /></a>
  <a href="https://kaos.sh/r/go-confluence"><img src="https://kaos.sh/r/go-confluence.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/w/go-confluence/ci"><img src="https://kaos.sh/w/go-confluence/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/go-confluence/codeql"><img src="https://kaos.sh/w/go-confluence/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="https://kaos.sh/b/go-confluence"><img src="https://kaos.sh/b/c367cff1-4b71-43de-9a47-9fb34e8c34df.svg" alt="Codebeat badge" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage-example">Usage example</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

`go-confluence` is a Go package for working with [Confluence REST API](https://docs.atlassian.com/ConfluenceServer/rest/7.3.3/).

**▲ Please take note that this package support only getting data from Confluence API (_i.e. you cannot create or modify data using this package_).**

### Installation

Make sure you have a working Go 1.18+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```bash
go get -u github.com/essentialkaos/go-confluence/v5
```

### Usage example

Authentication with username and password.

```go
package main

import (
  "fmt"
  cf "github.com/essentialkaos/go-confluence/v6"
)

func main() {
	api, err := cf.NewAPI(&Configuration{URL: "https://confl.domain.com", Auth: cf.AuthBasic{"john", "MySuppaPAssWOrd"}})

  api.SetUserAgent("MyApp", "1.2.3")

  if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
  }

  content, err := api.GetContentByID(
    "18173522", cf.ContentIDParameters{
      Version: 4,
      Expand:  []string{"space", "body.view", "version"},
    },
  )

  if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
  }

  fmt.Printf("ID: %s\n", content.ID)
}
```

Authentication with personal token. Please make sure your confluence 7.9 version and later. See [Using Personal Access Tokens guide](https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html)

```go
package main

import (
  "fmt"

  cf "github.com/essentialkaos/go-confluence/v6"
)

func main() {
	api, err := cf.NewAPI(&Configuration{URL: "https://confl.domain.com", Auth: cf.AuthToken{"avaMTxxxqKaxpFHpmwHPXhjmUFfAJMaU3VXUji73EFhf"}})

  api.SetUserAgent("MyApp", "1.2.3")

  if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
  }

  content, err := api.GetContentByID(
    "18173522", cf.ContentIDParameters{
      Version: 4,
      Expand:  []string{"space", "body.view", "version"},
    },
  )
  if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
  }

  fmt.Printf("ID: %s\n", content.ID)
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
