package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const debug = true

func findClosestMatch(names []string, value string) string {
	for _, name := range names {
		if strings.Contains(name, value) {
			return name
		}
	}
	return ""
}

func readValue(k registry.Key, name string) ([]byte, error) {
	b, _, err := k.GetBinaryValue(name)
	return b, err
}

func unmarshalJSON(b []byte, v interface{}) error {
	jsonb := bytes.Trim(b, "\x00")
	return json.Unmarshal(jsonb, v)
}

func read(k registry.Key, valueNames []string, name string, result interface{}) error {
	globalDataValueName := findClosestMatch(valueNames, name)

	b, err := readValue(k, globalDataValueName)
	if err != nil {
		return err
	}

	if debug {
		err = ioutil.WriteFile(fmt.Sprintf("%s.txt", globalDataValueName), b, 0644)
		if err != nil {
			return err
		}
	}

	return unmarshalJSON(b, result)
}

func main() {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Eleventh Hour Games\Last Epoch`, registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}

	defer k.Close()

	valueNames, err := k.ReadValueNames(-1)
	if err != nil {
		log.Fatal(err)
	}

	{
		var s struct {
			StashList []struct {
				StashType              int32 // currently 0
				SoloChallengeStashName string
				Gold                   int32
			}
		}
		if err := read(k, valueNames, "Epoch_Local_Global_Data_Beta", &s); err != nil {
			log.Fatal(err)
		}

		for _, stash := range s.StashList {
			if stash.StashType == 0 && stash.SoloChallengeStashName == "" {
				fmt.Println("Gold", stash.Gold)
			}
		}
	}

	{
		var s struct {
			CharacterName string
			Level         int32
		}
		if err := read(k, valueNames, "CHARACTERSLOT_BETA", &s); err != nil {
			log.Fatal(err)
		}

		if s.CharacterName != "" {
			fmt.Printf("%s (Level %d)\n", s.CharacterName, s.Level)
		}
	}
}
