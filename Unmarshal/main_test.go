// unit test就是從主要的code叫過來以後比較其他的是不是equal,如果是equal他就PASS
package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalSuccess(t *testing.T) {
	jsonData := []byte(`{
		"OSLRCDPath":"/Services/OSLRCD/AA.BBBB.CCC",
		"RSLKITID":"dd.ee.ff",
		"SCD":{
			"ISOPath":"/Services/SCD",
			"ISOName":"xx.yyyy.zzz"
		}
	}`)

	var oslrcd OSLRCD
	err := json.Unmarshal(jsonData, &oslrcd)
	assert.Nil(t, err)
	assert.Equal(t, "/Services/OSLRCD/AA.BBBB.CCC", oslrcd.OSLRCDPath)
	assert.Equal(t, "dd.ee.ff", oslrcd.RSLKitID)
	assert.Equal(t, "/Services/SCD", oslrcd.SCD.ISOPath)
	assert.Equal(t, "xx.yyyy.zzz", oslrcd.SCD.ISOName)
}

func TestUnmarshalInvalidJSON(t *testing.T) {
	jsonData := []byte(`{
		"OSLRCD":"/Services/OSLRCD/AA.BBBB.CCC",
		"RSLKITID":"dd.ee.ff",
		"SCD":{
			"ISOPath":"/Services/SCD",
			"ISOName":"xx.yyyy.zzz"
		`)

	var oslrcd OSLRCD
	err := json.Unmarshal(jsonData, &oslrcd)
	assert.NotNil(t, err)
}

func TestUnmarshalMissingFields(t *testing.T) {
	jsonData := []byte(`{
		"OSLRCDPath":"/Services/OSLRCD/AA.BBBB.CCC",
		"RSLKITID":"dd.ee.ff"
	}`)

	var oslrcd OSLRCD
	err := json.Unmarshal(jsonData, &oslrcd)
	assert.Nil(t, err)
	assert.Equal(t, "/Services/OSLRCD/AA.BBBB.CCC", oslrcd.OSLRCDPath)
	assert.Equal(t, "dd.ee.ff", oslrcd.RSLKitID)
	assert.Equal(t, "", oslrcd.SCD.ISOPath)
	assert.Equal(t, "", oslrcd.SCD.ISOName)
}
