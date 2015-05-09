package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// +build linux windows

const (
	VSSTRING_OFFSET = 0x50cb4

//	DLLPATH         = "Managed/Assembly-CSharp.dll"
)

var DLLPATH string

func findDLLPath() {
	// Sets the dll path
	dirContents, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}

	// Filter the directories
	dirs := make([]os.FileInfo, 0)
	for _, file := range dirContents {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
	}
	if len(dirs) == 0 {
		fmt.Println("Cant find airbrawl data folder, are you sure you but this patcher in the right place?")
		return
	} else if len(dirs) == 1 {
		// Simple, the directory is then the data folder
		DLLPATH = dirs[0].Name() + "/Managed/Assembly-CSharp.dll"
	} else {
		// A little more complicated then... so we prompt the user for which directory is the datafolder
		for index, item := range dirs {
			fmt.Printf("[%d]: %s\n", index+1, item.Name())
		}
		for {
			fmt.Println("What folder is the airbrawl data folder? (enter a number)")
			var folder int
			fmt.Scanln(&folder)
			if folder == 0 {
				continue
			}
			index := folder - 1
			if index >= len(dirs) {
				continue
			}
			dir := dirs[index]
			DLLPATH = dir.Name() + "/Managed/Assembly-CSharp.dll"
			break
		}
	}
}
