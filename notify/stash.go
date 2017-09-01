package notify

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

type cachedItem struct {
	item        NotifyItem
	createdTime int64
}

// NotifyStash stash notifies
type NotifyStash struct {
	sync.RWMutex
	notifyDurationSec int64
	collectDuration   time.Duration
	expired           int64
	cachedMap         map[string]cachedItem
}

// NewNotifyStash new instance
func NewNotifyStash(notifyDurationSec int64) *NotifyStash {
	if notifyDurationSec == 0 {
		notifyDurationSec = 300
	}
	return &NotifyStash{notifyDurationSec: notifyDurationSec,
		cachedMap:       make(map[string]cachedItem),
		collectDuration: 16 * time.Second,
		expired:         600,
	}
}

// Push new notify
func (p *NotifyStash) Push(key string, item NotifyItem) error {
	if key == "" {
		return nil
	}
	key = strings.ToLower(key)
	defer p.RUnlock()
	p.RLock()
	exitedItem, ok := p.cachedMap[key]
	var err error
	if ok {
		// do not send notification too quickly
		if time.Now().Unix()-exitedItem.createdTime > p.notifyDurationSec {
			p.cachedMap[key] = cachedItem{item: item, createdTime: time.Now().Unix()}
			err = item.MustNotify()
		}
		return err
	}
	err = item.MustNotify()
	if err != nil {
		log.Println(err)
	}
	p.cachedMap[key] = cachedItem{item: item, createdTime: time.Now().Unix()}
	return err
}

// Loop loop schedule
func (p *NotifyStash) Loop(ctx context.Context) {
	ticker := time.NewTicker(p.collectDuration)
	go func() {
		for {
			select {
			case <-ticker.C:
				for k, v := range p.cachedMap {
					if time.Now().Unix()-v.createdTime > p.expired {
						p.Lock()
						delete(p.cachedMap, k)
						p.Unlock()
					}
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
