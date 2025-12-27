package backend

type FileSystemNode struct {
	Name        string
	IsFolder    bool
	Size        int64
	
	// MULTI LINKED LIST (First Child - Next Sibling)
	FirstChild  *FileSystemNode 
	NextSibling *FileSystemNode 
}

func NewNode(name string, isFolder bool, size int64) *FileSystemNode {
	return &FileSystemNode{
		Name:        name,
		IsFolder:    isFolder,
		Size:        size,
		FirstChild:  nil,
		NextSibling: nil,
	}
}

// AddChild
// Menambahkan anak ke ujung antrean Sibling
func (n *FileSystemNode) AddChild(newChild *FileSystemNode) {
	if !n.IsFolder {
		return
	}

	if n.FirstChild == nil {
		// Jika belum punya anak, jadikan anak pertama
		n.FirstChild = newChild
	} else {
		// Jika sudah punya, cari saudara terakhir (traverse sibling)
		current := n.FirstChild
		for current.NextSibling != nil {
			current = current.NextSibling
		}
		// Sambungkan di ujung
		current.NextSibling = newChild
	}
}