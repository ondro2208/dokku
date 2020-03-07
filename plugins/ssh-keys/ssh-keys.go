package sshkeys

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func ListSshKeys() (keyNamesList []string) {
	cmd := exec.Command("bash", "-c", "dokku ssh-keys:list")
	output, err := cmd.Output()
	if err != nil {
		log.Println("ListSshKeys error")
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		keyNamesList = append(keyNamesList, parseSshKeyName(scanner.Text()))
	}
	return keyNamesList
}

func AddSshKey(sshKeyName string, sshKeyValue string) (string, error) {
	cmd := exec.Command("bash", "-c", "echo \""+sshKeyValue+"\" | sshcommand acl-add dokku "+sshKeyName)
	output, err := cmd.Output()
	if err != nil {
		log.Println("AddSshKey error")
		return "", errors.New("Something goes wrong")
	}

	return string(output), nil
}

func RemoveSshKey(sshKeyName string) (string, error) {
	cmd := exec.Command("bash", "-c", "sshcommand acl-remove dokku "+sshKeyName)
	output, err := cmd.Output()
	if err != nil {
		log.Println("RemoveSshKey error")
		return "", errors.New("Something goes wrong")
	}

	return string(output), nil
}

func parseSshKeyName(value string) string {
	re := regexp.MustCompile(`NAME=\".{0,20}\"`)
	tmp := re.FindString(value)
	tmp = strings.TrimPrefix(tmp, "NAME=\"")
	tmp = strings.TrimSuffix(tmp, "\"")
	return tmp
}
