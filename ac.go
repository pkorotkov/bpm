package bmp

type treeNode struct {
	parent  *treeNode
	char    byte
	failure *treeNode
	// See https://github.com/golang/go/issues/3512
	results     map[string][]byte
	transitions map[byte]*treeNode
}

func newTreeNode(p *treeNode, c byte) *treeNode {
	return &treeNode{
		parent:      p,
		char:        c,
		results:     make(map[string][]byte),
		transitions: make(map[byte]*treeNode),
	}
}

// Adds pattern ending in this node.
func (tn *treeNode) addResult(r []byte) {
	if _, ok := tn.results[string(r)]; ok {
		return
	}
	tn.results[string(r)] = r
}

// adds trabsition node.
func (tn *treeNode) addTransition(n *treeNode) {
	tn.transitions[n.char] = n
}

type _ACSearchEngine struct {
	bfr  *bufferedFileReader
	root *treeNode
}

func (ac *_ACSearchEngine) SetFile(fp string) (err error) {
	var bfr *bufferedFileReader
	if bfr, err = NewBufferedFileReader(fp); err != nil {
		return
	}
	ac.bfr = bfr
	return
}

// builds tree from specified patterns.
func (ac *_ACSearchEngine) buildTree(patterns [][]byte) {
	// Build pattern tree and transition function.
	ac.root = newTreeNode(nil, 0)
	for _, p := range patterns {
		n := ac.root
		for _, pc := range p {
			var nn *treeNode
			for c, tn := range n.transitions {
				if c == pc {
					nn = tn
					break
				}
			}
			if nn == nil {
				nn = newTreeNode(n, pc)
				n.addTransition(nn)
			}
			n = nn
		}
		n.addResult(p)
	}
	// Find failure functions.
	var nodes []*treeNode
	for _, tn := range ac.root.transitions {
		tn.failure = ac.root
		for _, tt := range tn.transitions {
			nodes = append(nodes, tt)
		}
	}
	// Deal with other nodes using BFS.
	for len(nodes) != 0 {
		var nnodes []*treeNode
		for _, n := range nodes {
			r, c := n.parent.failure, n.char
			for r != nil && r.transitions[c] == nil {
				r = r.failure
			}
			if r == nil {
				n.failure = ac.root
			} else {
				n.failure = r.transitions[c]
				for _, res := range n.failure.results {
					n.addResult(res)
				}
			}
			for _, tn := range n.transitions {
				nnodes = append(nnodes, tn)
			}
		}
		nodes = nnodes
	}
	ac.root.failure = ac.root
}

func (ac *_ACSearchEngine) FindAllOccurrences(patterns [][]byte) (srs SearchResults, err error) {
	var (
		index int64
		dl    int64 = ac.bfr.FileSize()
	)
	srs = newSearchResults()
	ac.buildTree(patterns)
	rt := ac.root
	for index < dl {
		var tn *treeNode
		for tn == nil {
			tn = rt.transitions[ac.bfr.ReadByteAt(index)]
			if rt == ac.root {
				break
			}
			if tn == nil {
				rt = rt.failure
			}
		}
		if tn != nil {
			rt = tn
		}
		for _, pattern := range rt.results {
			srs.putOne(pattern, index-int64(len(pattern))+1)
		}
		index++
	}
	return
}
