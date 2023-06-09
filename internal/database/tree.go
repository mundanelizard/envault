package database

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"sort"
)

const (
	ModeDir = "40000"
	//MODE = "100644" // 100755 for executables
)

type Tree struct {
	name    string
	id      string
	entries map[string]Enterable
}

type Enterable interface {
	Mode() string
	Name() string
	Id() string
	String() string
}

func NewTree(name string) *Tree {
	return &Tree{
		name:    name,
		entries: make(map[string]Enterable),
	}
}

func BuildTree(entries []Enterable) *Tree {
	root := NewTree("root")

	for _, entry := range entries {
		parents := getParentPaths(entry.Name())
		root.addEntry(parents, entry)
	}

	return root
}

func (t *Tree) Traverse(fn func(*Tree)) {
	for _, entry := range t.entries {
		if tree, ok := entry.(*Tree); ok {
			tree.Traverse(fn)
		}
	}
	fn(t)
}

func (t *Tree) addEntry(parents []string, entry Enterable) {
	if len(parents) == 0 {
		base := filepath.Base(entry.Name())
		t.entries[base] = entry
		return
	}

	parent := filepath.Base(parents[0])
	tree, ok := t.entries[parent].(*Tree)
	if !ok {
		tree = NewTree(parent)
		t.entries[parent] = tree
	}

	tree.addEntry(parents[1:], entry)
}

func getParentPaths(path string) []string {
	path = filepath.Clean(path)
	var parents []string

	parent := filepath.Dir(path)
	for parent != "." {
		parents = append(parents, parent)
		parent = filepath.Dir(parent)
	}

	sort.Strings(parents)

	return parents
}

// implementing enterable

func (t *Tree) Name() string {
	return t.name
}

func (t *Tree) Mode() string {
	return ModeDir
}

// implementing Enterable and database.Storable

func (t *Tree) Id() string {
	return t.id
}

// implementing methods for database.Storable

func (t *Tree) String() string {
	var buf bytes.Buffer

	keys := getSortedKeys(t.entries)

	for _, name := range keys {
		e := t.entries[name]
		_, err := fmt.Fprintf(&buf, "%s %s %s\n", e.Mode(), name, e.Id())
		if err != nil {
			panic(err)
		}
	}

	return buf.String()
}

func (t *Tree) SetId(id string) {
	t.id = id
}

func (t *Tree) Type() string {
	return "tree"
}

func getSortedKeys(m map[string]Enterable) []string {
	keys := make([]string, 0)

	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func hexDecode(h string) []byte {
	data, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	return data
}
