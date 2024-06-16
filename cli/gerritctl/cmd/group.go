package cmd

import (
	"fmt"
	"github.com/shijl0925/go-gerrit"
	"github.com/spf13/cobra"
	"os"
)

// Group Commands
var group = &cobra.Command{
	Use:   "group",
	Short: "group related commands",
}

var groupList = &cobra.Command{
	Use:   "list",
	Short: "List the groups.",
	Run: func(cmd *cobra.Command, args []string) {
		Limit, _ := cmd.Flags().GetInt("limit")
		Skip, _ := cmd.Flags().GetInt("skip")

		option := gerrit.ListGroupsOptions{}
		option.Skip = Skip
		option.Limit = Limit

		groups, _, err := gerritMod.Instance.Groups.List(gerritMod.Context, &option)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for name, group := range groups {
			fmt.Printf("✅ Group Name: %s, GroupID: %d.\n", name, group.GroupID)
			if Verbose {
				if out, err := ToIndentJSON(group); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%+v\n", out)
				}
			}
		}
	},
}

var groupCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a new group.",
	Run: func(cmd *cobra.Command, args []string) {
		groupName, _ := cmd.Flags().GetString("name")
		input := gerrit.GroupInput{
			Name: groupName,
		}
		if _, _, err := gerritMod.Instance.Groups.Create(gerritMod.Context, groupName, &input); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("✅ Create new group,Name: %s.\n", groupName)
	},
}

var groupGet = &cobra.Command{
	Use:   "show",
	Short: "Retrieve the group.",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("group_id")
		group, _, err := gerritMod.Instance.Groups.Get(gerritMod.Context, groupID)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific group: %s.\n %v", groupID, err)
		}
		fmt.Printf("✅ Group GroupID: %d.\n", group.Raw.GroupID)
		if Verbose {
			if out, err := ToIndentJSON(*group.Raw); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%+v\n", out)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(group)

	group.AddCommand(groupList)
	groupList.Flags().IntP("limit", "l", 25, "limit")
	groupList.Flags().IntP("skip", "s", 0, "skip")

	group.AddCommand(groupCreate)
	groupCreate.Flags().StringP("name", "n", "", "group name")
	groupCreate.MarkFlagRequired("name")

	group.AddCommand(groupGet)
	groupGet.Flags().StringP("group_id", "g", "", "group id")
	groupGet.MarkFlagRequired("group_id")
}
