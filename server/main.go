package main

import (
	"container/list"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type CacheItem struct {
	key        string
	value      string
	expiration time.Time
	ttl        int64
}

type LRUCache struct {
	capacity int
	mutex    sync.Mutex
	items    map[string]*list.Element
	order    *list.List
}

var broadcast = make(chan map[string]map[string]interface{})

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

func (cache *LRUCache) Get(key string) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if element, found := cache.items[key]; found {
		cache.order.MoveToFront(element)
		item := element.Value.(*CacheItem)
		if item.expiration.After(time.Now()) {
			return item.value, true
		}
		cache.order.Remove(element)
		delete(cache.items, key)
	}
	return nil, false
}

func (cache *LRUCache) Set(key string, value string, ttl time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if element, found := cache.items[key]; found {
		cache.order.MoveToFront(element)
		item := element.Value.(*CacheItem)
		item.value = value
		item.ttl = int64(ttl.Seconds())
		item.expiration = time.Now().Add(ttl)
	} else {
		if cache.order.Len() == cache.capacity {
			oldest := cache.order.Back()
			if oldest != nil {
				cache.order.Remove(oldest)
				delete(cache.items, oldest.Value.(*CacheItem).key)
			}
		}
		item := &CacheItem{key: key, value: value, ttl: int64(ttl.Seconds()), expiration: time.Now().Add(ttl)}
		element := cache.order.PushFront(item)
		cache.items[key] = element
	}
	broadcast <- cache.copyItems()
}

func (cache *LRUCache) Delete(key string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if element, found := cache.items[key]; found {
		cache.order.Remove(element)
		delete(cache.items, key)
	}
	broadcast <- cache.copyItems()
}

func (cache *LRUCache) copyItems() map[string]map[string]interface{} {
	copy := make(map[string]map[string]interface{})
	for key, element := range cache.items {
		item := element.Value.(*CacheItem)
		copy[key] = map[string]interface{}{
			"key":        item.key,
			"expiration": item.expiration,
		}
	}
	return copy
}

func main() {
	cache := NewLRUCache(100)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	r.GET("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		if value, found := cache.Get(key); found {
			c.JSON(http.StatusOK, gin.H{"value": value})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "no data found"})
		}
	})

	r.GET("/cache", func(c *gin.Context) {
		if len(cache.items) > 0 {
			cacheCopy := cache.copyItems()
			c.JSON(http.StatusOK, gin.H{"cache": cacheCopy})
		} else {
			c.JSON(http.StatusOK, gin.H{"error": "no data present"})
		}
	})

	r.POST("/cache", func(c *gin.Context) {
		var req struct {
			Key   string        `json:"key"`
			Value string        `json:"value"`
			TTL   time.Duration `json:"ttl"`
		}
		if err := c.BindJSON(&req); err == nil {
			cache.Set(req.Key, req.Value, req.TTL*time.Second)
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	r.DELETE("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		cache.Delete(key)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	go func() {
		for {
			time.Sleep(time.Second * 1)
			cache.mutex.Lock()
			for key, element := range cache.items {
				if element.Value.(*CacheItem).expiration.Before(time.Now()) {
					cache.order.Remove(element)
					delete(cache.items, key)
				}
			}
			cache.mutex.Unlock()
			broadcast <- cache.copyItems()
		}
	}()

	r.Run(":8080")
}
