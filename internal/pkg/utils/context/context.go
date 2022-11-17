package context

import (
	"context"

	"github.com/gin-gonic/gin"

	"farmer/internal/pkg/utils/logger"
)

type Context struct {
	*gin.Context

	ctx context.Context
}

func (c *Context) Ctx() context.Context {
	return c.ctx
}

func New(c *gin.Context) *Context {
	if c == nil {
		return NewDefault()
	}
	return &Context{
		Context: c,
		ctx:     c.Request.Context(),
	}
}

func NewDefault() *Context {
	return &Context{
		Context: &gin.Context{},
		ctx:     context.Background(),
	}
}

func (c *Context) New(ctx context.Context) *Context {
	return &Context{
		Context: c.Context,
		ctx:     ctx,
	}
}

func (c *Context) Abstract() context.Context {
	return c
}

func Parse(ctx context.Context) (*Context, bool) {
	if c, ok := ctx.(*Context); ok {
		return c, true
	}
	return nil, false
}

func Child(ctx context.Context, ID string) context.Context {
	return logger.BindLogger(ctx, map[string]string{
		"ID": ID,
	})
}
