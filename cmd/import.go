package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"lawsearch/internal/indexer"
	"lawsearch/internal/parser"
	"lawsearch/internal/store"
)

var importFile string
var importName string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "导入法规文本",
	RunE: func(cmd *cobra.Command, args []string) error {
		bundle, err := parser.ParseFile(importFile, importName)
		if err != nil {
			return err
		}
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		tokenIndex := indexer.Build(bundle.Articles)
		if err := s.SaveBundle(bundle, tokenIndex); err != nil {
			return err
		}
		fmt.Printf("导入完成：%s，条款 %d，章节 %d\n", bundle.Law.Name, bundle.Law.ArticleCount, bundle.Law.ChapterCount)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVar(&importFile, "file", "", "TXT/Markdown 文件路径")
	importCmd.Flags().StringVar(&importName, "name", "", "法规名称")
	_ = importCmd.MarkFlagRequired("file")
	_ = importCmd.MarkFlagRequired("name")
}
