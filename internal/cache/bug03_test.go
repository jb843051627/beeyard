package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/jb843051627/beeyard/internal/model"
)

// TestBug03_CacheUpdateConcurrentSafe 验证：并发调用 Update 不会触发 data race。
// 需以 go test -race -count=20 运行；RLock 写共享 map 会被 race detector 捕获。
func TestBug03_CacheUpdateConcurrentSafe(t *testing.T) {
	rc := NewReadingCache()
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			rc.Update(&model.SensorReading{
				HiveID:     int64(n % 20),
				SensorType: model.SensorTypeTemp,
				Value:      float64(n),
				RecordedAt: time.Now(),
			})
		}(i)
	}
	wg.Wait()
	snap := rc.Snapshot()
	if len(snap) == 0 {
		t.Fatal("expected non-empty cache after concurrent updates")
	}
}
