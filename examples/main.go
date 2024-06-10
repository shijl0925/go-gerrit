package main

import (
	"context"
	"fmt"
	"github.com/shijl0925/go-gerrit"
)

func main() {
	client, _ := gerrit.NewClient("https://review.lineageos.org", nil)

	ctx := context.Background()

	option := gerrit.ProjectOptions{}
	option.Limit = 25
	option.Skip = 0

	projects, resp, err := client.Projects.List(ctx, &option)
	if err != nil {
		panic(err)
	}
	for name, project := range projects {
		fmt.Printf("Project Name: %s, ID: %s, State: %s\n", name, project.ID, project.State)
	}
	fmt.Println(resp.StatusCode)

	//projectName := "Head-Developers"
	////projectName := "LineageOS/android"
	//project, _, err := client.Projects.Get(ctx, projectName)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Project: %+v\n", project.Raw)
}
