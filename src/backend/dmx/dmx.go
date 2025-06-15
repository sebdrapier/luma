package dmx

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"go.bug.st/serial"
)

const (
	DMXBaudRate     = 250000                 // DMX512 baud rate
	DMXChannels     = 512                    // Number of DMX channels
	DMXStartCode    = 0                      // Start code (0 for standard DMX)
	DMXFrameSize    = DMXChannels + 1        // Total frame size including start code
	DMXBreakTime    = 100 * time.Microsecond // Break signal duration (with margin)
	DMXMaBTime      = 10 * time.Microsecond  // Mark After Break duration (with margin)
	DMXMinFrameRate = 40 * time.Millisecond  // Minimum time between frames (25 fps)
	DMXMaxFrameRate = 1 * time.Millisecond   // Maximum frame rate (1000 fps)
)

type FadeMode int

const (
	FadeLinear FadeMode = iota
	FadeQuadratic
	FadeCubic
	FadeSine
	FadeExponential
)

type ChannelLimit struct {
	Min byte
	Max byte
}

type DMXController struct {
	port serial.Port
	data [DMXFrameSize]byte

	mu          sync.RWMutex
	flushRate   time.Duration
	stopSender  chan struct{}
	dataChanged chan struct{}

	masterDimmer  float64
	channelLimits map[int]*ChannelLimit

	fadeMu     sync.Mutex
	fadeCancel context.CancelFunc

	frameCount atomic.Uint64
	errorCount atomic.Uint64

	closed atomic.Bool
}

func NewDMXController(portName string) (*DMXController, error) {
	if portName == "" {
		return nil, fmt.Errorf("port name cannot be empty")
	}

	mode := &serial.Mode{
		BaudRate: DMXBaudRate,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.TwoStopBits,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open port %s: %w", portName, err)
	}

	d := &DMXController{
		port:          port,
		flushRate:     DMXMinFrameRate,
		stopSender:    make(chan struct{}),
		dataChanged:   make(chan struct{}, 2),
		masterDimmer:  1.0,
		channelLimits: make(map[int]*ChannelLimit),
	}

	d.data[0] = DMXStartCode
	for i := 1; i <= DMXChannels; i++ {
		d.data[i] = 0
	}

	go d.continuousSender()

	return d, nil
}

func (d *DMXController) continuousSender() {
	ticker := time.NewTicker(d.flushRate)
	defer ticker.Stop()

	for {
		select {
		case <-d.stopSender:
			return
		case <-d.dataChanged:
			d.sendFrame()
			ticker.Reset(d.flushRate)
		case <-ticker.C:
			d.sendFrame()
		}
	}
}

func (d *DMXController) sendFrame() {
	if d.closed.Load() {
		return
	}

	d.frameCount.Add(1)

	if err := d.port.Break(DMXBreakTime); err != nil {
		d.errorCount.Add(1)
		return
	}
	time.Sleep(DMXMaBTime)

	d.mu.RLock()
	frame := d.data
	master := d.masterDimmer
	limits := d.channelLimits
	d.mu.RUnlock()

	for i := 1; i <= DMXChannels; i++ {
		v := float64(frame[i]) * master
		b := byte(v)
		if lim, ok := limits[i]; ok {
			if b < lim.Min {
				b = lim.Min
			} else if b > lim.Max {
				b = lim.Max
			}
		}
		frame[i] = b
	}

	if _, err := d.port.Write(frame[:]); err != nil {
		d.errorCount.Add(1)
		return
	}
	d.port.Drain()
}

func (d *DMXController) SetChannel(ch int, value byte) error {
	if ch < 1 || ch > DMXChannels {
		return fmt.Errorf("channel must be 1-%d, got %d", DMXChannels, ch)
	}
	if d.closed.Load() {
		return fmt.Errorf("controller is closed")
	}

	d.mu.Lock()
	d.data[ch] = value
	d.mu.Unlock()

	d.signalChange()
	return nil
}

func (d *DMXController) SetChannels(vals map[int]byte) error {
	if d.closed.Load() {
		return fmt.Errorf("controller is closed")
	}

	d.mu.Lock()
	for ch, v := range vals {
		if ch < 1 || ch > DMXChannels {
			d.mu.Unlock()
			return fmt.Errorf("channel must be 1-%d, got %d", DMXChannels, ch)
		}
		d.data[ch] = v
	}
	d.mu.Unlock()

	d.signalChange()
	return nil
}

