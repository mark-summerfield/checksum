// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	_ "embed"
	//"fmt"
	//"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

//go:embed Version.dat
var Version string

const Hmargin = 6

func main() {
	gtk.Init(nil)
	app, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Failed to create window:", err)
	}
	app.Connect("destroy", func() { gtk.MainQuit() })
	app.SetTitle("Checksum")
	app.SetSizeRequest(360, 120)
	// TODO icon
	mainWindow := NewMainWindow()
	app.Add(mainWindow.container)
	app.ShowAll()
	gtk.Main()

}

type MainWindow struct {
	container        *gtk.Widget
	fileButton       *gtk.Button
	fileEntry        *gtk.Entry
	expectedLabel    *gtk.Label
	expectedEntry    *gtk.Entry
	md5LabelLabel    *gtk.Label
	md5Label         *gtk.Label
	sha1LabelLabel   *gtk.Label
	sha1Label        *gtk.Label
	sha256Label      *gtk.Label
	sha256LabelLabel *gtk.Label
	statusLabel      *gtk.Label
}

func NewMainWindow() *MainWindow {
	mainWindow := makeMainWindow()
	layoutMainWindow(mainWindow)
	return mainWindow
}

func makeMainWindow() *MainWindow {
	fileButton, err := gtk.ButtonNewWithMnemonic("_File...")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	fileButton.SetMarginEnd(Hmargin)
	fileEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	expectedLabel, err := gtk.LabelNewWithMnemonic("_Expected")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	expectedLabel.SetMarginEnd(Hmargin)
	expectedEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	md5LabelLabel, err := gtk.LabelNew("MD5")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	md5Label, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	sha1LabelLabel, err := gtk.LabelNew("SHA1")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	sha1Label, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	sha256LabelLabel, err := gtk.LabelNew("SHA256")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	sha256Label, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	statusLabel, err := gtk.LabelNew("Choose a file...")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	statusLabel.SetHExpand(true)
	for _, entry := range []*gtk.Entry{fileEntry, expectedEntry} {
		entry.SetHExpand(true)
	}
	for _, label := range []*gtk.Label{md5Label, sha1Label, sha256Label} {
		label.SetHExpand(true)
		label.SetMarginEnd(Hmargin)
	}
	return &MainWindow{nil, fileButton, fileEntry,
		expectedLabel, expectedEntry, md5LabelLabel, md5Label,
		sha1LabelLabel, sha1Label, sha256Label, sha256LabelLabel,
		statusLabel}
}

func layoutMainWindow(mainWindow *MainWindow) {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Failed to create grid:", err)
	}
	grid.Attach(mainWindow.fileButton, 0, 0, 1, 1)
	grid.Attach(mainWindow.fileEntry, 1, 0, 1, 1)
	grid.Attach(mainWindow.expectedLabel, 0, 1, 1, 1)
	grid.Attach(mainWindow.expectedEntry, 1, 1, 1, 1)
	grid.Attach(mainWindow.md5LabelLabel, 0, 2, 1, 1)
	grid.Attach(mainWindow.md5Label, 1, 2, 1, 1)
	grid.Attach(mainWindow.sha1LabelLabel, 0, 3, 1, 1)
	grid.Attach(mainWindow.sha1Label, 1, 3, 1, 1)
	grid.Attach(mainWindow.sha256LabelLabel, 0, 4, 1, 1)
	grid.Attach(mainWindow.sha256Label, 1, 4, 1, 1)
	grid.Attach(mainWindow.statusLabel, 0, 5, 2, 1)
	mainWindow.container = &grid.Container.Widget
}
