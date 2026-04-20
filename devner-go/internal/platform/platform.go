package platform

import (
	"os"
	"runtime"
)

type OS int

const (
	Unknown OS = iota
	MacOS
	Linux
	Windows
)

func Detect() OS {
	switch runtime.GOOS {
	case "darwin":
		return MacOS
	case "linux":
		return Linux
	case "windows":
		return Windows
	default:
		return Unknown
	}
}

func (o OS) String() string {
	switch o {
	case MacOS:
		return "macOS"
	case Linux:
		return "linux"
	case Windows:
		return "windows"
	default:
		return "unknown"
	}
}

func IsElevated() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	return os.Geteuid() == 0
}

func HostsFilePath() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}
