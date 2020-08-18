package g2j

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// Secrets a key-value map containing confidential data
type Secrets map[string]string

const (
	jsonExt string = ".json"
)

func loadSecrets(secretsLocation string) (Secrets, error) {
	fileBytes, err := ioutil.ReadFile(secretsLocation)
	if err != nil {
		return make(map[string]string), err
	}
	return interpretSecrets(secretsLocation, fileBytes)
}

func interpretSecrets(secretsLocation string, data []byte) (Secrets, error) {
	dest := make(map[string]string)
	switch ext := path.Ext(secretsLocation); ext {
	case jsonExt:
		err := json.Unmarshal(data, &dest)
		return dest, err
	default:
		err := fmt.Errorf("unrecognised extension: %s", ext)
		return dest, err
	}
}
