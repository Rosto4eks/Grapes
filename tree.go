package grapes

import (
	"strings"
)

type node struct {
	Handlers map[string] HandlerFunc // where key is method
	Children map[string]*node // where key is part of URL
}


func (n *node) insert(method string, url string, handler HandlerFunc) {
	parts := getArrPath(url)
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

func (n *node) search(url string) (*node, string) {
	parts := getArrPath(url)
	var treePath string
	return searchInTree(n, parts, &treePath, 0), treePath
}

func searchInTree(n *node, parts []string, treePath *string, index int) *node {
	// returns node when function reachs end point
	if index == len(parts) {
		return n
	}
	var nodeWithParam *node
	var key string
	// searching for * parameter
	for k, v := range n.Children {
		if k == "*" || k[0] == ':' {
			nodeWithParam = v
			key = k
		}
	}
	// when part has no matches with existing nodes
	// if * parameter exists, return * node, in other cases return nil
	if n.Children[parts[index]] == nil {
		if nodeWithParam != nil {
			*treePath += "/" + key
			return searchInTree(nodeWithParam, parts, treePath, index + 1)
		}
		return nil
	}
	*treePath += "/" + parts[index]
	return searchInTree(n.Children[parts[index]], parts, treePath, index + 1)
}

func getArrPath(path string) []string {
	return strings.Split(path, "/")[1:]
}