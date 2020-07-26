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
		gositemap.CompressOpt(false),
	)

	countSitemapFiles := 1
	maxPages := 135000
	for i := 0; i <= maxPages; i++ {
		u := &gositemap.URL{
			Loc:        fmt.Sprintf("product/inventory?page=%d", i),
			LastMod:    gositemap.XMLTime(time.Now()),
			ChangeFreq: gositemap.Always,
			Priority:   1,
		}
		err := sm.Add(u)
		if err != nil {
			if err == gositemap.ErrorAddEntity {
				name, err := sm.Build()
				if err != nil {
					// handle err
				}
				sitemapNames = append(sitemapNames, name)
				countSitemapFiles++
				sm = gositemap.NewSitemap(
					gositemap.DirOpt("./"),
					gositemap.FileNameOpt(fmt.Sprintf(sitemapName, countSitemapFiles)),
					gositemap.CompressOpt(false),
				)
				sm.Add(u)
				continue
			}
			log.Fatal("error", err)
		}

		if i == maxPages {
			name, err := sm.Build()
			if err != nil {
				// handle err
			}
			sitemapNames = append(sitemapNames, name)
		}
	}

	indexSm := gositemap.NewIndexSitemap(
		gositemap.DirIndexOpt("./"),
		gositemap.FileNameIndexOpt("index-sitemap.xml"),
		gositemap.CompressIndexOpt(false),
	)

	for _, v := range sitemapNames {
		if err := indexSm.Add(&gositemap.SitemapEntity{
			Loc:     "public" + v,
			LastMod: gositemap.XMLTime(time.Now()),
		}); err != nil {
			log.Fatal("error generate index sitemap", err)
		}

	}
	indexSm.Build()
}
