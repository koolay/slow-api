package notify

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/koolay/slow-api/notify/mocks"
)

func TestNewNotifyStash(t *testing.T) {
	NewNotifyStash(0)
}

func TestNotifyStash_Push(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		cachedMap map[string]cachedItem
	}
	type args struct {
		key  string
		item NotifyItem
	}
	m := &mocks.NotifyItem{}
	m.On("MustNotify").Return(nil)
	stash := NewNotifyStash(2)
	stash.Push("select * from user", m)
	m.AssertCalled(t, "MustNotify")

	m2 := &mocks.NotifyItem{}
	m2.On("MustNotify").Return(nil)
	stash.Push("SELECT * from user", m2)
	m2.AssertNotCalled(t, "MustNotify")

	fmt.Println("sleep 3 seconds")
	time.Sleep(3 * time.Second)
	m3 := &mocks.NotifyItem{}
	m3.On("MustNotify").Return(nil)
	stash.Push("select * from user", m3)
	m3.AssertCalled(t, "MustNotify")
}

func TestNotifyStash_Loop(t *testing.T) {
	stash := NewNotifyStash(2)
	stash.expired = 3
	stash.collectDuration = 100 * time.Millisecond
	now := time.Now().Unix()
	stash.cachedMap["k1"] = cachedItem{item: nil, createdTime: now + 1}
	stash.cachedMap["k2"] = cachedItem{item: nil, createdTime: now - 1}
	stash.cachedMap["k3"] = cachedItem{item: nil, createdTime: now - 10}
	rootCtx := context.Background()
	ctx, cancel := context.WithCancel(rootCtx)
	go func() {
		stash.Loop(ctx)
	}()
	time.Sleep(1 * time.Second)
	cancel()
	_, ok := stash.cachedMap["k2"]
	if !ok {
		t.Errorf("k2 should be existed")
	}
	_, ok = stash.cachedMap["k1"]
	if !ok {
		t.Errorf("k1 should be existed")
	}
	_, ok = stash.cachedMap["k3"]
	if ok {
		t.Errorf("k3 should be expired")
	}
}
