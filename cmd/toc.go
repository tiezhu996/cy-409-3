package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"lawsearch/internal/store"
)

var tocLaw string
var tocChapter string

var tocCmd = &cobra.Command{
	Use:   "toc",
	Short: "浏览法规章节目录",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		law, err := s.LawByName(tocLaw)
		if err != nil {
			return err
		}
		sections, err := s.Sections(law.ID)
		if err != nil {
			return err
		}
		for _, section := range sections {
			if tocChapter == "" || section.Title == tocChapter {
				fmt.Printf("- [%s] %s\n", section.Level, section.Title)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tocCmd)
	tocCmd.Flags().StringVar(&tocLaw, "law", "", "法规名称或 ID")
	tocCmd.Flags().StringVar(&tocChapter, "chapter", "", "指定章节")
	_ = tocCmd.MarkFlagRequired("law")
}
