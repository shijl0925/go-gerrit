package cmd

import (
	"fmt"
	"github.com/shijl0925/go-gerrit"
	"github.com/spf13/cobra"
	"os"
)

// Account Commands
var account = &cobra.Command{
	Use:   "account",
	Short: "account related commands",
}

var accountList = &cobra.Command{
	Use:   "list",
	Short: "list accounts",
	Run: func(cmd *cobra.Command, args []string) {
		Limit, _ := cmd.Flags().GetInt("limit")
		Start, _ := cmd.Flags().GetInt("start")
		AdditionalFields, _ := cmd.Flags().GetStringSlice("additional_fields")

		option := gerrit.QueryAccountOptions{}
		option.Start = Start
		option.Limit = Limit
		option.Query = []string{"is:active"}
		option.AdditionalFields = AdditionalFields //[]string{"DETAILS"}
		accounts, _, err := gerritMod.Instance.Accounts.Query(gerritMod.Context, &option)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, account := range *accounts {
			fmt.Printf("✅ Account AccountID: %d.\n", account.AccountID)
			if Verbose {
				if out, err := ToIndentJSON(account); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%+v\n", out)
				}
			}
		}
	},
}

var accountGet = &cobra.Command{
	Use:   "show",
	Short: "Retrieve the account.",
	Run: func(cmd *cobra.Command, args []string) {
		accountID, _ := cmd.Flags().GetString("account_id")
		account, _, err := gerritMod.Instance.Accounts.Get(gerritMod.Context, accountID)
		if err != nil {
			fmt.Printf("❌ Unable to find the specific account: %s.\n %v", accountID, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Account AccountID: %d.\n", account.Raw.AccountID)
		if Verbose {
			if out, err := ToIndentJSON(*account.Raw); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%+v\n", out)
			}
		}
	},
}

var accountCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a new account.",
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		username, _ := cmd.Flags().GetString("username")
		input := gerrit.AccountInput{
			Email:    email,
			Username: username,
		}
		if _, _, err := gerritMod.Instance.Accounts.Create(gerritMod.Context, username, &input); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("✅ Create new account,Email: %s.\n", email)
	},
}

func init() {
	rootCmd.AddCommand(account)

	account.AddCommand(accountList)
	accountList.Flags().IntP("limit", "l", 25, "Limit the number of accounts returned.")
	accountList.Flags().IntP("start", "s", 0, "Skip the first N accounts.")
	accountList.Flags().StringSliceP("additional_fields", "f", []string{}, "Additional fields to be returned.")

	account.AddCommand(accountGet)
	accountGet.Flags().StringP("account_id", "a", "", "The account ID.")
	accountGet.MarkFlagRequired("account_id")

	account.AddCommand(accountCreate)
	accountCreate.Flags().StringP("email", "e", "", "The email address of the new account.")
	accountCreate.Flags().StringP("username", "u", "", "The username of the new account.")
	accountCreate.MarkFlagRequired("email")
	accountCreate.MarkFlagRequired("username")
}
