package Grapes

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	http.ResponseWriter
	*http.Request
}
// sends file as a response
func (c *Context) File(filepath string) {
	http.ServeFile(c.ResponseWriter, c.Request, filepath)
}
// sends json as a response
func (c *Context) Json(message any) {
	json.NewEncoder(c.ResponseWriter).Encode(message)
}