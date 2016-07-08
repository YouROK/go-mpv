package main

import (
	"log"
	//	"time"

	"github.com/yourok/go-mpv/mpv"
)

func eventListener(m *mpv.Mpv) chan *mpv.Event {
	c := make(chan *mpv.Event)
	go func() {
		for {
			e := m.WaitEvent(1)
			c <- e
		}
	}()
	return c
}

func main() {
	m := mpv.Create()
	c := eventListener(m)
	log.Println("no-resume-playback", m.SetOption("no-resume-playback", mpv.FORMAT_FLAG, true))
	log.Println("volume", m.SetOption("volume", mpv.FORMAT_INT64, 20))
	log.Println("terminal", m.SetOptionString("terminal", "yes"))
	log.Println("softvol", m.SetOptionString("softvol", "yes"))

	// Disable video in three ways.
	//m.SetOption("no-video", mpv.FORMAT_FLAG, true)
	//m.SetOption("vo", mpv.FORMAT_STRING, "null")
	//m.SetOption("vid", mpv.FORMAT_STRING, "")

	//cache
	log.Println("cache-default", m.SetOption("cache-default", mpv.FORMAT_INT64, 160))  // 10 seconds
	log.Println("cache-seek-min", m.SetOption("cache-seek-min", mpv.FORMAT_INT64, 16)) // 1 second

	err := m.Initialize()
	if err != nil {
		log.Println("Mpv init:", err.Error())
		return
	}
	//Set video file
	log.Println("loadfile", m.Command([]string{"loadfile", "http://techslides.com/demos/sample-videos/small.webm"}))

	for {
		e := <-c
		log.Println(e)
		if e.Event_Id == mpv.EVENT_END_FILE {
			break
		}

		/*{
			pos, err := m.GetProperty("time-pos", mpv.FORMAT_DOUBLE)
			if err != nil {
				log.Println("time-pos", err)
				continue
			}
			vol, err := m.GetProperty("volume", mpv.FORMAT_DOUBLE)
			//log.Println(pos.(float64))
			if err != nil {
				log.Println("volume", err)
				continue
			}
			{
				position := pos.(float64)
				if position < 0 {
					position = 0
				}
				log.Println(time.Duration(position * float64(time.Second)))
				log.Println(vol)
			}
		}*/

	}
	m.TerminateDestroy()
}
