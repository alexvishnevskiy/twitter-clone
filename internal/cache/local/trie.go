package local

type Node struct {
	children map[string]*Node
}

type Trie struct {
	root *Node
}

func NewNode() *Node {
	p := &Node{}
	p.children = make(map[string]*Node)
	return p
}

func NewTrie() *Trie {
	return &Trie{root: NewNode()}
}

func (t *Trie) Insert(first_word string, second_word string) {
	temp := t.root

	if temp.children[first_word] == nil {
		temp.children[first_word] = NewNode()
	}

	temp = temp.children[first_word]
	temp.children[second_word] = NewNode()
}

func (t *Trie) Delete(first_word string, second_word string) {
	temp := t.root

	if temp.children[first_word] != nil {
		temp = temp.children[first_word]
		temp.children[second_word] = nil
	}

	temp = t.root
	if len(temp.children[first_word].children) == 0 {
		temp.children[first_word] = nil
	}
}

func (t *Trie) StartsWith(first_word string) []string {
	var result []string
	temp := t.root

	if temp.children[first_word] != nil {
		temp = temp.children[first_word]
		for child, _ := range temp.children {
			result = append(result, child)
		}
	}
	return result
}
