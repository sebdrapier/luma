package ws

import (
	"context"
	"log"
	"strconv"
	"time"

	"elano.fr/src/backend/dmx"
)

func runShowSequence(ctx context.Context, ctrl *dmx.DMXController, show ShowPayload, showID string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in runShowSequence: %v", r)
		}
		showMu.Lock()
		if currentShow != nil && currentShow.id == showID {
			currentShow = nil
		}
		showMu.Unlock()
		broadcast <- Message{Type: "show_stopped", Payload: mustMarshal(map[string]interface{}{"show_id": showID})}
	}()

	for {
		if err := ctx.Err(); err != nil {
			return
		}

		for i, step := range show.Steps {
			select {
			case <-ctx.Done():
				return
			default:
			}

			showMu.Lock()
			if currentShow == nil || currentShow.id != showID {
				showMu.Unlock()
				return
			}
			currentShow.currentStep = i
			showMu.Unlock()

			if err := ctrl.Blackout(); err != nil {
				log.Printf("Error during blackout: %v", err)
			}

			channels := make(map[int]byte)
			for addrStr, val := range step.Preset {
				addr, err := strconv.Atoi(addrStr)
				if err != nil || addr < 1 || addr > 512 || val < 0 || val > 255 {
					continue
				}
				channels[addr] = byte(val)
			}

			if step.FadeMs > 0 {
				if err := ctrl.FadeChannels(channels, time.Duration(step.FadeMs)*time.Millisecond, dmx.FadeLinear); err != nil {
					performManualFade(ctx, ctrl, channels, step.FadeMs)
				}
			} else {
				if err := ctrl.SetChannels(channels); err != nil {
					log.Printf("Error setting channels in show step %d: %v", i, err)
				}
			}

			broadcast <- Message{
				Type:    "show_step",
				Payload: mustMarshal(map[string]interface{}{"step": i, "total": len(show.Steps), "show_id": showID}),
			}

			if step.DelayMs > 0 {
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Duration(step.DelayMs) * time.Millisecond):
				}
			}
		}

		if !show.Loop {
			break
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func performManualFade(ctx context.Context, ctrl *dmx.DMXController, targetChannels map[int]byte, fadeMs int) {
	currentChannels, err := ctrl.GetAllChannels()
	if err != nil {
		ctrl.SetChannels(targetChannels)
		return
	}
	fadeSteps := max(fadeMs/20, 1)
	for f := 0; f <= fadeSteps; f++ {
		select {
		case <-ctx.Done():
			return
		default:
		}
		fadeChannels := make(map[int]byte)
		progress := float64(f) / float64(fadeSteps)
		for addr := 1; addr <= 512; addr++ {
			currentValue := currentChannels[addr-1]
			targetValue := targetChannels[addr]
			if currentValue != targetValue {
				newValue := byte(float64(currentValue) + (float64(targetValue)-float64(currentValue))*progress)
				fadeChannels[addr] = newValue
			}
		}
		if err := ctrl.SetChannels(fadeChannels); err != nil {
			log.Printf("Error during fade: %v", err)
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
}
