package core

import (
	"errors"
	"strings"

	ctx "saiang/framework/context"
)

type Tree struct {
	root *node
}

type node struct {
	isLast  bool            // Indicates whether it is the last one
	segment string          // A string in url that represents a string of segment in a route represented by this node
	handler ctx.HandlerFunc // Representes that the handler function contained in this node for the final local call
	childs  []*node         // Represents that all children node under this node
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
	}
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{root: root}
}

// isWildSegment is a function to determine whether a segment is generic, start with ':'.
func isWildSegment(seg string) bool {
	return strings.HasPrefix(seg, ":")
}

// filterChildNodes define a function that filters the next layer of child nodes that meet the segment rule.
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}

	// if a segment is generic and start with ':', all the next nodes will meet the requirements.
	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	// filter all the next layer of child nodes
	for _, node := range n.childs {
		if isWildSegment(node.segment) {
			nodes = append(nodes, node)
		} else if node.segment == segment {
			// If the next layer of child nodes does not have wildcards, but the text matches exactly,
			// the requirements are met.
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (n *node) matchNode(url string) *node {
	segments := strings.SplitN(url, "/", 2)
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}

	nodes := n.filterChildNodes(segment)
	if nodes == nil || len(nodes) == 0 {
		return nil
	}

	if len(segments) == 1 {
		for _, tn := range nodes {
			if tn.isLast {
				return tn
			}
		}
		return nil
	}

	for _, tn := range nodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}
	return nil
}

func (tr *Tree) AddRouter(url string, handler ctx.HandlerFunc) error {
	n := tr.root
	if n.matchNode(url) != nil {
		return errors.New("router exist: " + url)
	}

	segments := strings.Split(url, "/")
	for index, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := index == len(segment)-1
		var objNode *node
		childNodes := n.filterChildNodes(segment)
		if len(childNodes) > 0 {
			for _, node := range childNodes {
				if node.segment == segment {
					objNode = node
					break
				}
			}
		}

		if objNode == nil {
			node := newNode()
			node.segment = segment
			if isLast {
				node.isLast = true
				node.handler = handler
			}
			n.childs = append(n.childs, node)
			objNode = node
		}
		n = objNode
	}
	return nil
}

func (tr *Tree) FindHandler(url string) ctx.HandlerFunc {
	matchNode := tr.root.matchNode(url)
	if matchNode == nil {
		return nil
	}
	return matchNode.handler
}
