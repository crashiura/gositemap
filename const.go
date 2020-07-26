package gositemap

import "errors"

const (
	maxFileSize = 52428800
	maxUrls     = 50000
)

//ErrorAddEntity error add entity in sitemap
var ErrorAddEntity = errors.New("error max file size or max elements in file")

var (
	indexSitemapXMLHeader = []byte(`<?xml version="1.0" encoding="UTF-8"?>
      <sitemapindex
      xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
      xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9
        http://www.sitemaps.org/schemas/sitemap/0.9/siteindex.xsd"
      xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
    >`)
	indexSitemapXMLFooter = []byte("</sitemapindex>")
)

var (
	sitemapHeader = []byte(`<?xml version="1.0" encoding="UTF-8"?>
      <urlset
        xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9
          http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
        xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
        xmlns:xhtml="http://www.w3.org/1999/xhtml"
    >`)
	sitemapFooter = []byte("</urlset>")
)
