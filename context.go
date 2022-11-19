package Grapes

import "net/http"

type Context struct {
	http.ResponseWriter
	*http.Request
}

func (c *Context) SendFile(filepath string) {
	http.ServeFile(c.ResponseWriter, c.Request, filepath)
}