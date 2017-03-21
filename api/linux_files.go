/*
 * Github Authorized Keys - Use GitHub teams to manage system user accounts and authorized_keys
 *
 * Copyright 2016 Cloud Posse, LLC <hello@cloudposse.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/goruha/permbits"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type operationOnFileContent func(string) error

// FileExists - check if file exists
func (linux *Linux) FileExists(filePath string) bool {
	file, err := os.Open(linux.applyChroot(filePath))
	defer file.Close()
	return err == nil
}

// FileCreate - creates new file
func (linux *Linux) FileCreate(filePath string) error {
	if !linux.FileExists(filePath) {
		file, err := os.Create(linux.applyChroot(filePath))
		defer file.Close()
		return err
	}
	return os.ErrExist
}

// FileDelete - removes file
func (linux *Linux) FileDelete(filePath string) error {
	if linux.FileExists(filePath) {
		return os.Remove(linux.applyChroot(filePath))
	}
	return os.ErrNotExist
}

// FileEnsure - If file not exists create it file, update content in that file if it differs.
func (linux *Linux) FileEnsure(filePath, content string) error {
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

// FileGet - return file content
func (linux *Linux) FileGet(filePath string) (string, error) {
	if linux.FileExists(filePath) {
		buffer, err := ioutil.ReadFile(linux.applyChroot(filePath))
		content := string(buffer)
		return content, err
	}
	return "", os.ErrNotExist
}

// FileSet - set file content
func (linux *Linux) FileSet(filePath, content string) error {
	return ioutil.WriteFile(linux.applyChroot(filePath), []byte(content), 0777)
}

// FileEnsureLine - add line to file if it is not present in file content
func (linux *Linux) FileEnsureLine(filePath string, line string) error {
	return linux.FileEnsureLineMatch(filePath, "^"+line+"$", line)
}

// FileEnsureLineMatch - add line to file if there is no line that match regexp or that line differs
func (linux *Linux) FileEnsureLineMatch(filePath, matcher, line string) error {
	logger := log.WithFields(log.Fields{"class": "Linux", "method": "FileEnsureLineMatch"})

	return linux.doOnFileContent(func(fileContent string) error {
		re := regexp.MustCompile(matcher)

		matchedStrings := re.FindAllString(fileContent, -1)
		logger.Debugf("Matcher  %v", matcher)
		logger.Debugf("Matched strings count %v", len(matchedStrings))
		if len(matchedStrings) > 1 {
			return fmt.Errorf("Match regexp /%v/ is too wide - %v matches found", matcher, matchedStrings)
		}

		matchedString := ""
		if len(matchedStrings) == 1 {
			matchedString = matchedStrings[0]
		}

		if matchedString == "" {
			logger.Debugf("File %v does not contain target string", filePath)
			return linux.FileSet(filePath, fileContent+"\n"+line)
		} else if matchedString != line {
			newFileContent := re.ReplaceAllLiteralString(fileContent, line)
			return linux.FileSet(filePath, newFileContent)
		}
		logger.Debugf("File %v contains target string", filePath)
		return nil

	}, filePath, logger)
}

func (linux *Linux) doOnFileContent(f operationOnFileContent, filePath string, logger *log.Entry) (finalError error) {
	if linux.FileExists(filePath) {
		fileContent, err := linux.FileGet(filePath)
		if err == nil {
			finalError = f(fileContent)
		} else {
			finalError = err
			logger.Debugf("Can not read file %v", filePath)
		}
	} else {
		logger.Debugf("File %v not fould", filePath)
		finalError = os.ErrNotExist
	}

	return
}

// FileModeGet - get file permission mode
func (linux *Linux) FileModeGet(filePath string) (permbits.PermissionBits, error) {
	return permbits.Stat(linux.applyChroot(filePath))
}

// FileModeSet - set file permission mode
func (linux *Linux) FileModeSet(filePath string, mode permbits.PermissionBits) error {
	return permbits.Chmod(linux.applyChroot(filePath), mode)
}

// FileModeEnsure - set file permission mode if differs
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
	}
	return linux.root + path
}
