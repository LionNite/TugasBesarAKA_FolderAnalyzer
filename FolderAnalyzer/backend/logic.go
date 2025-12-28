package backend

import (
	"math/rand"
	"strconv"
	"time"
)

// GenerateDummyStructure: Multi Linked List + Mode Pilihan
func GenerateDummyStructure(totalItems int, mode string) (*FileSystemNode, int, int) {
	rand.Seed(time.Now().UnixNano())
	root := NewNode("Root", true, 0)

	countFiles := 0
	countFolders := 1 // Folder Root sudah dihitung

	if mode == "random" {
		// SKENARIO RANDOM (30% Folder : 70% File)
		// Buat slice bantu untuk memilih parent secara acak
		availableFolders := []*FileSystemNode{root}

		for i := 0; i < totalItems; i++ {
			parentIdx := rand.Intn(len(availableFolders))
			parent := availableFolders[parentIdx]

			if rand.Float32() > 0.3 { // 70% File
				size := int64(rand.Intn(100000) + 1024)
				file := NewNode("File_"+strconv.Itoa(i), false, size)
				parent.AddChild(file) // Menggunakan AddChild Linked List
				countFiles++
			} else { // 30% Folder
				folder := NewNode("Folder_"+strconv.Itoa(i), true, 0)
				parent.AddChild(folder)

				availableFolders = append(availableFolders, folder)
				countFolders++
			}
		}

	} else {
		// SKENARIO DEFAULT (1 Folder = 1 File, Deep Structure)
		// Ini akan membuat Linked List vertikal yang sangat dalam.

		current := root
		i := 0
		for i < totalItems {
			// 1. Buat Folder Baru
			folderName := "Folder_" + strconv.Itoa(countFolders)
			newFolder := NewNode(folderName, true, 0)
			current.AddChild(newFolder)
			countFolders++
			i++ // Node bertambah

			// Cek jika kuota N masih ada untuk File
			if i < totalItems {
				// 2. Buat 1 File di dalamnya
				size := int64(rand.Intn(100000) + 1024)
				file := NewNode("File_"+strconv.Itoa(countFiles), false, size)
				newFolder.AddChild(file)
				countFiles++
				i++ // Node bertambah
			}

			// 3. Masuk ke dalam folder baru
			current = newFolder
		}
	}

	return root, countFiles, countFolders
}

// HITUNG REKURSIF (Versi Multi Linked List)
func HitungRekursif(node *FileSystemNode) int64 {
	if !node.IsFolder {
		return node.Size
	}

	var total int64 = 0

	// Loop menelusuri Sibling (Kanan)
	child := node.FirstChild
	for child != nil {
		total += HitungRekursif(child)
		child = child.NextSibling
	}

	return total
}

// HITUNG ITERATIF (Versi Multi Linked List)
func HitungIteratif(root *FileSystemNode) int64 {
	if !root.IsFolder {
		return root.Size
	}

	var totalSize int64 = 0
	stack := []*FileSystemNode{root}

	for len(stack) > 0 {
		// Pop Stack
		index := len(stack) - 1
		current := stack[index]
		stack = stack[:index]

		if !current.IsFolder {
			totalSize += current.Size
		}

		// LOGIKA STACK UNTUK LINKED LIST:
		// Agar urutannya benar (Depth First), kita harus push NextSibling dulu, baru FirstChild.

		// 1. Push Saudara (NextSibling) agar diproses nanti setelah anak-anak selesai
		if current.NextSibling != nil {
			stack = append(stack, current.NextSibling)
		}

		// 2. Push Anak Pertama (FirstChild) agar diproses SEGERA (masuk ke dalam folder)
		if current.FirstChild != nil {
			stack = append(stack, current.FirstChild)
		}
	}
	return totalSize
}
