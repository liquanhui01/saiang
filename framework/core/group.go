package core

import (
	ctx "saiang/framework/context"
)

type GroupInterface interface {
	Get(string, ctx.HandlerFunc)
	Post(string, ctx.HandlerFunc)
	Put(string, ctx.HandlerFunc)
	Delete(string, ctx.HandlerFunc)
}

type Group struct {
	core   *Core
	prefix string
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}

func (g *Group) Get(path string, handler ctx.HandlerFunc) {
	url := g.prefix + path
	g.core.Get(url, handler)
}

func (g *Group) Post(path string, handler ctx.HandlerFunc) {
	url := g.prefix + path
	g.core.Post(url, handler)
}

func (g *Group) Put(path string, handler ctx.HandlerFunc) {
	url := g.prefix + path
	g.core.Put(url, handler)
}

func (g *Group) Delete(path string, handler ctx.HandlerFunc) {
	url := g.prefix + path
	g.core.Delete(url, handler)
}
