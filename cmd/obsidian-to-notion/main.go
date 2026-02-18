package main

import (
	"fmt"
	"os"

	"github.com/kfi/obsidian-to-notion/internal/notion"
	"github.com/spf13/cobra"
)

func main() {
	var token string

	root := &cobra.Command{
		Use:   "obsidian-to-notion",
		Short: "Migrate Obsidian vaults to Notion",
	}

	// --token is a persistent flag available to all subcommands.
	root.PersistentFlags().StringVar(&token, "token", "", "Notion integration token (or set NOTION_TOKEN)")

	root.AddCommand(newConnectCmd(&token))

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newConnectCmd(token *string) *cobra.Command {
	return &cobra.Command{
		Use:   "connect",
		Short: "Verify connection to the Notion API",
		RunE: func(cmd *cobra.Command, args []string) error {
			t := *token
			if t == "" {
				t = os.Getenv("NOTION_TOKEN")
			}
			if t == "" {
				return fmt.Errorf("notion token required: use --token or set NOTION_TOKEN")
			}

			client := notion.NewClient(t)
			name, err := client.Ping(cmd.Context())
			if err != nil {
				return fmt.Errorf("connection failed: %w", err)
			}

			fmt.Printf("Connected to Notion as: %s\n", name)
			return nil
		},
	}
}
