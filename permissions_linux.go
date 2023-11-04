//go:build linux && !android

package anping

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func getPermission() error {
	// root
	if uid := os.Geteuid(); uid == 0 {
		return nil
	}

	// cap
	path, err := os.Executable()
	if err != nil {
		return err
	}

	cap, err := getCap(path)
	if err != nil {
		return err
	}

	// 用正则表达式检索
	pattern := "cap_net_raw=[a-zA-Z]+"
	re := regexp.MustCompile(pattern)
	match := re.FindString(cap)
	// log.Println(match)
	equalIndex := strings.Index(match, "=")

	if equalIndex != -1 && equalIndex+1 < len(match) {
		// 获取等于号后面的子字符串
		afterEqual := match[equalIndex+1:]

		// 使用 strings.Contains 函数检查子字符串是否包含 "e" 和 "p"
		if strings.Contains(afterEqual, "e") && strings.Contains(afterEqual, "p") {
			return nil
		}
	}

	err = setCap("cap_net_raw=+ep", path)
	if err != nil {
		return err
	}

	return nil
}

func getCap(path string) (string, error) {
	cmd := exec.Command("getcap", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func setCap(cap, path string) error {
	log.Println("Requiring cap...")

	cmd := exec.Command("pkexec", "setcap", cap, path)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
