package sitemap

import (
	"encoding/xml"
	"fmt"
	"github.com/spf13/viper"
	"news/service/tpl"
	"news/utils"
	"os"
)

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`

	Links []Link `xml:"url"`
}

type Link struct {
	Loc string `xml:"loc"`
	//LastMod    string  `xml:"lastmod"`
	//ChangeFreq string  `xml:"changefreq"`
	Priority float64 `xml:"priority"`
}

func ListItem2Link(item utils.ListItem) Link {
	dateStr := fmt.Sprintf("%.0f", item.Score)
	return Link{
		Loc: tpl.GetUrlByDateStr(dateStr),
		//LastMod:  dateStr,
		Priority: 0.8,
	}
}

func New() *Sitemap {
	sitemap := &Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Links: make([]Link, 0),
	}

	sitemap.Links = append(sitemap.Links, Link{
		Loc:      viper.GetString("app.baseUrl"),
		Priority: 0.9,
	})

	return sitemap
}

func (s *Sitemap) Add(link Link) {
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
