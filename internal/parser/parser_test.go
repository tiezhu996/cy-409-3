package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "law.txt")
	text := "第一章 总则\n第一条 为了保护民事主体合法权益。\n第二条 民法调整平等主体之间的人身关系和财产关系。\n"
	if err := os.WriteFile(file, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
	bundle, err := ParseFile(file, "民法典")
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Law.ArticleCount != 2 {
		t.Fatalf("expected 2 articles, got %d", bundle.Law.ArticleCount)
	}
}
