package grapes

import (
	"encoding/json"
	"net/http"
)

// used to send json
// grapes.Obj{"id": 1, "fruits": grapes.Obj{"apple": "green", "count": 8}}
type Obj map[string]interface{}

type Context struct {
	Request *http.Request

	Response http.ResponseWriter
	
	TreePath string
}

func (c *Context) File(filepath string) {
	http.ServeFile(c.Response, c.Request, filepath)
}

func (c *Context) Json(message Obj) {
	json.NewEncoder(c.Response).Encode(message)
}

func (c *Context) String(message string) {
	json.NewEncoder(c.Response).Encode(message)
}

// function returns param from url 
// route /Home/:index -> /Home/credits will return "credits"
func (c *Context) NamedParam(param string) string {
	urlParts := getArrPath(c.Request.URL.Path)
	treeParts := getArrPath(c.TreePath)
	for i,part := range treeParts {
		if param == part[1:] {
			return urlParts[i]
		}
	}
	return ""
}

// returns query param from url
// route /Home?id=1 will return "1"
func (c *Context) Query(param string) string {
	return c.Request.URL.Query().Get(param)
}