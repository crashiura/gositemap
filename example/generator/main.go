package main

import (
	"log"
	"strconv"
	"time"

	"github.com/crashiura/gositemap"
)

func main() {
	g := gositemap.NewGenerator(
		gositemap.IndexSitemapGenOpt("index.xml"),
		gositemap.DirGenOpt("./web"),
		gositemap.CompressGenOpt(false),
		gositemap.HostSitemapGenOpt("https://www.test.com"),
		gositemap.PathSitemapGenOpt("assets/sitemap"),
	)

	productSitemap := g.GetSecondarySitemap("product%d.xml")

	for i := 0; i < 100005; i++ {
		err := productSitemap.Add(&gositemap.URL{
			Loc:        "product/id-" + strconv.Itoa(i),
			LastMod:    gositemap.XMLTime(time.Now()),
			ChangeFreq: gositemap.Always,
			Priority:   0,
		})
		if err != nil {
			log.Println(err)
		}
	}

	blogSitemap := g.GetSecondarySitemap("blog-secondary%d.xml")

	for i := 0; i < 33; i++ {
		err := blogSitemap.Add(&gositemap.URL{
			Loc:        "blog/number-" + strconv.Itoa(i),
			LastMod:    gositemap.XMLTime(time.Now()),
			ChangeFreq: gositemap.Hourly,
			Priority:   0.8,
		})
		if err != nil {
			log.Println(err)
		}
	}

	err := g.Build()
	if err != nil {
		log.Println(err)
	}
}
