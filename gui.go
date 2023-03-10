package main

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

var VERSION = "DEV"

func downloadHandler(window fyne.Window, input *widget.Entry, check *widget.Check, bottom *widget.Label) {
	if !strings.Contains(strings.ToLower(input.Text), "youtube.com") && !strings.Contains(input.Text, "youtu.be") {
		bottom.SetText("Najpierw podaj poprawny link do YouTube!")
		return
	}

	input.Disable()
	check.Disable()
	defer input.Enable()
	defer check.Enable()

	bottom.SetText("CHWILECZKĘ...")
	err := downloadVideo(input.Text, check.Checked)
	if err != nil {
		bottom.SetText("ERROR :(")
		dialog.ShowError(err, window)
		log.Println("Error:", err)
		// log.Panic
	} else {
		bottom.SetText("POBRANE :)")
	}
}

func showUpdateDialog(window fyne.Window, newVersion string, updateLink string) {
	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("%s -> %s", VERSION, newVersion)),
		widget.NewButton("Kopiuj link", func() {
			clipboard.WriteAll(updateLink)
		}),
	)
	dialog.ShowCustom("Dostępna aktualizacja", "OK", content, window)
}

func main() {
	a := app.New()
	w := a.NewWindow("yt-simple-dl-gui")

	hello := widget.NewLabel("YouTube Simple DL!")
	bottom := widget.NewLabel("")

	check := widget.NewCheck("Konwertuj do mp3", nil)
	check.SetChecked(true)

	input := widget.NewEntry()
	input.SetPlaceHolder("Wprowadź link do YouTube...")

	button := widget.NewButton("Pobierz!", func() {
		downloadHandler(w, input, check, bottom)
	})

	w.SetContent(container.NewVBox(
		hello,
		input,
		button,
		check,
		bottom,
	))

	w.Resize(fyne.Size{Width: 350})
	w.SetFixedSize(true)

	go func() {
		newVersion, updateLink, err := checkVersion(VERSION)
		if err != nil {
			log.Println("Error:", err)
		} else if newVersion != "" {
			showUpdateDialog(w, newVersion, updateLink)
		}
	}()

	w.ShowAndRun()
}
