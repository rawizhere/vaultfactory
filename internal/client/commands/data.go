package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tempizhere/vaultfactory/internal/client/service"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

// NewDataCommands создает команды для управления данными.
func NewDataCommands() *cobra.Command {
	dataCmd := &cobra.Command{
		Use:   "data",
		Short: "Data management commands",
	}

	addCmd := &cobra.Command{
		Use:   "add [type] [name] [data]",
		Short: "Add new data item",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			dataType := models.DataType(args[0])
			name := args[1]
			data := args[2]

			client := service.NewClientService()
			item, err := client.AddData(cmd.Context(), dataType, name, "", data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to add data: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Data added successfully: %s\n", item.ID)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all data items",
		Run: func(cmd *cobra.Command, args []string) {
			client := service.NewClientService()
			items, err := client.ListData(cmd.Context())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to list data: %v\n", err)
				os.Exit(1)
			}

			if len(items) == 0 {
				fmt.Println("No data items found")
				return
			}

			for _, item := range items {
				fmt.Printf("ID: %s, Type: %s, Name: %s\n", item.ID, item.Type, item.Name)
			}
		},
	}

	getCmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get data item by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			client := service.NewClientService()
			item, err := client.GetData(cmd.Context(), id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get data: %v\n", err)
				os.Exit(1)
			}

			data, _ := json.MarshalIndent(item, "", "  ")
			fmt.Println(string(data))
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete data item by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			client := service.NewClientService()
			err := client.DeleteData(cmd.Context(), id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to delete data: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Data deleted successfully: %s\n", id)
		},
	}

	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize data with server",
		Run: func(cmd *cobra.Command, args []string) {
			client := service.NewClientService()
			err := client.Sync(cmd.Context())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to sync data: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Data synchronized successfully")
		},
	}

	dataCmd.AddCommand(addCmd)
	dataCmd.AddCommand(listCmd)
	dataCmd.AddCommand(getCmd)
	dataCmd.AddCommand(deleteCmd)
	dataCmd.AddCommand(syncCmd)

	return dataCmd
}
