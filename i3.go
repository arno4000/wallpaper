package wallpaper

import (
	"errors"
	"os/exec"
	"strings"
)

func getI3() (string, error) {
	cmd := `tail -n 1 ~/.fehbg | cut -d" " -f 4 | sed -s "s/'//g"`
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(output), "\n", ""), nil
}

func setI3(file string) error {
	output, err := exec.Command("which", "feh").Output()
	if err != nil {
		err = errors.New("feh not found, please install it via package manager")
		return err
	}
	fehCmd := strings.ReplaceAll(string(output), "\n", "")
	err = exec.Command(fehCmd, "--bg-fill", file).Run()
	if err != nil {
		return err
	}
	return nil
}
