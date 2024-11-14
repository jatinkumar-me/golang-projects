package main

import (
	"fmt"
	"math"
	"sync"

	"github.com/eiannone/keyboard"
	"github.com/gordonklaus/portaudio"
)

const (
	sampleRate      = 44100
	framesPerBuffer = 64
	duration        = 0.3
)

var frequencies = map[rune]float64{
	'a': 261.63,
	's': 293.66,
	'd': 329.63,
	'f': 349.23,
	'g': 392.00,
	'h': 440.00,
	'j': 493.88,
	'k': 523.25,
	'l': 587.33,
	';': 659.25,
}

type AudioData struct {
	buffer []float32
	index  int
	mutex  sync.Mutex
}

func generateSineWave(freq float64, sampleRate int, duration float64) []float32 {
	samples := int(float64(sampleRate) * duration)
	wave := make([]float32, samples)
	for i := 0; i < samples; i++ {
		wave[i] = float32(math.Sin(2 * math.Pi * freq * float64(i) / float64(sampleRate)))
	}
	return wave
}

func generateSquareWave(freq, sampleRate, duration float64) []float32 {
	samples := int(sampleRate * duration)
	wave := make([]float32, samples)
	for i := 0; i < samples; i++ {
		if math.Sin(2*math.Pi*freq*float64(i)/sampleRate) >= 0 {
			wave[i] = 1
		} else {
			wave[i] = -1
		}
	}
	return wave
}

func generateSawtoothWave(freq, sampleRate, duration float64) []float32 {
	samples := int(sampleRate * duration)
	wave := make([]float32, samples)
	for i := 0; i < samples; i++ {
		wave[i] = float32(
			2 * (float64(i)*freq/sampleRate - math.Floor(0.5+float64(i)*freq/sampleRate)),
		)
	}
	return wave
}

func generateTriangleWave(freq, sampleRate, duration float64) []float32 {
	samples := int(sampleRate * duration)
	wave := make([]float32, samples)
	for i := 0; i < samples; i++ {
		wave[i] = float32(
			2*math.Abs(
				2*(float64(i)*freq/sampleRate-math.Floor(0.5+float64(i)*freq/sampleRate)),
			) - 1,
		)
	}
	return wave
}

func (data *AudioData) loadWave(wave []float32) {
	data.mutex.Lock()
	defer data.mutex.Unlock()
	data.buffer = wave
	data.index = 0
}

func (data *AudioData) playCallback(out []float32) {
	data.mutex.Lock()
	defer data.mutex.Unlock()

	for i := range out {
		if data.index < len(data.buffer) {
			out[i] = data.buffer[data.index]
			data.index++
		} else {
			out[i] = 0
		}
	}
}

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	if err := keyboard.Open(); err != nil {
		fmt.Println("Error initializing keyboard:", err)
		return
	}
	defer keyboard.Close()

	audioData := &AudioData{}

	outStream, err := portaudio.OpenDefaultStream(
		0,
		1,
		sampleRate,
		framesPerBuffer,
		audioData.playCallback,
	)
	if err != nil {
		fmt.Println("Error opening output stream:", err)
		return
	}
	defer outStream.Close()

	// Start the output stream
	if err := outStream.Start(); err != nil {
		fmt.Println("Error starting stream:", err)
		return
	}

	fmt.Println(
		"Simple Piano: Press 'A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L', or ';' to play notes. Press 'ESC' to exit.",
	)

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error reading key:", err)
			break
		}

		if key == keyboard.KeyEsc {
			fmt.Println("Exiting...")
			break
		}

		if freq, ok := frequencies[char]; ok {
			wave := generateSquareWave(freq, sampleRate, duration)

			audioData.loadWave(wave)
		}
	}

	if err := outStream.Stop(); err != nil {
		fmt.Println("Error stopping stream:", err)
		return
	}

	fmt.Println("Program exited.")
}
