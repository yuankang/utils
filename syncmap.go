package utils

import "sync"

//sync.Map 没有len()方法
func SyncMapLen(sm *sync.Map) int {
	l := 0
	sm.Range(func(k, v interface{}) bool {
		l++
		return true
	})
	return l
}
