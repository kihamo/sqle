package sqle

import (
	"strings"
	"sync"

	"github.com/lazada/sqle/strcase"
)

type NamingConvention interface {
	Name(string) string
	Reset() error
}

type CachedConvention struct {
	orig NamingConvention
	mu   sync.RWMutex
	name map[string]string
}

func NewCachedConvention(conv NamingConvention) *CachedConvention {
	if conv == nil {
		conv = new(NoopConvention)
	} else if c, ok := conv.(*CachedConvention); ok {
		return c
	}
	return &CachedConvention{orig: conv, name: make(map[string]string)}
}

func (c *CachedConvention) Reset() error {
	c.mu.Lock()
	c.name = make(map[string]string)
	c.mu.Unlock()
	return c.orig.Reset()
}

func (c *CachedConvention) Name(name string) string {
	if name == `` {
		return ``
	}
	c.mu.RLock()
	n, ok := c.name[name]
	c.mu.RUnlock()
	if ok {
		return n
	}
	n = c.orig.Name(name)
	c.mu.Lock()
	c.name[name] = n
	c.mu.Unlock()
	return n
}

type NoopConvention struct{}

func (n *NoopConvention) Reset() error            { return nil }
func (n *NoopConvention) Name(name string) string { return name }

type LowerConvention struct{}

func (n *LowerConvention) Reset() error            { return nil }
func (n *LowerConvention) Name(name string) string { return strings.ToLower(name) }

type UpperConvention struct{}

func (n *UpperConvention) Reset() error            { return nil }
func (n *UpperConvention) Name(name string) string { return strings.ToUpper(name) }

type SnakeConvention struct{}

func (n *SnakeConvention) Reset() error            { return nil }
func (n *SnakeConvention) Name(name string) string { return strcase.ToSnake(name) }

type CamelConvention struct{}

func (n *CamelConvention) Reset() error            { return nil }
func (n *CamelConvention) Name(name string) string { return strcase.ToCamel(name) }
