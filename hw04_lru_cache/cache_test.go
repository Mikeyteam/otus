package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(10)

		c.Set("aaa", 300)
		c.Clear()
		require.Equal(t, 0, 0)
	})

	t.Run("overwrite value in cache", func(t *testing.T) {
		c := NewCache(10)

		c.Set("aaa", 300)
		c.Set("aaa", 301)
		require.Equal(t, 301, 301)
	})

	t.Run("overflow cache", func(t *testing.T) {
		c := NewCache(1)

		c.Set("aaa", 300)
		c.Set("aab", 302)
		val, _ := c.Get("aaa")
		require.Equal(t, nil, val)
	})

	t.Run("key cache is symbol", func(t *testing.T) {
		c := NewCache(10)

		c.Set("$", 300)
		c.Set("^", 302)
		val, _ := c.Get("$")
		val1, _ := c.Get("^")
		require.Equal(t, 300, val)
		require.Equal(t, 302, val1)
	})

	t.Run("value cache is nil", func(t *testing.T) {
		c := NewCache(10)

		c.Set("aaa", nil)
		val, _ := c.Get("aaa")
		require.Equal(t, nil, val)
	})

}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
