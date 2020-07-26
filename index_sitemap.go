package gositemap

import (
	"bytes"
	"encoding/xml"
	"path/filepath"
)

//IndexSitemapOpt option for sitemap
type IndexSitemapOpt func(*IndexSitemap)

//IndexSitemap index sitemap entity representation
type IndexSitemap struct {
	fileName      string
	dir           string
	content       []byte
	publicPath    string
	host          string
	countSitemaps int
	compress      bool
}

//SitemapEntity sitemap item representation
type SitemapEntity struct {
	XMLName xml.Name `xml:"sitemap"`
	Loc     string   `xml:"loc"`
	LastMod XMLTime  `xml:"lastmod,omitempty"`
}

//NewIndexSitemap creates new index sitemap
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

//Add new sitemap entity in index sitemap
func (s *IndexSitemap) Add(e *SitemapEntity) error {
	e.Loc = s.host + "/" + s.publicPath + e.Loc
	b, err := xml.Marshal(e)
	if err != nil {
		return err
	}

	if s.validate(append(s.content, b...)) && s.countSitemaps+1 <= maxUrls {
		s.content = append(s.content, b...)
		s.countSitemaps++
		return nil
	}

	return ErrorAddEntity
}

func (s *IndexSitemap) validate(content []byte) bool {
	return len(content)+len(indexSitemapXMLHeader)+len(indexSitemapXMLFooter) < maxFileSize
}

//Build index sitemap
func (s *IndexSitemap) Build() error {
	fullPath := filepath.Join(
		s.dir,
		s.fileName,
	)

	return writeFile(s.dir, fullPath, s.GetXML(), s.compress)
}

//GetXML return xml output
func (s *IndexSitemap) GetXML() []byte {
	c := bytes.Join(bytes.Fields(indexSitemapXMLHeader), []byte(" "))
	c = append(append(c, s.content...), indexSitemapXMLFooter...)
	return c
}

//DirIndexOpt set dir for index sitemap
func DirIndexOpt(dir string) IndexSitemapOpt {
	return func(sitemap *IndexSitemap) {
		sitemap.dir = dir
	}
}

//FileNameIndexOpt set filename for index sitemap
func FileNameIndexOpt(fileName string) IndexSitemapOpt {
	return func(sitemap *IndexSitemap) {
		sitemap.fileName = fileName
	}
}

//HostIndexOpt set host for index sitemap
func HostIndexOpt(host string) IndexSitemapOpt {
	return func(sitemap *IndexSitemap) {
		sitemap.host = host
	}
}

//CompressIndexOpt compress to .gz for sitemap
func CompressIndexOpt(compress bool) IndexSitemapOpt {
	return func(sitemap *IndexSitemap) {
		sitemap.compress = compress
	}
}

//PublicPathIndexOpt set public path in sitemap
func PublicPathIndexOpt(path string) IndexSitemapOpt {
	return func(sitemap *IndexSitemap) {
		sitemap.publicPath = path
	}
}
