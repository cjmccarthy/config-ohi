package main

import (
	"encoding/json"
	"fmt"
	//"github.com/clbanning/mxj/x2j"
	"github.com/clbanning/mxj"
	"github.com/ghodss/yaml"
	"io/ioutil"
	//	"os"
	"reflect"
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

func sanitizeLeaf(leaf interface{}) string {
	/*
		fmt.Printf("leaf: %s\n", leaf)
		fmt.Printf("leaf: %s\n", reflect.TypeOf(leaf))
			stuff := make([]map[string]interface{}, 2)
			stuff1 := make(map[string]interface{})
			stuff2 := make(map[string]interface{})
			stuff1["qwer"] = "asdf"
			stuff2["uiop"] = "poui"
			stuff[0] = stuff1
			stuff[1] = stuff2
			fmt.Printf("stuff: %s\n", stuff)
			fmt.Printf("stuff: %s\n", reflect.TypeOf(stuff))
	*/
	marshalled, err := json.Marshal(leaf)
	check(err)
	/*
		fmt.Printf("marshalled: %s\n", marshalled)
		fmt.Printf("sanit: %s\n", string(marshalled))
	*/
	return string(marshalled)
}

func sanitizeMap(jdata interface{}) interface{} {
	jmap := reflect.ValueOf(jdata)
	sanitized := make(map[string]interface{})
	for _, key := range jmap.MapKeys() {
		jval := jmap.MapIndex(key)
		sanitized[key.String()] = sanitizeLeaf(jval.Elem().Interface())
	}
	return sanitized
}

func sanitizeSlice(jdata interface{}) interface{} {
	jslice := reflect.ValueOf(jdata)
	sanitized := make([]interface{}, jslice.Len())
	for i := 0; i < jslice.Len(); i++ {
		sanitized[i] = sanitizeLeaf(jslice.Index(i).Interface())
	}
	return sanitized
}

func extractConfs(fileDef FileDef) interface{} {

	fdata, err := ioutil.ReadFile(fileDef.Path)
	check(err)

	var invdata []byte
	switch ftype := fileDef.Type; ftype {
	case "yaml":
		invdata, err = yaml.YAMLToJSON(fdata)
		check(err)
	case "json":
		invdata = fdata
	case "xml":
		/*
			invdata, err := x2j.XmlToJson(fdata)
			check(err)
		*/
		xdata, err := mxj.NewMapXmlSeq(fdata)
		check(err)
		invdata, err = xdata.Json()
		check(err)
	//case "properties":
	default:
		fmt.Printf("Unsupported file adapter type %s:", ftype)
	}
	var jdata interface{}
	err = json.Unmarshal(invdata, &jdata)
	check(err)

	// Encode any leaf maps as scalars to adhere to infra sdk spec
	// TODO: WTF
	var sanitized interface{}
	switch reflect.TypeOf(jdata).Kind() {
	case reflect.Map:
		sanitized = sanitizeMap(jdata)
	case reflect.Slice:
		sanitized = sanitizeSlice(jdata)
	default:
		sanitized = sanitizeLeaf(jdata)
	}

	return sanitized
}

func main() {
	/*
		args := os.Args[1:]
		confData, err := ioutil.ReadFile(args[0])
	*/
	confData, err := ioutil.ReadFile("./test.yml")
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
