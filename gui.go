package main

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

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

func main() {
	a := app.New()
	w := a.NewWindow("yt-simple-dl-gui")

	w.SetOnClosed(func() {
		log.Println("GUI window closed")
		// killProcessByName("yt-dlp")
	})

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
	w.ShowAndRun()
}
