package nano

import (
	"sync"
	"time"

	"github.com/FloatTech/ttl"
)

var (
	triggeredMessages   = ttl.NewCache[string, []string](time.Minute * 5)
	triggeredMessagesMu = sync.Mutex{}
)

func logtriggeredmessages(id, reply string) {
	triggeredMessagesMu.Lock()
	defer triggeredMessagesMu.Unlock()
	triggeredMessages.Set(id, append(triggeredMessages.Get(id), reply))
}

// GetTriggeredMessages 获取被 id 消息触发的回复消息 id
func GetTriggeredMessages(id string) []string {
	triggeredMessagesMu.Lock()
	defer triggeredMessagesMu.Unlock()
	return triggeredMessages.Get(id)
}
