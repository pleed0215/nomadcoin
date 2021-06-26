package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	HandleError(encoder.Encode(i))
	return blockBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	HandleError(gob.NewDecoder(bytes.NewReader(data)).Decode(i))
}

func Hash(i interface {}) string {
	s:=fmt.Sprintf("%v", i)
	hash := sha256.Sum256(([]byte(s)))
	return fmt.Sprintf("%x", hash)
}