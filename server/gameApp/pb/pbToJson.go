package pb

import "encoding/json"

type TestItemJson struct {
	Index int32
	Item  *TestItem
}

func (item *TestItem) ToJson(index int32) ([]byte, error) {
	jsonObj := TestItemJson{index, item}
	return json.Marshal(jsonObj)
}
