package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getMList(c *colly.Collector, db *sql.DB)  {
	start := time.Now()
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Something went wrong:", err, r.StatusCode)
	})

	c.OnRequest(func(r *colly.Request) {
		for key, val := range doubanConfig {
			r.Headers.Set(key, val)
		}
		fmt.Println(r.ProxyURL," : Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("get response at:%v\n", time.Since(start))
		log.Printf("Proxy Address: %s\n", r.Request.ProxyURL)
		start = time.Now()
		fmt.Printf("%s\n\n", string(r.Body))
		json.Unmarshal(r.Body, movieArr)
		var insertStr bytes.Buffer
		foo := bufio.NewWriter(&insertStr)
		foo.WriteString("insert into gomovie ( movieId, url, rate, title, playable, cover, IsNew) values ")
		for idx, movie := range movieArr.Subjects {
			createStr(movie, []string{"id", "url", "rate", "title", "playable", "cover", "is_new"}, foo)
			if idx < len(movieArr.Subjects)-1 {
				foo.WriteString(",")
			}
		}
		foo.WriteString(" on duplicate key update rate=values(rate),playable=values(playable), IsNew=values(IsNew)")
		foo.Flush()
		str := insertStr.String()
		fmt.Printf("get : %v ", str)
		stmt, err := db.Prepare(str)
		checkErr(err)
		sqlres, err := stmt.Exec()
		checkErr(err)
		affect, err := sqlres.RowsAffected()
		checkErr(err)
		fmt.Println(affect)
		fmt.Printf("insert finished for: %v\n", time.Since(start))
		start = time.Now()
		//rows, err := db.Query(`select COST from products_tbl`)
		//checkErr(err)
		//for rows.Next() {
		//	var cost string
		//	err := rows.Scan(&cost)
		//	checkErr(err)
		//	fmt.Printf("get: %v", cost)
		//}
		//fmt.Printf("%v,%v", json.Valid(r.Body), *res)
	})

	c.OnScraped(func(r *colly.Response) {
	})
}

func getMInfo( cInfo *colly.Collector, infoIserte chan string)  {
	start := time.Now()
	cInfo.OnRequest(func(r *colly.Request) {
		println(r.ProxyURL)
	})

	cInfo.OnError(func(r *colly.Response, err error) {
		log.Println("Something went wrong:", err, r.StatusCode)
	})

	cInfo.OnResponse(func(r *colly.Response) {
		fmt.Printf("get response at:%v\n", time.Since(start))
		log.Printf("Proxy Address: %s\n", r.Request.ProxyURL)
	})

	cInfo.OnHTML("#content", func(e *colly.HTMLElement) {
		//strings.Replace(e.Text, " ", "", -1)
		println("onHtml")
		var infoArr [][]string
		var MInfo MovieInfo
		var RInfo RateInfo
		var RatePer = make(map[string]string)
		MInfo.MovieId = strings.Split(e.Request.URL.Path, "/")[2]
		reg := regexp.MustCompile(`(\s?\S+)+`)
		params := reg.FindAllString(strings.Replace(e.ChildText("#info"), "\n", "", -1), -1)
		e.DOM.Find(".ratings-on-weight .item .rating_per").Each(func(i int, sel *goquery.Selection) {
			RatePer[strconv.Itoa(5-i)+"star"] = sel.Text()
		})
		MInfo.Summary = e.ChildText("#link-report span[property='v:summary']")
		if e.DOM.IsNodes(e.DOM.Has("#link-report .a_show_full").Nodes...) {
			MInfo.Summary = e.ChildText("#link-report .all.hidden")
		}
		RInfo.Total = e.ChildText(".ll.rating_num")
		RInfo.RatingSum = e.ChildText(".rating_people span[property='v:votes']")
		RInfo.RatePer = RatePer
		jsonStr, _ := json.Marshal(RInfo)
		MInfo.RateInfo = string(jsonStr)
		for _, param := range params {
			strArr := strings.Split(param, ":")
			if strArr[0] != "又名" {
				strArr = strings.Split(strings.Replace(param, " ", "", -1), ":")
			}
			strArr[1] = strings.Trim(strings.Replace(strArr[1], " / ", "/", -1), " ")
			//fmt.Printf("\nget %d: %v\n", idx, strArr[0])
			infoArr = append(infoArr, strArr)
		}
		for _, val := range infoArr {
			switch val[0] {
			case "又名":
				MInfo.NTitle = val[1]
			case "上映日期":
				MInfo.ShowDate = val[1]
			case "片长":
				MInfo.Duration = val[1]
			case "IMDb链接":
				MInfo.Imbd = val[1]
			}
		}
		reg = regexp.MustCompile(`\d+`)
		MInfo.CommentNum = reg.FindString(e.ChildText(".mod-hd span.pl"))
		MInfo.HasSeen = reg.FindString(e.ChildText(".subject-others-interests-ft a:first-child"))
		MInfo.WantSee = reg.FindString(e.ChildText(".subject-others-interests-ft a:last-child"))
		var insertStr bytes.Buffer
		foo := bufio.NewWriter(&insertStr)
		foo.WriteString(`insert into movieInfo ( movieId, rateInfo, showDate, duration, imbd, nTitle, summary, commentNum, wantSee, hasSeen) VALUES `)
		createStr(MInfo, []string{"MovieId", "RateInfo", "ShowDate", "Duration", "Imbd", "NTitle", "Summary", "CommentNum", "WantSee", "HasSeen"}, foo)
		foo.WriteString(" on duplicate key update rateInfo=values(rateInfo), commentNum=values(commentNum), wantSee=values(wantSee), hasSeen=values(hasSeen);")
		foo.Flush()
		str := insertStr.String()
		go func(s string) {
			println("insert")
			infoIserte <- s
		}(str)
		//fmt.Printf("\n%s\n", strings.Replace(e.Text, "\n", "", -1))
	})
}

