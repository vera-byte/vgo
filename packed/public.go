package packed

import "github.com/gogf/gf/v2/os/gres"

func init() {
	if err := gres.Add("H4sIAAAAAAAC/wrwZmYRYeBgYGC4KmgQzYAE+Bg4GQpKk3Iyk/VLC3LyE1OKQ0NYGRiFT5zICPBm50BWijBkpq0ViiFscEPAmotcrmHRzMgkwozbDRAgwPDWEUTjdBHCEGxugBny39EVbgiSi1jZQNJMDEwMTQwMYMzAAAgAAP//JZ1s1RoBAAA="); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
