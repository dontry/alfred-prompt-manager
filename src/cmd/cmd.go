package cmd

import (
	"github.com/dontry/alfred-prompt-manager/src/service"
	"github.com/spf13/cobra"
)

var svc *service.Service

const (
	CUSTOM_PROMPTS_FILE_NAME      = "custom_prompts.json"
	AWESOME_PROMPTS_DOWNLOAD_LINK = "https://raw.githubusercontent.com/f/awesome-chatgpt-prompts/main/prompts.csv"
	AWESOME_PROMPTS_FILE_NAME     = "awesome_prompts.json"
)

var rootCmd = &cobra.Command{
	Use:   "prompt",
	Short: "prompt is a CLI tool for managing prompts for chatbots",
	Long:  "prompt is a CLI tool for managing prompts for chatbots",
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a prompt",
	Long:  "add a prompt to custom_prompts.json",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		content := args[1]
		cb := svc.Add(name, content)
		svc.Run(cb)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a prompt",
	Long:  "delete a prompt from custom_prompts.json",
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		cb := svc.Delete(input)
		svc.Run(cb)
	},
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query a prompt",
	Long:  "query a prompt from custom_prompts.json & awesome_prompts.json",
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		all, _ := cmd.Flags().GetBool("all")
		action, _ := cmd.Flags().GetString("action")
		cb := svc.Query(input, all, action)
		svc.Run(cb)
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download awesome_prompts.json",
	Long:  "download awesome_prompts.json from github",
	Run: func(cmd *cobra.Command, args []string) {
		cb := svc.Download()
		svc.Run(cb)
	},
}

func Init() {
	svc = service.NewService(
		CUSTOM_PROMPTS_FILE_NAME,
		AWESOME_PROMPTS_DOWNLOAD_LINK,
		AWESOME_PROMPTS_FILE_NAME,
	)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(deleteCmd)
	queryCmd.Flags().BoolP("all", "a", false, "show all prompts")
	queryCmd.Flags().StringP("action", "c", "", "query action")
	rootCmd.AddCommand(queryCmd)

}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
