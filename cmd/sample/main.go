package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
)

const (
	sampleRate      = 44100
	channelCount    = 1
	framesPerBuffer = 64
)

var (
	recordDuration int
)

func init() {
	flag.IntVar(&recordDuration, "duration", 10, "Recording duration in seconds")
	flag.Parse()
}

func main() {
	err := portaudio.Initialize()
	if err != nil {
		fmt.Println("Error initializing PortAudio:", err)
		return
	}
	defer portaudio.Terminate()

	inputParameters := portaudio.LowLatencyParameters(nil, nil)
	inputParameters.Input.Channels = channelCount
	inputParameters.SampleRate = sampleRate
	inputParameters.FramesPerBuffer = framesPerBuffer

	file, err := os.Create("tmp/output.wav")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file.Close()

	wavEncoder := wav.NewEncoder(file, sampleRate, 16, channelCount, 1)
	defer wavEncoder.Close()

	stream, err := portaudio.OpenStream(inputParameters, nil, callback(wavEncoder))
	if err != nil {
		fmt.Println("Error opening stream:", err)
		return
	}
	defer stream.Close()

	err = stream.Start()
	if err != nil {
		fmt.Println("Error starting stream:", err)
		return
	}

	fmt.Println("Recording for", recordDuration, "seconds...")
	time.Sleep(time.Duration(recordDuration) * time.Second)

	err = stream.Stop()
	if err != nil {
		fmt.Println("Error stopping stream:", err)
		return
	}

	fmt.Println("Recording saved to output.wav")
}

func callback(wavEncoder *wav.Encoder) func(in []float32) {
	return func(in []float32) {
		buf := &audio.IntBuffer{
			Format: &audio.Format{
				NumChannels: channelCount,
				SampleRate:  sampleRate,
			},
			Data: make([]int, len(in)),
		}

		for i, f := range in {
			buf.Data[i] = int(f * 32767.0)
		}

		err := wavEncoder.Write(buf)
		if err != nil {
			fmt.Println("Error encoding audio data:", err)
			return
		}
	}
}
