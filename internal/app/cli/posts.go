package cli

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"redditclone/internal/domain"
	"redditclone/internal/pkg/apperror"
)

// statusCmd represents the status command
var postsCmd = &cobra.Command{
	Use:   "posts",
	Short: "Returns all posts",
	Long:  `Return all posts of all users`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("posts called")
		ctx := context.Background()
		items, err := app.Domain.Post.Service.Query(ctx, domain.DBQueryConditions{})
		if err != nil {
			if err == apperror.ErrNotFound {
				app.Logger.With(ctx).Info(err)
			}
			app.Logger.With(ctx).Error(err)
		}
		fmt.Println(items)
	},
}

func init() {
	app.rootCmd.AddCommand(postsCmd)
}
