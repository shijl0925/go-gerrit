# go-gerrit

go-gerrit is a [Go](https://golang.org/) client library for the [Gerrit Code Review](https://www.gerritcodereview.com/) system.

See https://pkg.go.dev/github.com/shijl0925/go-gerrit for supported endpoints.

## Installation

```shell
go get github.com/shijl0925/go-gerrit
```

## Usage

Use `gerrit.NewClient` to create a new client. It needs a Gerrit baseUrl and username / password, and optionally accepts
an existing `*http.Client`.

After creating the client, create a new request from one of the endpoint-specific packages, i.e.
`gerrit.ProjectOptions` from `github.com/shijl0925/go-gerrit/projects`.

```go
package main
import (
    "context"
    "fmt"
    "github.com/shijl0925/go-gerrit"
    "log"
)

func main() {
    ctx := context.Background()
    baseUrl := "https://review.lineageos.org"
    client, _ := gerrit.NewClient(baseUrl, gerrit.DefaultClient)

    option := gerrit.ProjectOptions{}
    option.Limit = 25
    option.Skip = 0

    projects, _, err := client.Projects.List(ctx, &option)
    if err != nil {
        log.Fatalf("could not search projects: %+v", err)
    }
    for name, project := range projects {
        fmt.Printf("Project Name: %s, ID: %s, State: %s\n", name, project.ID, project.State)
    }
}
```
