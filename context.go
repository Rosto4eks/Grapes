package grapes

import (
	"encoding/json"
	"net/http"
)

type Obj map[string]interface{}

type Context struct {
	http.ResponseWriter
	*http.Request
}
// sends file as a response
func (c *Context) File(filepath string) {
	http.ServeFile(c.ResponseWriter, c.Request, filepath)
}
// sends json as a response
func (c *Context) Json(message Obj) {
	json.NewEncoder(c.ResponseWriter).Encode(message)
}

func (c *Context) String(message string) {
	json.NewEncoder(c.ResponseWriter).Encode(message)
}