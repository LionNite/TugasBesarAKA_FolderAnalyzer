package main

import (
	"fmt"
	"strconv"
	"time"

	// Import package backend yang baru kita buat
	"TugasBesar_AKA/backend"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Tugas Besar: Folder Size Analyzer")
	myWindow.Resize(fyne.NewSize(500, 450))

	// --- UI COMPONENTS ---

	labelTitle := widget.NewLabel("Analisis Perbandingan Algoritma")
	labelTitle.TextStyle = fyne.TextStyle{Bold: true}
	labelTitle.Alignment = fyne.TextAlignCenter

	entryN := widget.NewEntry()
	entryN.SetText("10000")
	entryN.PlaceHolder = "Masukkan Jumlah N"

	labelResult := widget.NewLabel("Klik tombol untuk mulai...")
	labelResult.Wrapping = fyne.TextWrapWord

	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()

	// --- TOMBOL PROSES ---

	btnProcess := widget.NewButton("Mulai Analisis (Iteratif vs Rekursif)", func() {
		// Validasi Input
		n, err := strconv.Atoi(entryN.Text)
		if err != nil {
			labelResult.SetText("Error: Masukkan angka yang valid!")
			return
		}

		// Disable tombol & Show Progress
		progressBar.Show()
		labelResult.SetText("Sedang memproses " + entryN.Text + " data...")

		// Jalankan proses berat di Goroutine (Thread terpisah)
		go func() {
			// 1. Generate Data (Panggil dari package backend)
			startGen := time.Now()
			root := backend.GenerateDummyStructure(n)
			durationGen := time.Since(startGen)

			// 2. Hitung Rekursif (Panggil dari package backend)
			startRec := time.Now()
			sizeR := backend.HitungRekursif(root)
			durationRec := time.Since(startRec)

			// 3. Hitung Iteratif (Panggil dari package backend)
			startIter := time.Now()
			sizeI := backend.HitungIteratif(root)
			durationIter := time.Since(startIter)

			// Format Hasil Output
			resultText := fmt.Sprintf("Data Generated (%d items): %s\n", n, durationGen)
			resultText += "--------------------------------------------------\n"
			resultText += fmt.Sprintf("REKURSIF:\n   Waktu: %s\n   Total: %d bytes\n\n", durationRec, sizeR)
			resultText += fmt.Sprintf("ITERATIF:\n   Waktu: %s\n   Total: %d bytes\n", durationIter, sizeI)

			// Kesimpulan
			resultText += "--------------------------------------------------\n"
			if durationIter < durationRec {
				resultText += ">> Iteratif lebih cepat."
			} else {
				resultText += ">> Rekursif lebih cepat/setara."
			}

			// Update UI
			labelResult.SetText(resultText)
			progressBar.Hide()
		}()
	})

	// --- LAYOUT ---

	content := container.NewVBox(
		labelTitle,
		widget.NewSeparator(),
		widget.NewLabel("Ukuran Input (N):"),
		entryN,
		btnProcess,
		progressBar,
		widget.NewSeparator(),
		labelResult,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
