package model

var BoardMap map[string][]node = make(map[string][]node)

type node struct {
	User string
	Px,Py int
}

func BoardStart(group string) {
	var nodes []node
	BoardMap[group] = nodes
}

func Play(group string,user string,px int,py int) error {
	nodes := BoardMap[group]
	nodes = append(nodes, node{
		User: user,
		Px:   px,
		Py:   py,
	})
	return nil
}

func Judge() error {


	return nil
}

