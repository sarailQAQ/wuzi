package data_form

import (
	"bytes"
	"encoding/gob"
)

func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

func TypeString(data interface{}) (string) {
	if str,ok := data.(string); ok {
		return str
	}
	return ""
}

func TypeInt(data interface{}) (int) {
	if x,ok := data.(int); ok {
		return x
	}
	return 0
}

func TypeUint(data interface{}) (uint) {
	if x,ok := data.(uint); ok {
		return x
	}
	return 0
}
