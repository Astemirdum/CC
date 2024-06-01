package parser

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
)

func newGraph() *gographviz.Graph {
	g := gographviz.NewGraph()
	_ = g.SetName("Parsing")
	_ = g.SetDir(true)
	return g
}

type Node struct {
	Data     string
	Children []*Node
}

func NewNode(data string) *Node {
	return &Node{
		Data:     data,
		Children: []*Node{},
	}
}

func (node *Node) AddChild(child *Node) {
	node.Children = append(node.Children, child)
}

var vv = map[string]bool{}

func (node *Node) ToAst(g *gographviz.Graph, parentID string) string {
	nodeID := fmt.Sprintf(`"%s"`, node.Data)
	_ = g.AddNode("G", nodeID, map[string]string{"label": fmt.Sprintf("\"%s\"", node.Data)})
	if parentID != "" && !vv[nodeID] {
		vv[nodeID] = true
		_ = g.AddEdge(parentID, nodeID, true, nil)
	}
	for _, child := range node.Children {
		child.ToAst(g, nodeID)
	}
	return nodeID
}
