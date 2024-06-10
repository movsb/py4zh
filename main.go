package main

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

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
	n := utf8.RuneCountInString(words)
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
	ruby := flag.Bool(`ruby`, false, `生成 <ruby> HTML 内容`)
	flag.Parse()

	if flag.NArg() <= 0 {
		return
	}

	words := flag.Args()
	allPys := getPinyins(strings.Join(words, ""))

	if !*ruby {
		for _, pys := range allPys {
			fmt.Printf("%s: %s\n", pys.char, pys.pinyins)
		}
	} else {
		for _, pys := range allPys {
			for _, p := range pys.pinyins {
				fmt.Printf("<ruby>%s<rp>(</rp><rt>%s</rt><rp>)</rp></ruby>\n", pys.char, p)
			}
		}
	}
}
