package jwt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Payload struct{
	Iss string `json:"iss"`
	Exp string `json:"exp"`
	Username string `json:"user"`
	Id uint `json:"id"`
	Iat    string `json:"iat"`
}

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func NewHeader() Header{
	return Header{
		Alg: "HS256",
		Typ: "JWT",
	}
}

func Creat(user string,id uint)string{
	header := NewHeader()
	payload := Payload{
		Iss:      "thhbmz",
		Exp:      strconv.FormatInt(time.Now().Add(3*time.Hour).Unix(), 10),
		Username: user,
		Id:       id,
		Iat:      strconv.FormatInt(time.Now().Unix(), 10),
	}

	h , e1 := json.Marshal(header)
	p , e2 := json.Marshal(payload)
	if e1 != nil || e2 != nil {
		fmt.Println(e1.Error(),e2.Error())
		return ""
	}
	header64 := base64.StdEncoding.EncodeToString(h)
	payload64 := base64.StdEncoding.EncodeToString(p)
	str1 := strings.Join([]string{header64,payload64},".")

	key := "thhbmz"
	mac := hmac.New(sha256.New,[]byte(key))
	mac.Write([]byte(str1))
	s := mac.Sum(nil)

	signature := base64.StdEncoding.EncodeToString(s)
	return str1+"."+signature
}

func Check(token string)(uid uint,user string,err error){
	arr := strings.Split(token,".")
	//fmt.Println(token)
	if len(arr) != 3 {
		err =errors.New("token length error")
		return
	}

	_ , err = base64.StdEncoding.DecodeString(arr[0])
	if err != nil {
		fmt.Println("token:header errror",err)
		return
	}

	pay , err := base64.StdEncoding.DecodeString(arr[1])
	if err != nil {
		err=errors.New("token:payload error")
		return
	}

	sign , err :=base64.StdEncoding.DecodeString(arr[2])
	if err != nil {
		err = errors.New("tken:signature error")
		return
	}

	str1 := arr[0] + "." + arr[1]

	key := []byte("thhbmz")
	mac := hmac.New(sha256.New,key)
	mac.Write([]byte(str1))
	s := mac.Sum(nil)
	fmt.Println(sign,s)
	if res := bytes.Compare(sign,s); res != 0 {
		err = errors.New("token:signature error")
		return
	}

	var payload Payload
	json.Unmarshal(pay,&payload)
	uid = payload.Id
	user = payload.Username
	return
}
