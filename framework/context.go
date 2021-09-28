package framework

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type (
	HandlerFunc func(w http.ResponseWriter, r http.Request)
	IContext    interface {
		GetRequest() *http.Request
		GetResponse() http.ResponseWriter
		HasTimeout() bool
	}
	Context struct {
		response   http.ResponseWriter
		request    *http.Request
		ctx        context.Context
		mux        *sync.RWMutex
		hasTimeout bool
		handler    HandlerFunc
	}
)

const (
	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + ";" + charsetUTF8
	MIMEApplicationXML             = "application/xml"
	MIMEApplicationXMLCharsetUTF8  = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextHTML                   = "text/html"
	MIMETextHTMLCharsetUTF8        = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                  = "text/plain"
	MIMETextPlainCharsetUTF8       = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm              = "multipart/form-data"
)

const (
	charsetUTF8 = "charset=UTF-8"
)

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		response: w,
		request:  r,
		mux:      &sync.RWMutex{},
	}
}

// RWMux returns the mux of Context
func (c *Context) RWMux() *sync.RWMutex {
	return c.mux
}

// GetRequest returns the request of Context
func (c *Context) GetRequest() *http.Request {
	return c.request
}

// GetResponse return the response of Context
func (c *Context) GetResponse() http.ResponseWriter {
	return c.response
}

// HasTimeout returns the hasTimeout of Context
func (c *Context) HasTimeout() bool {
	return c.hasTimeout
}

// SethasTimeout set the hasTimeout of Context to true
func (c *Context) SethasTimeout() {
	c.hasTimeout = true
}

// BaseContext returns base context from request
func (c *Context) BaseContext() context.Context {
	return c.request.Context()
}

// Deadline returns two variables, deadline is time type, ok is bool type
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.BaseContext().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.BaseContext().Done()
}

// Value defines a function that both the parameters and return values are
// interface types
func (c *Context) Value(key interface{}) interface{} {
	return c.BaseContext().Value(key)
}

// #region form post

func (c *Context) FormInt(key string, def int) int {
	param := c.formAll()
	if val, ok := param[key]; ok && len(val) > 0 {
		intVal, err := strconv.Atoi(val[len(val)-1])
		if err != nil {
			return def
		}
		return intVal
	}
	return def
}

func (c *Context) FormString(key string, def []string) []string {
	param := c.formAll()
	if vals, ok := param[key]; ok {
		return vals
	}
	return def
}

func (c *Context) formAll() map[string][]string {
	if c.request != nil {
		return map[string][]string(c.request.PostForm)
	}
	return map[string][]string{}
}

// #endregion

// #region response
// Json defines a method that set header and transfer obj to JSON
func (c *Context) Json(status int, obj interface{}) error {
	byt, err := json.Marshal(obj)
	if err != nil {
		c.response.WriteHeader(500)
		return err
	}
	c.Blob(status, byt, MIMEApplicationJSON)
	return nil
}

// HTML
func (c *Context) HTML(status int, html string) {
	c.Blob(status, []byte(html), MIMETextHTML)
}

// Blob perform specific settings, parameters include status, data and contentType
func (c *Context) Blob(status int, data []byte, contentType string) {
	if c.HasTimeout() {
		return
	}
	c.response.Header().Set("Content-Type", contentType)
	c.response.WriteHeader(status)
	c.response.Write(data)
}

// # endregion
