package gositemap

import (
	"fmt"
	"time"
)

//GeneratorOpt option for generator sitemaps
type GeneratorOpt func(*Generator)

//SecondarySitemap represent seconadary sitemap entity for additional sitemap in main sitemap
type SecondarySitemap struct {
	placeholder string
	count       int
	sm          *Sitemap
	fileNames   map[string]struct{}
}

//Generator sitemap generator
type Generator struct {
	secondarySitemap map[string]*SecondarySitemap
	compress         bool
	indexSitemapName string
	dir              string
	path             string
	host             string
}

//DirGenOpt set dir for sitemaps
func DirGenOpt(dir string) GeneratorOpt {
	return func(generator *Generator) {
		generator.dir = dir
	}
}

//CompressGenOpt compress sitemap option
func CompressGenOpt(compress bool) GeneratorOpt {
	return func(generator *Generator) {
		generator.compress = compress
	}
}

//IndexSitemapGenOpt set sitemap file name
func IndexSitemapGenOpt(fileName string) GeneratorOpt {
	return func(generator *Generator) {
		generator.indexSitemapName = fileName
	}
}

//PathSitemapGenOpt web path for sitemap
func PathSitemapGenOpt(path string) GeneratorOpt {
	return func(generator *Generator) {
		generator.path = path
	}
}

//HostSitemapGenOpt set host for sitemap
func HostSitemapGenOpt(hostName string) GeneratorOpt {
	return func(generator *Generator) {
		generator.host = hostName
	}
}

//NewGenerator creates new sitemap generator
func NewGenerator(opts ...GeneratorOpt) *Generator {
	g := &Generator{
		secondarySitemap: make(map[string]*SecondarySitemap),
		compress:         false,
		indexSitemapName: "index-sitemap.xml",
		dir:              "./web",
		path:             "public",
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

//GetSecondarySitemap create secondary sitemap example: placeholder product%d.xml
func (g *Generator) GetSecondarySitemap(placeholder string) *SecondarySitemap {
	sm := NewSitemap(
		DirOpt(g.dir),
		FileNameOpt(fmt.Sprintf(placeholder, 0)),
		CompressOpt(g.compress),
		HostOpt(g.host),
	)

	ssm := &SecondarySitemap{
		placeholder: placeholder,
		sm:          sm,
		fileNames:   make(map[string]struct{}),
	}

	g.secondarySitemap[placeholder] = ssm

	return ssm
}

//Build build sitemap
func (g *Generator) Build() error {
	indexSm := NewIndexSitemap(
		DirIndexOpt(g.dir),
		FileNameIndexOpt(g.indexSitemapName),
		CompressIndexOpt(g.compress),
		HostIndexOpt(g.host),
		PublicPathIndexOpt(g.path),
	)

	for _, v := range g.secondarySitemap {
		if _, ok := v.fileNames[v.sm.fileName]; !ok {
			fileName, err := v.sm.Build()
			if err != nil {
				return err
			}
			v.fileNames[fileName] = struct{}{}
		}

		for ssm := range v.fileNames {
			if err := indexSm.Add(&SitemapEntity{
				Loc:     "/" + ssm,
				LastMod: XMLTime(time.Now()),
			}); err != nil {
				return err
			}
		}
	}

	if err := indexSm.Build(); err != nil {
		return err
	}

	return nil
}

//Add adds entity in sitemap
func (ssm *SecondarySitemap) Add(u *URL) error {
	if err := ssm.sm.Add(u); err != nil {
		if err != ErrorAddEntity {
			return err
		}

		fileName, err := ssm.sm.Build()
		if err != nil {
			return err
		}

		ssm.count++
		ssm.fileNames[fileName] = struct{}{}

		newSm := NewSitemap(
			DirOpt(ssm.sm.dir),
			FileNameOpt(fmt.Sprintf(ssm.placeholder, ssm.count)),
			CompressOpt(ssm.sm.compress),
			HostOpt(ssm.sm.host),
		)

		ssm.sm = newSm
	}

	return nil
}
