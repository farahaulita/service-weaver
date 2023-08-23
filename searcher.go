package main


import (
	"context"
	"sort"
	"strings"

	"github.com/ServiceWeaver/weaver" 
 	"golang.org/x/exp/slices"
)

// sesuai dokumentasi menambah komponen service weaver, return dan input disesuaikan
type Searcher interface{
	Search(ctx context.Context, query string)([]string,error)
}

// sesuai dokumentasi menambah komponen service weaver
type searcher struct{
	weaver.Implements[Searcher]
	cache weaver.Ref[Cache]
}

// sesuai dokumentasi menambah komponen service weaver, return dan input disesuaikan
func (s * searcher) Search(ctx context.Context, query string) ([]string,error) {
	s.Logger(ctx).Debug("Search", "query",query)

	if emojis, err:= s.cache.Get().Get(ctx,query); err !=nil {
		s.Logger(ctx).Error("cache.Get","query",query, "err", err)
	} else if len(emojis) > 0{
		return emojis, nil
	}

	
	
	// change to lowercase
	words := strings.Fields(strings.ToLower(query))

	// store results
	var results  []string

	// loop through emojis file, append resultskarena isinya list emoji yg matching
	for emoji, labels := range emojis {
		if matches(labels,words){
			results = append(results, emoji)
		}

	}
	// sort for better results
	sort.Strings(results)
	// nil is the error

	if err := s.cache.Get().Put(ctx,query,results); err != nil{
		s.Logger(ctx).Error("cache.Put","query",query, "err", err)
	}
	return results, nil
}

// fucktion to check if the ada yg contains, returns bool
func matches(labels, words []string) bool {
	for _, word := range words{
		if !slices.Contains(labels,word){
			return false
		}
	}
	return true
}