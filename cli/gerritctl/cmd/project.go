package cmd

import (
	"fmt"
	"github.com/shijl0925/go-gerrit"
	"github.com/spf13/cobra"
	"os"
)

// Project Commands
var project = &cobra.Command{
	Use:   "project",
	Short: "project related commands",
}

// Branch Commands
var branch = &cobra.Command{
	Use:   "branch",
	Short: "project branch related commands",
}

// Tag Commands
var tag = &cobra.Command{
	Use:   "tag",
	Short: "project tag related commands",
}

// projectList Command
var projectList = &cobra.Command{
	Use:   "list",
	Short: "List all projects accessible.",
	Run: func(cmd *cobra.Command, args []string) {
		All, _ := cmd.Flags().GetBool("all")
		Limit, _ := cmd.Flags().GetInt("limit")
		Skip, _ := cmd.Flags().GetInt("skip")
		Description, _ := cmd.Flags().GetBool("description")
		Prefix, _ := cmd.Flags().GetString("prefix")
		Regex, _ := cmd.Flags().GetString("regex")
		State, _ := cmd.Flags().GetString("state")
		Tree, _ := cmd.Flags().GetBool("tree")
		Substring, _ := cmd.Flags().GetString("substring")
		Type, _ := cmd.Flags().GetString("type")
		Branch, _ := cmd.Flags().GetString("branch")

		option := gerrit.ProjectOptions{}
		option.All = All
		option.Limit = Limit
		option.Skip = Skip
		option.Description = Description
		option.Prefix = Prefix
		option.Regex = Regex
		option.State = State
		option.Tree = Tree
		option.Substring = Substring
		option.Type = Type
		option.Branch = Branch

		projects, _, err := gerritMod.Instance.Projects.List(gerritMod.Context, &option)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for name, project := range projects {
			fmt.Printf("✅ Project Name: %s.\n", name)
			if Verbose {
				if out, err := ToIndentJSON(project); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%+v\n", out)
				}
			}
		}
	},
}

// projectGet Command
var projectGet = &cobra.Command{
	Use:   "show",
	Short: "Retrieve a project.",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("name")
		project, _, err := gerritMod.Instance.Projects.Get(gerritMod.Context, projectName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific project: %s.\n %v", projectName, err)
			os.Exit(1)
		}

		fmt.Printf("✅ Project Name: %s, Id: %s, State: %s\n", projectName, project.Raw.ID, project.Raw.State)
		if Verbose {
			if out, err := ToIndentJSON(*project.Raw); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%+v\n", out)
			}
		}
	},
}

// projectCreate Command
var projectCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a project.",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		parent, _ := cmd.Flags().GetString("parent")
		description, _ := cmd.Flags().GetString("description")

		option := gerrit.ProjectInput{
			Name:        name,
			Parent:      parent,
			Description: description,
		}
		_, _, err := gerritMod.Instance.Projects.Create(gerritMod.Context, name, &option)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("✅ Create new project,Name: %s.\n", name)
	},
}

// projectDelete Command
var projectDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a project.",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		project, _, err := gerritMod.Instance.Projects.Get(gerritMod.Context, name)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific project: %s.\n %v", name, err)
			os.Exit(1)
		}

		input := gerrit.DeleteOptionsInfo{Force: true, Preserve: true}
		if _, _, err := project.Delete(gerritMod.Context, &input); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("✅ Delete project,Name: %s.\n", name)
	},
}

// branchList Command
var branchList = &cobra.Command{
	Use:   "list",
	Short: "List the branches of a project.",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project_name")
		project, _, err := gerritMod.Instance.Projects.Get(gerritMod.Context, projectName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific project: %s.\n %v", projectName, err)
			os.Exit(1)
		}

		Limit, _ := cmd.Flags().GetInt("limit")
		Skip, _ := cmd.Flags().GetInt("skip")
		Substring, _ := cmd.Flags().GetString("substring")
		Regex, _ := cmd.Flags().GetString("regex")
		option := gerrit.BranchOptions{
			Limit:     Limit,
			Skip:      Skip,
			Substring: Substring,
			Regex:     Regex,
		}

		branches, _, err := project.Branches.List(gerritMod.Context, &option)

		for _, branch := range *branches {
			fmt.Printf("✅ Branch Name: %s.\n", branch.Ref)
			if Verbose {
				if out, err := ToIndentJSON(branch); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%+v\n", out)
				}
			}
		}
	},
}

