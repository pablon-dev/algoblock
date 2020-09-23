package accounts

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var path string

func init() {
	path = os.Getenv("algoMnemonicPath")
}

//SaveMnemonics Save mnemonics to an external file
func SaveMnemonics(mnemonics ...string) error {
	arr, err := ReadMnemonics()
	if !os.IsNotExist(err) {
		return err
	}
	for _, x := range mnemonics {
		arr = append(arr, x)
	}

	js, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, js, 0644); err != nil {
		return err
	}
	return nil
}

//ReadMnemonics read mnemonics from external file
func ReadMnemonics() ([]string, error) {
	var arr []string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return arr, err
	}
	s, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(s, &arr); err != nil {
		return nil, err
	}
	return arr, nil
}
