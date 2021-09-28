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
	HandlerFunc func(c *Context) error
	IContext    struct {
	}

	Context struct {
		response   http.ResponseWriter
		request    *http.Request
		path       string
		hasTimeout bool
		mux        *sync.RWMutex
		ctx        context.Context
		handlers   HandlerFunc
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
		path:     r.URL.Path,
		mux:      &sync.RWMutex{},
	}
}

// GetRequest returns the request of Context
func (c *Context) GetRequest() *http.Request {
	return c.request
}

// GetRequest returns the response of Context
func (c *Context) GetResponse() http.ResponseWriter {
	return c.response
}

// GetPath return the path of Context
func (c *Context) GetPath() string {
	return c.path
}

// GetTimeout return the hasTimeout of Context
func (c *Context) GetTimeout() bool {
	return c.hasTimeout
}

// SetTimeout set hasTimeout to true
func (c *Context) SetTimeout() {
	c.hasTimeout = true
}

// #region context
func (c *Context) BaseContext() context.Context {
	return c.request.Context()
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.BaseContext().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.BaseContext().Done()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.BaseContext().Value(key)
}

// #endregion

// #region form post
func (c *Context) FormInt(key string, def int) int {
	param := c.formAll()
	if vals, ok := param[key]; ok && len(vals) > 0 {
		intVal, err := strconv.Atoi(vals[len(vals)-1])
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
	if c.request == nil {
		return map[string][]string{}
	}
	return map[string][]string(c.request.PostForm)
}

// #endregion

// #region response
// Json defines a method that set header and transfer data to JSON
func (c *Context) Json(data interface{}, status int) error {
	byt, err := json.Marshal(data)
	if err != nil {
		c.response.WriteHeader(500)
		return err
	}
	c.Blob(MIMEApplicationJSON, byt, status)
	return nil
}

func (c *Context) HTML(html string, status int) {
	c.Blob(MIMETextHTML, []byte(html), status)
}

// Blob perform specific settings, parameters include status, data and contentType
func (c *Context) Blob(contentType string, data []byte, status int) {
	if c.hasTimeout {
		return
	}
	c.response.Header().Set("Content-Type", contentType)
	c.response.WriteHeader(status)
	c.response.Write(data)
}

// #endregion
