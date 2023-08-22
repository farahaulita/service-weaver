package main

import (
	"context"
	"fmt"
	
	"github.com/ServiceWeaver/weaver"
)

func main(){
	if err := weaver.Run(context.Background(), run); err != nil{
		panic(err)
	}
}

type app struct{
	weaver.Implements[weaver.Main]
	// Add searcher component
	searcher weaver.Ref[Searcher]
}

func run(ctx context.Context, a *app) error {
	// modify to print out  emojis
	emojis, err := a.searcher.Get().Search(ctx,"pig")
	if err != nil {
		return err
	}
	fmt.Println(emojis)

	return nil
}