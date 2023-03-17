package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gordonklaus/portaudio"
)

var (
	sigs = make(chan os.Signal, 1)
	done = make(chan bool, 1)
)

func init() {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println("###############")
		fmt.Println(sig)
		done <- true
		fmt.Println("exiting")
	}()
}

func main() {

	fmt.Println("#### start audio ####")

	// portauido
	portaudio.Initialize()
	defer portaudio.Terminate()
	e := newEcho(time.Second / 3)
	defer e.Close()

	chk(e.Start())

	<-done

	chk(e.Stop())
	fmt.Println("#### echo stopping ####")

}

type echo struct {
	*portaudio.Stream
	buffer []float32
	i      int
}

func newEcho(delay time.Duration) *echo {
	h, err := portaudio.DefaultHostApi()
	chk(err)
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
	p.Input.Channels = 1
	p.Output.Channels = 1
	e := &echo{buffer: make([]float32, int(p.SampleRate*delay.Seconds()))}
	e.Stream, err = portaudio.OpenStream(p, e.processAudio)
	chk(err)
	return e
}

func (e *echo) processAudio(in, out []float32) {
	for i := range out {
		out[i] = .7 * e.buffer[e.i]
		e.buffer[e.i] = in[i]
		e.i = (e.i + 1) % len(e.buffer)
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
