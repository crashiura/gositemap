package gositemap

import (
	"bytes"
	"encoding/xml"
	"log"
	"path/filepath"
)

type IndexSitemapOpt func(*IndexSitemap)

type IndexSitemap struct {
	fileName      string
	dir           string
	content       []byte
	publicPath    string
	host          string
	countSitemaps int
	compress      bool
}

type SitemapEntity struct {
	XMLName xml.Name `xml:"sitemap"`
	Loc     string   `xml:"loc"`
	LastMod XmlTime  `xml:"lastmod,omitempty"`
}

func NewIndexSitemap(opts ...IndexSitemapOpt) *IndexSitemap {
	s := &IndexSitemap{
		fileName: "sitemap.xml",
		dir:      "./",
		content:  make([]byte, 0, 256),
		host:     "https://example.com",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *IndexSitemap) Add(e *SitemapEntity) error {
	e.Loc = s.host + "/" + s.publicPath + e.Loc
	b, err := xml.Marshal(e)
	if err != nil {
		log.Println("Error marshal url ", err)
		return err
	}

	if s.validate(append(s.content, b...)) && s.countSitemaps+1 <= maxUrls {
		s.content = append(s.content, b...)
		s.countSitemaps++
		return nil
	}

	return ErrorFileValidation
}

func (s *IndexSitemap) validate(content []byte) bool {
	return len(content)+len(indexSitemapXMLHeader)+len(indexSitemapXMLFooter) < MaxFileSize
}

func (s *IndexSitemap) Build() {
	fullPath := filepath.Join(
		s.dir,
		s.fileName,
	)
	writeFile(s.dir, fullPath, s.GetXml(), s.compress)
}

func (s *IndexSitemap) GetXml() []byte {
	c := bytes.Join(bytes.Fields(indexSitemapXMLHeader), []byte(" "))
	c = append(append(c, s.content...), indexSitemapXMLFooter...)
	return c
}

func DirIndexOpt(dir string) IndexSitemapOpt {
	return func(sitemap *IndexSitemap) {
		sitemap.dir = dir
	}
}

func FileNameIndexOpt(fileName string) IndexSitemapOpt {
	return func(sitemap *IndexSitemap) {
		sitemap.fileName = fileName
	}
}

func HostIndexOpt(host string) IndexSitemapOpt {
	return func(sitemap *IndexSitemap) {
		sitemap.host = host
	}
}

func CompressIndexOpt(compress bool) IndexSitemapOpt {
	return func(sitemap *IndexSitemap) {
		sitemap.compress = compress
	}
}
