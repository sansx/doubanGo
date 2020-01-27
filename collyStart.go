package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type IpRes struct {
	Ip     string
	Port   string
	IpType string `json:"ip_type"`
	Recall string
}

type IpCheck struct {
	IpRes
	CanUse bool
}

type MovieInfo struct {
	MovieId    string
	RateInfo   string
	ShowDate   string
	Duration   string
	Imbd       string
	NTitle     string
	Summary    string
	CommentNum string
	WantSee    string
	HasSeen    string
}

type InfoItem struct {
	Rate     string `json:"rate"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	Playable bool   `json:"playable"`
	Cover    string `json:"cover"`
	Id       string `json:"id"`
	IsNew    bool   `json:"is_new"`
}

type commInfo struct {
	Cid     string
	MovieId string
	User    string
	ULink   string
	UImg    string
	Date    string
	Rate    string
	Votes   string
	Comment string
}

type RateInfo struct {
	Total     string
	RatingSum string
	RatePer   map[string]string
}

type NetReturn struct {
	Subjects []InfoItem `json:"subjects"`
}

var (
	IpFile       = "./connect.json"
	doubanConfig = map[string]string{
		"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		//"Accept-Encoding":  "gzip, deflate, br",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Connection":      "keep-alive",
		"Host":            "movie.douban.com",
		//"Content-Type": 		"application/json;charset=UTF-8",
		"Sec-Fetch-Mode":   "navigate",
		"Sec-Fetch-Site":   "none",
		"Sec-Fetch-User":   "?1",
		"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36",
		"X-Requested-With": "XMLHttpRequest",
		"Cookie":           `bid=ushAA-2w-7Q; ll="118174"; __utmc=30149280; __utmc=223695111; __yadk_uid=upEYssVluItDlvwW3J9GdG9lLO8hxdjN; _vwo_uuid_v2=DF2B44256E1AC02453A49B9648A40100E|f561403fb009dba0f8cc78a30b7e2c70; viewed="27021790"; gr_user_id=36686aa0-46bb-49f2-b1fa-eba892225df8; trc_cookie_storage=taboola%2520global%253Auser-id%3D13252cda-6cfa-4214-9cc8-a6ca7bd4df0e-tuct4ff4bab; ct=y; __utmz=30149280.1579137437.8.5.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __utmz=223695111.1579137437.7.5.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1579590689%2C%22https%3A%2F%2Fwww.google.com%2F%22%5D; _pk_ses.100001.4cf6=*; ap_v=0,6.0; __utma=30149280.1573436667.1578367462.1579566610.1579590689.21; __utma=223695111.29153026.1578367462.1579566610.1579590689.20; __utmb=223695111.0.10.1579590689; dbcl2="209565533:9wyG4GF57go"; ck=4mOa; push_noty_num=0; push_doumail_num=0; __utmv=30149280.20956; __utmb=30149280.4.10.1579590689; _pk_id.100001.4cf6=bb96bb66a5b250d0.1578367462.20.1579594912.1579566609.`,
		//`ll="118174"; bid=iixAdsHav7g; __utmc=30149280; __utmc=223695111; __yadk_uid=1kCWbAOjqecDW3oPsC1w1gIAAbNMP0gF; _vwo_uuid_v2=D9889C3872072D24349633B722DCE59C9|68435c7981e40979aeee3e8570450858; trc_cookie_storage=taboola%2520global%253Auser-id%3Dfc4e983d-5c62-4d5f-bda4-f55698840db0-tuct50a29ac; ct=y; __utmz=30149280.1578919579.10.3.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __utmz=223695111.1578919579.10.3.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1579315650%2C%22https%3A%2F%2Fwww.google.com%2F%22%5D; _pk_ses.100001.4cf6=*; __utma=30149280.1092024888.1578148770.1579312712.1579315650.18; __utmb=30149280.0.10.1579315650; __utma=223695111.1173361526.1578148770.1579312712.1579315651.18; __utmb=223695111.0.10.1579315651; _pk_id.100001.4cf6=c42d1db054733a4e.1578148770.18.1579319697.1579312721.`,
		//`ll="118174"; bid=iixAdsHav7g; __utmc=30149280; __utmz=30149280.1578148770.1.1.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __utmc=223695111; __utmz=223695111.1578148770.1.1.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __yadk_uid=1kCWbAOjqecDW3oPsC1w1gIAAbNMP0gF; _vwo_uuid_v2=D9889C3872072D24349633B722DCE59C9|68435c7981e40979aeee3e8570450858; trc_cookie_storage=taboola%2520global%253Auser-id%3Dfc4e983d-5c62-4d5f-bda4-f55698840db0-tuct50a29ac; ap_v=0,6.0; ct=y; _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1578202711%2C%22https%3A%2F%2Fwww.google.com%2F%22%5D; _pk_id.100001.4cf6=c42d1db054733a4e.1578148770.4.1578202711.1578196026.; _pk_ses.100001.4cf6=*; __utma=30149280.1092024888.1578148770.1578195755.1578202711.4; __utmb=30149280.0.10.1578202711; __utma=223695111.1173361526.1578148770.1578195755.1578202711.4; __utmb=223695111.0.10.1578202711`,
		"Upgrade-Insecure-Requests": "1",
	}
	movieArr   = new(NetReturn)
	infoIserte = make(chan string, 2)
	wp         = &sync.WaitGroup{}
	infoWp     = &sync.WaitGroup{}
	startCom   = 0
	countNum   = make(chan int)
	commEnd    = false
	MLoop      = 1
)

func collyStart() {
	rand.Seed(time.Now().Unix())
	proxysArr := GetProxys("http://163.172.147.94:8811", "1",3,2e9)
	start := time.Now()
	AStart := time.Now()
	ACount := 0
	mCount := countDown(&startCom)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"),
		colly.AllowURLRevisit(),
		colly.Async(true),
	)
	db := connSql()
	defer db.Close()
	//proxysArr := []string{
	//	//"http://140.255.186.40:9999",
	//	"http://222.95.144.43:3000",
	//	//"http://123.160.1.96:9999",
	//}
	if p, err := proxy.RoundRobinProxySwitcher(
		proxysArr...,
	); err == nil {
		fmt.Printf("Set success")
		c.SetProxyFunc(p)
	}else {
		log.Fatal(err)
	}
	c.SetRequestTimeout(5e9)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*.douban.*",
		Parallelism: 2,
		RandomDelay: 3 * time.Second,
		//Delay:      5 * time.Second,
	})
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	getMList(c, db)
	url := "https://movie.douban.com/j/search_subjects?type=movie&tag=%E7%83%AD%E9%97%A8&sort=time&page_limit=2&page_start=10"
	c.Visit(url)
	c.Wait()
	fmt.Printf("end sys: %v\n", time.Since(start))
	//fmt.Printf("finail movieArr : %v", movieArr)
	cInfo := c.Clone()
	cInfo.SetRequestTimeout(10e9)

	getMInfo(cInfo, infoIserte)
	go sqlIn(infoIserte, db, wp)
	wp.Add(len(movieArr.Subjects))
	for idx, val := range movieArr.Subjects {
		fmt.Printf("\n%d:%v\n", idx, val.Url)
		err := cInfo.Visit(val.Url)
		checkErr(err)
		if idx+1%5 == 0 {
			println("take a break")
			time.Sleep(1e9)
			println("continue to work!")
		}
	}
	cInfo.Wait()
	wp.Wait()
	cCom := cInfo.Clone()
	//cCom.Async = false
	var insertStr bytes.Buffer
	itemCount := 0
	commUrl := ""
	foo := bufio.NewWriter(&insertStr)
	//insert into movieComm (cid, movieId, user, userLink, userImg, date, rate, votes, comment) values () on duplicate key update votes=values(votes), rate=values(rate), comment=values(comment);
	foo.WriteString(`insert into movieComm (cid, movieId, user, userLink, userImg, date, rate, votes, comment) values `)

	cCom.OnRequest(func(r *colly.Request) {
		for key, val := range doubanConfig {
			r.Headers.Set(key, val)
		}
		fmt.Println(r.ProxyURL, " : Visiting", r.URL)
	})

	cCom.OnHTML(".article", func(e *colly.HTMLElement) {
		MId := strings.Split(e.Request.URL.Path, "/")[2]
		itemCount = itemCount + e.DOM.Find(`.comment-item`).Size()
		println(e.DOM.Find(`.comment-item`).Length())
		if e.DOM.Find(`.avatar`).Length() == 0 {
			commEnd = true
			infoWp.Done()
			return
		}
		println("onHtml startCom:", startCom)
		if startCom > MLoop -1  {
			foo.WriteString(",")
		}

		e.DOM.Find(`.comment-item`).Each(func(i int, s *goquery.Selection) {
			var CInfo commInfo
			if val, ok := s.Attr("data-cid"); ok {
				CInfo.Cid = val
			}
			CInfo.MovieId = MId
			user := s.Find(".avatar a")
			name, _ := user.Attr("title")
			link, _ := user.Attr("href")
			img, _ := user.Find("img").Attr("src")
			CInfo.User = name
			CInfo.ULink = link
			CInfo.UImg = img
			date, _ := s.Find(".comment-time").Attr("title")
			CInfo.Date = date
			CInfo.Votes = s.Find(".comment-vote .votes").Text()
			CInfo.Comment = s.Find("p .short").Text()
			switch true {
			case s.Find(".rating").HasClass("allstar10"):
				CInfo.Rate = "1"
			case s.Find(".rating").HasClass("allstar20"):
				CInfo.Rate = "2"
			case s.Find(".rating").HasClass("allstar30"):
				CInfo.Rate = "3"
			case s.Find(".rating").HasClass("allstar40"):
				CInfo.Rate = "4"
			case s.Find(".rating").HasClass("allstar50"):
				CInfo.Rate = "5"
			}
			//cid, movieId, user, userLink, userImg, date, rate, votes, comment
			createStr(CInfo, []string{"Cid", "MovieId", "User", "ULink", "UImg", "Date", "Rate", "Votes", "Comment"}, foo)
			if e.DOM.Find(`.comment-item`).Size() != i+1 {
				foo.WriteString(",")
			}
			//fmt.Printf("\n%v\n", CInfo)
		})
		rege := regexp.MustCompile(`\d+`)
		num, _ := strconv.Atoi(rege.FindString(e.ChildText("li:first-child")))
		mCount()
		go func(n int, startCom int) {
			//startCom := <-countNum
			println("\n current:", startCom)
			if n > startCom*20 {
				infoWp.Add(1)
				println("get:!!", startCom, n)
				comUrl := strings.Join([]string{commUrl, strconv.Itoa(startCom * 20), "&limit=20&sort=new_score&status=P"}, "")
				sleepNum, _ := strconv.ParseFloat(strconv.Itoa(5+rand.Intn(22))+"e8", 64)
				time.Sleep(time.Duration(sleepNum))
				err := cCom.Visit(comUrl)
				checkErr(err)
			}
			infoWp.Done()
		}(num, startCom)
	})
	for idx, val := range movieArr.Subjects {
		start = time.Now()
		for i := 0; i < MLoop; i++ {
			infoWp.Add(1)
			//startCom := <-countNum
			println("link num:", startCom)
			commUrl = strings.Join([]string{"https://movie.douban.com/subject/", val.Id, "/comments?start="}, "")
			comUrl := strings.Join([]string{commUrl, strconv.Itoa(startCom * 20), "&limit=20&sort=new_score&status=P"}, "")
			err := cCom.Visit(comUrl)
			checkErr(err)
			if i < MLoop-1 {
				mCount()
			}
		}
		infoWp.Wait()
		foo.WriteString(" on duplicate key update rate=values(rate), votes=values(votes), comment=values(comment);")
		foo.Flush()
		sqlStr := insertStr.String()
		//fmt.Printf("\n%s\n", sqlStr)
		fmt.Printf("No.%d %s-%s finish:%v\ntotal get: %d ", idx, val.Title, val.Id, time.Since(start), itemCount)
		ACount += itemCount
		println("ACount += itemCount : ", ACount-itemCount, " += ", itemCount, " is ", ACount)
		//f, _ := os.OpenFile("./comm.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		//_, err := f.WriteString(sqlStr)
		//checkErr(err)
		//defer f.Close()
		stmt, err := db.Prepare(sqlStr)
		checkErr(err)
		sqlres, err := stmt.Exec()
		checkErr(err)
		affect, err := sqlres.RowsAffected()
		checkErr(err)
		fmt.Println("affect:", affect)
		insertStr = bytes.Buffer{}
		foo = bufio.NewWriter(&insertStr)
		foo.WriteString(`insert into movieComm (cid, movieId, user, userLink, userImg, date, rate, votes, comment) values `)
		startCom = 0
		itemCount = 0
		commEnd = false
	}
	fmt.Printf("All Done for: %v \nAll get: %d",  time.Since(AStart), ACount)
}

func sqlIn(ch chan string, d *sql.DB, wp *sync.WaitGroup) {
	count := 0
	for str := range ch {
		count++
		println("count:", count)
		stmt, err := d.Prepare(str)
		checkErr(err)
		sqlres, err := stmt.Exec()
		checkErr(err)
		affect, err := sqlres.RowsAffected()
		checkErr(err)
		fmt.Println("affect:", affect)
		wp.Done()
	}
}

func countDown(n *int) func() {
	return func() {
		*n++
	}
}

func createStr(movie interface{}, queryArr []string, str *bufio.Writer) {
	var movieMap map[string]interface{}
	inrec, _ := json.Marshal(movie)
	json.Unmarshal(inrec, &movieMap)
	//fmt.Printf("%d %v\n", len(queryArr), movieMap)
	str.WriteString("(")
	for idx, val := range queryArr {
		//fmt.Printf("%d %v\n", len(queryArr), movieMap[val])
		switch t := movieMap[val].(type) {
		case bool:
			str.WriteString(strconv.FormatBool(t))
		case string:
			if len(t) > 0 {
				str.WriteString("'" + strings.Replace(t, "'", "\\'", -1) + "'")
				break
			}
			str.WriteString("null")
		case nil:
			panic("val not exist :" + val)
		}
		if idx != len(queryArr)-1 {
			str.WriteString(",")
		}
	}
	str.WriteString(")")

}

