package xmlx

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Node is a generic XML node.
type Node struct {

	// The name of the node
	Name string

	// The attribute list of the node
	Attrs map[string]string

	// The data located within the node
	Data string

	// The subnodes within the node
	Nodes []Node

	prefix string

	IsInvalid bool
}

// UnmarshalXML takes the content of an XML node and puts it into the Node structure.
func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	if len(start.Attr) != 0 {
		n.Attrs = map[string]string{}
		for _, v := range start.Attr {
			n.Attrs[v.Name.Local] = v.Value
		}
	}
	n.Name = start.Name.Local

	balance := 1

	for balance != 0 {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch t := token.(type) {

		case xml.CharData:
			v := bytes.TrimSpace(t)
			if len(v) == 0 {
				continue
			}
			n.Data = string(t)

		case xml.StartElement:
			if t.Name.Local == start.Name.Local {
				balance++
				continue
			}

			node := Node{}
			err := node.UnmarshalXML(d, t)
			if err != nil {
				return err
			}

			n.Nodes = append(n.Nodes, node)

		case xml.EndElement:
			if t.Name.Local == start.Name.Local {
				balance--
			}
		}
	}

	return nil
}

func GetXMLNode(xmlA string, labelsA ...string) (*Node, error) {
	var nodeT Node

	errT := xml.Unmarshal([]byte(xmlA), &nodeT)

	if errT != nil {
		return nil, errT
	}

	if labelsA == nil || len(labelsA) == 0 {
		return &nodeT, nil
	}

	if len(labelsA) == 1 {
		if labelsA[0] != nodeT.Name {
			return nil, nil
		}

		return &nodeT, nil
	}

	if labelsA[0] != nodeT.Name {
		return nil, nil
	}

	return nodeT.GetSubNode(labelsA[1:]...), nil

}

func (n Node) GetSubNode(labelsA ...string) *Node {
	if labelsA == nil || len(labelsA) == 0 {
		return &n
	}

	currentNodeT := n

	for _, term := range labelsA {
		// fmt.Printf("currentNodeT: %v, currentNodeT.name: %v\n", currentNodeT, currentNodeT.Name)
		if term == "" {
			continue
		}

		foundT := false
		for _, node := range currentNodeT.Nodes {
			if node.Name == term {
				currentNodeT = node
				foundT = true
				break
			}
		}

		if !foundT {
			return nil
		}
	}

	return &currentNodeT

}

func (n Node) GetSubNodeX(labelsA string) *Node {
	labelsT := strings.Split(labelsA, "/")

	if labelsT == nil || len(labelsT) == 0 {
		return &n
	}

	currentNodeT := n

	for _, term := range labelsT {
		// fmt.Printf("currentNodeT: %v, currentNodeT.name: %v\n", currentNodeT, currentNodeT.Name)
		if term == "" {
			continue
		}

		foundT := false
		for _, node := range currentNodeT.Nodes {
			if node.Name == term {
				currentNodeT = node
				foundT = true
				break
			}
		}

		if !foundT {
			return nil
		}
	}

	return &currentNodeT

}

// GetSubNodeString default/error will return empty
func (n Node) GetSubNodeString(labelsA ...string) string {
	if labelsA == nil || len(labelsA) == 0 {
		return ""
	}

	currentNodeT := n

	for _, term := range labelsA {
		// fmt.Printf("currentNodeT: %v, currentNodeT.name: %v\n", currentNodeT, currentNodeT.Name)
		if term == "" {
			continue
		}

		foundT := false
		for _, node := range currentNodeT.Nodes {
			if node.Name == term {
				currentNodeT = node
				foundT = true
				break
			}
		}

		if !foundT {
			return ""
		}
	}

	return currentNodeT.Data

}

func (n Node) GetSubNodeStringX(labelsA string) string {
	labelsT := strings.Split(labelsA, "/")

	if labelsT == nil || len(labelsT) == 0 {
		return ""
	}

	currentNodeT := n

	for _, term := range labelsT {
		// fmt.Printf("currentNodeT: %v, currentNodeT.name: %v\n", currentNodeT, currentNodeT.Name)

		foundT := false
		for _, node := range currentNodeT.Nodes {
			if node.Name == term {
				currentNodeT = node
				foundT = true
				break
			}
		}

		if !foundT {
			return ""
		}
	}

	return currentNodeT.Data

}

