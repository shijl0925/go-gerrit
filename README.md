# go-gerrit

![Alt](https://repobeats.axiom.co/api/embed/cc7f0e92adb80793ac5636c37392e3fb0a7e4f7d.svg "Repobeats analytics image")

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=shijl0925_go-gerrit&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=shijl0925_go-gerrit)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/9fed6a9c3725480db1aa6187b6926ca1)](https://app.codacy.com/gh/shijl0925/go-gerrit?utm_source=github.com&utm_medium=referral&utm_content=shijl0925/go-gerrit&utm_campaign=Badge_Grade)
[![DeepSource](https://app.deepsource.com/gh/shijl0925/go-gerrit.svg/?label=active+issues&show_trend=true&token=gTZMEaVQMah8hOul0B3mw_RG)](https://app.deepsource.com/gh/shijl0925/go-gerrit/)

go-gerrit is a [Go](https://golang.org/) client library for the [Gerrit Code Review](https://www.gerritcodereview.com/) system.

See https://pkg.go.dev/github.com/shijl0925/go-gerrit for supported endpoints.

## Installation

```shell
go get github.com/shijl0925/go-gerrit
```

## Usage

Use `gerrit.NewClient` to create a new Gerrit client. It needs a Gerrit baseUrl and username / password, and optionally accepts
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

Use `gerrit.NewGitilesClient` to create a new Gitiles client. It needs a Gitiles baseUrl and username / password, and optionally accepts
an existing `*http.Client`.

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
    // baseUrl := "http://127.0.0.1:8080/plugins/gitiles/"
    baseUrl := "https://gerrit.googlesource.com/"
    client, _ := gerrit.NewGitilesClient(baseUrl, nil)

    projectName := "gerrit"
    commitID := "ec36cba6080bac72790c7875c36f5b86fc55372c"
    commit, _, err := client.GetCommit(ctx, projectName, commitID)

    if err != nil {
        log.Panicf("Gitiles.GetCommit returned error: %v", err)
    }
    log.Printf("Commit: %v", commit)
}
```
