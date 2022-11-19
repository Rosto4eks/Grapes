package Grapes

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(Context)

type Route struct {
	Method  string
	Url     string
	Handler HandlerFunc
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) add(method, url string, handler HandlerFunc) {
	r.routes = append(r.routes, Route{
		Method:  method,
		Url:     url,
		Handler: handler,
	})
}

func (r *Router) Get(adress string, handler HandlerFunc) {
	r.add(http.MethodGet, adress, handler)
}

func (r *Router) Head(adress string, handler HandlerFunc) {
	r.add(http.MethodHead, adress, handler)
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	
	for _, route := range r.routes {
		if (route.Url == req.URL.Path && route.Method == req.Method) || (strings.Contains(route.Url, "*") && strings.Contains(req.URL.Path, route.Url[:len(route.Url) - 1])) {
			route.Handler(Context{
				writer,
				req,
			})
			return
		}
	}
	http.NotFound(writer, req)
}

func (r *Router) Static(path, pattern string) {
	sf := func(c Context) {
		fileServer := http.StripPrefix(pattern, http.FileServer(http.Dir(path)))
		fileServer.ServeHTTP(c.ResponseWriter, c.Request)
	}
	
	r.Get(pattern + "/*", sf)
	r.Head(pattern + "/*", sf)
}

func (r *Router) Run(port int) {
	addr := fmt.Sprintf(":%d", port)
	fmt.Println("Server started\nPORT: ", port)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalln(err.Error())
	}
}
