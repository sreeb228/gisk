package gisk

import "sync"

type Context struct {
	variateMutex sync.RWMutex
	variateData  map[string]Value
}

func (ctx *Context) GetVariate(key string) (Value Value, ok bool) {
	ctx.variateMutex.RLock()
	defer ctx.variateMutex.RUnlock()
	if Value, ok = ctx.variateData[key]; ok {
	}
	return
}

func (ctx *Context) GetVariates() map[string]Value {
	ctx.variateMutex.RLock()
	defer ctx.variateMutex.RUnlock()
	return ctx.variateData
}

func (ctx *Context) SetVariate(key string, Value Value) {
	ctx.variateMutex.Lock()
	defer ctx.variateMutex.Unlock()
	ctx.variateData[key] = Value
}
