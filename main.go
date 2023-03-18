package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Error loading .env file")
	}

	apiKey := os.Getenv("OPEN_AI_API_KEY")

	if len(os.Args) < 2 {
		os.Args[1] = time.Now().Format(time.RFC3339Nano)
	}
	fmt.Println("Recording.  Press Ctrl-C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	fileName := os.Args[1]
	if !strings.HasSuffix(fileName, ".aiff") {
		fileName += ".aiff"
	}
	f, err := os.Create("tmp/" + fileName)
	if err != nil {
		panic(err)
	}

	// form chunk
	_, err = f.WriteString("FORM")
	if err != nil {
		panic(err)
	}
	if err := binary.Write(f, binary.BigEndian, int32(0)); err != nil {
		panic(err)
	}
	//total bytes
	_, err = f.WriteString("AIFF")
	if err != nil {
		panic(err)
	}

	// common chunk
	_, err = f.WriteString("COMM")
	if err != nil {
		panic(err)
	}
	err = binary.Write(f, binary.BigEndian, int32(18)) //size
	if err != nil {
		panic(err)
	}
	err = binary.Write(f, binary.BigEndian, int16(1)) //channels
	if err != nil {
		panic(err)
	}
	err = binary.Write(f, binary.BigEndian, int32(0)) //number of samples
	if err != nil {
		panic(err)
	}
	err = binary.Write(f, binary.BigEndian, int16(32)) //bits per sample
	if err != nil {                                    //sample rate
		panic(err) //80-bit sample rate 44100
	} //80-bit sample rate 44100
	_, err = f.Write([]byte{0x40, 0x0e, 0xac, 0x44, 0, 0, 0, 0, 0, 0}) //80-bit sample rate 44100
	if err != nil {
		panic(err)
	}

	// sound chunk
	_, err = f.WriteString("SSND")
	if err != nil {
		panic(err)
	}
	err = binary.Write(f, binary.BigEndian, int32(0)) //size
	if err != nil {
		panic(err)
	}
	err = binary.Write(f, binary.BigEndian, int32(0)) //offset
	if err != nil {
		panic(err)
	}
	err = binary.Write(f, binary.BigEndian, int32(0)) //block
	if err != nil {
		panic(err)
	}

	nSamples := 0
	defer func() {
		// fill in missing sizes
		totalBytes := 4 + 8 + 18 + 8 + 8 + 4*nSamples
		_, err = f.Seek(4, 0)
		if err != nil {
			panic(err)
		}
		err = binary.Write(f, binary.BigEndian, int32(totalBytes))
		if err != nil {
			panic(err)
		}
		_, err = f.Seek(22, 0)
		if err != nil {
			panic(err)
		}
		err = binary.Write(f, binary.BigEndian, int32(nSamples))
		if err != nil {
			panic(err)
		}
		_, err = f.Seek(42, 0)
		if err != nil {
			panic(err)
		}
		err = binary.Write(f, binary.BigEndian, int32(4*nSamples+8))
		if err != nil {
			panic(err)
		}
		f.Close()
	}()

	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	stream.Start()
	for {
		stream.Read()
		binary.Write(f, binary.BigEndian, in)
		nSamples += len(in)
		select {
		case <-sig:
			return
		default:
		}
	}

}
