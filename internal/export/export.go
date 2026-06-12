package export

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"lawsearch/pkg/models"
)

func Results(results []models.SearchResult, format string, file string) error {
	var payload []byte
	var err error
	switch format {
	case "json":
		payload, err = json.MarshalIndent(results, "", "  ")
	default:
		lines := []string{"| 法规 | 条款 | 内容 |", "| --- | --- | --- |"}
		for _, result := range results {
			content := strings.ReplaceAll(result.Snippet, "\n", "<br>")
			lines = append(lines, fmt.Sprintf("| %s | 第%d条 | %s |", result.LawName, result.Article.Number, content))
		}
		payload = []byte(strings.Join(lines, "\n") + "\n")
	}
	if err != nil {
		return err
	}
	return os.WriteFile(file, payload, 0644)
}
