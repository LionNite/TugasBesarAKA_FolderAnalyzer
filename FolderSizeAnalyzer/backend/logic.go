package backend

import (
	"math/rand"
	"strconv"
	"time"
)

// GenerateDummyStructure membuat struktur folder acak
func GenerateDummyStructure(totalItems int) *FileSystemNode {
	rand.Seed(time.Now().UnixNano())
	root := NewNode("Root", true, 0)
	availableFolders := []*FileSystemNode{root}

	for i := 0; i < totalItems; i++ {
		// Pilih folder induk acak
		parentIdx := rand.Intn(len(availableFolders))
		parent := availableFolders[parentIdx]

		// 70% File, 30% Folder
		if rand.Float32() > 0.3 {
			// Buat File (Size 1KB - 100KB)
			size := int64(rand.Intn(100000) + 1024)
			file := NewNode("File_"+strconv.Itoa(i), false, size)
			parent.AddChild(file)
		} else {
			// Buat Folder
			folder := NewNode("Folder_"+strconv.Itoa(i), true, 0)
			parent.AddChild(folder)
			availableFolders = append(availableFolders, folder)
		}
	}
	return root
}

// HitungRekursif: Menghitung ukuran menggunakan fungsi rekursif
func HitungRekursif(node *FileSystemNode) int64 {
	if !node.IsFolder {
		return node.Size
	}

	var total int64 = 0
	for _, child := range node.Children {
		total += HitungRekursif(child)
	}
	return total
}

// HitungIteratif: Menghitung ukuran menggunakan Stack (Manual)
func HitungIteratif(root *FileSystemNode) int64 {
	if !root.IsFolder {
		return root.Size
	}

	var totalSize int64 = 0
	// Stack inisialisasi
	stack := []*FileSystemNode{root}

	for len(stack) > 0 {
		// Pop item terakhir (LIFO - Last In First Out)
		index := len(stack) - 1
		current := stack[index]
		stack = stack[:index] // Hapus dari stack

		if !current.IsFolder {
			totalSize += current.Size
		} else {
			// Push semua anak ke stack
			stack = append(stack, current.Children...)
		}
	}
	return totalSize
}