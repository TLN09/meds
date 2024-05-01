package seedTree

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"golang.org/x/crypto/sha3"
)

type SeedTreeNode struct {
	seed   []byte
	salt   []byte
	i, j   int
	parent *SeedTreeNode
	left   *SeedTreeNode
	right  *SeedTreeNode
}

func New(seed, salt []byte, i, j int, parent, left, right *SeedTreeNode) *SeedTreeNode {
	return &SeedTreeNode{seed, salt, i, j, parent, left, right}
}

func (node *SeedTreeNode) Seed() []byte {
	return node.seed
}

func (node *SeedTreeNode) CreateSeedTree(maxHeight int, leafs *[]*SeedTreeNode) error {
	// fmt.Printf("(%v, %v)\n", node.i, node.j)
	if node.j >= len(*leafs) {
		// fmt.Printf("(%v, %v): too far down\n", node.i, node.j)
		return nil
	}

	if maxHeight == 0 {
		// fmt.Printf("(%v, %v): inserting this into leafs\n", node.i, node.j)
		(*leafs)[node.j] = node
		return nil
	}

	err := node.createChildren()
	if err != nil {
		return err
	}

	// fmt.Printf("(%v, %v): left.CreateSeedTree\n", node.i, node.j)
	err = node.left.CreateSeedTree(maxHeight-1, leafs)
	// fmt.Printf("(%v, %v): left.CreateSeedTree created\n", node.i, node.j)
	if err != nil {
		return err
	}

	if (*leafs)[len(*leafs)-1] == nil {
		// fmt.Printf("(%v, %v): right.CreateSeedTree\n", node.i, node.j)
		err = node.right.CreateSeedTree(maxHeight-1, leafs)
		// fmt.Printf("(%v, %v): right.CreateSeedTree created\n", node.i, node.j)
	} else {
		node.right = nil
	}

	return err
}

// Creates the children for this node
// Precondition: Node does not have any children currently
func (node *SeedTreeNode) createChildren() error {
	G := sha3.NewShake256()
	G.Write(node.salt)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(node.i*node.i-1+node.j))
	if err != nil {
		return err
	}
	G.Write(buf.Bytes())
	G.Write(node.seed)
	seedLength := len(node.seed)
	node.left = New(make([]byte, seedLength), node.salt, node.i+1, 2*node.j, node, nil, nil)
	node.right = New(make([]byte, seedLength), node.salt, node.i+1, 2*node.j+1, node, nil, nil)
	G.Read(node.left.seed)
	G.Read(node.right.seed)

	return nil
}

func (node *SeedTreeNode) SeedTreeToPath(path *[]byte, idx *int) {
	if node.HasLabel() {
		seed := node.Seed()
		for i := 0; i < len(seed); i++ {
			(*path)[i+(*idx)] = seed[i]
		}
		(*idx) += len(seed)
	} else {
		if node.left != nil {
			node.left.SeedTreeToPath(path, idx)
		}
		if node.right != nil {
			node.right.SeedTreeToPath(path, idx)
		}
	}
}

func (node *SeedTreeNode) RemoveSeedLabel() {
	node.i = -1
	node.j = -1
	if node.parent != nil {
		node.parent.RemoveSeedLabel()
	}
}

func (node *SeedTreeNode) HasLabel() bool {
	return node.i >= 0 || node.j >= 0
}

func EmptyTree(h int, salt []byte, leafs *[]*SeedTreeNode) *SeedTreeNode {
	root := New([]byte{}, salt, 0, 0, nil, nil, nil)
	root.createEmptyTree(h, leafs)
	return root
}

func (node *SeedTreeNode) createEmptyTree(maxHeight int, leafs *[]*SeedTreeNode) {
	if node.j >= len(*leafs) {
		return
	}

	if maxHeight == 0 {
		(*leafs)[node.j] = node
		return
	}

	node.left = New([]byte{}, node.salt, node.i+1, 2*node.j, node, nil, nil)
	node.left.createEmptyTree(maxHeight-1, leafs)

	if (*leafs)[len(*leafs)-1] == nil {
		node.right = New([]byte{}, node.salt, node.i+1, 2*node.j+1, node, nil, nil)
		node.right.CreateSeedTree(maxHeight-1, leafs)
	} else {
		node.right = nil
	}
}

func (node *SeedTreeNode) isLeaf() bool {
	return node.left == nil && node.right == nil
}

func (node *SeedTreeNode) pathToSeedTreeComputeSeeds() error {
	G := sha3.NewShake256()
	G.Write(node.salt)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(node.i*node.i-1+node.j))
	if err != nil {
		return err
	}
	G.Write(buf.Bytes())
	G.Write(node.seed)
	if node.left != nil {
		node.left.seed = make([]byte, len(node.seed))
		G.Read(node.left.seed)
		err = node.left.pathToSeedTreeComputeSeeds()
		if err != nil {
			return err
		}
	}

	if node.right != nil {
		node.right.seed = make([]byte, len(node.seed))
		G.Read(node.right.seed)
		err = node.right.pathToSeedTreeComputeSeeds()
	}
	return err

}

func (node *SeedTreeNode) addLeafSeeds(seeds *[][]byte, idx *int, l_tree_seed int) {
	if node.isLeaf() {
		// fmt.Printf("(%v, %v): %v\n", node.i, node.j, node.seed)
		// fmt.Printf("idx: %3v\n", *idx)
		if len(node.seed) != 0 {
			for i := 0; i < l_tree_seed; i++ {
				(*seeds)[(*idx)][i] = node.seed[i]
			}
			(*idx)++
		}
	} else {
		if node.left != nil {
			node.left.addLeafSeeds(seeds, idx, l_tree_seed)
		}
		if node.right != nil {
			node.right.addLeafSeeds(seeds, idx, l_tree_seed)
		}
	}
}

func (node *SeedTreeNode) PathToSeedTree(path *[]byte, l_tree_seed int, seeds *[][]byte, seeds_idx, path_idx *int) error {
	var err error = nil
	if node.HasLabel() {
		node.seed = (*path)[(*path_idx) : (*path_idx)+l_tree_seed]
		(*path_idx) += l_tree_seed
		err = node.pathToSeedTreeComputeSeeds()
		if err != nil {
			return err
		}

		node.addLeafSeeds(seeds, seeds_idx, l_tree_seed)
	} else {
		if node.isLeaf() {
			(*seeds_idx)++
		} else {
			if node.left != nil {
				err = node.left.PathToSeedTree(path, l_tree_seed, seeds, seeds_idx, path_idx)
			}
			if err != nil {
				return err
			}
			if node.right != nil {
				err = node.right.PathToSeedTree(path, l_tree_seed, seeds, seeds_idx, path_idx)
			}
		}
	}
	return err
}

func (node *SeedTreeNode) String() string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			fmt.Sprintf(
				"{\"name\": \"(%v, %v)\", \"seed\": \"%v\", \"left\": %v, \"right\": %v}",
				node.i,
				node.j,
				node.seed,
				node.left,
				node.right,
			),
			"<",
			"\"",
		),
		">",
		"\"",
	)
}
