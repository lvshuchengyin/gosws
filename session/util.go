// util
package session

import (
	"encoding/json"
	"gosws/util"
)

// return base64
func SessionEncode(key string, data map[string]interface{}) (string, error) {
	bs, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// encrypt
	cryptData, err := util.AesEncrypt(bs, []byte(key))
	if err != nil {
		return "", err
	}

	encoded := util.UrlBase64Encode(cryptData)
	return string(encoded), nil
}

// data is base64
func SessionDecode(key, data string) (map[string]interface{}, error) {
	// decrypt
	data64, err := util.UrlBase64Decode([]byte(data))
	if err != nil {
		return nil, err
	}

	originData, err := util.AesDecrypt(data64, []byte(key))
	if err != nil {
		return nil, err
	}

	var sessMap map[string]interface{}
	err = json.Unmarshal(originData, &sessMap)
	if err != nil {
		return nil, err
	}

	return sessMap, nil
}
