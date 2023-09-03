//go:build linux
// +build linux

package wallpaper

import (
	"errors"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// Get returns the current wallpaper.
func Get() (string, error) {
	if isGNOMECompliant() {
		return parseDconf("gsettings", "get", "org.gnome.desktop.background", "picture-uri")
	}

	switch Desktop {
	case "KDE":
		return getKDE()
	case "X-Cinnamon":
		return parseDconf("dconf", "read", "/org/cinnamon/desktop/background/picture-uri")
	case "MATE":
		return parseDconf("dconf", "read", "/org/mate/desktop/background/picture-filename")
	case "XFCE":
		return getXFCE()
	case "LXDE":
		return getLXDE()
	case "Deepin":
		return parseDconf("dconf", "read", "/com/deepin/wrap/gnome/desktop/background/picture-uri")
	case "i3":
		return getI3()
	case "Hyprland":
		return getHyprland()
	default:
		return "", ErrUnsupportedDE
	}
}

// SetFromFile sets wallpaper from a file path.
func SetFromFile(file string) error {
	if isGNOMECompliant() {
		theme, err := getGnomeTheme()
		if err != nil {
			return err
		}
		if theme == "dark" {
			return exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", strconv.Quote("file://"+file)).Run()
		} else if theme == "light" {
			return exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", strconv.Quote("file://"+file)).Run()

		}
	}

	switch Desktop {
	case "KDE":
		return setKDE(file)
	case "X-Cinnamon":
		return exec.Command("dconf", "write", "/org/cinnamon/desktop/background/picture-uri", strconv.Quote("file://"+file)).Run()
	case "MATE":
		return exec.Command("dconf", "write", "/org/mate/desktop/background/picture-filename", strconv.Quote(file)).Run()
	case "XFCE":
		return setXFCE(file)
	case "LXDE":
		return exec.Command("pcmanfm", "-w", file).Run()
	case "Deepin":
		return exec.Command("dconf", "write", "/com/deepin/wrap/gnome/desktop/background/picture-uri", strconv.Quote("file://"+file)).Run()
	case "i3":
		return setI3(file)
	case "Hyprland":
		return setHyprland(file)
	default:
		err := exec.Command("swaybg", "-i", file).Start()
		// if the command completed successfully, return
		if err == nil {
			return nil
		}

		return exec.Command("feh", "-bg-fill", file).Run()
	}
}

// SetMode sets the wallpaper mode.
func SetMode(mode Mode) error {
	if isGNOMECompliant() {
		return exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-options", strconv.Quote(mode.getGNOMEString())).Run()
	}

	switch Desktop {
	case "KDE":
		return setKDEMode(mode)
	case "X-Cinnamon":
		return exec.Command("dconf", "write", "/org/cinnamon/desktop/background/picture-options", strconv.Quote(mode.getGNOMEString())).Run()
	case "MATE":
		return exec.Command("dconf", "write", "/org/mate/desktop/background/picture-options", strconv.Quote(mode.getGNOMEString())).Run()
	case "XFCE":
		return setXFCEMode(mode)
	case "LXDE":
		return exec.Command("pcmanfm", "--wallpaper-mode", mode.getLXDEString()).Run()
	case "Deepin":
		return exec.Command("dconf", "write", "/com/deepin/wrap/gnome/desktop/background/picture-options", strconv.Quote(mode.getGNOMEString())).Run()
	default:
		return ErrUnsupportedDE
	}
}

func getCacheDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".cache"), nil
}

// function to get if gnome is in light or dark mode, because in gnome 40 and higher you need to set the wallpaper separately
func getGnomeTheme() (string, error) {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	if strings.Contains(outputStr, "dark") {
		return "dark", err
	} else if strings.Contains(outputStr, "light") {
		return "light", err
	} else {
		err = errors.New("Could not get gnome theme")
		return "", err
	}
}
