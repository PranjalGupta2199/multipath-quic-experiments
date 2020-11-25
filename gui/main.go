package main

import (
    "log"
    "fmt"
    "net"
    // "os"
    // "github.com/gotk3/gotk3/glib"
    "github.com/gotk3/gotk3/gtk"
)

func setUpWindow(title string, width, height int) (*gtk.Window){
    win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    if err != nil {
        log.Fatal("Unable to create window: ", err)
    }

    win.SetTitle(title)
    win.SetDefaultSize(width, height)
    win.SetPosition(gtk.WIN_POS_CENTER)
    win.Connect("destroy", func(){
        gtk.MainQuit()
    })  
    return win 
}

func createLabel(title string, wrapText bool) (*gtk.Label){
    label, err := gtk.LabelNew(title)
    if err != nil {
        log.Fatal("Unable to add label: ", err)
    }
    label.SetLineWrap(wrapText)
    return label
}

func createEntry() (*gtk.Entry){
    entry, err := gtk.EntryNew()
    if err != nil {
        log.Fatal("Unable to create entry: ", err)
    }
    return entry
}

func createButton(label string, onClick func(), args ...interface{}) (*gtk.Button){
    button, err := gtk.ButtonNewWithLabel(label)
    if err != nil {
        log.Fatal("Unable to create button: ", err)
    }
    button.Connect("clicked", onClick, args)
    return button
}

func createBox(orient gtk.Orientation, spacing int) (*gtk.Box){
    box, err := gtk.BoxNew(orient, spacing)
    if err != nil{
        log.Fatal("Unable to create a Box")
    }
    return box
}

func createStackSwitcher() (*gtk.StackSwitcher){
    stackSwitcher,err := gtk.StackSwitcherNew()
    if err != nil {
        log.Fatal("Unable to add StackSwitcher: ", err)
    }
    return stackSwitcher
}

func createStack() (*gtk.Stack){
    stack, err := gtk.StackNew()
    if err != nil {
        log.Fatal("Unable to add Stack: ", err)
    }
    return stack
}

func createGrid(columnHomogeneous, rowHomogeneous bool, colSpacing, rowSpacing uint) (*gtk.Grid) {
    grid, err := gtk.GridNew()
    if err != nil {
        log.Fatal("Unable to add grid: ", err)
    }
    grid.SetColumnHomogeneous(columnHomogeneous)
    grid.SetRowHomogeneous(rowHomogeneous)
    grid.SetColumnSpacing(colSpacing)
    grid.SetRowSpacing(rowSpacing)
    return grid
}

func createDialog(title string, width, height int) (*gtk.Dialog){
    dialog, err := gtk.DialogNew()
    if err != nil {
        log.Fatal("Unable to add Dialog: ", err)
    }
    dialog.SetDefaultSize(width, height)
    dialog.SetTitle(title)
    return dialog
}

func createDialogWithButtons(title string, parent gtk.IWindow, flag gtk.DialogFlags, buttons []interface{}) (*gtk.Dialog){
    dialog, err := gtk.DialogNewWithButtons(title, parent, flag, buttons)
    if err != nil {
        log.Fatal("Unable to create Dialog with Buttons: ", err)
    }
    return dialog
}

func getOutboundIPAddr() (string){
    var allIPAddr string

    ifaces, err := net.Interfaces()
    if err != nil {
        log.Fatal("Error retrieving IP Addrs:",err)
        return ""
    }
    for _, i := range ifaces {
        addrs, err := i.Addrs()
        if err != nil {
            log.Fatal("Error retrieving IP Addrs:",err)
            return ""
        }
        for _, addr := range addrs {
            ipnet, ok := addr.(*net.IPNet)

            if !ok{
                continue
            }
            ipv4 := ipnet.IP.To4()
            if ipv4 == nil || ipv4[0] == 127{
                continue
            }
            
            allIPAddr += ipv4.String() + " / "
        }
    }

    return allIPAddr
}

