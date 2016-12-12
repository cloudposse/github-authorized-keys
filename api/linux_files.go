package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/goruha/permbits"
	"io/ioutil"
	"os"
	"strings"
)

func (linux *Linux) FileExists(filePath string) bool {
	file, err := os.Open(linux.applyChroot(filePath))
	defer file.Close()
	return err == nil
}

func (linux *Linux) FileCreate(filePath string) error {
	if !linux.FileExists(filePath) {
		file, err := os.Create(linux.applyChroot(filePath))
		defer file.Close()
		return err
	}
	return os.ErrExist
}

func (linux *Linux) FileDelete(filePath string) error {
	if linux.FileExists(filePath) {
		return os.Remove(linux.applyChroot(filePath))
	}
	return os.ErrNotExist
}

func (linux *Linux) FileEnsure(filePath string, content string) error {
	logger := log.WithFields(log.Fields{"class": "Linux", "method": "FileEnsure"})

	if !linux.FileExists(filePath) {
		logger.Debugf("File %v not found", filePath)
		linux.FileCreate(filePath)
	}
	fileContent, err := linux.FileGet(filePath)
	if err == nil {
		if strings.Compare(fileContent, content) != 0 {
			logger.Debugf("File %v content differs from expected", filePath)
			return linux.FileSet(filePath, content)
		}
	} else {
		logger.Debugf("Can not read file %v", filePath)
	}
	return err
}

func (linux *Linux) FileGet(filePath string) (string, error) {
	if linux.FileExists(filePath) {
		buffer, err := ioutil.ReadFile(linux.applyChroot(filePath))
		content := string(buffer)
		return content, err
	}
	return "", os.ErrNotExist
}

func (linux *Linux) FileSet(filePath string, content string) error {
	return ioutil.WriteFile(linux.applyChroot(filePath), []byte(content), 0777)
}

func (linux *Linux) FileEnsureLine(filePath string, line string) (err error) {
	logger := log.WithFields(log.Fields{"class": "Linux", "method": "FileEnsureLine"})
	if linux.FileExists(filePath) {
		fileContent, err := linux.FileGet(filePath)
		if err == nil {
			if !strings.Contains(fileContent, line) {
				logger.Debugf("File %v miss target string", filePath)
				err = linux.FileSet(filePath, fileContent+"\n"+line)
			} else {
				logger.Debugf("File %v contains target string", filePath)
			}
		} else {
			logger.Debugf("Can not read file %v", filePath)
		}
	} else {
		logger.Debugf("File %v not fould", filePath)
		err = os.ErrNotExist
	}

	return
}

func (linux *Linux) FileModeGet(filePath string) (permbits.PermissionBits, error) {
	return permbits.Stat(linux.applyChroot(filePath))
}

func (linux *Linux) FileModeSet(filePath string, mode permbits.PermissionBits) error {
	return permbits.Chmod(linux.applyChroot(filePath), mode)
}

func (linux *Linux) FileModeEnsure(filePath string, mode permbits.PermissionBits) error {
	logger := log.WithFields(log.Fields{"class": "Linux", "method": "FileModeEnsure"})
	currentMode, err := linux.FileModeGet(filePath)
	if err == nil {
		linux.FileModeSet(filePath, currentMode|mode)
	} else {
		logger.Debugf("Cannot get permissions of file %v", filePath)
	}
	return err
}

func (linux *Linux) applyChroot(path string) string {
	if linux.root == "/" {
		return path
	} else {
		return linux.root + path
	}
}
