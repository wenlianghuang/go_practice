package main

import (
	"encoding/json"
	"fmt"
)

type SCD struct {
	ISOName string `json:"ISOName"`
	ISOPath string `json:"ISOPath"`
}

type OSLRCD struct {
	OSLRCDPath string `json:"OSLRCDPath"`
	RSLKitID   string `json:"RSLKitID"`
	SCD        SCD    `json:"SCD"`
}

func main() {
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
	if err != nil {
		fmt.Println("JSON decode failed: ", err)
		return
	}

	fmt.Println("ISOName: ", oslrcd.SCD.ISOName)
	fmt.Println("ISOPath: ", oslrcd.SCD.ISOPath)
	fmt.Println("OSLRCDPath: ", oslrcd.OSLRCDPath)
	fmt.Println("RSLKitID: ", oslrcd.RSLKitID)
}
