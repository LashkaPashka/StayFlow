package encode

import "encoding/json"

func Encode(payload any) []byte {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	return payloadBytes
}