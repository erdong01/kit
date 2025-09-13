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
	var (
		buffer    = bytes.NewBuffer(nil)
		jsonBytes = jsonContent
	)
	if err := json.Indent(buffer, jsonBytes, "", "    "); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(buffer.String())
}
