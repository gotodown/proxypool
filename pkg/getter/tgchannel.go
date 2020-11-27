package getter

import (
	"fmt"
	logger "log"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/zu1k/proxypool/pkg/proxy"
	"github.com/zu1k/proxypool/pkg/tool"
)

var log *logger.Logger

func init() {
	log = tool.Logger
	Register("tgchannel", NewTGChannelGetter)
}

type TGChannelGetter struct {
	c         *colly.Collector
	NumNeeded int
	results   []string
	Url       string
}

func NewTGChannelGetter(options tool.Options) (getter Getter, err error) {
	num, found := options["num"]
	t := 200
	switch num.(type) {
	case int:
		t = num.(int)
	case float64:
		t = int(num.(float64))
	}

	if !found || t <= 0 {
		t = 200
	}
	urlInterface, found := options["channel"]
	if found {
		url, err := AssertTypeStringNotNull(urlInterface)
		if err != nil {
			return nil, err
		}
		return &TGChannelGetter{
			c:         tool.GetColly(),
			NumNeeded: t,
			Url:       "https://t.me/s/" + url,
		}, nil
	}
	return nil, ErrorUrlNotFound
}

func (g *TGChannelGetter) Get() proxy.ProxyList {
	result := make(proxy.ProxyList, 0)
	g.results = make([]string, 0)
	// 找到所有的文字消息
	g.c.OnHTML("div.tgme_widget_message_text", func(e *colly.HTMLElement) {

		text, _ := e.DOM.Html()
		lines := strings.Split(text, "<br/>")
		for _, line := range lines {
			links := GrepLinksFromString(line)
			g.results = append(g.results, links...)
		}

		// 抓取到http链接，有可能是订阅链接或其他链接，无论如何试一下
		// 打印--
		// log.Println(e.Text, g.Url, "========\n")
		subUrls := make([]string, 0)
		e.DOM.Find("a").Each(func(i int, selection *goquery.Selection) {
			u := selection.Text()
			subUrls = append(subUrls, u)
		})
		// subUrls := urlRe.FindAllString(e.Text, -1)
		for _, url := range subUrls {
			tgName := strings.Split(g.Url, "/")
			if strings.Contains(url, "t.me") || strings.HasPrefix(url, "@") || len(url) < 30 {
				continue
			}
			log.Printf("来自%s频道的可能的链接%s \n", tgName[len(tgName)-1], url)
			result = append(result, (&Subscribe{Url: url}).Get()...)
		}
	})
	// 找到之前消息页面的链接，加入访问队列
	g.c.OnHTML("link[rel=prev]", func(e *colly.HTMLElement) {
		if len(g.results) < g.NumNeeded {
			_ = e.Request.Visit(e.Attr("href"))
		}
	})

	g.results = make([]string, 0)
	err := g.c.Visit(g.Url)
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
	}
	return append(result, StringArray2ProxyArray(g.results, g.Url)...)
}

func (g *TGChannelGetter) Get2Chan(pc chan proxy.Proxy, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := g.Get()
	log.Printf("STATISTIC: TGChannel\tcount=%d\turl=%s\n", len(nodes), g.Url)
	for _, node := range nodes {
		pc <- node
	}
}
