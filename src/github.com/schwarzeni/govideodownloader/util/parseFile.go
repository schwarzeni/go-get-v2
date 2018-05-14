package util

import (
	"io/ioutil"

	"log"
	"os"

	"encoding/json"

	"github.com/schwarzeni/govideodownloader/model"
)

func ParseConfigFile() model.Config {
	file, e := ioutil.ReadFile("./data.json")
	if e != nil {
		log.Fatal(e)
		os.Exit(1)
	}

	var jsontype model.Config
	json.Unmarshal(file, &jsontype)
	return jsontype
}
