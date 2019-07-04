package main

import (
	"fmt"
	"github.com/dawidd6/go-appindicator"
	"github.com/dawidd6/go-spotify-dbus"
	"github.com/godbus/dbus"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

const (
	id       = "indicator-spotify"
	icon     = id
	category = appindicator.CategoryApplicationStatus
)

func main() {
	gtk.Init(nil)
	defer gtk.Main()

	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}

	menu, err := gtk.MenuNew()
	if err != nil {
		log.Fatal(err)
	}
	defer menu.ShowAll()

	itemQuit, err := gtk.MenuItemNewWithLabel("Quit")
	if err != nil {
		log.Fatal(err)
	}
	defer menu.Add(itemQuit)

	_, err = itemQuit.Connect("activate", func() {
		gtk.MainQuit()
	})
	if err != nil {
		log.Fatal(err)
	}

	indicator := appindicator.New(id, icon, category)
	indicator.SetMenu(menu)

	song := ""
	status := ""

	listeners := &spotify.Listeners{
		OnMetadata: func(metadata *spotify.Metadata) {
			song = fmt.Sprintf("%s - %s", metadata.Artist, metadata.Title)

			set(song, status, indicator)
		},
		OnPlaybackStatus: func(playbackStatus spotify.PlaybackStatus) {
			switch playbackStatus {
			case spotify.StatusPlaying:
				status = "\u25B6"
			case spotify.StatusPaused:
				status = "\u23F8"
			}

			set(song, status, indicator)
		},
		OnServiceStart: func() {
			indicator.SetStatus(appindicator.StatusActive)
			indicator.SetLabel("", "")
		},
		OnServiceStop: func() {
			indicator.SetStatus(appindicator.StatusPassive)
			indicator.SetLabel("", "")
		},
	}

	go spotify.Listen(conn, nil, listeners)
}

func set(song, status string, indicator *appindicator.Indicator) {
	label := fmt.Sprintf("%s %s", status, song)
	indicator.SetLabel(label, "")
}
