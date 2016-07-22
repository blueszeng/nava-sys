package models

import ()

type Node struct {
	ID       int     `json:"-"`
	ParentID int     `json:"-"`
	Text     string  `json:"text"`
	Icon     string  `json:"icon"`
	Path     string  `json:"href"`
	Note string `json:"note"`
	Child    []*Node `json:"nodes,omitempty"`
}

func (this *Node) Size() int {
	var size int = len(this.Child)
	for _, c := range this.Child {
		size += c.Size()
	}
	return size
}

func (this *Node) Add(nodes ...*Node) bool {
	var size = this.Size()
	for _, n := range nodes {
		if n.ParentID == this.ID {
			this.Child = append(this.Child, n)
		} else {
			for _, c := range this.Child {
				if c.Add(n) {
					break
				}
			}
		}
	}
	return this.Size() == size+len(nodes)
}
