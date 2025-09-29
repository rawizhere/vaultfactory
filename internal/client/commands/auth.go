package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tempizhere/vaultfactory/internal/client/service"
)

func NewAuthCommands() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
	}

	registerCmd := &cobra.Command{
		Use:   "register [email] [password]",
		Short: "Register a new user",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			email := args[0]
			password := args[1]

			client := service.NewClientService()
			user, err := client.Register(cmd.Context(), email, password)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Registration failed: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("User registered successfully: %s\n", user.Email)
		},
	}

	loginCmd := &cobra.Command{
		Use:   "login [email] [password]",
		Short: "Login user",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			email := args[0]
			password := args[1]

			client := service.NewClientService()
			user, accessToken, refreshToken, err := client.Login(cmd.Context(), email, password)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Login failed: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Login successful: %s\n", user.Email)
			fmt.Printf("Access token: %s\n", accessToken)
			fmt.Printf("Refresh token: %s\n", refreshToken)
		},
	}

	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout user",
		Run: func(cmd *cobra.Command, args []string) {
			client := service.NewClientService()
			err := client.Logout()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Logout failed: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Logout successful")
		},
	}

	authCmd.AddCommand(registerCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)

	return authCmd
}
