package cache

import (
	"sync"
)

type Cache struct {
	posts []Post
	mu    sync.RWMutex
}

func NewCache() Cache {
	return Cache{posts: make([]Post, 0, 100)}
}

func (c *Cache) AddPost(post Post) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.posts = append([]Post{post}, c.posts...)
	if len(c.posts) > 100 {
		c.posts = c.posts[:100]
	}
}

func (c *Cache) GetRecentPosts(limit int) []Post {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if limit > len(c.posts) {
		limit = len(c.posts)
	}
	return c.posts[:limit]
}
