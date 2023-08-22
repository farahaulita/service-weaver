package main

import (
    "context"
    "testing"
	"fmt"

    "github.com/google/go-cmp/cmp"
    "github.com/ServiceWeaver/weaver/weavertest"
)

func TestSearch(t *testing.T) {

	// test struct consists of query utk input and a list of strings for the expected output
	type test struct{
		query string
		want []string

	}

	//init	each query and expected output for every test
	for _, test := range []test{
		{"pig", []string{"🐖", "🐗", "🐷", "🐽"}},
		{"PiG", []string{"🐖", "🐗", "🐷", "🐽"}},
		{"black cat", []string{"🐈\u200d⬛"}}, 
		{"foo bar baz", nil},

	} {
	// run each tests
	 for _, runner := range weavertest.AllRunners() {
		// name, not necessary
		runner.Name = fmt.Sprintf("%s%q", runner.Name, test.query)
		// running the tests
		runner.Test(t, func(t *testing.T, searcher Searcher){
			// get output from query
			got, err := searcher.Search(context.Background(), test.query)

			// Check errors
			if err != nil {
				t.Fatalf("Search: %v", err)
			} 
			// Find Difference between expected and actual ouput
			if diff := cmp.Diff(test.want, got); diff != ""{
				t.Fatalf("Search: (-want,+got):\n%s", diff)
			}
		})
	}


	}
}