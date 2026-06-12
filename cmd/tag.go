package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"lawsearch/internal/store"
)

var tagLaw string
var tagNames []string

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "管理法规标签",
}

var tagAddCmd = &cobra.Command{
	Use:   "add",
	Short: "为法规添加标签",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		law, err := s.LawByName(tagLaw)
		if err != nil {
			return err
		}
		if err := s.AddTags(law.ID, tagNames); err != nil {
			return err
		}
		law, err = s.LawByName(tagLaw)
		if err != nil {
			return err
		}
		fmt.Printf("已添加标签，当前标签：[%s]\n", strings.Join(law.Tags, ", "))
		return nil
	},
}

var tagRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "为法规移除标签",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		law, err := s.LawByName(tagLaw)
		if err != nil {
			return err
		}
		if err := s.RemoveTags(law.ID, tagNames); err != nil {
			return err
		}
		law, err = s.LawByName(tagLaw)
		if err != nil {
			return err
		}
		tagStr := "无"
		if len(law.Tags) > 0 {
			tagStr = fmt.Sprintf("[%s]", strings.Join(law.Tags, ", "))
		}
		fmt.Printf("已移除标签，当前标签：%s\n", tagStr)
		return nil
	},
}

var tagSetCmd = &cobra.Command{
	Use:   "set",
	Short: "设置法规标签（覆盖原有标签）",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		law, err := s.LawByName(tagLaw)
		if err != nil {
			return err
		}
		if err := s.SetTags(law.ID, tagNames); err != nil {
			return err
		}
		law, err = s.LawByName(tagLaw)
		if err != nil {
			return err
		}
		tagStr := "无"
		if len(law.Tags) > 0 {
			tagStr = fmt.Sprintf("[%s]", strings.Join(law.Tags, ", "))
		}
		fmt.Printf("已设置标签，当前标签：%s\n", tagStr)
		return nil
	},
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有已使用的标签",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		tags, err := s.AllTags()
		if err != nil {
			return err
		}
		if len(tags) == 0 {
			fmt.Println("暂无标签")
			return nil
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.Header("标签")
		for _, tag := range tags {
			_ = table.Append(tag)
		}
		return table.Render()
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.AddCommand(tagAddCmd)
	tagCmd.AddCommand(tagRemoveCmd)
	tagCmd.AddCommand(tagSetCmd)
	tagCmd.AddCommand(tagListCmd)

	for _, subCmd := range []*cobra.Command{tagAddCmd, tagRemoveCmd, tagSetCmd} {
		subCmd.Flags().StringVar(&tagLaw, "law", "", "法规名称或 ID")
		subCmd.Flags().StringSliceVar(&tagNames, "tags", nil, "标签列表，多个用逗号分隔")
		_ = subCmd.MarkFlagRequired("law")
		_ = subCmd.MarkFlagRequired("tags")
	}
}
