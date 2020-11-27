package file_transfer

import (
    "log"

    "github.com/gotk3/gotk3/gtk"

    widgets "../widgets"
    utils "../utils"
)


func SetupSenderUI(win *gtk.Window) (*gtk.Grid){
    grid := widgets.GridNew(true, false, 5, 20)

    grid.Attach(widgets.LabelNew("Server IP Address:", false), 0, 0, 2, 1)
    grid.Attach(widgets.EntryNew(), 2, 0, 2, 1)
    grid.Attach(widgets.LabelNew("Client IP Address:", false), 0, 1, 2, 1)
    grid.Attach(widgets.LabelNew(utils.GetOutboundIPAddr(), false), 2, 1, 2, 1)

    // senderChannel := make(chan bool)

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

    sendFileButton := widgets.ButtonNew("Send File", func (){
        button, err := grid.GetChildAt(0, 3)
        utils.HandleError(err)

        path, err := pathLabel.GetText()
        utils.HandleError(err)

        // start client_multipath.go
        go SenderRoutine(path, button)

        button.ToWidget().SetSensitive(false)
    }, grid, pathLabel)

    grid.Attach(widgets.LabelNew("Select file for transfer: ", false), 0, 2, 2, 1)
    grid.Attach(pathLabel, 2, 2, 2, 1)

    grid.Attach(sendFileButton, 0, 3, 2, 1)
    grid.Attach(fileChooserButton, 2, 3, 2, 1)

    return grid
}

func SetupReceiverUI(win *gtk.Window) (*gtk.Grid){
    grid := widgets.GridNew(true, false, 5, 20)

    receiverChannel := make(chan bool)

    statusLabel := widgets.LabelNew("", true)
    startButton := widgets.ButtonNew("Start", func (){
        button, err := grid.GetChildAt(1, 2)
        utils.HandleError(err)

        // start sender-multipath.go
        go ReceiverRoutine(receiverChannel, button)

        button.ToWidget().SetSensitive(false)
    }, grid)

    stopButton := widgets.ButtonNew("Stop", func (){
        button, err := grid.GetChildAt(1, 2)
        utils.HandleError(err)

        if !button.ToWidget().GetSensitive() {
            receiverChannel <- true
        }

        button.ToWidget().SetSensitive(true)
    }, grid)


    grid.Attach(widgets.LabelNew("Server IP Address:", false), 0, 0, 2, 1)
    grid.Attach(widgets.LabelNew(utils.GetOutboundIPAddr(), true), 2, 0, 2, 1)
    grid.Attach(widgets.LabelNew("Server Status:", false), 0, 1, 2, 1)
    grid.Attach(statusLabel, 2, 1, 2, 1)
    grid.Attach(startButton, 1, 2, 1, 1)
    grid.Attach(stopButton, 2, 2, 1, 1)
    return grid
}