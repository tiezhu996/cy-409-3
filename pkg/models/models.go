package models

import "time"

type Law struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	SourceFile   string    `json:"source_file"`
	ArticleCount int       `json:"article_count"`
	ChapterCount int       `json:"chapter_count"`
	ImportedAt   time.Time `json:"imported_at"`
}

type Section struct {
	ID       string    `json:"id"`
	LawID    string    `json:"law_id"`
	Title    string    `json:"title"`
	Level    string    `json:"level"`
	ParentID string    `json:"parent_id"`
	Children []Section `json:"children,omitempty"`
}

type Article struct {
	ID      string `json:"id"`
	LawID   string `json:"law_id"`
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Chapter string `json:"chapter"`
}

type ImportBundle struct {
	Law      Law       `json:"law"`
	Sections []Section `json:"sections"`
	Articles []Article `json:"articles"`
}

type SearchResult struct {
	LawName string  `json:"law_name"`
	Article Article `json:"article"`
	Score   int     `json:"score"`
	Snippet string  `json:"snippet"`
}
