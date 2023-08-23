package main

import(
	"context"
	"sync"
	"github.com/ServiceWeaver/weaver"

)
// Making cache

type cacheRouter struct{}
func (cacheRouter) Get(_ context.Context, query string) string {return query}
func (cacheRouter) Put(_ context.Context, query string, _ []string) string {return query}

type Cache interface {
	// Utk return emoji di cahce dari query yang diberikan
	Get(context.Context, string) ([]string, error)

	//Untuk memasukkan query and emoji pair di cache
	Put(context.Context, string, []string) error
}



//implementasi
type cache struct {
	weaver.Implements[Cache]
	weaver.WithRouter[cacheRouter]
	mu	sync.Mutex
	emojis map[string][]string

}

func (c *cache) Init(context.Context) error{
	c.emojis = map[string][]string{}
	return nil
}

func (c *cache) Get(ctx context.Context, query string) ([]string,error){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Logger(ctx).Debug("Get", "query", query)
	return c.emojis[query], nil
}

func (c *cache) Put(ctx context.Context, query string, emojis []string) error{
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Logger(ctx).Debug("Put", "query", query)
	c.emojis[query] = emojis
	return nil

}