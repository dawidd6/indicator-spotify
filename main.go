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
    id = "indicator-spotify"
    icon = id
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
    indicator.SetStatus(appindicator.StatusActive)

    song := ""
    symbol := ""

    listeners := &spotify.Listeners{
        OnMetadata: func(metadata *spotify.Metadata) {
            song = fmt.Sprintf("%s - %s", metadata.Artist, metadata.Title)
            label := fmt.Sprintf("%s %s", symbol, song)
            indicator.SetLabel(label, "")
        },
        OnPlaybackStatus: func(status spotify.PlaybackStatus) {
            switch status {
            case spotify.StatusPlaying:
                symbol = "\u25B6"
            case spotify.StatusPaused:
                symbol = "\u23F8"
            }
            label := fmt.Sprintf("%s %s", symbol, song)
            indicator.SetLabel(label, "")
        },
        OnServiceStart: func() {
            indicator.SetStatus(appindicator.StatusActive)
            indicator.SetLabel("","")
        },
        OnServiceStop: func() {
            indicator.SetStatus(appindicator.StatusPassive)
            indicator.SetLabel("", "")
        },
    }

    go spotify.Listen(conn, listeners)
}


