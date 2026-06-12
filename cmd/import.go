package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"lawsearch/internal/indexer"
	"lawsearch/internal/parser"
	"lawsearch/internal/store"
)

var importFile string
var importName string
var importTags []string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "导入法规文本",
	RunE: func(cmd *cobra.Command, args []string) error {
		bundle, err := parser.ParseFile(importFile, importName, importTags)
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
		tagStr := ""
		if len(bundle.Law.Tags) > 0 {
			tagStr = fmt.Sprintf("，标签 [%s]", strings.Join(bundle.Law.Tags, ", "))
		}
		fmt.Printf("导入完成：%s，条款 %d，章节 %d%s\n", bundle.Law.Name, bundle.Law.ArticleCount, bundle.Law.ChapterCount, tagStr)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVar(&importFile, "file", "", "TXT/Markdown 文件路径")
	importCmd.Flags().StringVar(&importName, "name", "", "法规名称")
	importCmd.Flags().StringSliceVar(&importTags, "tags", nil, "业务标签，多个用逗号分隔")
	_ = importCmd.MarkFlagRequired("file")
	_ = importCmd.MarkFlagRequired("name")
}
