package indexer

import (
	"strings"

	"github.com/yanyiwu/gojieba"

	"lawsearch/pkg/models"
)

func Build(articles []models.Article) map[string][]string {
	jieba := gojieba.NewJieba()
	defer jieba.Free()
	index := map[string][]string{}
	for _, article := range articles {
		tokens := jieba.CutForSearch(article.Content, true)
		tokens = append(tokens, strings.Fields(article.Content)...)
		for _, token := range tokens {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			index[token] = appendUnique(index[token], article.ID)
		}
	}
	return index
}

func appendUnique(items []string, item string) []string {
	for _, existing := range items {
		if existing == item {
			return items
		}
	}
	return append(items, item)
}
