package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func splitDirAndNameAndExt(path string) (string, string, string) {
	dir, fileNameOnly := filepath.Split(path)
	ext := filepath.Ext(fileNameOnly)
	name := strings.TrimSuffix(fileNameOnly, filepath.Ext(fileNameOnly))
	return dir, name, ext
}

func visit(files *[]string, allowFilePrefix string, allowFileExtsion string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() == false {
			_, name, ext := splitDirAndNameAndExt(path)
			//print("dir, file, ext ", dir, name, ext, "\n")
			if strings.HasPrefix(name, allowFilePrefix) && ext == allowFileExtsion {
				*files = append(*files, path)
			}
		}
		return nil
	}
}

func main() {
	var files []string
	allowPrefix := "demo_tc_"
	allowExt := ".xaml"
	configExt := ".config.json"
	baseLogDir := "C:/temp/uitest-logs/"
	logDir := baseLogDir + time.Now().Format("20060102-150405")

	exepath := "C:/apps/uipath/app-20.4.1/uirobot.exe"
	root := "C:/data/projects_uipath/mm-ui-tests"

	err := filepath.Walk(root, visit(&files, allowPrefix, allowExt))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		// fmt.Println(file)
		dir, name, _ := splitDirAndNameAndExt(file)
		jsonFileName := dir + name + configExt
		// fmt.Println(jsonFileName)
		jsonFile, err := os.Open(jsonFileName)
		if err != nil {
			fmt.Println(err)
		}
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println(err)
		}
		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(byteValue), &jsonMap)
		jsonMap["logDirectory"] = logDir
		fmt.Println("\n", jsonMap["cityName"])
		jsonString, err := json.Marshal(jsonMap)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(exepath)
		fmt.Println(string(jsonString))

		cmd := exec.Command(exepath, "execute", "--file", file, "--input", string(jsonString))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err1 := cmd.Run()
		if err1 != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err1)
		}
	}
}
