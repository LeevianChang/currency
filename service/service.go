package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go-currency/model"
	"go-currency/tool"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var CurrencyService *Service

type Service struct {
	Currency *Currency
}

const (
	accessKey = "a3f9c33d76d640112a2c3d3fd83e456d"
	url       = "http://api.currencylayer.com/live"
)

type Currency struct {
	Success   bool               `json:"success"`
	Terms     string             `json:"terms"`
	Privacy   string             `json:"privacy"`
	Timestamp int64              `json:"timestamp"`
	Source    string             `json:"source"`
	Quotes    map[string]float64 `json:"quotes"`
}

func Init() {
	CurrencyService = &Service{}
	ctx := context.Background()
	CurrencyService.HttpGetCurrency(ctx)
	go CurrencyService.LoopGetCurrency(ctx)
}

func (s *Service) LoopGetCurrency(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		go CurrencyService.HttpGetCurrency(ctx)
	}
}

func (s *Service) GetIV(timeStamp int64) (res string) {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%d", timeStamp))
	md := fmt.Sprintf("%x", h.Sum(nil))
	list := make([][]string, 0)
	item := make([]string, 0)
	for i := 0; i < len(md); i++ {
		if (i+1)%8 == 0 {
			list = append(list, item)
			item = make([]string, 0)
		} else {
			if len(item) < 3 {
				item = append(item, string(md[i]))
			}
		}
	}
	for _, v := range list {
		res += strings.Join(v, "")
	}
	return
}

func (s *Service) GetCurrency(ctx context.Context, iv string) (res *model.Reply) {
	data := make(map[string]interface{})
	currency, err := json.Marshal(s.Currency)
	if err != nil {
		return
	}
	b, err := tool.GcmEncrypt(string(currency), "6143ec9acb9160154306ffb7d12ee141", []byte(iv))
	if err != nil {
		return &model.Reply{
			Errno:   model.CodeErrEncrypt,
			Message: model.CodeErrEncryptMessage,
		}
	}
	data["ed"] = b
	res = &model.Reply{
		Data: data,
	}
	return
}
func (s *Service) HttpGetCurrency(ctx context.Context) (code int, res *Currency) {
	// 处理query
	query := fmt.Sprintf("?access_key=%s", accessKey)

	req, _ := http.NewRequest("GET", url+query, nil)
	// 处理header
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "api.currencylayer.com")
	req.Header.Set("Accept", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	resp, err := (&http.Client{Timeout: time.Millisecond * 5000}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &s.Currency)
	fmt.Println(s.Currency)
	res = s.Currency
	return
}
