// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	_ "embed"
	//"fmt"
	//"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"path"
)

//go:embed Version.dat
var Version string

const (
	Hmargin = 6
	Icon    = "checksum.svg"
)

func main() {
	gtk.Init(nil)
	app := getApp()
	mainWindow := NewMainWindow()
	app.Add(mainWindow.container)
	app.ShowAll()
	gtk.Main()

}

func getApp() *gtk.Window {
	app, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Failed to create window:", err)
	}
	app.Connect("destroy", func() { gtk.MainQuit() })
	app.SetTitle("Checksum")
	app.SetSizeRequest(360, 120)
	addIcon(app)
	return app
}

func addIcon(app *gtk.Window) {
	filename, err := os.Executable()
	if err == nil {
		filename = path.Join(path.Dir(filename), Icon)
		if !PathExists(filename) {
			filename = ""
			folder, err := os.Getwd()
			if err == nil {
				filename = path.Join(folder, Icon)
			}
		}
		if filename != "" {
			err := app.SetIconFromFile(filename)
			if err != nil {
				log.Println("Failed to load icon:", err)
			}
		}
	}
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
	mainWindow := &MainWindow{}
	mainWindow.makeWidgets()
	mainWindow.makeLayout()
	mainWindow.makeConnections()
	return mainWindow
}

func (me *MainWindow) makeWidgets() {
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
	me.fileButton = fileButton
	me.fileEntry = fileEntry
	me.expectedLabel = expectedLabel
	me.expectedEntry = expectedEntry
	me.md5LabelLabel = md5LabelLabel
	me.md5Label = md5Label
	me.sha1LabelLabel = sha1LabelLabel
	me.sha1Label = sha1Label
	me.sha256Label = sha256Label
	me.sha256LabelLabel = sha256LabelLabel
	me.statusLabel = statusLabel
}

func (me *MainWindow) makeLayout() {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Failed to create grid:", err)
	}
	grid.Attach(me.fileButton, 0, 0, 1, 1)
	grid.Attach(me.fileEntry, 1, 0, 1, 1)
	grid.Attach(me.expectedLabel, 0, 1, 1, 1)
	grid.Attach(me.expectedEntry, 1, 1, 1, 1)
	grid.Attach(me.md5LabelLabel, 0, 2, 1, 1)
	grid.Attach(me.md5Label, 1, 2, 1, 1)
	grid.Attach(me.sha1LabelLabel, 0, 3, 1, 1)
	grid.Attach(me.sha1Label, 1, 3, 1, 1)
	grid.Attach(me.sha256LabelLabel, 0, 4, 1, 1)
	grid.Attach(me.sha256Label, 1, 4, 1, 1)
	grid.Attach(me.statusLabel, 0, 5, 2, 1)
	me.container = &grid.Container.Widget
}

// TODO Esc → quit
func (me *MainWindow) makeConnections() {
	me.fileButton.Connect("activate", func(_ *gtk.Button) bool {
		log.Println("TODO show file choose dialog") // TODO
		// If user chooses then compute hashes each using a function set in
		// a goroutine for glib.IdleAdd: see
		// ~/zip/gtk-examples/goroutines/goroutines.go
		return true
	})
	me.expectedLabel.Connect("mnemonic-activate", func(_ *gtk.Label) bool {
		me.expectedEntry.GrabFocus()
		return true
	})
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
