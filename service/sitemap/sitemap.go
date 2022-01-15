package sitemap

import (
	"encoding/xml"
	"fmt"
	"news/service/tpl"
	"news/utils"
	"os"
)

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`

	Links []Url `xml:"url"`
}

type Url struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
}

func ListItem2Link(item utils.ListItem) Url {
	dateStr := fmt.Sprintf("%.0f", item.Score)
	return Url{
		Loc:        tpl.GetUrlByDateStr(dateStr),
		LastMod:    dateStr,
		ChangeFreq: "monthly",
		Priority:   0.8,
	}
}

func New() *Sitemap {
	return &Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Links: make([]Url, 0),
	}
}

func (s *Sitemap) Add(link Url) {
	s.Links = append(s.Links, link)
}

func (s *Sitemap) Save() {
	content, err := xml.Marshal(s)
	if err != nil {
		panic(err)
	}
	sitemapContent := append([]byte(xml.Header), content...)
	err = os.WriteFile(utils.AbsolutPath("/cache/sitemap.xml"), sitemapContent, 0777)
	if err != nil {
		panic(err)
	}
}
