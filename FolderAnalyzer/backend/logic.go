package backend

import (
	"math/rand"
	"strconv"
	"time"
)

// GenerateDummyStructure: Logika sama, tapi fungsi AddChild di struct sudah berubah
func GenerateDummyStructure(totalItems int) (*FileSystemNode, int, int) {
	rand.Seed(time.Now().UnixNano())
	root := NewNode("Root", true, 0)
	
	// Kita butuh cara untuk menyimpan daftar folder agar bisa dipilih acak.
	// Karena Linked List susah diakses acak (harus traverse), 
	// kita BANTU dengan slice sementara hanya untuk proses pembuatan dummy ini.
	availableFolders := []*FileSystemNode{root}

	countFiles := 0
	countFolders := 1

	for i := 0; i < totalItems; i++ {
		parentIdx := rand.Intn(len(availableFolders))
		parent := availableFolders[parentIdx]

		if rand.Float32() > 0.3 {
			// Buat File
			size := int64(rand.Intn(100000) + 1024)
			file := NewNode("File_"+strconv.Itoa(i), false, size)
			parent.AddChild(file)
			countFiles++
		} else {
			// Buat Folder
			folder := NewNode("Folder_"+strconv.Itoa(i), true, 0)
			parent.AddChild(folder)
			
			availableFolders = append(availableFolders, folder)
			countFolders++
		}
	}
	return root, countFiles, countFolders
}

// --- HITUNG REKURSIF (Versi Multi Linked List) ---
func HitungRekursif(node *FileSystemNode) int64 {
	if !node.IsFolder {
		return node.Size
	}

	var total int64 = 0
	
	// Loop menelusuri Linked List Sibling
	// Mulai dari anak pertama, geser ke kanan terus sampai nil
	child := node.FirstChild
	for child != nil {
		total += HitungRekursif(child)
		child = child.NextSibling // Geser ke saudara berikutnya
	}
	
	return total
}

// --- HITUNG ITERATIF (Versi Multi Linked List) ---
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

		// Proses node saat ini
		if !current.IsFolder {
			totalSize += current.Size
		}

		// LOGIKA STACK UNTUK LINKED LIST:
		// Jika punya anak, masukkan anak pertama ke Stack.
		// Nanti anak pertama akan menarik saudara-saudaranya (NextSibling) 
		// ATAU kita masukkan semua sibling ke stack sekarang juga?
		
		// Cara Paling Aman: Masukkan anak pertama saja ke Stack?
		// TAPI tunggu, jika kita hanya push FirstChild, bagaimana NextSibling diproses?
		
		// Strategi Iteratif yang benar untuk First Child/Next Sibling:
		// 1. Proses Current
		// 2. Jika punya NextSibling, Push ke Stack (agar diproses nanti)
		// 3. Jika punya FirstChild, Push ke Stack (agar diproses duluan/DFS)
		
		if current.FirstChild != nil {
			stack = append(stack, current.FirstChild)
		}
		
		// PERHATIAN:
		// Struktur data kita mencampur File dan Folder sebagai Node.
		// Jika 'current' adalah Folder, dia punya FirstChild.
		// TAPI 'current' sendiri mungkin adalah FirstChild dari bapaknya, 
		// jadi dia punya NextSibling.
		
		// KOREKSI LOGIKA ITERATIF:
		// Di looping awal (stack root), root tidak punya sibling.
		// Tapi anak-anaknya punya sibling.
		
		// Mari kita ubah strateginya:
		// Stack berisi Node yang perlu "dijelajahi".
		// Saat kita pop Folder, kita iterasi semua children-nya MANUAL menggunakan linked list,
		// lalu push ke stack.
		
		if current.IsFolder {
			child := current.FirstChild
			for child != nil {
				// Push semua anak ke stack
				stack = append(stack, child)
				child = child.NextSibling
			}
		}
	}
	return totalSize
}