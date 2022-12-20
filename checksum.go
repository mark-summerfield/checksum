// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	_ "embed"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"hash"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

//go:embed Version.dat
var Version string

const (
	Margin = 6
	Icon   = "checksum.svg"
	MD5    = "MD5"
	SHA1   = "SHA1"
	SHA256 = "SHA256"
)

func main() {
	gtk.Init(&os.Args)
	filename := ""
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	mainWindow := NewMainWindow("Checksum", filename)
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
	md5Frame         *gtk.Frame
	md5Entry         *gtk.Entry
	sha1LabelLabel   *gtk.Label
	sha1Frame        *gtk.Frame
	sha1Entry        *gtk.Entry
	sha256LabelLabel *gtk.Label
	sha256Frame      *gtk.Frame
	sha256Entry      *gtk.Entry
	statusLabel      *gtk.Label
}

func NewMainWindow(title, filename string) *MainWindow {
	mainWindow := &MainWindow{}
	mainWindow.makeWidgets()
	mainWindow.makeLayout()
	mainWindow.makeConnections(filename)
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
	me.fileButton.SetMarginEnd(Margin)
	me.fileEntry, err = gtk.EntryNew()
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.fileEntry.SetHExpand(true)
	me.expectedLabel, err = gtk.LabelNewWithMnemonic("_Expected")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.expectedLabel.SetMarginEnd(Margin)
	me.expectedEntry, err = gtk.EntryNew()
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.expectedEntry.SetHExpand(true)
	me.md5LabelLabel, err = gtk.LabelNew("MD5")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.md5Frame, err = gtk.FrameNew("")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	prepareFrame(me.md5Frame)
	me.md5Entry, err = gtk.EntryNew()
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	prepareEntry(me.md5Entry)
	me.md5Frame.Add(me.md5Entry)
	me.sha1LabelLabel, err = gtk.LabelNew("SHA1")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.sha1Frame, err = gtk.FrameNew("")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	prepareFrame(me.sha1Frame)
	me.sha1Entry, err = gtk.EntryNew()
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	prepareEntry(me.sha1Entry)
	me.sha1Frame.Add(me.sha1Entry)
	me.sha256LabelLabel, err = gtk.LabelNew("SHA256")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.sha256Frame, err = gtk.FrameNew("")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	prepareFrame(me.sha256Frame)
	me.sha256Entry, err = gtk.EntryNew()
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	prepareEntry(me.sha256Entry)
	me.sha256Frame.Add(me.sha256Entry)
	me.statusLabel, err = gtk.LabelNew("Choose a file...")
	if err != nil {
		log.Fatal("Failed to create widget:", err)
	}
	me.statusLabel.SetHAlign(gtk.ALIGN_START)
}

func prepareFrame(frame *gtk.Frame) {
	frame.SetHExpand(true)
	frame.SetBorderWidth(Margin / 2)
	frame.SetShadowType(gtk.SHADOW_IN)
}

func prepareEntry(entry *gtk.Entry) {
	entry.SetHExpand(true)
	entry.SetEditable(false)
}

func (me *MainWindow) makeLayout() {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Failed to create grid:", err)
	}
	grid.SetRowSpacing(Margin)
	grid.Attach(me.fileButton, 0, 0, 1, 1)
	grid.Attach(me.fileEntry, 1, 0, 1, 1)
	grid.Attach(me.expectedLabel, 0, 1, 1, 1)
	grid.Attach(me.expectedEntry, 1, 1, 1, 1)
	grid.Attach(me.md5LabelLabel, 0, 2, 1, 1)
	grid.Attach(me.md5Frame, 1, 2, 1, 1)
	grid.Attach(me.sha1LabelLabel, 0, 3, 1, 1)
	grid.Attach(me.sha1Frame, 1, 3, 1, 1)
	grid.Attach(me.sha256LabelLabel, 0, 4, 1, 1)
	grid.Attach(me.sha256Frame, 1, 4, 1, 1)
	grid.Attach(me.statusLabel, 1, 5, 1, 1)
	me.container = &grid.Container.Widget
}

