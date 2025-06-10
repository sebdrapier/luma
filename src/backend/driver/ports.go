package driver

import (
	"errors"
	"runtime"

	"strings"

	"go.bug.st/serial"
)

func ListDMXPorts() ([]string, error) {
	ports, err := serial.GetPortsList()

	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, errors.New("no serial ports found")
	}

	var result []string
	for _, port := range ports {
		if isLikelyDMXPort(port) {
			result = append(result, port)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no DMX USB interfaces found")
	}
	return result, nil
}

func isLikelyDMXPort(port string) bool {
	switch runtime.GOOS {
	case "windows":
		return strings.HasPrefix(strings.ToUpper(port), "COM")
	case "darwin":
		return strings.Contains(port, "usb")
	case "linux":
		return strings.HasPrefix(port, "/dev/ttyUSB") || strings.HasPrefix(port, "/dev/ttyACM")
	default:
		return false
	}
}
