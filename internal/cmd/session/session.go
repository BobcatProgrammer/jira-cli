package session

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"

	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

// NewCmdSession is a session command.
func NewCmdSession() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage session cookie authentication",
		Long:  "Manage session cookie authentication for Jira.",
	}

	cmd.AddCommand(NewCmdSessionSet())

	return cmd
}

// NewCmdSessionSet creates the set subcommand.
func NewCmdSessionSet() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "Set a session cookie from browser export",
		Long:  `Set a session cookie from an exported browser session.

Export your Jira session cookie from your browser and paste it when prompted.
The cookie will be securely stored in your system's keyring.`,
		Run: sessionSet,
	}
}

func sessionSet(*cobra.Command, []string) {
	login := viper.GetString("login")
	if login == "" {
		cmdutil.Failed("No login configured. Run 'jira init' first.")
	}

	// Prompt for session cookie with masked input
	var cookie string
	qs := &survey.Password{
		Message: "Paste your session cookie (cloud.session.token value):",
		Help:    "Export the cloud.session.token cookie from your browser and paste it here. It will be masked as you type.",
	}

	if err := survey.AskOne(qs, &cookie); err != nil {
		cmdutil.Failed("Error reading session cookie: %v", err)
	}

	if cookie == "" {
		cmdutil.Failed("Session cookie cannot be empty.")
	}

	// Confirm the cookie
	confirmMsg := "Do you want to save this session cookie?"
	confirmed := false
	if err := survey.AskOne(&survey.Confirm{
		Message: confirmMsg,
		Default: true,
	}, &confirmed); err != nil {
		cmdutil.Failed("Error: %v", err)
	}

	if !confirmed {
		fmt.Println("Session cookie not saved.")
		os.Exit(0)
	}

	// Store the cookie in the keyring
	if err := keyring.Set("jira-cli", login, cookie); err != nil {
		cmdutil.Failed("Error storing session cookie: %v", err)
	}

	fmt.Println("Session cookie saved successfully!")
	fmt.Println("You can now use jira commands with session authentication.")
}
