package gositemap

import (
	"bytes"
	"encoding/xml"
	"log"
	"path/filepath"
	"time"
)

type ChangeFreq string

const (
	Always  ChangeFreq = "always"
	Hourly  ChangeFreq = "hourly"
	Daily   ChangeFreq = "daily"
	Weekly  ChangeFreq = "weekly"
	Monthly ChangeFreq = "monthly"
	Yearly  ChangeFreq = "yearly"
	Never   ChangeFreq = "never"
)

type Sitemap struct {
	fileName string
	dir      string
	content  []byte
	countUrl int
	compress bool
	host     string
}

type XmlTime time.Time

type SitemapOpt func(*Sitemap)

func (xt XmlTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := time.Time(xt)
	if t.IsZero() {
		t = time.Now()
	}
	v := t.Format(time.RFC3339)

	return e.EncodeElement(v, start)
}

type URL struct {
	XMLName    xml.Name   `xml:"url"`
	Loc        string     `xml:"loc"`
	LastMod    XmlTime    `xml:"lastmod,omitempty"`
	ChangeFreq ChangeFreq `xml:"changefreq,omitempty"`
	Priority   float32    `xml:"priority,omitempty"`
}

func NewSitemap(opts ...SitemapOpt) *Sitemap {
	s := &Sitemap{
		fileName: "sitemap.xml",
		dir:      "./",
		content:  make([]byte, 0, 256),
		compress: false,
		host:     "https://example.com",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Sitemap) Add(u *URL) error {
	u.Loc = s.host + "/" + u.Loc
	b, err := xml.Marshal(u)
	if err != nil {
		log.Println("Error marshal url ", err)
		return err
	}

	if s.validate(append(s.content, b...)) && s.countUrl+1 <= maxUrls {
		s.content = append(s.content, b...)
		s.countUrl++
		return nil
	}

	return ErrorFileValidation
}

func (s *Sitemap) Build() string {
	fullPath := filepath.Join(s.dir, s.fileName)
	writeFile(s.dir, fullPath, s.GetXml(), s.compress)
	return s.fileName
}

func (s *Sitemap) GetXml() []byte {
	c := bytes.Join(bytes.Fields(sitemapHeader), []byte(" "))
	c = append(append(c, s.content...), sitemapFooter...)
	return c
}

func (s *Sitemap) validate(content []byte) bool {
	return len(content)+len(sitemapHeader)+len(sitemapFooter) < MaxFileSize
}

func DirOpt(dir string) SitemapOpt {
	return func(sitemap *Sitemap) {
		sitemap.dir = dir
	}
}

func FileNameOpt(fileName string) SitemapOpt {
	return func(sitemap *Sitemap) {
		sitemap.fileName = fileName
	}
}

func HostOpt(host string) SitemapOpt {
	return func(sitemap *Sitemap) {
		sitemap.host = host
	}
}

func CompressOpt(compress bool) SitemapOpt {
	return func(sitemap *Sitemap) {
		sitemap.compress = compress
	}
}
