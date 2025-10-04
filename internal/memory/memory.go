package memory

import "ethra-go/internal/types"

var store = map[string]struct {
	History []types.MemoryEntry
}{}

func AddToMemory(sessionId, role string, content interface{}, agent string) {
	if _, ok := store[sessionId]; !ok {
		store[sessionId] = struct{ History []types.MemoryEntry }{}
	}
	h := store[sessionId]
	h.History = append(h.History, types.MemoryEntry{Role: role, Content: content, Agent: agent})
	store[sessionId] = h
}

func GetMemory(sessionId string) []types.MemoryEntry {
	if h, ok := store[sessionId]; ok {
		return h.History
	}
	return nil
}

func TruncatedMemory(sessionId string, maxEntries int) []types.MemoryEntry {
	h := GetMemory(sessionId)
	if len(h) > maxEntries {
		return h[len(h)-maxEntries:]
	}
	return h
}
