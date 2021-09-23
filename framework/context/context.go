package context

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handler        HandlerFunc

	// mark whether the timeout
	hasTimeout bool
	// the protection mechanism for writing
	writerMux *sync.Mutex
}

// NewContext initialize the Context sturct
func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		writerMux:      &sync.Mutex{},
	}
}

// GetResponse returns the writerMux field in the Context structure
func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

// GetResponse returns the request field in the Context structure
func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

// GetResponse returns the responseWriter field in the Context structure
func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

// SetHasTimeout set the hashTimeout field in the Context structure to true
func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

// HasTimeout returns the hashTimeout field in the Context structure
func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

// #region form post

func (ctx *Context) FormInt(key string, def int) int {
	params := ctx.formAll()
	if vals, ok := params[key]; ok {
		le := len(vals)
		if le > 0 {
			intval, err := strconv.Atoi(vals[le-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) FormString(key string, def string) string {
	params := ctx.formAll()
	if vals, ok := params[key]; ok {
		le := len(vals)
		if le > 0 {
			return vals[le-1]
		}
	}
	return def
}

func (ctx *Context) FormArray(key string, def []string) []string {
	params := ctx.formAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

func (ctx *Context) formAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}

// #endregion

// #region response
// Json define a method that set header and tansfer obj to json
func (ctx *Context) Json(status int, obj interface{}) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.responseWriter.WriteHeader(500)
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

// HTML define a method that set header and tansfer obj to html
func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "text/html")
	ctx.responseWriter.WriteHeader(status)
	ctx.responseWriter.Write([]byte(template))
	return nil
}

// Text
func (ctx *Context) Text(status int, obj interface{}) error {
	if ctx.HasTimeout() {
		return nil
	}
	return nil
}

// #endregion
