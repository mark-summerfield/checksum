// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	_ "embed"
	//"fmt"
	//"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gdk"
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
	mainWindow := NewMainWindow("Checksum")
	mainWindow.window.ShowAll()
	gtk.Main()

}

type MainWindow struct {
	window           *gtk.Window
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

func NewMainWindow(title string) *MainWindow {
	mainWindow := &MainWindow{}
	mainWindow.makeWidgets()
	mainWindow.makeLayout()
	mainWindow.makeConnections()
	mainWindow.window.SetTitle(title)
	mainWindow.window.SetSizeRequest(360, 120)
	mainWindow.window.Add(mainWindow.container)
	mainWindow.addIcon()
	return mainWindow
}

func (me *MainWindow) makeWidgets() {
	var err error
	me.window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Failed to create window:", err)
	}
	me.fileButton, err = gtk.ButtonNewWithMnemonic("_File...")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.fileButton.SetMarginEnd(Hmargin)
	me.fileEntry, err = gtk.EntryNew()
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.fileEntry.SetHExpand(true)
	me.expectedLabel, err = gtk.LabelNewWithMnemonic("_Expected")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.expectedLabel.SetMarginEnd(Hmargin)
	me.expectedEntry, err = gtk.EntryNew()
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.expectedEntry.SetHExpand(true)
	me.md5LabelLabel, err = gtk.LabelNew("MD5")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.md5Label, err = gtk.LabelNew("")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.md5Label.SetHExpand(true)
	me.md5Label.SetMarginEnd(Hmargin)
	me.sha1LabelLabel, err = gtk.LabelNew("SHA1")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.sha1Label, err = gtk.LabelNew("")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.sha1Label.SetHExpand(true)
	me.sha1Label.SetMarginEnd(Hmargin)
	me.sha256LabelLabel, err = gtk.LabelNew("SHA256")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.sha256Label, err = gtk.LabelNew("")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.sha256Label.SetHExpand(true)
	me.sha256Label.SetMarginEnd(Hmargin)
	me.statusLabel, err = gtk.LabelNew("Choose a file...")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
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

func (me *MainWindow) makeConnections() {
	me.window.Connect("destroy", func(_ *gtk.Window) { me.onQuit() })
	me.window.Connect("key-press-event", func(_ *gtk.Window,
		event *gdk.Event) {
		keyEvent := &gdk.EventKey{Event: event}
		me.onKey(keyEvent)
	})
	me.fileButton.Connect("clicked", func() {
		fileChooserDlg, err := gtk.FileChooserNativeDialogNew("Open",
			me.window, gtk.FILE_CHOOSER_ACTION_OPEN, "_Open", "_Cancel")
		if err != nil {
			log.Fatal("Failed to create file chooser dialog:", err)
		}
		reply := fileChooserDlg.NativeDialog.Run()
		if gtk.ResponseType(reply) == gtk.RESPONSE_ACCEPT {
			filename := fileChooserDlg.GetFilename()
			me.onNewFile(filename)
		}
	})
	me.expectedLabel.Connect("mnemonic-activate", func(_ *gtk.Label) bool {
		me.expectedEntry.GrabFocus()
		return true
	})
}

func (me *MainWindow) addIcon() {
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
			err := me.window.SetIconFromFile(filename)
			if err != nil {
				log.Println("Failed to load icon:", err)
			}
		}
	}
}

func (me *MainWindow) onQuit() {
	// Here is where any state would be saved.
	gtk.MainQuit()
}

func (me *MainWindow) onKey(event *gdk.EventKey) {
	if event.KeyVal() == gdk.KEY_Escape {
		me.onQuit()
	}
}

func (me *MainWindow) onNewFile(filename string) {
	me.fileEntry.SetText(filename)
	me.expectedEntry.GrabFocus()
	log.Println("onNewFile", filename) // TODO read hashes using glib.IdleAdd()
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
