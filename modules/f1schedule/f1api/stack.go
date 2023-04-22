package f1api

import (
	"errors"

	"golang.org/x/net/html"
)

type NodeStack struct {
	nodes []*html.Node
}

func (s *NodeStack) Push(node *html.Node) {
	s.nodes = append(s.nodes, node)
}

func (s *NodeStack) Len() int {
	return len(s.nodes)
}

func (s *NodeStack) IsEmpty() bool {
	return len(s.nodes) < 1
}

func (s *NodeStack) Pop() (*html.Node, error) {
	if len(s.nodes) < 1 {
		return nil, errors.New("nodestack_empty")
	}
	c := s.nodes[0]
	s.nodes = s.nodes[1:]
	return c, nil
}
