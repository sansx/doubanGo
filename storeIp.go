package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

//f, err := os.OpenFile(IpFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
//defer f.Close()
//if err != nil {
//	fmt.Println("getErr", err.Error())
//	return
//}
//inputFile := json.NewEncoder(f)
//if err := inputFile.Encode(res); err != nil {
//	log.Println("URLStore: ", err)
//}

func getIp(fileName string) (IpArr []IpRes, err error) {
	f, err := os.Open(fileName)
	defer f.Close()
	d := json.NewDecoder(f)
	for err == nil {
		var r IpRes
		if err = d.Decode(&r); err == nil {
			IpArr = append(IpArr, r)
		}
	}
	if err == io.EOF {
		return IpArr, nil
	}
	log.Println("Error decoding:", err)
	return
}
