package backend

// FileSystemNode merepresentasikan node dalam tree (bisa file atau folder)
// Huruf besar di awal (Exported) agar bisa dibaca oleh main.go
type FileSystemNode struct {
	Name     string
	IsFolder bool
	Size     int64
	Children []*FileSystemNode
}

// NewNode adalah factory function untuk membuat node baru
func NewNode(name string, isFolder bool, size int64) *FileSystemNode {
	return &FileSystemNode{
		Name:     name,
		IsFolder: isFolder,
		Size:     size,
		Children: []*FileSystemNode{},
	}
}

// AddChild menambahkan anak ke node folder
func (n *FileSystemNode) AddChild(child *FileSystemNode) {
	if n.IsFolder {
		n.Children = append(n.Children, child)
	}
}