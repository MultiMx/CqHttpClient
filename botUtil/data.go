package botUtil

import "sync"

type botData struct {
	Groups     map[int]bool
	GroupsLock *sync.RWMutex
}

var BotData = botData{
	GroupsLock: &sync.RWMutex{},
}

func (a *botData) WriteGroups(k int, v bool) {
	a.GroupsLock.Lock()
	a.Groups[k] = v
	a.GroupsLock.Unlock()
}

func (a *botData) ReadGroups(k int) bool {
	a.GroupsLock.RLock()
	t := a.Groups[k]
	a.GroupsLock.RUnlock()
	return t
}
