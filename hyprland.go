package wallpaper

import (
	"errors"
	"os/exec"
	"strings"
)

func getHyprland() (string, error) {
	cmd := `ps aux | grep 'swaybg -i' | awk '{match($0, /-i ([^ ]+)/, arr); if (arr[1]) print arr[1]}'`
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(output), "\n", ""), nil
}

func setHyprland(file string) error {
	output, err := exec.Command("which", "swaybg").Output()
	if err != nil {
		err = errors.New("swaybg not found, please install it via package manager")
		return err
	}
	swaybgCmd := strings.ReplaceAll(string(output), "\n", "")
	err = exec.Command(swaybgCmd, "-i", file).Run()
	if err != nil {
		return err
	}
	return nil
}
