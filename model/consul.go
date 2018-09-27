package model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// ConsulValue contains value returned by consul
type ConsulValue struct {
	Value string
}

// GetValueWithKey get value from consul with given key and url
func GetValueWithKey(key, url string) (string, error) {
	client := &http.Client{}
	res, err := client.Get(url + key)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if len(body) == 0 {
		return "", errors.New("can not get satisfied value")
	}
	ret := []ConsulValue{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return "", err
	}
	if len(ret) == 0 {
		return "", errors.New("Can not get value")
	}
	decoded, err := base64.StdEncoding.DecodeString(ret[0].Value)

	return string(decoded), err
}
