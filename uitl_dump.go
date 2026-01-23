package kit

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func DumpJson(value any) {
	switch result := value.(type) {
	case []byte:
		doDumpJSON(result)
	case string:
		doDumpJSON([]byte(result))
	default:
		jsonContent, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		doDumpJSON(jsonContent)
	}
}

func doDumpJSON(jsonContent []byte) {
	if len(jsonContent) == 0 {
		fmt.Println("")
		return
	}
	if !json.Valid(jsonContent) {
		fmt.Println(string(jsonContent))
		return
	}
	var buffer bytes.Buffer
	if err := json.Indent(&buffer, jsonContent, "", "    "); err != nil {
		fmt.Println(err.Error())
		fmt.Println(string(jsonContent))
		return
	}
	fmt.Println(buffer.String())
}