func clientFileTransfer(win *gtk.Window) (*gtk.Grid){
    grid := createGrid(true, false, 5, 20)

    grid.Attach(createLabel("Server IP Address:", false), 0, 0, 2, 1)
    grid.Attach(createEntry(), 2, 0, 2, 1)
    grid.Attach(createLabel("Client IP Address:", false), 0, 1, 2, 1)
    grid.Attach(createLabel(getOutboundIPAddr(), false), 2, 1, 2, 1)


    pathLabel := createLabel("<Path will appear hear>", true)
    fileChooserButton := createButton("Click to Select File", func(){
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


    grid.Attach(createLabel("Select file for transfer: ", false), 0, 2, 2, 1)
    grid.Attach(pathLabel, 2, 2, 2, 1)

    grid.Attach(createButton("Send File", func(){}), 0, 3, 2, 1)
    grid.Attach(fileChooserButton, 2, 3, 2, 1)

    return grid
}

func addClientSide(win *gtk.Window) {
    stackSwitcher,err := gtk.StackSwitcherNew()
    if err != nil {
        log.Fatal("Unable to add StackSwitcher: ", err)
    }   
    stack, err := gtk.StackNew()
    if err != nil {
        log.Fatal("Unable to add stack: ", err)
    }

    grid := clientFileTransfer(win)
    stack.AddTitled(grid, "Page1", "File Transfer")
    stack.AddTitled(createLabel("Hello World", false), "Page2", "Video Stream")
    stackSwitcher.SetStack(stack)

    box := createBox(gtk.ORIENTATION_VERTICAL, 0)
    box.PackStart(stackSwitcher, false, false, 0)
    box.PackStart(stack, true, true, 0)

    win.Add(box)
}

func serverFileTransfer(win *gtk.Window) (*gtk.Grid){
    grid := createGrid(true, false, 5, 20)

    grid.Attach(createLabel("Server IP Address:", false), 0, 0, 2, 1)
    grid.Attach(createLabel(getOutboundIPAddr(), true), 2, 0, 2, 1)
    grid.Attach(createLabel("Server Status:", false), 0, 1, 2, 1)
    statusLabel := createLabel("", true)
    grid.Attach(statusLabel, 2, 1, 2, 1)
    return grid
}


func addServerSide(win *gtk.Window){
    stackSwitcher := createStackSwitcher()
    stack := createStack()

    grid := serverFileTransfer(win)
    stack.AddTitled(grid, "Page1", "File Transfer")
    stack.AddTitled(createLabel("Hello World", false), "Page2", "Video Stream")
    stackSwitcher.SetStack(stack)

    box := createBox(gtk.ORIENTATION_VERTICAL, 0)
    box.PackStart(stackSwitcher, false, false, 0)
    box.PackStart(stack, true, true, 0)

    win.Add(box)
}

func createRadioButton(grpMember *gtk.RadioButton, label string, callback func()) (*gtk.RadioButton){
    radio, err := gtk.RadioButtonNewWithLabelFromWidget(grpMember, label)
    if err != nil {
        log.Fatal("Unable to generate radio button: ", err)
    }
    radio.Connect("toggled", callback)
    return radio
}

func setupDialog() (){
    dialog := createDialog("MPQUIC Experiment", 300, 150)

    dialog.AddButton("OK", gtk.RESPONSE_OK)
    dialog.AddButton("Cancel", gtk.RESPONSE_CLOSE)

    contentArea, err := dialog.GetContentArea()
    if err != nil {
        log.Fatal("Unable to fetch contentArea: ", err)
    }

    box := createBox(gtk.ORIENTATION_VERTICAL, 0)
    clientButton := createRadioButton(nil, "Server", func(){
        role = "server"
    })
    serverButton := createRadioButton(clientButton, "Client", func(){
        role = "client"
    })
    box.PackStart(clientButton, false, false, 0)
    box.PackStart(serverButton, false, false, 0)


    contentArea.PackStart(createLabel("Which role do you want to start?", false), false, false, 0)
    contentArea.PackStart(box, false, false, 0)
    dialog.ShowAll()

    reply := dialog.Run()
    if reply == gtk.RESPONSE_OK {
        fmt.Println("OK")
    }
    dialog.Destroy()
}

var (
    role = "NULL"
    win *gtk.Window
)

func main(){
    gtk.Init(nil)
    setupDialog()

    log.Printf("Selected Role: %s\n", role)

    if role == "client"{
        win = setUpWindow("Client", 800, 200)
        addClientSide(win)
    } else if role == "server" {
        win = setUpWindow("Server", 800, 200)
        addServerSide(win)
    }

    if role != "NULL" {
        win.ShowAll()
        gtk.Main()
    }
}