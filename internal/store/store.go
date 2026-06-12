package store

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go.etcd.io/bbolt"

	"lawsearch/pkg/models"
)

var (
	lawsBucket     = []byte("laws")
	articlesBucket = []byte("articles")
	sectionsBucket = []byte("sections")
	indexBucket    = []byte("index")
)

type Store struct {
	db *bbolt.DB
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	return s, db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{lawsBucket, articlesBucket, sectionsBucket, indexBucket} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) SaveBundle(bundle models.ImportBundle, tokenIndex map[string][]string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		if err := putJSON(tx.Bucket(lawsBucket), bundle.Law.ID, bundle.Law); err != nil {
			return err
		}
		for _, section := range bundle.Sections {
			if err := putJSON(tx.Bucket(sectionsBucket), section.ID, section); err != nil {
				return err
			}
		}
		for _, article := range bundle.Articles {
			if err := putJSON(tx.Bucket(articlesBucket), article.ID, article); err != nil {
				return err
			}
		}
		for token, ids := range tokenIndex {
			if err := putJSON(tx.Bucket(indexBucket), bundle.Law.ID+"::"+token, ids); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Laws() ([]models.Law, error) {
	var result []models.Law
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(lawsBucket).ForEach(func(_, value []byte) error {
			var law models.Law
			if err := json.Unmarshal(value, &law); err != nil {
				return err
			}
			result = append(result, law)
			return nil
		})
	})
	return result, err
}

func (s *Store) LawByName(name string) (models.Law, error) {
	laws, err := s.Laws()
	if err != nil {
		return models.Law{}, err
	}
	for _, law := range laws {
		if law.Name == name || law.ID == name {
			return law, nil
		}
	}
	return models.Law{}, fmt.Errorf("未找到法规: %s", name)
}

func (s *Store) Articles(lawID string) ([]models.Article, error) {
	var articles []models.Article
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(articlesBucket).ForEach(func(_, value []byte) error {
			var article models.Article
			if err := json.Unmarshal(value, &article); err != nil {
				return err
			}
			if article.LawID == lawID {
				articles = append(articles, article)
			}
			return nil
		})
	})
	return articles, err
}

func (s *Store) Sections(lawID string) ([]models.Section, error) {
	var sections []models.Section
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(sectionsBucket).ForEach(func(_, value []byte) error {
			var section models.Section
			if err := json.Unmarshal(value, &section); err != nil {
				return err
			}
			if section.LawID == lawID {
				sections = append(sections, section)
			}
			return nil
		})
	})
	return sections, err
}

func (s *Store) ArticleByNumber(lawID string, number int) (models.Article, error) {
	var article models.Article
	err := s.db.View(func(tx *bbolt.Tx) error {
		raw := tx.Bucket(articlesBucket).Get([]byte(lawID + "-article-" + strconv.Itoa(number)))
		if raw == nil {
			return fmt.Errorf("未找到第 %d 条", number)
		}
		return json.Unmarshal(raw, &article)
	})
	return article, err
}

func (s *Store) DeleteLaw(lawID string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		_ = tx.Bucket(lawsBucket).Delete([]byte(lawID))
		for _, bucket := range []*bbolt.Bucket{tx.Bucket(articlesBucket), tx.Bucket(sectionsBucket), tx.Bucket(indexBucket)} {
			var keys [][]byte
			_ = bucket.ForEach(func(key, value []byte) error {
				if len(key) >= len(lawID) && string(key[:len(lawID)]) == lawID {
					keys = append(keys, append([]byte(nil), key...))
				}
				return nil
			})
			for _, key := range keys {
				_ = bucket.Delete(key)
			}
		}
		return nil
	})
}

func (s *Store) UpdateLaw(law models.Law) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		return putJSON(tx.Bucket(lawsBucket), law.ID, law)
	})
}

func (s *Store) AddTags(lawID string, tags []string) error {
	law, err := s.getLaw(lawID)
	if err != nil {
		return err
	}
	law.Tags = mergeUnique(law.Tags, tags)
	return s.UpdateLaw(law)
}

func (s *Store) RemoveTags(lawID string, tags []string) error {
	law, err := s.getLaw(lawID)
	if err != nil {
		return err
	}
	law.Tags = removeItems(law.Tags, tags)
	return s.UpdateLaw(law)
}

func (s *Store) SetTags(lawID string, tags []string) error {
	law, err := s.getLaw(lawID)
	if err != nil {
		return err
	}
	law.Tags = dedupe(tags)
	return s.UpdateLaw(law)
}

func (s *Store) LawsByTags(tags []string) ([]models.Law, error) {
	laws, err := s.Laws()
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return laws, nil
	}
	var result []models.Law
	for _, law := range laws {
		if hasAllTags(law.Tags, tags) {
			result = append(result, law)
		}
	}
	return result, nil
}

func (s *Store) AllTags() ([]string, error) {
	laws, err := s.Laws()
	if err != nil {
		return nil, err
	}
	tagSet := make(map[string]bool)
	for _, law := range laws {
		for _, tag := range law.Tags {
			tagSet[tag] = true
		}
	}
	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	return tags, nil
}

func (s *Store) getLaw(lawID string) (models.Law, error) {
	var law models.Law
	err := s.db.View(func(tx *bbolt.Tx) error {
		raw := tx.Bucket(lawsBucket).Get([]byte(lawID))
		if raw == nil {
			return fmt.Errorf("未找到法规: %s", lawID)
		}
		return json.Unmarshal(raw, &law)
	})
	return law, err
}

func mergeUnique(existing, newTags []string) []string {
	set := make(map[string]bool)
	for _, tag := range existing {
		set[tag] = true
	}
	for _, tag := range newTags {
		set[tag] = true
	}
	var result []string
	for tag := range set {
		result = append(result, tag)
	}
	return result
}

func removeItems(items, toRemove []string) []string {
	removeSet := make(map[string]bool)
	for _, tag := range toRemove {
		removeSet[tag] = true
	}
	var result []string
	for _, tag := range items {
		if !removeSet[tag] {
			result = append(result, tag)
		}
	}
	return result
}

func dedupe(tags []string) []string {
	set := make(map[string]bool)
	var result []string
	for _, tag := range tags {
		if !set[tag] {
			set[tag] = true
			result = append(result, tag)
		}
	}
	return result
}

func hasAllTags(lawTags, requiredTags []string) bool {
	tagSet := make(map[string]bool)
	for _, tag := range lawTags {
		tagSet[tag] = true
	}
	for _, required := range requiredTags {
		if !tagSet[required] {
			return false
		}
	}
	return true
}

func putJSON(bucket *bbolt.Bucket, key string, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return bucket.Put([]byte(key), payload)
}
