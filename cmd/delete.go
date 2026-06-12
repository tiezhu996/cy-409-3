package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"lawsearch/internal/store"
)

var deleteLaw string

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除已导入法规",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		law, err := s.LawByName(deleteLaw)
		if err != nil {
			return err
		}
		if err := s.DeleteLaw(law.ID); err != nil {
			return err
		}
		fmt.Printf("已删除：%s\n", law.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deleteLaw, "law", "", "法规名称或 ID")
	_ = deleteCmd.MarkFlagRequired("law")
}
