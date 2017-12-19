package main

import (
	"log"
	"net/url"
	"os"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

var logf = log.Printf
var logn = log.Println

var runelen = utf8.RuneCountInString

type xResult struct {
	index   int
	char    string
	pinyins []string
}

func getPinyin(r rune, i int, rch chan<- xResult) {
	u := "http://hanyu.baidu.com/s?wd="
	u += url.QueryEscape(string(r))
	u += "&ptype=zici"

	result := xResult{index: i, char: string(r)}
	doc, err := goquery.NewDocument(u)

	if err == nil {
		doc.Find("#pinyin > span > b").Each(func(_ int, s *goquery.Selection) {
			result.pinyins = append(result.pinyins, s.Text())
		})
	}

	rch <- result
}

func getPinyins(words string) (pys []xResult) {
	n := runelen(words)
	if n <= 0 {
		return nil
	}

	pys = make([]xResult, n)
	rch := make(chan xResult)

	index := 0

	for _, word := range words {
		go getPinyin(word, index, rch)
		index++
	}

	for i := 0; i < n; i++ {
		ret := <-rch
		pys[ret.index] = ret
	}

	return pys
}

func main() {
	if len(os.Args) <= 1 {
		return
	}

	words := os.Args[1]
	allPys := getPinyins(words)

	for _, pys := range allPys {
		logf("%s:\t%s\n", pys.char, pys.pinyins)
	}
}
