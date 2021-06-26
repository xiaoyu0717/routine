package routine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func init() {
	storageGCInterval = time.Millisecond * 50 // for faster test
}

func TestStorage(t *testing.T) {
	var s storage
	key := "k1"

	for i := 0; i < 100; i++ {
		src := "hello"
		s.Set(key, src)
		p := s.Get(key)
		assert.True(t, p.(string) == src)
	}

	for i := 0; i < 1000; i++ {
		num := rand.Int()
		s.Set(strconv.Itoa(num), num)
		num2 := s.Get(strconv.Itoa(num))
		assert.True(t, num2.(int) == num)
	}

	v := s.Del(key)
	assert.True(t, v != nil)

	s.Clear()
	v = s.Get(key)
	assert.True(t, v == nil)
}

func TestStorageConcurrency(t *testing.T) {
	const concurrency = 100
	const loopTimes = 100000

	var s storage

	waiter := new(sync.WaitGroup)
	waiter.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			v := rand.Uint64()
			k := fmt.Sprint()
			for i := 0; i < loopTimes; i++ {
				s.Set(k, v)
				tmp := s.Get(k)
				assert.True(t, tmp.(uint64) == v)
			}
			waiter.Done()
		}()
	}
	waiter.Wait()
}

func TestStorageGC(t *testing.T) {
	var s1, s2, s3, s4, s5 storage

	// use LocalStorage in multi goroutines
	for i := 0; i < 10; i++ {
		for i := 0; i < 1000; i++ {
			go func() {
				s1.Set("s1", "hello world")
				s2.Set("s2", true)
				s3.Set("s3", &s3)
				s4.Set("s4", rand.Int())
				s5.Set("s5", time.Now())
			}()
		}
		assert.True(t, gcRunning(), "#%v, timer not running?", i)

		// wait for a while
		time.Sleep(storageGCInterval + time.Second)
		assert.True(t, !gcRunning(), "#%v, timer not stoped?", i)
		storeMap := storages.Load().(map[int64]*store)
		assert.True(t, len(storeMap) == 0, "#%v, storeMap not empty - %d", i, len(storeMap))
	}

	//time.Sleep(time.Minute)
}

// BenchmarkLoadCurrentStore-12    	 9630090	       118.2 ns/op	      16 B/op	       1 allocs/op
func BenchmarkStorage(b *testing.B) {
	var s storage
	var variable = "hello world"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Get(variable)
		s.Set(variable, variable)
		s.Del(variable)
	}
}
