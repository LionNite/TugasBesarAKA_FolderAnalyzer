package backend

import (
	"math/rand"
	"strconv"
	"time"
)

// GenerateDummyStructure: Versi Random Tree (Natural)
// Mengembalikan 3 nilai: Root Node, Jumlah File, Jumlah Folder
func GenerateDummyStructure(totalItems int) (*FileSystemNode, int, int) {
	rand.Seed(time.Now().UnixNano())
	
	// Root Node
	root := NewNode("Root", true, 0)
	
	// Slice untuk menampung folder yang tersedia sebagai induk
	availableFolders := []*FileSystemNode{root}

	countFiles := 0
	countFolders := 1 // Root dihitung 1 folder

	for i := 0; i < totalItems; i++ {
		// Pilih induk secara acak (Membuat struktur pohon yang lebar/alami)
		parentIdx := rand.Intn(len(availableFolders))
		parent := availableFolders[parentIdx]

		// Logika Probabilitas: 70% File, 30% Folder
		if rand.Float32() > 0.3 {
			// --- BUAT FILE (70%) ---
			// Ukuran acak antara 1KB - 100KB
			size := int64(rand.Intn(100000) + 1024)
			name := "File_" + strconv.Itoa(i)
			file := NewNode(name, false, size)
			
			parent.AddChild(file)
			countFiles++
		} else {
			// --- BUAT FOLDER (30%) ---
			name := "Folder_" + strconv.Itoa(i)
			folder := NewNode(name, true, 0)
			
			parent.AddChild(folder)
			
			// Masukkan folder baru ke daftar 'availableFolders' agar bisa dipilih jadi induk nanti
			availableFolders = append(availableFolders, folder)
			countFolders++
		}
	}

	return root, countFiles, countFolders
}

// HitungRekursif: Menghitung total size menggunakan metode Rekursif
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

// HitungIteratif: Menghitung total size menggunakan metode Iteratif
func HitungIteratif(root *FileSystemNode) int64 {
	if !root.IsFolder {
		return root.Size
	}

	var totalSize int64 = 0
	stack := []*FileSystemNode{root}

	for len(stack) > 0 {
		index := len(stack) - 1
		current := stack[index]
		stack = stack[:index] 

		if !current.IsFolder {
			totalSize += current.Size
		} else {
			stack = append(stack, current.Children...)
		}
	}
	return totalSize
}