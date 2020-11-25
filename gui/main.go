package main

import (
    "log"
    "fmt"
    // "net"
    // "os"
    // "github.com/gotk3/gotk3/glib"
    "github.com/gotk3/gotk3/gtk"

    utils "./utils"
    widgets "./widgets"
)


func clientFileTransfer(win *gtk.Window) (*gtk.Grid){
    grid := widgets.GridNew(true, false, 5, 20)

    grid.Attach(widgets.LabelNew("Server IP Address:", false), 0, 0, 2, 1)
    grid.Attach(widgets.EntryNew(), 2, 0, 2, 1)
    grid.Attach(widgets.LabelNew("Client IP Address:", false), 0, 1, 2, 1)
    grid.Attach(widgets.LabelNew(utils.GetOutboundIPAddr(), false), 2, 1, 2, 1)


    pathLabel := widgets.LabelNew("<Path will appear hear>", true)
    fileChooserButton := widgets.ButtonNew("Click to Select File", func(){
        dialog, err := gtk.FileChooserDialogNewWith2Buttons(
                            "Choose file to send from client to server",
                            win,
                            gtk.FILE_CHOOSER_ACTION_OPEN,
                            "Cancel",
                            gtk.RESPONSE_CLOSE,
                            "Select",
                            gtk.RESPONSE_OK)
    
        if err != nil {
            log.Fatal("Unable to create FileChooserDialog: ", err)
        }

        reply := dialog.Run()
        if reply == gtk.RESPONSE_OK {
            pathLabel.SetText(dialog.GetFilename())
        }
        dialog.Destroy()
    }, pathLabel)


    grid.Attach(widgets.LabelNew("Select file for transfer: ", false), 0, 2, 2, 1)
    grid.Attach(pathLabel, 2, 2, 2, 1)

    grid.Attach(widgets.ButtonNew("Send File", func(){}), 0, 3, 2, 1)
    grid.Attach(fileChooserButton, 2, 3, 2, 1)

    return grid
}

func addClientSide(win *gtk.Window) {
    stackSwitcher := widgets.StackSwitcherNew()  
    stack := widgets.StackNew()

    grid := clientFileTransfer(win)
    stack.AddTitled(grid, "Page1", "File Transfer")
    stack.AddTitled(widgets.LabelNew("Hello World", false), "Page2", "Video Stream")
    stackSwitcher.SetStack(stack)

    box := widgets.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    box.PackStart(stackSwitcher, false, false, 0)
    box.PackStart(stack, true, true, 0)

    win.Add(box)
}

func serverFileTransfer(win *gtk.Window) (*gtk.Grid){
    grid := widgets.GridNew(true, false, 5, 20)

    grid.Attach(widgets.LabelNew("Server IP Address:", false), 0, 0, 2, 1)
    grid.Attach(widgets.LabelNew(utils.GetOutboundIPAddr(), true), 2, 0, 2, 1)
    grid.Attach(widgets.LabelNew("Server Status:", false), 0, 1, 2, 1)
    statusLabel := widgets.LabelNew("", true)
    grid.Attach(statusLabel, 2, 1, 2, 1)
    return grid
}


func addServerSide(win *gtk.Window){
    stackSwitcher := widgets.StackSwitcherNew()
    stack := widgets.StackNew()

    grid := serverFileTransfer(win)
    stack.AddTitled(grid, "Page1", "File Transfer")
    stack.AddTitled(widgets.LabelNew("Hello World", false), "Page2", "Video Stream")
    stackSwitcher.SetStack(stack)

    box := widgets.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    box.PackStart(stackSwitcher, false, false, 0)
    box.PackStart(stack, true, true, 0)

    win.Add(box)
}


func setupDialog() (){
    dialog := widgets.DialogNew("MPQUIC Experiment", 300, 150)

    dialog.AddButton("OK", gtk.RESPONSE_OK)
    dialog.AddButton("Cancel", gtk.RESPONSE_CLOSE)

    contentArea, err := dialog.GetContentArea()
    if err != nil {
        log.Fatal("Unable to fetch contentArea: ", err)
    }

    box := widgets.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    clientButton := widgets.RadioButtonNew(nil, "Server", func(){
        role = "server"
    })
    serverButton := widgets.RadioButtonNew(clientButton, "Client", func(){
        role = "client"
    })
    box.PackStart(clientButton, false, false, 0)
    box.PackStart(serverButton, false, false, 0)


    contentArea.PackStart(widgets.LabelNew("Which role do you want to start?", false), false, false, 0)
    contentArea.PackStart(box, false, false, 0)
    dialog.ShowAll()

    reply := dialog.Run()
    if reply == gtk.RESPONSE_OK {
        fmt.Println("OK")
    } else {
        close = true
    }
    dialog.Destroy()
}

var (
    role = "server"
    close = false
    win *gtk.Window
)

func main(){
    gtk.Init(nil)
    setupDialog()

    log.Printf("Selected Role: %s\n", role)

    if role == "client"{
        win = widgets.WindowNew("Client", 800, 200)
        addClientSide(win)
    } else {
        win = widgets.WindowNew("Server", 800, 200)
        addServerSide(win)
    }
    
    if !close {
        win.ShowAll()
        gtk.Main()
    }

}