package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"lawsearch/internal/store"
)

var infoLaw string

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "查看法规元信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		law, err := s.LawByName(infoLaw)
		if err != nil {
			return err
		}
		fmt.Printf("名称：%s\n来源：%s\n条款数：%d\n章节数：%d\n导入时间：%s\n", law.Name, law.SourceFile, law.ArticleCount, law.ChapterCount, law.ImportedAt.Format("2006-01-02 15:04:05"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVar(&infoLaw, "law", "", "法规名称或 ID")
	_ = infoCmd.MarkFlagRequired("law")
}
