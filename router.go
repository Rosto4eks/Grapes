package grapes

import (
	"fmt"
	"io/fs"
	"os"
	"log"
	"net/http"
)

type HandlerFunc func(Context)

type Router struct {
	tree *node
}
// use this function to create new router
func NewRouter() *Router {
	router := Router{
		tree: &node{
			Children: make(map[string]*node),
		},
	}
	return &router
}

func (r *Router) add(method, url string, handler HandlerFunc) {
	r.tree.insert(method, url, handler)
}

func (r *Router) Get(adress string, handler HandlerFunc) {
	r.add(http.MethodGet, adress, handler)
}

func (r *Router) Post(adress string, handler HandlerFunc) {
	r.add(http.MethodPost, adress, handler)
}

func (r *Router) Put(adress string, handler HandlerFunc) {
	r.add(http.MethodPut, adress, handler)
}

func (r *Router) Patch(adress string, handler HandlerFunc) {
	r.add(http.MethodPatch, adress, handler)
}

func (r *Router) Delete(adress string, handler HandlerFunc) {
	r.add(http.MethodDelete, adress, handler)
}

func (r *Router) Head(adress string, handler HandlerFunc) {
	r.add(http.MethodHead, adress, handler)
}
// serves all http requests
// searchs in tree node with requested url
func (r *Router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	node := r.tree.search(req.URL.Path)
	if node != nil && node.Handlers[req.Method] != nil {
		node.Handlers[req.Method](Context{
			writer,
			req,
		})
		return
	}
	http.NotFound(writer, req)
}
// use for serving your static files
func (r *Router) Static(path, pattern string) {
    r.addStaticRoutes(path, pattern)
}

func (r *Router) addStaticRoutes(path, pattern string) {
	files, err := os.ReadDir(path)
    if err != nil {
        log.Fatal(err)
    }
	for _,file := range files {
		r.addStaticRoute(path, pattern, file)
	}
}

func (r *Router) addStaticRoute(path, pattern string, file fs.DirEntry) {
	if file.IsDir() {
		r.addStaticRoutes(path + "/" + file.Name(), pattern + "/" + file.Name())
	} else {
		sf := func(c Context) {
			c.File(path + "/" + file.Name())
		}
		r.Get(pattern + "/" + file.Name(), sf)
		r.Head(pattern + "/" + file.Name(), sf)
	}
}

func (r *Router) Run(port int) {
	addr := fmt.Sprintf(":%d", port)
	fmt.Println("Server started\nPORT: ", port)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalln(err.Error())
	}
}
