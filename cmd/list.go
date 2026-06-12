package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"lawsearch/internal/store"
)

var listTags []string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出已导入法规",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		laws, err := s.LawsByTags(listTags)
		if err != nil {
			return err
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.Header("名称", "条款数", "章节数", "标签", "导入时间")
		for _, law := range laws {
			tagStr := "-"
			if len(law.Tags) > 0 {
				tagStr = strings.Join(law.Tags, ", ")
			}
			_ = table.Append(law.Name, fmt.Sprintf("%d", law.ArticleCount), fmt.Sprintf("%d", law.ChapterCount), tagStr, law.ImportedAt.Format("2006-01-02 15:04"))
		}
		return table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringSliceVar(&listTags, "tags", nil, "按标签过滤，多个用逗号分隔")
}
