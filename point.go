package grapes

import "net/http"

// point allows you to create nodes for more comfortable using
// router := grapes.NewRouter()
// point := router.Point("/Home/Info")
// point.Get("/", ...)         handles /Home/Info
// point.Get("/feedback", ...) handles /Home/Info/feedback
type point struct {
	r    *Router
	path string
}

func newPoint(r *Router, path string) *point {
	point := &point{
		r:    r,
		path: path,
	}
	return point
}
// use this function to create point
func (r *Router) Point(path string) *point {
	point := newPoint(r, path)
	return point
}

func (p *point) Get(adress string, handler HandlerFunc) {
	p.r.add(http.MethodGet, p.path + adress, handler)
}

func (p *point) Post(adress string, handler HandlerFunc) {
	p.r.add(http.MethodPost, p.path + "/" + adress, handler)
}

func (p *point) Put(adress string, handler HandlerFunc) {
	p.r.add(http.MethodPut, p.path + "/" + adress, handler)
}

func (p *point) Patch(adress string, handler HandlerFunc) {
	p.r.add(http.MethodPatch, p.path + "/" + adress, handler)
}

func (p *point) Delete(adress string, handler HandlerFunc) {
	p.r.add(http.MethodDelete, p.path + "/" + adress, handler)
}

func (p *point) Head(adress string, handler HandlerFunc) {
	p.r.add(http.MethodHead, p.path + "/" + adress, handler)
}