// branchCreate Command
var branchCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a branch.",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project_name")
		branchName, _ := cmd.Flags().GetString("name")
		revision, _ := cmd.Flags().GetString("revision")
		input := gerrit.BranchInput{
			Ref:      branchName,
			Revision: revision,
		}
		project, _, err := gerritMod.Instance.Projects.Get(gerritMod.Context, projectName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific project: %s.\n %v", projectName, err)
			os.Exit(1)
		}
		if _, _, err := project.Branches.Create(gerritMod.Context, branchName, &input); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("✅ Create new branch,Name: %s.\n", branchName)
	},
}

// branchGet Command
var branchGet = &cobra.Command{
	Use:   "show",
	Short: "Retrieve the branch of a project.",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project_name")
		project, _, err := gerritMod.Instance.Projects.Get(gerritMod.Context, projectName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific project: %s.\n %v", projectName, err)
			os.Exit(1)
		}

		branchName, _ := cmd.Flags().GetString("name")
		branch, _, err := project.Branches.Get(gerritMod.Context, branchName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific branch: %s.\n %v", branchName, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Branch Name: %s.\n", branch.Raw.Ref)
		if Verbose {
			if out, err := ToIndentJSON(*branch.Raw); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%+v\n", out)
			}
		}
	},
}

// branchDelete Command
var branchDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete the branch of a project.",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project_name")
		project, _, err := gerritMod.Instance.Projects.Get(gerritMod.Context, projectName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific project: %s.\n %v", projectName, err)
			os.Exit(1)
		}

		branchName, _ := cmd.Flags().GetString("name")
		branch, _, err := project.Branches.Get(gerritMod.Context, branchName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific branch: %s.\n %v", branchName, err)
			os.Exit(1)
		}
		if _, _, err := branch.Delete(gerritMod.Context); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("✅ Delete branch, name: %s.\n", branch.Raw.Ref)
	},
}

// tagList Command
var tagList = &cobra.Command{
	Use:   "list",
	Short: "List the tags of a project.",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project_name")
		project, _, err := gerritMod.Instance.Projects.Get(gerritMod.Context, projectName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific project: %s.\n %v", projectName, err)
			os.Exit(1)
		}

		Limit, _ := cmd.Flags().GetInt("limit")
		Skip, _ := cmd.Flags().GetInt("skip")
		Substring, _ := cmd.Flags().GetString("substring")
		Regex, _ := cmd.Flags().GetString("regex")
		option := gerrit.TagOptions{
			Limit:     Limit,
			Skip:      Skip,
			Substring: Substring,
			Regex:     Regex,
		}
		tags, _, err := project.Tags.List(gerritMod.Context, &option)

		for _, tag := range *tags {
			fmt.Printf("✅ Tag Name: %s.\n", tag.Ref)
			if Verbose {
				if out, err := ToIndentJSON(tag); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%+v\n", out)
				}
			}
		}
	},
}

// tagCreate Command
var tagCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a tag.",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project_name")
		tagName, _ := cmd.Flags().GetString("name")
		revision, _ := cmd.Flags().GetString("revision")
		input := gerrit.TagInput{
			Revision: revision,
		}
		project, _, err := gerritMod.Instance.Projects.Get(gerritMod.Context, projectName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific project: %s.\n %v", projectName, err)
			os.Exit(1)
		}
		if _, _, err := project.Tags.Create(gerritMod.Context, tagName, &input); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("✅ Create new tag,Name: %s.\n", tagName)
	},
}

// tagGet Command
var tagGet = &cobra.Command{
	Use:   "show",
	Short: "Retrieve the tag of a project.",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project_name")
		project, _, err := gerritMod.Instance.Projects.Get(gerritMod.Context, projectName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific project: %s.\n %v", projectName, err)
			os.Exit(1)
		}

		tagName, _ := cmd.Flags().GetString("name")
		tag, _, err := project.Tags.Get(gerritMod.Context, tagName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific tag: %s.\n %v", tagName, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Tag Name: %s.\n", tag.Raw.Ref)
		if Verbose {
			if out, err := ToIndentJSON(*tag.Raw); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%+v\n", out)
			}
		}
	},
}

// tagDelete Command
var tagDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete the tag of a project.",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project_name")
		project, _, err := gerritMod.Instance.Projects.Get(gerritMod.Context, projectName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific project: %s.\n %v", projectName, err)
			os.Exit(1)
		}

		tagName, _ := cmd.Flags().GetString("name")
		tag, _, err := project.Tags.Get(gerritMod.Context, tagName)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific tag: %s.\n %v", tagName, err)
			os.Exit(1)
		}
		if _, _, err := tag.Delete(gerritMod.Context); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("✅ Delete Tag,Name: %s.\n", tag.Raw.Ref)
	},
}

