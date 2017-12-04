package utils

/*
*新浪短链接转长链接接口
*/
import (
	"encoding/json"
	"fmt"
	"zcm_tools/http"
)

type Url struct {
	Result    bool `json:"result"`
	UrlShort  string `json:"url_short"`
	UrlLong   string `json:"url_long"`
	Type      int `json:"type"`
	TranScode int `json:"transcode"`
}
type LongUrl struct {
	Urls []Url `json:"urls"`
}

//解析短链接到长链接
func SortToLong(accessToken, urlShort string) (newUrl string, err error) {
	url := SHORTURLTOLANG + "?access_token=" + accessToken + "&url_short=" + urlShort
	b, err := http.Get(url)
	if err != nil {
		fmt.Println(29,err.Error())
		return "", err
	} else {
		var m LongUrl
		if err := json.Unmarshal(b, &m); err == nil {
			if len(m.Urls) > 0 {
				return m.Urls[0].UrlLong, err
			}
		}
		return "", err
	}
}