func (n Node) SubNodes(labelA string) []Node {
	if labelA == "" {
		return n.Nodes
	}

	bufT := make([]Node, 0, len(n.Nodes))

	for _, node := range n.Nodes {
		if node.Name == labelA {
			bufT = append(bufT, node)
		}
	}

	return bufT
}

// GetSubNodeBy get a sub-node value of a node, while one other sub-node has the proper value
func (n Node) GetSubNodeBy(labelA string, subLabelA string, valueA string) *Node {

	nodesT := n.SubNodes(labelA)

	for _, node := range nodesT {
		subNodeT := node.GetSubNode(subLabelA)
		if subNodeT == nil {
			continue
		}

		if subNodeT.Data == valueA {
			return &node
		}
	}

	return nil
}

// GetSubNodeStringBy if 1 sub-node has proper value, return the other sub-nodes' value
func (n Node) GetSubNodeStringBy(labelA string, subLabel1A string, value1A string, subLabel2A string) string {

	nodesT := n.SubNodes(labelA)

	for _, node := range nodesT {
		subNode1T := node.GetSubNode(subLabel1A)
		if subNode1T == nil {
			continue
		}

		if subNode1T.Data != value1A {
			continue
		}

		subNode2T := node.GetSubNode(subLabel2A)
		if subNode2T == nil {
			continue
		}

		return subNode2T.Data

	}

	return ""
}

// GetSubNodeBy2 get a sub-node value of a node, while two other sub-nodes have the proper value
func (n Node) GetSubNodeBy2(labelA string, subLabel1A string, value1A string, subLabel2A string, value2A string) *Node {

	nodesT := n.SubNodes(labelA)

	for _, node := range nodesT {
		subNode1T := node.GetSubNode(subLabel1A)
		if subNode1T == nil {
			continue
		}

		if subNode1T.Data != value1A {
			continue
		}

		subNode2T := node.GetSubNode(subLabel2A)
		if subNode2T == nil {
			continue
		}

		if subNode2T.Data != value2A {
			continue
		}

		return &node
	}

	return nil
}

// GetSubNodeStringBy2 if 2 sub-nodes have proper values, return the other sub-nodes' value
func (n Node) GetSubNodeStringBy2(labelA string, subLabel1A string, value1A string, subLabel2A string, value2A string, subLabel3A string) string {

	nodesT := n.SubNodes(labelA)

	for _, node := range nodesT {
		subNode1T := node.GetSubNode(subLabel1A)
		if subNode1T == nil {
			continue
		}

		if subNode1T.Data != value1A {
			continue
		}

		subNode2T := node.GetSubNode(subLabel2A)
		if subNode2T == nil {
			continue
		}

		if subNode2T.Data != value2A {
			continue
		}

		subNode3T := node.GetSubNode(subLabel3A)
		if subNode3T == nil {
			continue
		}

		return subNode3T.Data
	}

	return ""
}

func (n Node) GetSubNodeByX(rootLabesA string, labelA string, labelValuePairA ...string) *Node {
	rootLabelsT := strings.Split(rootLabesA, "/")
	rootNodeT := n.GetSubNode(rootLabelsT...)

	if rootNodeT == nil {
		return nil
	}

	nodesT := rootNodeT.SubNodes(labelA)

	lenT := len(labelValuePairA) / 2

	for _, node := range nodesT {
		foundT := false
		for i := 0; i < lenT; i++ {
			subNodeT := node.GetSubNodeX(labelValuePairA[i*2])
			if subNodeT == nil {
				foundT = true
				break
			}

			if subNodeT.Data != labelValuePairA[i*2+1] {
				foundT = true
				break
			}

		}

		if foundT {
			continue
		}

		return &node

	}

	return nil
}

