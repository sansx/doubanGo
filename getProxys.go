package main

import (
	"bytes"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	IpArr     []IpRes
	fmtIp     []string
	spiderUrl string = "https://www.kuaidaili.com/free/inha/"
	page      string
)

func GetProxys(proxyIp string, start string, loop int, sleep time.Duration) []string {
	innerIp := make(chan []IpRes)
	filterRes := make(chan []IpRes)
	page = start
	go filterIp(innerIp, filterRes)
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"),
		colly.AllowURLRevisit(),
		//colly.Async(true),
	)
	if p, err := proxy.RoundRobinProxySwitcher(
		proxyIp); err == nil {
		c.SetProxyFunc(p)
	}else {
		log.Fatal(err)
	}
	c.SetRequestTimeout(10e9)

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.OnRequest(func(r *colly.Request) {
		webConfig := map[string]string{
			"Accept":                    "*/*",
			"Accept-Encoding":           "gzip, deflate, br",
			"Accept-Language":           "zh-CN,zh;q=0.9",
			"Cache-Control":             "max-age=0",
			"Connection":                "keep-alive",
			"Host":                      "www.kuaidaili.com",
			"Sec-Fetch-Mode":            "navigate",
			"Sec-Fetch-Site":            "none",
			"Sec-Fetch-User":            "?1",
			"Upgrade-Insecure-Requests": "1",
			"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36",
			"X-Requested-With":          "XMLHttpRequest",
		}
		//var cookieFront = "channelid=0; sid=1578140149928365; _ga=GA1.2.1446772863.1578141329; _gid=GA1.2.1755857370.1578141329; Hm_lvt_7ed65b1cc4b810e9fd37959c9bb51b31=1578141329; _gat=1;"
		//webConfig["Cookie"] = strings.Join([]string{cookieFront,"Hm_lpvt_7ed65b1cc4b810e9fd37959c9bb51b31=",strconv.FormatInt(time.Now().Unix(),10)},"")
		webConfig["Referer"] = strings.Join([]string{spiderUrl, page, "/"}, "")
		for key, val := range webConfig {
			r.Headers.Set(key, val)
		}
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("\n-Something went wrong:%s %v %v %v %v", err, r.Body,r.Request.Body,r.Ctx,r.Headers)
		panic("err")
	})
	// Find and visit all links
	c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
		//fmt.Printf("%v\n",e)
		res := IpRes{}
		e.ForEach("td", func(i int, e *colly.HTMLElement) {
			switch e.Attr("data-title") {
			case "IP":
				res.Ip = e.Text
			case "PORT":
				res.Port = e.Text
			case "类型":
				res.IpType = e.Text
			case "响应速度":
				res.Recall = e.Text
			}
		})
		IpArr = append(IpArr, res)
	})

	c.OnScraped(func(r *colly.Response) {
		innerIp <- IpArr
		IpArr = []IpRes{}
	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Printf("%s",r.Body)
		//log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
	})
	for i := 0; i < loop; i++ {
		num, _ := strconv.Atoi(start)
		page = strconv.Itoa(num + i)
		str := strings.Join([]string{spiderUrl, page, "/"}, "")
		//if p, err := proxy.RoundRobinProxySwitcher(
		//	fmtIp...); err == nil && len(fmtIp)>0 {
		//	c = c.Clone()
		//	c.SetProxyFunc(p)
		//}
		fmt.Printf("\nstr:%s\n", str)
		err := c.Visit(str)
		if err != nil {
			fmt.Printf("get Err:%v", err.Error())
		}
		time.Sleep(sleep)
		getRes := <-filterRes
		if len(getRes) > 0 {
			for _, v := range getRes {
				var str bytes.Buffer
				str.WriteString(`"` + strings.ToLower(v.IpType))
				str.WriteString("://")
				str.WriteString( v.Ip)
				str.WriteString(":")
				str.WriteString(v.Port + `",`)
				fmtIp = append(fmtIp, str.String())
			}
			fmt.Printf("\nfmt res:%v \n", fmtIp)
		}
	}
	return fmtIp
}
