package main

import (
	"bytes"
	"fmt"
	"net"
)

func filterIp(inner chan []IpRes, res chan []IpRes) {
	for {
		checkRes := make(chan IpCheck, 5)
		newArr := []IpRes{}
		ipArr := <-inner
		for k, val := range ipArr {
			go func(a IpRes, num int) {
				canUse := dailTest(a)
				checkRes <- IpCheck{a, canUse}
			}(val, k)
		}
		for i := 0; i < len(ipArr); i++ {
			res := <-checkRes
			if res.CanUse {
				newArr = append(newArr, res.IpRes)
			}
		}
		fmt.Printf("filter res : %v \n", newArr)
		res <- newArr
	}
}

func dailTest(ip IpRes) bool {
	var buffer bytes.Buffer
	buffer.WriteString(ip.Ip)
	buffer.WriteString(":")
	buffer.WriteString(ip.Port)
	res := buffer.String()
	conn, err := net.DialTimeout("tcp", res, 5e9)
	if err != nil {
		//fmt.Println("err" + err.Error() + "\n")
		return false
	}
	fmt.Printf("ip can use: %s \n", res)
	conn.Close()
	return true
}
