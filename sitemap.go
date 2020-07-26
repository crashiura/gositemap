package gositemap

import (
	"bytes"
	"encoding/xml"
	"path/filepath"
	"time"
)

//ChangeFreq enum for frequency
type ChangeFreq string

const (
	//Always ...
	Always ChangeFreq = "always"
	//Hourly ...
	Hourly ChangeFreq = "hourly"
	//Daily ...
	Daily ChangeFreq = "daily"
	//Weekly ...
	Weekly ChangeFreq = "weekly"
	//Monthly ...
	Monthly ChangeFreq = "monthly"
	//Yearly ...
	Yearly ChangeFreq = "yearly"
	//Never ...
	Never ChangeFreq = "never"
)

//Sitemap secondary sitemap struct
type Sitemap struct {
	fileName string
	dir      string
	content  []byte
	countURL int
	compress bool
	host     string
}

//XMLTime alias for xml time format
type XMLTime time.Time

//SitemapOpt option for sitemap
type SitemapOpt func(*Sitemap)

//MarshalXML decorate sitemap time format
func (xt XMLTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := time.Time(xt)
	if t.IsZero() {
		t = time.Now()
	}
	v := t.Format(time.RFC3339)

	return e.EncodeElement(v, start)
}

//URL representation for url sitemap
type URL struct {
	XMLName    xml.Name   `xml:"url"`
	Loc        string     `xml:"loc"`
	LastMod    XMLTime    `xml:"lastmod,omitempty"`
	ChangeFreq ChangeFreq `xml:"changefreq,omitempty"`
	Priority   float32    `xml:"priority,omitempty"`
}

//NewSitemap creates new secondary sitemap
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

//Add new element in secondary sitemap
func (s *Sitemap) Add(u *URL) error {
	u.Loc = s.host + "/" + u.Loc
	b, err := xml.Marshal(u)
	if err != nil {
		return err
	}

	if s.validate(append(s.content, b...)) && s.countURL+1 <= maxUrls {
		s.content = append(s.content, b...)
		s.countURL++
		return nil
	}

	return ErrorAddEntity
}

//Build create secondary sitemap
func (s *Sitemap) Build() (string, error) {
	fullPath := filepath.Join(s.dir, s.fileName)
	err := writeFile(s.dir, fullPath, s.GetXML(), s.compress)

	return s.fileName, err
}

//GetXML return xml byte representation
func (s *Sitemap) GetXML() []byte {
	c := bytes.Join(bytes.Fields(sitemapHeader), []byte(" "))
	c = append(append(c, s.content...), sitemapFooter...)
	return c
}

func (s *Sitemap) validate(content []byte) bool {
	return len(content)+len(sitemapHeader)+len(sitemapFooter) < maxFileSize
}

//DirOpt add directory for sitemap
func DirOpt(dir string) SitemapOpt {
	return func(sitemap *Sitemap) {
		sitemap.dir = dir
	}
}

//FileNameOpt set file name for sitemap
func FileNameOpt(fileName string) SitemapOpt {
	return func(sitemap *Sitemap) {
		sitemap.fileName = fileName
	}
}

//HostOpt set host in sitemap
func HostOpt(host string) SitemapOpt {
	return func(sitemap *Sitemap) {
		sitemap.host = host
	}
}

//CompressOpt compress sitemap gz
func CompressOpt(compress bool) SitemapOpt {
	return func(sitemap *Sitemap) {
		sitemap.compress = compress
	}
}
