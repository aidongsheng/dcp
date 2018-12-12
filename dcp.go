package main

import (
	"database/sql"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

const dataSourceName  = "root:ads---@/wcc"

func insertIntoDCP(id string, category string, p_at_book string) {
	mydb,openErr := sql.Open("mysql",dataSourceName)
	if openErr != nil {
		log.Fatal(openErr)
	}

	defer mydb.Close()
	insertStmt, insertErr := mydb.Prepare("insert into dcp (t,dcpid,category,patbook) values (?,?,?,?)")
	if insertErr != nil {
		log.Fatal(insertErr)
	}else {
		t := time.Now().Format("2006-01-02 15:04:05")
		_,resultErr := insertStmt.Exec(t,id,category,p_at_book)
		if resultErr != nil{
			log.Fatal(resultErr)
		}
	}
}

func insertIntoDCP1(id string, content string, indexTb string) {
	mydb,openErr := sql.Open("mysql",dataSourceName)
	if openErr != nil {
		log.Fatal(openErr)
	}

	defer mydb.Close()
	insertStmt, insertErr := mydb.Prepare("update dcp set content="+"'"+content+"'" + "," + "indexTb="+"'"+indexTb+"'"+"where dcpid="+"'"+id+"'")
	if insertErr != nil {
		log.Fatal(insertErr)
	}else {
		_,resultErr := insertStmt.Exec()
		if resultErr != nil{
			log.Fatal(resultErr)
		}
	}
}

func main() {
	c := colly.NewCollector()
	c.OnError(func(response *colly.Response, e error) {

	})
	c.OnRequest(func(request *colly.Request) {
		log.Printf("请求ID:%d 链接:%s",request.ID,request.URL)
	})

	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36"

	c.OnHTML("div[class]", func(element *colly.HTMLElement) {

		element.ForEach("div[class]", func(i int, element *colly.HTMLElement) {

			var category,position,content,indexTb string

			if element.Attr("class") == "panel-heading" {
				element.ForEach("span[id]", func(i int, element *colly.HTMLElement) {
					if element.Attr("id") == "ContentPlaceHolder1_lbl_name" {
						log.Printf("门类:%s",element.Text)
						category = element.Text
					}
				})
				element.ForEach("ol[class]", func(i int, element *colly.HTMLElement) {
					if element.Attr("class") == "breadcrumb nomargin" {

						element.ForEach("li", func(i int, element *colly.HTMLElement) {
							position = position + element.Text
						})
						log.Printf("位于植物志的位置:%s",position)
					}
				})
				id := element.Request.URL.String()[len(element.Request.URL.String())-5:]
				insertIntoDCP(id,category,position)
			}

			if element.Attr("class") == "panel-body" {
				element.ForEach("table[class]", func(i int, element *colly.HTMLElement) {
					if element.Attr("class") == "contentPadding" {	//	正文简介
						element.ForEach("span[id]", func(i int, element *colly.HTMLElement) {
							if element.Attr("id") == "ContentPlaceHolder1_lbl_fulltext" {
								log.Printf("正文:%s",element.Text)
								content = element.Text
							}
						})
					}
					if element.Attr("class") == "IndexTb" {	//	检索表
						element.ForEach("tr", func(i int, element *colly.HTMLElement) {
							if element.ChildAttrs("span","class")[0] == "padding-right5" {
								indexTb = indexTb + element.Text + ","
							}
							if element.ChildAttrs("span","class")[1] == "padding-left2 pull-right" {
								url := element.Request.AbsoluteURL(element.ChildAttr("a","href"))
								url = strings.Replace(url," ","%20",-1)
								indexTb = indexTb + url + ";"
							}
						})

					}
				})
				log.Printf("检索表子项内容和链接:%s",indexTb)
				id := element.Request.URL.String()[len(element.Request.URL.String())-5:]
				insertIntoDCP1(id,content,indexTb)
			}


		})
	})
	var i int64
	for  i = 49520; i <= 56144;i++  {
		url := "http://db.kib.ac.cn/XZFlora/SearchResult.aspx?id="+strconv.FormatInt(i,10)
		log.Printf("%s",url)
		c.Visit(url)
	}
}