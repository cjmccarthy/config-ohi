package main

import (
	"encoding/json"
	"fmt"
	"github.com/clbanning/mxj/x2j"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	//"reflect"
)

type InfraPayload struct {
	Name               string                 `json:"name"`
	ProtocolVersion    string                 `json:"protocol_version"`
	IntegrationVersion string                 `json:"integration_version"`
	Inventory          map[string]interface{} `json:"inventory"`
}

type FileDef struct {
	Inventory string
	Path      string
	Type      string
}

func check(e error) {
	if e != nil {
		fmt.Printf("err: %v\n", e)
		panic(e)
	}
}

func extractConfs(fileDef FileDef) interface{} {

	fdata, err := ioutil.ReadFile(fileDef.Path)
	check(err)

	switch ftype := fileDef.Type; ftype {
	case "yaml":
		invdata, err := yaml.YAMLToJSON(fdata)
		check(err)
		var jdata interface{}
		err = json.Unmarshal(invdata, &jdata)
		check(err)
		return jdata
	case "json":
		/*
			Array vs Node top level object?

			rt := reflect.TypeOf(jdata)
			if rt.Kind() == reflect.Slice {

			}
		*/
		var jdata interface{}
		err = json.Unmarshal(fdata, &jdata)
		check(err)
		return jdata
	case "xml":
		invdata, err := x2j.XmlToJson(fdata)
		check(err)
		var jdata interface{}
		err = json.Unmarshal(invdata, &jdata)
		return jdata
	case "properties":
	default:
		fmt.Printf("Unsupported file adapter type %s:", ftype)
	}
	return nil
}

func main() {
	args := os.Args[1:]
	confData, err := ioutil.ReadFile(args[0])
	check(err)

	confs := []FileDef{}

	err = yaml.Unmarshal([]byte(confData), &confs)
	check(err)

	var payload InfraPayload
	payload.Name = "haus.chris.nringest"
	payload.ProtocolVersion = "1"
	payload.IntegrationVersion = "1.0.0"
	payload.Inventory = make(map[string]interface{})

	for _, fileDef := range confs {
		confs := extractConfs(fileDef)
		payload.Inventory[fileDef.Inventory] = confs
	}

	payloadJson, err := json.Marshal(payload)
	check(err)

	fmt.Printf(string(payloadJson))

}