func (n Node) GetSubNodeStringByX(rootLabesA string, labelA string, subLabelA string, labelValuePairA ...string) string {
	rootLabelsT := strings.Split(rootLabesA, "/")
	rootNodeT := n.GetSubNode(rootLabelsT...)

	if rootNodeT == nil {
		return ""
	}

	nodesT := rootNodeT.SubNodes(labelA)

	lenT := len(labelValuePairA) / 2

	for _, node := range nodesT {
		foundT := false
		for i := 0; i < lenT; i++ {
			subNodeT := node.GetSubNodeX(labelValuePairA[i*2])
			if subNodeT == nil {
				foundT = true
				break
			}

			if subNodeT.Data != labelValuePairA[i*2+1] {
				foundT = true
				break
			}

		}

		if foundT {
			continue
		}

		subNodeT := node.GetSubNodeX(subLabelA)
		if subNodeT == nil {
			return ""
		}

		return subNodeT.Data

	}

	return ""
}

func (p Node) Text() string {
	return p.Data
}

func (p Node) IsValid() bool {
	return !p.IsInvalid
}

func (p Node) FindNodeRecursively(label string) Node {
	if len(label) == 0 {
		return p
	}

	for _, node := range p.Nodes {
		if node.Name == label {
			return node
		}

		if node.Nodes != nil && len(node.Nodes) > 0 {
			subResult := node.FindNodeRecursively(label)

			if subResult.IsInvalid == false {
				return subResult
			}
		}
	}

	return Node{IsInvalid: true}

}

func (p Node) FindNode(label string) *Node {
	if len(label) == 0 {
		return &p
	}

	for _, node := range p.Nodes {
		if node.Name == label {
			return &node
		}

		if node.Nodes != nil && len(node.Nodes) > 0 {
			subResult := node.FindNode(label)

			if subResult != nil {
				return subResult
			}
		}
	}

	return nil

}

// Split the node into many: each time the split label is encountered within a subnode of the node,
// a new node is created.
func (n Node) Split(label string) []Node {

	// Return the node itself if no label is specified: there is no split to do.
	if len(label) == 0 {
		return []Node{n}
	}

	// Create a leveled array of children.
	terms := strings.Split(label, ".")
	gen := len(terms)
	children := make([][]Node, gen+1, gen+1)
	children[0] = n.Nodes

	// Explore the leveled array. On each level, put the children of the matching node
	// to the upper level.
	for i, term := range terms {
		for _, node := range children[i] {
			if node.Name == term {
				children[i+1] = append(children[i+1], node.Nodes...)
			}
		}
	}

	// Create a node for each child. Also rename the child with label.
	var nodes []Node
	for _, child := range children[gen] {
		node := n.clone()
		var i int
		for i < len(node.Nodes) {
			if node.Nodes[i].Name == terms[0] {
				node.Nodes = append(node.Nodes[:i], node.Nodes[i+1:]...)
			}
			i++
		}
		child.Name = terms[gen-1]
		node.Nodes = append(node.Nodes, child)
		nodes = append(nodes, node)
	}

	return nodes
}

// Map returns a flatten representation of the node. If a node contains nodes
// having the same name, only the last node will exist in the map.
func (n Node) Map() map[string]string {

	out := n.flatten()

	var toProcess []Node
	for _, node := range n.Nodes {
		node.prefix = fmt.Sprintf("#nodes.%s", node.Name)
		toProcess = append(toProcess, node)
	}

	var i int
	for i < len(toProcess) {
		node := toProcess[i]
		for k, v := range node.flatten() {
			name := fmt.Sprintf("%s.%s", node.prefix, k)
			out[name] = v
		}

		for _, child := range node.Nodes {
			child.prefix = fmt.Sprintf("%s.#nodes.%s", node.prefix, child.Name)
			toProcess = append(toProcess, child)
		}
		i++
	}

	return out
}

// flatten takes a node a generates a map with it.
func (n Node) flatten() map[string]string {

	var t = map[string]string{}

	// put simple values into transcient.
	if len(n.Name) != 0 {
		t["#name"] = n.Name
	}
	if len(n.Data) != 0 {
		t["#data"] = n.Data
	}

	// put attributes into transcient.
	for k, v := range n.Attrs {
		name := fmt.Sprintf("#attr.%s", k)
		t[name] = v
	}

	return t
}

// clone creates a copy of the node and returns it.
func (n Node) clone() Node {

	node := Node{
		Name:      n.Name,
		Attrs:     n.Attrs,
		Data:      n.Data,
		prefix:    n.prefix,
		IsInvalid: n.IsInvalid,
	}
	for i := range n.Nodes {
		node.Nodes = append(node.Nodes, n.Nodes[i].clone())
	}

	return node
}
