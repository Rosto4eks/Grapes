package grapes

import (
	"strings"
)

type node struct {
	Handlers map[string] HandlerFunc // where key is method
	Children map[string]*node // where key is part of URL
}


func (n *node) insert(method string, url string, handler HandlerFunc) {
	parts := strings.Split(url, "/")[1:]
	insertIntoTree(n, method, parts, 0, handler)
}

// recursive function that inserts node into the tree
func insertIntoTree(n *node, method string, parts []string, index int, handler HandlerFunc) {
	// check if current node has a child with same name like part of url path
	// when it is end point, add handler to node
	if n.Children[parts[index]] != nil && index == len(parts) - 1 { 
		n.Children[parts[index]].Handlers[method] = handler
		return
	} else if n.Children[parts[index]] != nil { // when it isn't end point
		insertIntoTree(n.Children[parts[index]], method, parts, index + 1, handler)
		return
	}
	// if node is not found, and it isn't end point
	// just create empty node and process the next point
	if index < len(parts) - 1 {
		n.Children[parts[index]] = &node{
			Handlers: make(map[string]HandlerFunc),
			Children: make(map[string]*node),
		}
		insertIntoTree(n.Children[parts[index]], method, parts, index + 1, handler)
		return
	} else { // if it is end point, create node and add handler
		n.Children[parts[index]] = &node{
			Handlers: make(map[string]HandlerFunc),
			Children: make(map[string]*node),
		}
		n.Children[parts[index]].Handlers[method] = handler
	}
} 

func (n *node) search(url string) *node {
	parts := strings.Split(url, "/")[1:]
	return searchInTree(n, parts, 0)
}

func searchInTree(n *node, parts []string, index int) *node {
	// returns node when function reachs end point
	if index == len(parts) {
		return n
	}
	var next *node
	// searching for * parameter
	for k, v := range n.Children {
		if k == "*" {
			next = v
		}
	}
	// when part has no matches with existing nodes
	// if * parameter exists, return * node, in other cases return nil
	if n.Children[parts[index]] == nil {
		if next != nil {
			return searchInTree(next, parts, index + 1)
		}
		return nil
	}
	return searchInTree(n.Children[parts[index]], parts, index + 1)
}
