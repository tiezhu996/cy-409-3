package cmd

import (
	"fmt"
	"strings"

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
		tagStr := "无"
		if len(law.Tags) > 0 {
			tagStr = strings.Join(law.Tags, ", ")
		}
		fmt.Printf("名称：%s\n来源：%s\n条款数：%d\n章节数：%d\n标签：%s\n导入时间：%s\n", law.Name, law.SourceFile, law.ArticleCount, law.ChapterCount, tagStr, law.ImportedAt.Format("2006-01-02 15:04:05"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVar(&infoLaw, "law", "", "法规名称或 ID")
	_ = infoCmd.MarkFlagRequired("law")
}
