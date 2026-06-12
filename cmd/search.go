package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	exporter "lawsearch/internal/export"
	"lawsearch/internal/search"
	"lawsearch/internal/store"
)

var searchQuery string
var searchLaw string
var searchFuzzy bool
var searchMode string
var searchOutput string
var searchOutputFile string

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "全文检索法规条文",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		law, err := s.LawByName(searchLaw)
		if err != nil {
			return err
		}
		articles, err := s.Articles(law.ID)
		if err != nil {
			return err
		}
		results := search.Engine{Articles: articles, LawName: law.Name}.Search(searchQuery, searchFuzzy, searchMode)
		if searchOutputFile != "" {
			format := searchOutput
			if format == "" {
				format = "markdown"
			}
			return exporter.Results(results, format, searchOutputFile)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.Header("法规", "条款", "内容")
		for _, result := range results {
			_ = table.Append(result.LawName, fmt.Sprintf("第%d条", result.Article.Number), result.Snippet)
		}
		return table.Render()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVar(&searchQuery, "query", "", "检索关键词")
	searchCmd.Flags().StringVar(&searchLaw, "law", "", "法规名称或 ID")
	searchCmd.Flags().BoolVar(&searchFuzzy, "fuzzy", false, "启用模糊/拼音检索")
	searchCmd.Flags().StringVar(&searchMode, "mode", "and", "多关键词组合：and/or")
	searchCmd.Flags().StringVar(&searchOutput, "output", "table", "输出格式：table/json/markdown")
	searchCmd.Flags().StringVar(&searchOutputFile, "file", "", "导出文件路径")
	_ = searchCmd.MarkFlagRequired("query")
	_ = searchCmd.MarkFlagRequired("law")
}