func (d *DMXController) FadeChannels(targets map[int]byte, duration time.Duration, mode FadeMode) error {
	if d.closed.Load() {
		return fmt.Errorf("controller is closed")
	}

	d.fadeMu.Lock()
	if d.fadeCancel != nil {
		d.fadeCancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	d.fadeCancel = cancel
	d.fadeMu.Unlock()

	d.mu.RLock()
	start := make(map[int]byte, len(targets))
	for ch := range targets {
		start[ch] = d.data[ch]
	}
	d.mu.RUnlock()

	go func() {
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()
		startTime := time.Now()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				elapsed := time.Since(startTime)
				f := float64(elapsed) / float64(duration)
				if f >= 1 {
					d.SetChannels(targets)
					return
				}
				adj := applyFadeCurve(f, mode)
				curr := make(map[int]byte, len(targets))
				for ch, tgt := range targets {
					s := float64(start[ch])
					curr[ch] = byte(s + (float64(tgt)-s)*adj)
				}
				d.SetChannels(curr)
			}
		}
	}()

	return nil
}

func (d *DMXController) SetMasterDimmer(val float64) error {
	if val < 0 || val > 1 {
		return fmt.Errorf("master dimmer must be 0.0-1.0, got %f", val)
	}
	d.mu.Lock()
	d.masterDimmer = val
	d.mu.Unlock()
	d.signalChange()
	return nil
}

func (d *DMXController) SetChannelLimit(ch int, min, max byte) error {
	if ch < 1 || ch > DMXChannels {
		return fmt.Errorf("channel must be 1-%d, got %d", DMXChannels, ch)
	}
	if min > max {
		return fmt.Errorf("min %d cannot exceed max %d", min, max)
	}
	d.mu.Lock()
	d.channelLimits[ch] = &ChannelLimit{Min: min, Max: max}
	d.mu.Unlock()
	return nil
}

func (d *DMXController) RemoveChannelLimit(ch int) error {
	if ch < 1 || ch > DMXChannels {
		return fmt.Errorf("channel must be 1-%d, got %d", DMXChannels, ch)
	}
	d.mu.Lock()
	delete(d.channelLimits, ch)
	d.mu.Unlock()
	return nil
}

func (d *DMXController) Blackout() error {
	if d.closed.Load() {
		return fmt.Errorf("controller is closed")
	}
	d.fadeMu.Lock()
	if d.fadeCancel != nil {
		d.fadeCancel()
	}
	d.fadeMu.Unlock()

	d.mu.Lock()
	for i := 1; i <= DMXChannels; i++ {
		d.data[i] = 0
	}
	d.mu.Unlock()

	d.signalChange()
	return nil
}

func (d *DMXController) BlackoutWithFade(duration time.Duration, mode FadeMode) error {
	targets := make(map[int]byte, DMXChannels)
	for i := 1; i <= DMXChannels; i++ {
		targets[i] = 0
	}
	return d.FadeChannels(targets, duration, mode)
}

func (d *DMXController) GetStatistics() (frames, errors uint64) {
	return d.frameCount.Load(), d.errorCount.Load()
}

func (d *DMXController) ResetStatistics() {
	d.frameCount.Store(0)
	d.errorCount.Store(0)
}

func (d *DMXController) Close() error {
	if d.closed.Swap(true) {
		return nil
	}
	d.fadeMu.Lock()
	if d.fadeCancel != nil {
		d.fadeCancel()
	}
	d.fadeMu.Unlock()

	close(d.stopSender)

	d.mu.Lock()
	for i := 1; i <= DMXChannels; i++ {
		d.data[i] = 0
	}
	d.mu.Unlock()
	d.port.Write(d.data[:])
	d.port.Drain()

	return d.port.Close()
}

func (d *DMXController) GetAllChannels() ([]byte, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("controller is closed")
	}
	d.mu.RLock()
	defer d.mu.RUnlock()

	out := make([]byte, DMXChannels)
	copy(out, d.data[1:])
	return out, nil
}

func (d *DMXController) signalChange() {
	select {
	case d.dataChanged <- struct{}{}:
	default:
	}
}

func applyFadeCurve(progress float64, mode FadeMode) float64 {
	switch mode {
	case FadeLinear:
		return progress
	case FadeQuadratic:
		return progress * progress
	case FadeCubic:
		return progress * progress * progress
	case FadeSine:
		return (1 - math.Cos(progress*math.Pi)) / 2
	case FadeExponential:
		return math.Pow(2, 10*(progress-1))
	default:
		return progress
	}
}
