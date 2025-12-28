package main

import (
	"FolderAnalyzer/backend"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type PageData struct {
	InputN       int
	SelectedMode string // Menampung pilihan user (default/random)
	Result       template.HTML
	IsProcessed  bool

	ChartLabels   string
	ChartDataRec  string
	ChartDataIter string

	TotalFiles   int
	TotalFolders int
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			InputN:       10000,
			SelectedMode: "default", // Default awal
			IsProcessed:  false,
		}

		if r.Method == "POST" {
			nStr := r.FormValue("inputN")
			mode := r.FormValue("scenario") // Ambil dari Dropdown HTML

			n, err := strconv.Atoi(nStr)
			if err == nil && n > 0 {
				data.InputN = n
				data.SelectedMode = mode
				data.IsProcessed = true
				analyzeWithCheckpoints(n, mode, &data) // Kirim mode ke fungsi analisis
			}
		}

		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, "Error: index.html missing", 500)
			return
		}
		tmpl.Execute(w, data)
	})

	fmt.Println("Dashboard Berjalan: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func analyzeWithCheckpoints(n int, mode string, data *PageData) {
	var labels []int
	initials := []int{1, 2, 4}
	for _, val := range initials {
		if val < n {
			labels = append(labels, val)
		}
	}
	current := 10
	for current < n {
		labels = append(labels, current)
		current *= 2
	}
	if len(labels) == 0 || labels[len(labels)-1] != n {
		labels = append(labels, n)
	}

	var timesRec, timesIter []float64
	var finalSizeR, finalSizeI int64
	var finalTimeR, finalTimeI float64
	var currentFiles, currentFolders int

	for i, currentN := range labels {
		// PANGGIL GENERATOR DENGAN MODE YANG DIPILIH
		var root *backend.FileSystemNode
		root, currentFiles, currentFolders = backend.GenerateDummyStructure(currentN, mode)

		// Ukur Rekursif
		start := time.Now()
		sizeR := backend.HitungRekursif(root)
		durRec := float64(time.Since(start).Nanoseconds()) / 1000.0
		timesRec = append(timesRec, durRec)

		// Ukur Iteratif
		start = time.Now()
		sizeI := backend.HitungIteratif(root)
		durIter := float64(time.Since(start).Nanoseconds()) / 1000.0
		timesIter = append(timesIter, durIter)

		if i == len(labels)-1 {
			finalSizeR = sizeR
			finalSizeI = sizeI
			finalTimeR = durRec
			finalTimeI = durIter
			data.TotalFiles = currentFiles
			data.TotalFolders = currentFolders
		}
	}

	labelsJSON, _ := json.Marshal(labels)
	recJSON, _ := json.Marshal(timesRec)
	iterJSON, _ := json.Marshal(timesIter)

	data.ChartLabels = string(labelsJSON)
	data.ChartDataRec = string(recJSON)
	data.ChartDataIter = string(iterJSON)

	// --- FORMAT OUTPUT ---
	res := fmt.Sprintf("<b>Hasil Akhir (N=%d, Mode=%s):</b><br><hr>", n, mode)

	res += fmt.Sprintf("<b>REKURSIF:</b> %.4f µs (Size: %d bytes)<br>", finalTimeR, finalSizeR)
	res += fmt.Sprintf("<b>ITERATIF:</b> %.4f µs (Size: %d bytes)<br><hr>", finalTimeI, finalSizeI)

	if finalTimeI < finalTimeR {
		res += "<b style='color:green'>Iteratif Lebih Cepat.</b>"
	} else if finalTimeI > finalTimeR {
		res += "<b style='color:green'>Rekursif Lebih Cepat.</b>"
	} else {
		res += "<b style='color:orange'>Setara.</b>"
	}
	data.Result = template.HTML(res)
}
