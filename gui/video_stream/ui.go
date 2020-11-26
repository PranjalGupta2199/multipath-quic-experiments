package video_stream

import (
    "fmt"
    "time"
    "github.com/gotk3/gotk3/gtk"

    config "../config"
    utils "../utils"
    widgets "../widgets"
)


func SetupVideo(quit chan bool, img *gtk.Image){
    time.Sleep(3 * time.Second)

    counter := 0
    path := config.VIDEO_DIR + "/img%d.jpg"
    
    for {
        select{
        case <- quit:
            img.Clear()
            return
        default:
            currentPath := fmt.Sprintf(path, counter)

            if utils.PathExists(currentPath) {
                img.SetFromFile(currentPath)
                // time.Sleep(200 * time.Millisecond)
                counter += 1
                img.Show()
            }                
        }
    }
}


func SetupSenderUI(win *gtk.Window) (*gtk.Grid) {
    grid := widgets.GridNew(true, false, 5, 20)
    img := widgets.ImageNew()
    
    streamReceiverChannel := make(chan bool)
    videoGUIChannel := make(chan bool)
    
    startButton := widgets.ButtonNew("Start", func (){
        button, err := grid.GetChildAt(1, 0)
        utils.HandleError(err)
        
        // start stream_receiver.go
        go ReceiverRoutine(streamReceiverChannel)
        // start videoGUIChannel
        go SetupVideo(videoGUIChannel, img)

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

func SetupReceiverUI(win *gtk.Window) (*gtk.Grid) {
    grid := widgets.GridNew(true, false, 5, 20)
    img := widgets.ImageNew()


    streamGeneratorChannel := make(chan bool)
    streamSenderChannel := make(chan bool)
    // videoGUIChannel := make(chan bool)

    startButton := widgets.ButtonNew("Start", func(){
        button, err := grid.GetChildAt(1, 0)
        utils.HandleError(err)

        // start stream_generator.py
        go GeneratorRoutine(streamGeneratorChannel)
        // start stream_sender.go
        go SenderRoutine(streamSenderChannel)
        // start video stream on GUI
        // go streamVideoOnGUI(videoGUIChannel, img)

        button.ToWidget().SetSensitive(false)
    }, streamGeneratorChannel, streamSenderChannel, img)

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
    }, streamGeneratorChannel, streamSenderChannel, grid)

    grid.Attach(startButton, 1, 0, 1, 1)
    grid.Attach(stopButton, 2, 0, 1, 1)
    grid.Attach(widgets.LabelNew("FPS : ", false), 1, 1, 1, 1)
    grid.Attach(widgets.LabelNew("--", false), 2, 1, 1, 1)
    grid.Attach(img, 0, 3, 4, 4)
    return grid
}
