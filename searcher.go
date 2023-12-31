package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
	"github.com/rivo/uniseg"
	"golang.org/x/exp/slices"
)

var (
	cacheHits = metrics.NewCounter(
		"search_cache hits",
		"Number of search cache hits",
	)
	cacheMisses = metrics.NewCounter(
		"search_cache misses",
		"Number of search cache misses",
	)
)

// sesuai dokumentasi menambah komponen service weaver, return dan input disesuaikan
type Searcher interface {
	Search(ctx context.Context, query string) ([]string, error)
	SearchChatGPT(ctx context.Context, query string) ([]string, error)
}

// sesuai dokumentasi menambah komponen service weaver
type searcher struct {
	weaver.Implements[Searcher]
	cache   weaver.Ref[Cache]
	chatgpt weaver.Ref[ChatGPT]
}

// sesuai dokumentasi menambah komponen service weaver, return dan input disesuaikan
func (s *searcher) Search(ctx context.Context, query string) ([]string, error) {
	s.Logger(ctx).Debug("Search", "query", query)

	if emojis, err := s.cache.Get().Get(ctx, query); err != nil {
		s.Logger(ctx).Error("cache.Get", "query", query, "err", err)
	} else if len(emojis) > 0 {
		cacheHits.Add(1)
		return emojis, nil
	} else {
		cacheMisses.Add(1)
	}

	// change to lowercase
	words := strings.Fields(strings.ToLower(query))

	// store results
	var results []string

	// loop through emojis file, append resultskarena isinya list emoji yg matching
	for emoji, labels := range emojis {
		if matches(labels, words) {
			results = append(results, emoji)
		}

	}
	// sort for better results
	sort.Strings(results)
	// nil is the error

	if err := s.cache.Get().Put(ctx, query, results); err != nil {
		s.Logger(ctx).Error("cache.Put", "query", query, "err", err)
	}
	return results, nil
}

func (s *searcher) SearchChatGPT(ctx context.Context, query string) ([]string, error) {
	prompt := fmt.Sprintf("Give me a list of emojis that related to the query %q. Don't give an explanation.", query)
	completion, err := s.chatgpt.Get().Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("chatgpt: %w", err)
	}
	s.Logger(ctx).Debug("ChatGPT completion", "query", query, "completion", completion)

	var results []string
	seen := map[string]bool{}
	graphemes := uniseg.NewGraphemes(completion)
	for graphemes.Next() {
		emoji := graphemes.Str()
		if _, ok := emojis[emoji]; ok && !seen[emoji] {
			results = append(results, emoji)
		}
		seen[emoji] = true
	}
	return results, nil

}

// fucktion to check if the ada yg contains, returns bool
func matches(labels, words []string) bool {
	for _, word := range words {
		if !slices.Contains(labels, word) {
			return false
		}
	}
	return true
}
