package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"lawsearch/pkg/models"
)

var (
	articleRe = regexp.MustCompile(`^第([零一二三四五六七八九十百千万0-9]+)条[\s　]*(.*)$`)
	chapterRe = regexp.MustCompile(`^第[零一二三四五六七八九十百千万0-9]+[编章节].*`)
)

func ParseFile(file string, name string) (models.ImportBundle, error) {
	f, err := os.Open(file)
	if err != nil {
		return models.ImportBundle{}, err
	}
	defer f.Close()

	lawID := normalizeID(name)
	bundle := models.ImportBundle{
		Law: models.Law{ID: lawID, Name: name, SourceFile: file, ImportedAt: time.Now()},
	}
	currentChapter := ""
	var current *models.Article
	flush := func() {
		if current != nil {
			current.Content = strings.TrimSpace(current.Content)
			bundle.Articles = append(bundle.Articles, *current)
		}
		current = nil
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if chapterRe.MatchString(line) {
			flush()
			currentChapter = line
			bundle.Sections = append(bundle.Sections, models.Section{
				ID:    fmt.Sprintf("%s-section-%d", lawID, len(bundle.Sections)+1),
				LawID: lawID,
				Title: line,
				Level: detectLevel(line),
			})
			continue
		}
		if matches := articleRe.FindStringSubmatch(line); len(matches) == 3 {
			flush()
			number := chineseNumberToInt(matches[1])
			current = &models.Article{
				ID:      fmt.Sprintf("%s-article-%d", lawID, number),
				LawID:   lawID,
				Number:  number,
				Title:   fmt.Sprintf("第%d条", number),
				Content: strings.TrimSpace(matches[2]),
				Chapter: currentChapter,
			}
			continue
		}
		if current == nil {
			currentChapter = line
			bundle.Sections = append(bundle.Sections, models.Section{
				ID:    fmt.Sprintf("%s-section-%d", lawID, len(bundle.Sections)+1),
				LawID: lawID,
				Title: line,
				Level: "heading",
			})
		} else {
			current.Content += "\n" + line
		}
	}
	flush()
	if err := scanner.Err(); err != nil {
		return models.ImportBundle{}, err
	}
	bundle.Law.ArticleCount = len(bundle.Articles)
	bundle.Law.ChapterCount = len(bundle.Sections)
	return bundle, nil
}

func detectLevel(line string) string {
	if strings.Contains(line, "编") {
		return "part"
	}
	if strings.Contains(line, "章") {
		return "chapter"
	}
	if strings.Contains(line, "节") {
		return "section"
	}
	return "heading"
}

func normalizeID(name string) string {
	replacer := strings.NewReplacer(" ", "-", "　", "-", "中华人民共和国", "")
	id := replacer.Replace(name)
	if id == "" {
		return "law"
	}
	return id
}

func chineseNumberToInt(raw string) int {
	if n, err := strconv.Atoi(raw); err == nil {
		return n
	}
	digits := map[rune]int{'零': 0, '一': 1, '二': 2, '两': 2, '三': 3, '四': 4, '五': 5, '六': 6, '七': 7, '八': 8, '九': 9}
	units := map[rune]int{'十': 10, '百': 100, '千': 1000, '万': 10000}
	total, section, number := 0, 0, 0
	for _, r := range raw {
		if value, ok := digits[r]; ok {
			number = value
			continue
		}
		if unit, ok := units[r]; ok {
			if number == 0 {
				number = 1
			}
			section += number * unit
			number = 0
			if unit == 10000 {
				total += section
				section = 0
			}
		}
	}
	return total + section + number
}
