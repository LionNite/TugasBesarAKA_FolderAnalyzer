package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
	"FolderAnalyzer/backend"
)

// Struktur Data untuk dikirim ke HTML
type PageData struct {
	InputN      int
	Result      template.HTML
	IsProcessed bool
	
	// Data Grafik & Tabel
	ChartLabels   string
	ChartDataRec  string
	ChartDataIter string
	
	// Statistik Struktur (Baru)
	TotalFiles   int
	TotalFolders int
}

func main() {
	// 1. Setup Handler Server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Default N awal
		data := PageData{
			InputN:      10000, 
			IsProcessed: false,
		}

		// Jika tombol ditekan (POST request)
		if r.Method == "POST" {
			nStr := r.FormValue("inputN")
			n, err := strconv.Atoi(nStr)
			if err == nil && n > 0 {
				data.InputN = n
				data.IsProcessed = true
				// Jalankan analisis
				analyzeWithCheckpoints(n, &data)
			}
		}

		// 2. Baca file index.html
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, "Error: File index.html tidak ditemukan. Pastikan ada di folder yang sama.", 500)
			fmt.Println("Error loading template:", err)
			return
		}
		
		// 3. Tampilkan ke Browser
		tmpl.Execute(w, data)
	})

	// Info Terminal
	port := ":8080"
	fmt.Println("==============================================")
	fmt.Println("  Aplikasi Dashboard Berjalan!")
	fmt.Println("  Buka browser dan ketik: http://localhost" + port)
	fmt.Println("  Tekan Ctrl+C di sini untuk berhenti.")
	fmt.Println("==============================================")
	
	// Jalankan Server
	http.ListenAndServe(port, nil)
}

// Fungsi Analisis Utama
func analyzeWithCheckpoints(n int, data *PageData) {
	// --- 1. MEMBUAT TITIK GRAFIK (Labels) ---
	// Menggunakan Pola Eksponensial agar terlihat bagus di data kecil maupun besar
	var labels []int
	
	// Awal: 1, 2, 4 (Pola 2^x)
	initials := []int{1, 2, 4}
	for _, val := range initials {
		if val < n {
			labels = append(labels, val)
		}
	}

	// Lanjut: 10, 20, 40... (Doubling)
	current := 10
	for current < n {
		labels = append(labels, current)
		current *= 2 
	}

	// Akhir: Pastikan N user masuk
	if len(labels) == 0 || labels[len(labels)-1] != n {
		labels = append(labels, n)
	}
	
	// --- 2. PERSIAPAN VARIABEL ---
	var timesRec, timesIter []float64
	
	// Variabel untuk hasil akhir (Snapshot terakhir)
	var finalSizeR, finalSizeI int64
	var finalTimeR, finalTimeI float64

	// Variabel sementara untuk menangkap jumlah file/folder dari generator
	var currentFiles, currentFolders int

	// --- 3. LOOPING ANALISIS ---
	for i, currentN := range labels {
		// Generate Data (Tangkap 3 return value: root, files, folders)
		var root *backend.FileSystemNode
		root, currentFiles, currentFolders = backend.GenerateDummyStructure(currentN)

		// --- UKUR REKURSIF ---
		start := time.Now()
		sizeR := backend.HitungRekursif(root)
		durRNano := time.Since(start).Nanoseconds()
		// Konversi Nano ke Micro (Float)
		valRec := float64(durRNano) / 1000.0 
		timesRec = append(timesRec, valRec)

		// --- UKUR ITERATIF ---
		start = time.Now()
		sizeI := backend.HitungIteratif(root)
		durINano := time.Since(start).Nanoseconds()
		valIter := float64(durINano) / 1000.0
		timesIter = append(timesIter, valIter)

		// Simpan data jika ini adalah titik terakhir (N User)
		if i == len(labels)-1 {
			finalSizeR = sizeR
			finalSizeI = sizeI
			finalTimeR = valRec
			finalTimeI = valIter
			
			// Simpan statistik struktur untuk ditampilkan di panel kiri
			data.TotalFiles = currentFiles
			data.TotalFolders = currentFolders
		}
	}

	// --- 4. FORMAT DATA KE JSON (Untuk Frontend) ---
	labelsJSON, _ := json.Marshal(labels)
	recJSON, _ := json.Marshal(timesRec)
	iterJSON, _ := json.Marshal(timesIter)

	data.ChartLabels = string(labelsJSON)
	data.ChartDataRec = string(recJSON)
	data.ChartDataIter = string(iterJSON)

	// --- 5. BUAT KESIMPULAN TEKS ---
	// Menggunakan 6 digit desimal agar sangat presisi
	res := fmt.Sprintf("<b>Hasil Akhir (N = %d):</b><br><hr>", n)
	res += fmt.Sprintf("<b>REKURSIF:</b> %.6f µs (Size: %d bytes)<br>", finalTimeR, finalSizeR)
	res += fmt.Sprintf("<b>ITERATIF:</b> %.6f µs (Size: %d bytes)<br><hr>", finalTimeI, finalSizeI)
	
	if finalTimeI < finalTimeR {
		res += "<b style='color:green'>Kesimpulan: Iteratif lebih cepat.</b>"
	} else if finalTimeI > finalTimeR {
		res += "<b style='color:green'>Kesimpulan: Rekursif lebih cepat.</b>"
	} else {
		res += "<b style='color:orange'>Kesimpulan: Performa setara.</b>"
	}
	data.Result = template.HTML(res)
}