func init() {
	rootCmd.AddCommand(project)

	project.AddCommand(projectList)
	projectList.Flags().BoolP("all", "a", false, "List all projects")
	projectList.Flags().IntP("limit", "l", 25, "Limit the number of projects to be included in the results")
	projectList.Flags().IntP("skip", "S", 0, "Skip the first N projects in the results")
	projectList.Flags().BoolP("description", "d", false, "Include the project description in the results")
	projectList.Flags().StringP("prefix", "p", "", "Only include projects with the given prefix")
	projectList.Flags().StringP("regex", "r", "", "Only include projects matching the given regular expression")
	projectList.Flags().StringP("state", "s", "", "Only include projects with the given state")
	projectList.Flags().BoolP("tree", "t", false, "Include the project tree in the results")
	projectList.Flags().StringP("substring", "u", "", "Only include projects with the given substring")
	projectList.Flags().StringP("type", "T", "", "Only include projects with the given type")
	projectList.Flags().StringP("branch", "b", "", "Only include projects with the given branch")

	project.AddCommand(projectGet)
	projectGet.Flags().StringP("name", "n", "", "The name of the project (required)")
	projectGet.MarkFlagRequired("name")

	project.AddCommand(projectCreate)
	projectCreate.Flags().StringP("name", "n", "", "The name of the project")
	projectCreate.Flags().StringP("parent", "P", "", "The name of the parent project")
	projectCreate.Flags().StringP("description", "D", "", "The description of the project")

	project.AddCommand(projectDelete)
	projectDelete.Flags().StringP("name", "n", "", "The name of the project (required)")
	projectDelete.MarkFlagRequired("name")

	project.AddCommand(branch)
	branch.AddCommand(branchList)
	branchList.Flags().StringP("project_name", "p", "", "The name of the project (required)")
	branchList.MarkFlagRequired("project_name")
	branchList.Flags().IntP("limit", "l", 0, "Limit the number of branches to be included in the results")
	branchList.Flags().IntP("skip", "S", 0, "Skip the first N branches in the results")
	branchList.Flags().StringP("substring", "u", "", "Only include branches with the given substring")
	branchList.Flags().StringP("regex", "r", "", "Only include branches matching the given regular expression")

	branch.AddCommand(branchCreate)
	branchCreate.Flags().StringP("project_name", "p", "", "The name of the project (required)")
	branchCreate.MarkFlagRequired("project_name")
	branchCreate.Flags().StringP("name", "n", "", "The name of the branch (required)")
	branchCreate.MarkFlagRequired("name")
	branchCreate.Flags().StringP("revision", "r", "", "The revision of the branch")
	branchCreate.MarkFlagRequired("revision")

	branch.AddCommand(branchGet)
	branchGet.Flags().StringP("project_name", "p", "", "The name of the project (required)")
	branchGet.MarkFlagRequired("project_name")
	branchGet.Flags().StringP("name", "n", "", "The name of the branch (required)")
	branchGet.MarkFlagRequired("name")

	branch.AddCommand(branchDelete)
	branchDelete.Flags().StringP("project_name", "p", "", "The name of the project (required)")
	branchDelete.MarkFlagRequired("project_name")
	branchDelete.Flags().StringP("name", "n", "", "The name of the branch (required)")
	branchDelete.MarkFlagRequired("name")

	project.AddCommand(tag)
	tag.AddCommand(tagList)
	tagList.Flags().StringP("project_name", "p", "", "The name of the project (required)")
	tagList.MarkFlagRequired("project_name")
	tagList.Flags().IntP("limit", "l", 0, "Limit the number of tags to be included in the results")
	tagList.Flags().IntP("skip", "S", 0, "Skip the first N tags in the results")
	tagList.Flags().StringP("substring", "u", "", "Only include tags with the given substring")
	tagList.Flags().StringP("regex", "r", "", "Only include tags matching the given regular expression")

	tag.AddCommand(tagCreate)
	tagCreate.Flags().StringP("project_name", "p", "", "The name of the project (required)")
	tagCreate.MarkFlagRequired("project_name")
	tagCreate.Flags().StringP("name", "n", "", "The name of the tag (required)")
	tagCreate.MarkFlagRequired("name")
	tagCreate.Flags().StringP("revision", "r", "", "The revision of the tag")
	tagCreate.MarkFlagRequired("revision")

	tag.AddCommand(tagGet)
	tagGet.Flags().StringP("project_name", "p", "", "The name of the project (required)")
	tagGet.MarkFlagRequired("project_name")
	tagGet.Flags().StringP("name", "n", "", "The name of the tag (required)")
	tagGet.MarkFlagRequired("name")

	tag.AddCommand(tagDelete)
	tagDelete.Flags().StringP("project_name", "p", "", "The name of the project (required)")
	tagDelete.MarkFlagRequired("project_name")
	tagDelete.Flags().StringP("name", "n", "", "The name of the tag (required)")
	tagDelete.MarkFlagRequired("name")
}
