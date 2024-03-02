package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/service"
	"github.com/vagnerclementino/bragdoc/internal/usercase"
	"time"

	"github.com/spf13/cobra"
)

var userCase usercase.BragUserCase

// bragCmd represents the brag command
var bragCmd = &cobra.Command{
	Use:   "brag",
	Short: "The brag related commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("brag called")
	},
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "AddBrag a new brag",
	RunE: func(cmd *cobra.Command, args []string) error {
		details := "Some details about the brag"
		createdAt := time.Now()

		brag := &domain.Brag{
			ID:          uuid.NewString(),
			Description: "Exciting Achievement",
			Details:     &details,
			CreatedAt:   createdAt,
			UpdatedAt:   &createdAt,
		}

		if err := userCase.AddBrag(brag); err != nil {
			return err
		}
		fmt.Printf("brag created with success. Details:\n%s", *brag)
		return nil
	},
}

func init() {
	userCase = service.NewBragService()
	bragCmd.AddCommand(addCmd)
}
