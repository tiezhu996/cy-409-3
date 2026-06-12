package search

import (
	"testing"

	"lawsearch/pkg/models"
)

func TestSearchExact(t *testing.T) {
	engine := Engine{LawName: "民法典", Articles: []models.Article{{Number: 1, Content: "合同解除应当符合法律规定"}}}
	results := engine.Search("合同解除", false, "and")
	if len(results) != 1 {
		t.Fatalf("expected result")
	}
}
