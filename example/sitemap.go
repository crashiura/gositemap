package main

import (
	"fmt"
	"log"
	"time"

	"github.com/crashiura/gositemap"
)

func main() {
	sitemapNames := make([]string, 0)
	sitemapName := "sitemap-product%d.xml"

	sm := gositemap.NewSitemap(
		gositemap.DirOpt("./"),
		gositemap.FileNameOpt(fmt.Sprintf(sitemapName, 1)),
		gositemap.CompressOpt(true),
	)

	countSitemapFiles := 1
	maxPages := 135000
	for i := 0; i <= maxPages; i++ {
		u := &gositemap.URL{
			Loc:        fmt.Sprintf("product/inventory?page=%d", i),
			LastMod:    gositemap.XmlTime(time.Now()),
			ChangeFreq: gositemap.Always,
			Priority:   1,
		}
		err := sm.Add(u)
		if err != nil {
			if err == gositemap.ErrorFileValidation {
				sitemapNames = append(sitemapNames, sm.Build())
				countSitemapFiles++
				sm = gositemap.NewSitemap(
					gositemap.DirOpt("./"),
					gositemap.FileNameOpt(fmt.Sprintf(sitemapName, countSitemapFiles)),
					gositemap.CompressOpt(true),
				)
				sm.Add(u)
				continue
			}
			log.Fatal("error", err)
		}
		if i == maxPages {
			sitemapNames = append(sitemapNames, sm.Build())
		}
	}

	indexSm := gositemap.NewIndexSitemap(
		gositemap.DirIndexOpt("./"),
		gositemap.FileNameIndexOpt("index-sitemap.xml"),
		gositemap.CompressIndexOpt(true),
	)

	for _, v := range sitemapNames {
		if err := indexSm.Add(&gositemap.SitemapEntity{
			Loc:     "public" + v,
			LastMod: gositemap.XmlTime(time.Now()),
		}); err != nil {
			log.Fatal("error generate index sitemap", err)
		}

	}
	indexSm.Build()
}
