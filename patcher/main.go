package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	VSLENGTH         = 31
	ORIGINAL_VERSION = "Air Brawl 0.701"
	//VSSERVER         = "http://localhost:7448/version" // Server which we get our version string from
	VSSERVER = "http://home.jonas747.com:7448/version" // Server which we get our version string from
)

var (
	teamCode string
)

func main() {
	fmt.Println("Starting abpatcher for airbrawl tournament!")
	findDLLPath()
	fmt.Println("Set DLL Path to ", DLLPATH)
	actionLoop()
}

func getTeamCode() string {
	fmt.Println("Please enter your teamcode that you received from the admins (case sensitive)")
	var tc string
	fmt.Scanln(&tc)
	return tc
}

func actionLoop() {
	for {
		fmt.Println("=======================")
		fmt.Println("What do you want to do?")
		fmt.Println("(1) Restore original version")
		fmt.Println("(2) Patch for tournament usage")
		var action int
		fmt.Scanln(&action)
		switch action {
		case 1:
			fmt.Println("Restoring original version...")
			actionRestore()
		case 2:
			fmt.Println("Patching for tournament usage...")
			actionPatch()
		default:
			fmt.Println("Invalid option, returning to main menu")
			continue
		}
	}
}

func actionRestore() {
	applyVersion(ORIGINAL_VERSION)
}

func actionPatch() {
	var version string
	if teamCode == "" {
		// Get the teamcode from user input
		for {
			code := getTeamCode()
			v, err := getVersionFromServer(code)
			if err != nil {
				fmt.Println("Error getting version from server: ", err)
				continue
			}
			version = v
			teamCode = code
			break
		}
	} else {
		v, err := getVersionFromServer(teamCode)
		if err != nil {
			fmt.Println("Error getting version from server: ", err)
			return
		}
		version = v
	}
	if version != "" {
		applyVersion(version)
	}
}

func applyVersion(version string) {
	// Check if Assembly-CSharp.dll is where it's supposed to be
	info, err := os.Stat(DLLPATH)
	if err != nil {
		fmt.Println("Error applying version: ", err)
		fmt.Println("Are you sure you put this next to the airbrawl executeable?")
		return
	}
	if info.IsDir() {
		fmt.Println("This really should'nt happen, the \"Assembly-CSharp.dll\" is a directory?")
		fmt.Println("Are you sure you put this next to the airbrawl executable?")
		return
	}

	perm := info.Mode()

	// open the file
	file, err := os.OpenFile(DLLPATH, os.O_WRONLY, perm)
	defer file.Close()
	if err != nil {
		fmt.Println("Error applying version: ", err)
		fmt.Println("Are you sure you put this next to the airbrawl executeable?")
		return
	}

	encoded := encodeVSString(version)
	n, err := file.WriteAt(encoded, VSSTRING_OFFSET)
	if n < 31 {
		fmt.Println("Didn't write full 31 bytes!? (", n, ")")
		return
	}
	if err != nil {
		fmt.Println("Error applying version: ", err)
		fmt.Println("Are you sure you put this next to the airbrawl executeable?")
		return
	}
	fmt.Println("Successfully applied version: ", version)
	fmt.Println("You can now (re)start airbrawl!")
}

type Response struct {
	Error   string
	Version string
}

func getVersionFromServer(teamCode string) (string, error) {
	url := VSSERVER + "?tc=" + teamCode
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var decoded Response
	err = json.Unmarshal(body, &decoded)
	if err != nil {
		fmt.Println(string(body))
		return "", err
	}

	if decoded.Error != "" {
		return "", errors.New(decoded.Error)
	}

	return decoded.Version, nil
}

func encodeVSString(in string) []byte {
	out := make([]byte, VSLENGTH)
	for i := 0; i < VSLENGTH; i++ {
		if i >= len(in)*2 {
			break
		}
		if i%2 == 0 && i < 30 {
			sIndex := i / 2
			out[i] = in[sIndex]
			fmt.Printf("@ index[%d]: 0x%X\n", i, in[sIndex])
		} else {
			continue
		}
	}
	return out
}
