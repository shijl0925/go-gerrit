package cmd

import (
	"fmt"
	"github.com/shijl0925/go-gerrit"
	"github.com/spf13/cobra"
	"os"
)

// Change Commands
var change = &cobra.Command{
	Use:   "change",
	Short: "change related commands",
}

var changeQuery = &cobra.Command{
	Use:   "query",
	Short: "Query changes.",
	Run: func(cmd *cobra.Command, args []string) {
		Limit, _ := cmd.Flags().GetInt("limit")
		Start, _ := cmd.Flags().GetInt("start")
		Query, _ := cmd.Flags().GetStringSlice("query")
		AdditionalFields, _ := cmd.Flags().GetStringSlice("additional_fields")

		option := gerrit.QueryChangeOptions{}
		option.Start = Start
		option.Limit = Limit

		option.Query = Query
		option.AdditionalFields = AdditionalFields

		changes, _, err := gerritMod.Instance.Changes.Query(gerritMod.Context, &option)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, change := range *changes {
			fmt.Printf("✅ Change ChangeID: %s.\n", change.ID)
			if Verbose {
				if out, err := ToIndentJSON(change); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%+v\n", out)
				}
			}
		}
	},
}

var changeCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a new change.",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project_name")
		branchName, _ := cmd.Flags().GetString("branch_name")
		subject, _ := cmd.Flags().GetString("subject")
		input := gerrit.ChangeInput{
			Project: projectName,
			Branch:  branchName,
			Subject: subject,
		}
		if change, _, err := gerritMod.Instance.Changes.Create(gerritMod.Context, &input); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("✅ Create new change,ChangeID: %s.\n", change.Raw.ID)
			if Verbose {
				if out, err := ToIndentJSON(change); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%+v\n", out)
				}
			}
		}
	},
}

var changeGet = &cobra.Command{
	Use:   "show",
	Short: "Retrieve a change.",
	Run: func(cmd *cobra.Command, args []string) {
		changeID, _ := cmd.Flags().GetString("change_id")
		change, _, err := gerritMod.Instance.Changes.Get(gerritMod.Context, changeID)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific change: %s.\n %v", changeID, err)
		}
		fmt.Printf("✅ Change ChangeID: %s.\n", change.Raw.ID)
		if Verbose {
			if out, err := ToIndentJSON(change); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%+v\n", out)
			}
		}
	},
}

var changeDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a change.",
	Run: func(cmd *cobra.Command, args []string) {
		changeID, _ := cmd.Flags().GetString("change_id")
		if _, _, err := gerritMod.Instance.Changes.Delete(gerritMod.Context, changeID); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("✅ Delete change,ChangeID: %s.\n", changeID)
	},
}

func init() {
	rootCmd.AddCommand(change)

	change.AddCommand(changeQuery)
	changeQuery.Flags().IntP("limit", "l", 25, "limit")
	changeQuery.Flags().IntP("start", "s", 0, "start")
	changeQuery.Flags().StringSliceP("query", "q", []string{"is:open"}, "query")
	changeQuery.Flags().StringSliceP("additional_fields", "a", []string{}, "additional fields")

	change.AddCommand(changeCreate)
	changeCreate.Flags().StringP("project_name", "p", "", "project name")
	changeCreate.Flags().StringP("branch_name", "b", "", "branch name")
	changeCreate.Flags().StringP("subject", "s", "", "subject")
	changeCreate.MarkFlagRequired("project_name")
	changeCreate.MarkFlagRequired("branch_name")
	changeCreate.MarkFlagRequired("subject")

	change.AddCommand(changeGet)
	changeGet.Flags().StringP("change_id", "c", "", "change id")
	changeGet.MarkFlagRequired("change_id")

	change.AddCommand(changeDelete)
	changeDelete.Flags().StringP("change_id", "c", "", "change id")
	changeDelete.MarkFlagRequired("change_id")
}
