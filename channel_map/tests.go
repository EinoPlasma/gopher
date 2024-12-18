package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestSafeMap_Concurrent(t *testing.T) {
	const (
		numGoroutines = 1000
		numOperations = 1000
	)

	sm := NewSafeMap[int, int]()
	var wg sync.WaitGroup

	// 测试并发写入
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := id*numOperations + j
				sm.Set(key, key)
			}
		}(i)
	}
	wg.Wait()

	// 验证写入的正确性
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numOperations; j++ {
			key := i*numOperations + j
			value, exists := sm.Get(key)
			if !exists || value != key {
				t.Errorf("Expected key %d to exist with value %d, got value %d", key, key, value)
			}
		}
	}

	// 测试并发读写
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := id*numOperations + j
				sm.Set(key, key*2)
				value, exists := sm.Get(key)
				if !exists || value != key*2 {
					t.Errorf("Expected key %d to exist with value %d, got value %d", key, key*2, value)
				}
			}
		}(i)
	}
	wg.Wait()

	// 验证并发读写的正确性
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numOperations; j++ {
			key := i*numOperations + j
			value, exists := sm.Get(key)
			if !exists || value != key*2 {
				t.Errorf("Expected key %d to exist with value %d, got value %d", key, key*2, value)
			}
		}
	}

	// 测试并发删除
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := id*numOperations + j
				sm.Delete(key)
			}
		}(i)
	}
	wg.Wait()

	// 验证删除的正确性
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numOperations; j++ {
			key := i*numOperations + j
			_, exists := sm.Get(key)
			if exists {
				t.Errorf("Expected key %d to be deleted, but it still exists", key)
			}
		}
	}

	// 验证长度
	if sm.Len() != 0 {
		t.Errorf("Expected map length to be 0, got %d", sm.Len())
	}
}

func BenchmarkSafeMap_Set(b *testing.B) {
	sm := NewSafeMap[int, int]()

	for i := 0; i < b.N; i++ {
		sm.Set(i, i)
	}
}

func BenchmarkSafeMap_Get(b *testing.B) {
	sm := NewSafeMap[int, int]()
	for i := 0; i < b.N; i++ {
		sm.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Get(i)
	}
}

func BenchmarkSafeMap_Delete(b *testing.B) {
	sm := NewSafeMap[int, int]()
	for i := 0; i < b.N; i++ {
		sm.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Delete(i)
	}
}

func main() {
	// 创建一个 testing.T 对象，用于手动调用测试函数
	t := &testing.T{}

	// 调用测试函数
	fmt.Println("Running TestSafeMap_Concurrent...")
	TestSafeMap_Concurrent(t)

	// 输出测试结果
	if !t.Failed() {
		fmt.Println("All tests passed!")
	} else {
		fmt.Println("Some tests failed!")
	}

	// 运行基准测试（注意：基准测试不会在 main 中运行，需要使用 go test -bench=.）
	fmt.Println("To run benchmarks, use 'go test -bench=.'")
}
