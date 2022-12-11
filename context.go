package grapes

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"mime/multipart"
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

// used to set status code to the response
func (c *Context) Status(status int) {
	c.Response.WriteHeader(status)
}

// used to send file of any type, also can be used for sending html files
func (c *Context) SendFile(filepath string) {
	http.ServeFile(c.Response, c.Request, filepath)
}

func (c *Context) SendJson(message Obj) {
	json.NewEncoder(c.Response).Encode(message)
}

func (c *Context) SendString(message string) {
	json.NewEncoder(c.Response).Encode(message)
}

// redirect incoming request to another url
func (c *Context) Redirect(url string) {
	http.Redirect(c.Response, c.Request, url, http.StatusMovedPermanently)
}

// used to send html with parameters, all parameters must start with Uppercase symbol
func (c *Context) Template(path string, data Obj) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(c.Response, data)
}

// function returns param from url 
// route /Home/:index -> /Home/credits will return "credits"
func (c *Context) GetNamedParam(param string) string {
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
func (c *Context) GetQueryParam(param string) string {
	return c.Request.URL.Query().Get(param)
}

// used to get file from post form
func (c *Context) GetFormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, nil, err
	}
	file, header, err := c.Request.FormFile(key)
	if err != nil {
		return nil, nil, err
	}
	return file, header, nil
}

// used to get value from post form
func (c *Context) GetFormValue(key string) (string, error) {
	value := c.Request.FormValue(key)
	if value == "" {
		return "", errors.New("incorrect key")
	}
	return value, nil
}