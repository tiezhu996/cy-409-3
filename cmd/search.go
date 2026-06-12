package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	exporter "lawsearch/internal/export"
	"lawsearch/internal/search"
	"lawsearch/internal/store"
	"lawsearch/pkg/models"
)

var searchQuery string
var searchLaw string
var searchTags []string
var searchFuzzy bool
var searchMode string
var searchOutput string
var searchOutputFile string

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "全文检索法规条文",
	RunE: func(cmd *cobra.Command, args []string) error {
		if searchLaw == "" && len(searchTags) == 0 {
			return fmt.Errorf("必须指定 --law 或 --tags 其中之一")
		}
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()

		var laws []models.Law
		if searchLaw != "" {
			law, err := s.LawByName(searchLaw)
			if err != nil {
				return err
			}
			laws = []models.Law{law}
		} else {
			laws, err = s.LawsByTags(searchTags)
			if err != nil {
				return err
			}
			if len(laws) == 0 {
				return fmt.Errorf("未找到匹配标签的法规")
			}
		}

		var allResults []models.SearchResult
		for _, law := range laws {
			articles, err := s.Articles(law.ID)
			if err != nil {
				return err
			}
			results := search.Engine{Articles: articles, LawName: law.Name}.Search(searchQuery, searchFuzzy, searchMode)
			allResults = append(allResults, results...)
		}

		if searchOutputFile != "" {
			format := searchOutput
			if format == "" {
				format = "markdown"
			}
			return exporter.Results(allResults, format, searchOutputFile)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.Header("法规", "条款", "内容")
		for _, result := range allResults {
			_ = table.Append(result.LawName, fmt.Sprintf("第%d条", result.Article.Number), result.Snippet)
		}
		return table.Render()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVar(&searchQuery, "query", "", "检索关键词")
	searchCmd.Flags().StringVar(&searchLaw, "law", "", "法规名称或 ID（与 tags 二选一）")
	searchCmd.Flags().StringSliceVar(&searchTags, "tags", nil, "按标签过滤法规，多个用逗号分隔（与 law 二选一）")
	searchCmd.Flags().BoolVar(&searchFuzzy, "fuzzy", false, "启用模糊/拼音检索")
	searchCmd.Flags().StringVar(&searchMode, "mode", "and", "多关键词组合：and/or")
	searchCmd.Flags().StringVar(&searchOutput, "output", "table", "输出格式：table/json/markdown")
	searchCmd.Flags().StringVar(&searchOutputFile, "file", "", "导出文件路径")
	_ = searchCmd.MarkFlagRequired("query")
}
