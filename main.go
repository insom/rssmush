package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Item struct {
	Raw string `xml:",innerxml"`
}

type Rss struct {
	XMLName     xml.Name `xml:"rss"`
	Title       string   `xml:"channel>title"`
	Link        string   `xml:"channel>link"`
	Description string   `xml:"channel>description"`
	Items       []Item   `xml:"channel>item"`
}

type Blog struct {
	Url string
	Rss Rss
}

func (b *Blog) Read() {
	p, _ := http.Get(b.Url)
	defer p.Body.Close()
	data, _ := ioutil.ReadAll(p.Body)
	xml.Unmarshal(data, &b.Rss)
}

type Blogs []*Blog

func (bs Blogs) Read() {
	for _, b := range bs {
		b.Read()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: rssmush [input urls...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	output := flag.String("filename", "output.xml", "where to write the finished RSS feed")
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}
	bs := make(Blogs, flag.NArg())
	for i, x := range flag.Args() {
		bs[i] = &Blog{Url: x}
	}
	bs.Read()
	mainBlogRss := &bs[0].Rss
	for i, b := range bs {
		if i > 0 {
			mainBlogRss.Items = append(mainBlogRss.Items, b.Rss.Items...)
		}
	}
	data, _ := xml.Marshal(mainBlogRss)
	ioutil.WriteFile(*output, data, 0644)
}
