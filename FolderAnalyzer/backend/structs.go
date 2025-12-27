package backend

type FileSystemNode struct {
	Name        string
	IsFolder    bool
	Size        int64
	
	// MULTI LINKED LIST POINTERS:
	FirstChild  *FileSystemNode // Pointer ke anak paling kiri (isi folder)
	NextSibling *FileSystemNode // Pointer ke saudara di sebelah kanan
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

// AddChild Linked List
// Menyambungkan node baru ke ujung rantai sibling
func (n *FileSystemNode) AddChild(newChild *FileSystemNode) {
	if !n.IsFolder {
		return
	}

	// Jika belum punya anak, jadikan anak pertama
	if n.FirstChild == nil {
		n.FirstChild = newChild
	} else {
		// Jika sudah punya anak, telusuri (traverse) sampai anak terakhir (sibling paling ujung)
		current := n.FirstChild
		for current.NextSibling != nil {
			current = current.NextSibling
		}
		// Sambungkan di ujung
		current.NextSibling = newChild
	}
}