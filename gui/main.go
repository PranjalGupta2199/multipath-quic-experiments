package main

import (
    "log"
    "fmt"
    //"os"
    "syscall"
    "bytes"
    "os/exec"
    "time"
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

func streamGeneratorRoutine(quit chan bool){
    c := exec.Command("python3", "../live-video-stream/stream_generator.py")
    var out bytes.Buffer
    var stderr bytes.Buffer
    
    c.Stdout = &out
    c.Stderr = &stderr
    
    pythonChannel := make(chan bool)

    go func(quit chan bool){
        err := c.Start()
        log.Println("stream_generator.py started")

        for {
            select{
            case <- quit:
                if err := c.Process.Kill(); err != nil{
                    log.Println("Error occurred while killing stream_generator.py")
                } else {
                    log.Println("Killed stream_generator.py process")
                }
                return
            default:
                if err != nil {
                    log.Fatal("Error running stream_generator.py")
                    log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
                } 
            }
        }
    }(pythonChannel)

    for {
        select{
        case <- quit:
            pythonChannel <- true
            // close stream_generator.py process
            return
        }
    }
}

func streamSenderRoutine(quit chan bool){
    c := exec.Command("go","run", "../live-video-stream/stream_sender.go")
    var out bytes.Buffer
    var stderr bytes.Buffer
    
    c.Stdout = &out
    c.Stderr = &stderr
    
    goChannel := make(chan bool)
    
    go func(quit chan bool){ 
        err := c.Start()
        log.Println("stream_sender.go started")
        
        for {
            select{
            case <- quit:
                if err := c.Process.Kill(); err != nil{
                    log.Println("Error occurred while killing stream_sender.go")
                } else {
                    log.Println("Killed stream_sender.go process")
                }
                return
            default:
                if err != nil {
                    log.Fatal("Error running stream_sender.go")
                    log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
                } 
            }
        }
    }(goChannel)

    for {
        select{
        case <- quit:
            // close stream_sender.go process
            goChannel <- true
            return
        }
    }
}

func streamReceiverRoutine(quit chan bool){
    c := exec.Command("go", "run", "../live-video-stream/stream_receiver.go")
    c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
    var out bytes.Buffer
    var stderr bytes.Buffer

    c.Stdout = &out
    c.Stderr = &stderr

    goChannel := make(chan bool)
    
    go func (quit chan bool){
        err := c.Start()
        log.Println("stream_receiver.go started")

        for {
            select{
            case <- quit:
                pgid, err := syscall.Getpgid(c.Process.Pid)
                if err == nil {
                    syscall.Kill(-pgid, 15)  // note the minus sign
                }

                err = c.Wait()
                log.Println(err)
                return
            default:
                if err != nil {
                    log.Fatal("Error running stream_receiver.go")
                    log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
                } else {
                    // log.Println(out.String())
                } 
            }
        }
    }(goChannel)

    for {
        select{
        case <- quit:
            // close stream_receiver.go process
            goChannel <- true
            return
        }
    }
}

func streamVideoOnGUI(quit chan bool, img *gtk.Image){
    // to make sure that atleast one frame is present
    time.Sleep(3 * time.Second)

    counter := 0
    // missing_frame_count := 0
    path := "../live-video-stream/sample/img%d.jpg"
    
    for {
        select{
        case <- quit:
            img.Clear()
            return
        default:
            // Start showing images
            currentPath := fmt.Sprintf(path, counter)

            if utils.PathExists(currentPath) {
                // log.Printf("Setting Image file %s", currentPath)
                img.SetFromFile(currentPath)
                // time.Sleep(200 * time.Millisecond)
                counter += 1
                img.Show()
            }                
            // } else if missing_frame_count <= 100{
            //     // log.Println("Missing Frame Count: %d", missing_frame_count)
            //     time.Sleep(200 * time.Millisecond)
            //     missing_frame_count += 1
            // } else {
            //     log.Println("Exiting")
            //     img.Clear()
            //     return
                // quit <- true
            // }

        }
    }

}

func clientVideoStream(win *gtk.Window) (*gtk.Grid) {
    grid := widgets.GridNew(true, false, 5, 20)
    img := widgets.ImageNew()


    streamGeneratorChannel := make(chan bool)
    streamSenderChannel := make(chan bool)
    videoGUIChannel := make(chan bool)

    startButton := widgets.ButtonNew("Start", func(){
        button, err := grid.GetChildAt(1, 0)
        utils.HandleError(err)

        // start stream_generator.py
        go streamGeneratorRoutine(streamGeneratorChannel)
        // start stream_sender.go
        go streamSenderRoutine(streamSenderChannel)
        // start video stream on GUI
        // go streamVideoOnGUI(videoGUIChannel, img)

        button.ToWidget().SetSensitive(false)
    }, streamGeneratorChannel, streamSenderChannel, videoGUIChannel, img)

    stopButton := widgets.ButtonNew("Stop", func (){
        button, err := grid.GetChildAt(1, 0)
        utils.HandleError(err)
        
        if sensitive := button.ToWidget().GetSensitive(); !sensitive{
            // stop stream_generator.py
            streamSenderChannel <- true
            // stop stream_sender.go
            streamGeneratorChannel <- true
            // stop video stream on GUI
            // videoGUIChannel <- true
        }
        button.ToWidget().SetSensitive(true)
    }, streamGeneratorChannel, streamSenderChannel, videoGUIChannel, grid)

    grid.Attach(startButton, 1, 0, 1, 1)
    grid.Attach(stopButton, 2, 0, 1, 1)
    grid.Attach(widgets.LabelNew("FPS : ", false), 1, 1, 1, 1)
    grid.Attach(widgets.LabelNew("--", false), 2, 1, 1, 1)
    grid.Attach(img, 0, 3, 4, 4)
    return grid
}


func addClientSide(win *gtk.Window) {
    stackSwitcher := widgets.StackSwitcherNew()  
    stack := widgets.StackNew()

    gridFileTransfer := clientFileTransfer(win)
    gridVideoStream := clientVideoStream(win)
    stack.AddTitled(gridFileTransfer, "Page1", "File Transfer")
    stack.AddTitled(gridVideoStream, "Page2", "Video Stream")
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


func serverVideoStream(win *gtk.Window) (*gtk.Grid) {
    grid := widgets.GridNew(true, false, 5, 20)
    img := widgets.ImageNew()
    
    streamReceiverChannel := make(chan bool)
    videoGUIChannel := make(chan bool)
    
    startButton := widgets.ButtonNew("Start", func (){
        button, err := grid.GetChildAt(1, 0)
        utils.HandleError(err)
        
        // start stream_receiver.go
        go streamReceiverRoutine(streamReceiverChannel)
        // start videoGUIChannel
        go streamVideoOnGUI(videoGUIChannel, img)

        button.ToWidget().SetSensitive(false)
    }, grid, img)


    stopButton := widgets.ButtonNew("Stop", func (){
        button, err := grid.GetChildAt(1,0)
        utils.HandleError(err)

        if sensitive := button.ToWidget().GetSensitive(); !sensitive{
            // stop stream_receiver.go
            streamReceiverChannel <- true
            // stop video stream on GUI
            videoGUIChannel <- true
        }

        button.ToWidget().SetSensitive(true)
    }, grid)

    
    grid.Attach(startButton, 1, 0, 1, 1)
    grid.Attach(stopButton, 2, 0, 1, 1)
    grid.Attach(widgets.LabelNew("FPS : ", false), 1, 1, 1, 1)
    grid.Attach(widgets.LabelNew("--", false), 2, 1, 1, 1)
    grid.Attach(img, 0, 3, 4, 4)
    return grid

}


func addServerSide(win *gtk.Window){
    stackSwitcher := widgets.StackSwitcherNew()
    stack := widgets.StackNew()

    gridFileTransfer := serverFileTransfer(win)
    gridVideoStream := serverVideoStream(win)
    stack.AddTitled(gridFileTransfer, "Page1", "File Transfer")
    stack.AddTitled(gridVideoStream, "Page2", "Video Stream")
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
        win = widgets.WindowNew("Server", 800, 600)
        addServerSide(win)
    }
    
    if !close {
        win.ShowAll()
        gtk.Main()
    }

}