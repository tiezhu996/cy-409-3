package search

import (
	"sort"
	"strings"

	pinyin "github.com/mozillazg/go-pinyin"

	"lawsearch/pkg/models"
)

type Engine struct {
	Articles []models.Article
	LawName  string
}

func (e Engine) Search(query string, fuzzy bool, mode string) []models.SearchResult {
	terms := strings.Fields(query)
	if len(terms) == 0 {
		terms = []string{query}
	}
	var results []models.SearchResult
	for _, article := range e.Articles {
		matches := 0
		for _, term := range terms {
			if term == "" {
				continue
			}
			if matchArticle(article, term, fuzzy) {
				matches++
			}
		}
		if (mode == "or" && matches > 0) || (mode != "or" && matches == len(terms)) {
			results = append(results, models.SearchResult{
				LawName: e.LawName,
				Article: article,
				Score:   matches,
				Snippet: highlight(article.Content, terms),
			})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Article.Number < results[j].Article.Number
	})
	return results
}

func matchArticle(article models.Article, term string, fuzzy bool) bool {
	content := article.Title + article.Content
	if strings.Contains(content, term) {
		return true
	}
	if !fuzzy {
		return false
	}
	queryPY := strings.Join(pinyin.LazyConvert(term, nil), "")
	contentPY := strings.Join(pinyin.LazyConvert(content, nil), "")
	if queryPY != "" && strings.Contains(contentPY, queryPY) {
		return true
	}
	for _, word := range strings.Fields(content) {
		if levenshtein(word, term) <= 2 {
			return true
		}
	}
	return false
}

func highlight(content string, terms []string) string {
	snippet := content
	if len([]rune(snippet)) > 120 {
		snippet = string([]rune(snippet)[:120]) + "..."
	}
	for _, term := range terms {
		if term != "" {
			snippet = strings.ReplaceAll(snippet, term, "**"+term+"**")
		}
	}
	return snippet
}

func levenshtein(a, b string) int {
	ar, br := []rune(a), []rune(b)
	dp := make([][]int, len(ar)+1)
	for i := range dp {
		dp[i] = make([]int, len(br)+1)
	}
	for i := range ar {
		dp[i+1][0] = i + 1
	}
	for j := range br {
		dp[0][j+1] = j + 1
	}
	for i := 1; i <= len(ar); i++ {
		for j := 1; j <= len(br); j++ {
			cost := 0
			if ar[i-1] != br[j-1] {
				cost = 1
			}
			dp[i][j] = min(dp[i-1][j]+1, dp[i][j-1]+1, dp[i-1][j-1]+cost)
		}
	}
	return dp[len(ar)][len(br)]
}

func min(values ...int) int {
	result := values[0]
	for _, value := range values[1:] {
		if value < result {
			result = value
		}
	}
	return result
}
