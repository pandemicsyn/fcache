/*
Package fcache implements a basic fast in memory k/v cache

Its a simple cache with TTL support that does what *I needed*.
It probably doesn't do what you need, but thats okay...
Inspired by go-cache and others (that where missing bits and pieces).

*/
package fcache

import (
	"fmt"
	"sync"
	"time"
)

const (
	// NoExpiration constant for NO ttl's
	NoExpiration time.Duration = -1
	// DefaultExpiration constant for the default ttl
	DefaultExpiration time.Duration = 0
)

// FCache instance
type FCache struct {
	*fcache
}

type fcache struct {
	sync.RWMutex
	items map[string]*Item
	p     bool
	s     int
	e     time.Duration
}

// Item is our cache object and its expiration time
type Item struct {
	Object     interface{}
	Expiration *time.Time
}

//Expired returns whether an item is expired or not
func (i *Item) Expired() bool {
	if i.Expiration == nil {
		return false
	}
	return i.Expiration.Before(time.Now())
}

//Set a key/value/ttl
func (c *fcache) Set(k string, v interface{}, ttl time.Duration) {
	c.Lock()
	c.setWithTTL(k, v, ttl)
	c.Unlock()
}

func (c *fcache) setWithTTL(k string, v interface{}, ttl time.Duration) {
	var e *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		e = &t
	}
	c.items[k] = &Item{
		Object:     v,
		Expiration: e,
	}
}

//IncrementInt increment's an int, create the key initialized to 0 and
//then increment it by the given if it doesn't exist or is expired.
func (c *fcache) IncrementInt(k string, n int) error {
	c.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.setWithTTL(k, 0+n, NoExpiration)
		c.Unlock()
		return nil
	}
	rv, ok := v.Object.(int)
	if !ok {
		c.Unlock()
		return fmt.Errorf("Value not an int")
	}
	v.Object = rv + n
	c.Unlock()
	return nil
}

//IncrementFloat64 increments an in64, create the key initialized to 0 and
//then increment it by the given if it doesn't exist or is expired.
func (c *fcache) IncrementInt64(k string, n int64) error {
	c.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.setWithTTL(k, 0+n, NoExpiration)
		c.Unlock()
		return nil
	}
	rv, ok := v.Object.(int64)
	if !ok {
		c.Unlock()
		return fmt.Errorf("Value not an int64")
	}
	v.Object = rv + n
	c.Unlock()
	return nil
}

//IncrementFloat64 increment a float64, create the key initialized to 0 and
//then increment it by the given if it doesn't exist or is expired.
func (c *fcache) IncrementFloat64(k string, n float64) error {
	c.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		if found && v.Expired() {
			go c.asyncExpiredDel(k)
		}
		c.Unlock()
		return fmt.Errorf("No such key")
	}
	rv, ok := v.Object.(float64)
	if !ok {
		c.Unlock()
		return fmt.Errorf("Value not an float64")
	}
	v.Object = rv + n
	c.Unlock()
	return nil
}

// Get a key from the cache, if an expired key is encountered
// we'll spawn a goroutine to reap the key and immediately return
// that it doesn't exist.
func (c *fcache) Get(k string) (interface{}, bool) {
	c.RLock()
	v, found := c.get(k)
	c.RUnlock()
	return v, found
}

func (c *fcache) get(k string) (interface{}, bool) {
	item, found := c.items[k]
	if !found || item.Expired() {
		if found && item.Expired() {
			go c.asyncExpiredDel(k)
		}
		return nil, false
	}
	return item.Object, true
}

// async delete a key if expired
func (c *fcache) asyncExpiredDel(k string) {
	c.Lock()
	item, found := c.items[k]
	if found && item.Expired() {
		delete(c.items, k)
	}
	c.Unlock()
}

// Delete a key from the cache
func (c *fcache) Delete(k string) {
	c.Lock()
	delete(c.items, k)
	c.Unlock()
}

// Empty the cache completely
func (c *fcache) Empty() {
	c.Lock()
	if c.p {
		c.items = make(map[string]*Item, c.s)
	} else {
		c.items = make(map[string]*Item)
	}
	c.Unlock()
}

//NewPreAllocated cache creates a new cache with preallocated map of given size
func NewPreAllocated(size int) *FCache {
	items := make(map[string]*Item, size)
	return &FCache{newCache(items, true, size, NoExpiration)}
}

func newCache(items map[string]*Item, p bool, s int, e time.Duration) *fcache {
	c := &fcache{
		items: items,
		p:     p,
		s:     s,
		e:     e,
	}
	return c
}