func (me *MainWindow) makeConnections(filename string) {
	if filename != "" {
		me.window.Connect("map", func(_ *gtk.Window) {
			me.onNewFile(filename)
		})
	}
	me.window.Connect("destroy", func(_ *gtk.Window) { me.onQuit() })
	me.window.Connect("key-press-event", me.onKeyPress)
	me.fileButton.Connect("clicked", me.onFileButton)
	for _, signal := range []string{"changed", "delete-text", "insert-text",
		"activate", "paste-clipboard"} {
		me.expectedEntry.Connect(signal, func(_ *gtk.Entry) bool {
			me.onChange()
			return true
		})
	}
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
	// Here is where any state would be saved (e.g., window size & pos).
	gtk.MainQuit()
}

func (me *MainWindow) onKeyPress(_ *gtk.Window, event *gdk.Event) {
	keyEvent := &gdk.EventKey{Event: event}
	if keyEvent.KeyVal() == gdk.KEY_Escape {
		me.onQuit()
	}
}

func (me *MainWindow) onFileButton() {
	fileChooserDlg, err := gtk.FileChooserNativeDialogNew(
		"Checksum Choose File", me.window, gtk.FILE_CHOOSER_ACTION_OPEN,
		"_Open", "_Cancel")
	if err != nil {
		log.Fatal("Failed to create file chooser dialog:", err)
	}
	reply := fileChooserDlg.NativeDialog.Run()
	if gtk.ResponseType(reply) == gtk.RESPONSE_ACCEPT {
		filename := fileChooserDlg.GetFilename()
		me.onNewFile(filename)
	}
}

func (me *MainWindow) onNewFile(filename string) {
	me.fileEntry.SetText(filename)
	for _, entry := range []*gtk.Entry{me.md5Entry, me.sha1Entry,
		me.sha256Entry} {
		entry.SetText("")
	}
	me.statusLabel.SetText(fmt.Sprintf("Computing hashes for %s", filename))
	me.expectedEntry.GrabFocus()
	me.setCursor("wait")
	go func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			glib.IdleAdd(func() bool {
				calcHash(filename, MD5, me.md5Entry)
				return false
			})
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			glib.IdleAdd(func() bool {
				calcHash(filename, SHA1, me.sha1Entry)
				return false
			})
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			glib.IdleAdd(func() bool {
				calcHash(filename, SHA256, me.sha256Entry)
				return false
			})
		}()
		wg.Wait()
		// We do this in the idle time since only then have all the widgets
		// been updated with the completed calculations.
		go func() {
			glib.IdleAdd(func() bool {
				me.onChange()
				me.setCursor("default")
				return false
			})
		}()
	}()
}

func (me *MainWindow) onChange() {
	me.statusLabel.SetText("")
	expected, err := me.expectedEntry.GetText()
	if err != nil {
		me.statusLabel.SetText(fmt.Sprintf("error: %s", err))
	}
	if expected != "" {
		expected = strings.TrimSpace(strings.ToLower(expected))
		if h, err := me.md5Entry.GetText(); err == nil && h == expected {
			me.statusLabel.SetText("Expected equals MD5")
			return
		}
		if h, err := me.sha1Entry.GetText(); err == nil && h == expected {
			me.statusLabel.SetText("Expected equals SHA1")
			return
		}
		if h, err := me.sha256Entry.GetText(); err == nil && h == expected {
			me.statusLabel.SetText("Expected equals SHA256")
			return
		}
		me.statusLabel.SetText("Expected doesn't equal any hash")
	} else {
		me.statusLabel.SetText("Enter or Paste Expected to check...")
	}
}

func (me *MainWindow) setCursor(name string) {
	display, err := gdk.DisplayGetDefault()
	if err == nil {
		cursor, err := gdk.CursorNewFromName(display, name)
		if err == nil {
			window, err := me.container.GetWindow()
			if err == nil {
				window.SetCursor(cursor)
			}
		}
	}
}

func calcHash(filename, algorithm string, entry *gtk.Entry) {
	file, err := os.Open(filename)
	if err != nil {
		entry.SetText(fmt.Sprintf("Failed to open %s: %s", filename, err))
		return
	}
	defer file.Close()
	var h hash.Hash
	if algorithm == MD5 {
		h = md5.New()
	} else if algorithm == SHA1 {
		h = sha1.New()
	} else {
		h = sha256.New()
	}
	if _, err := io.Copy(h, file); err != nil {
		entry.SetText(fmt.Sprintf("Failed to compute %s: %s", algorithm,
			err))
		return
	}
	entry.SetText(fmt.Sprintf("%x", h.Sum(nil)))
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
