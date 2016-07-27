package model

type Node struct {
	ID       int     `json:"-"`
	ParentID int     `json:"-"`
	Text     string  `json:"text"`
	Icon     string  `json:"icon"`
	Path     string  `json:"href"`
	Note     string  `json:"note"`
	Child    []*Node `json:"nodes,omitempty"`
}

func (n *Node) Size() int {
	var size int = len(n.Child)
	for _, c := range n.Child {
		size += c.Size()
	}
	return size
}

func (n *Node) Add(nodes ...*Node) bool {
	var size = n.Size()
	for _, n := range nodes {
		if n.ParentID == n.ID {
			n.Child = append(n.Child, n)
		} else {
			for _, c := range n.Child {
				if c.Add(n) {
					break
				}
			}
		}
	}
	return n.Size() == size+len(nodes)
}
