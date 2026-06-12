package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"lawsearch/internal/store"
)

var articleLaw string
var articleNumber int
var articleFrom int
var articleTo int

var articleCmd = &cobra.Command{
	Use:   "article",
	Short: "按条款号查询",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		law, err := s.LawByName(articleLaw)
		if err != nil {
			return err
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.Header("条款", "内容")
		if articleFrom > 0 && articleTo >= articleFrom {
			for i := articleFrom; i <= articleTo; i++ {
				article, err := s.ArticleByNumber(law.ID, i)
				if err == nil {
					_ = table.Append(fmt.Sprintf("第%d条", article.Number), article.Content)
				}
			}
		} else {
			article, err := s.ArticleByNumber(law.ID, articleNumber)
			if err != nil {
				return err
			}
			_ = table.Append(fmt.Sprintf("第%d条", article.Number), article.Content)
		}
		return table.Render()
	},
}

func init() {
	rootCmd.AddCommand(articleCmd)
	articleCmd.Flags().StringVar(&articleLaw, "law", "", "法规名称或 ID")
	articleCmd.Flags().IntVar(&articleNumber, "number", 0, "条款号")
	articleCmd.Flags().IntVar(&articleFrom, "from", 0, "范围起始条款号")
	articleCmd.Flags().IntVar(&articleTo, "to", 0, "范围结束条款号")
	_ = articleCmd.MarkFlagRequired("law")
}
