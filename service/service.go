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
	CurrencyLive *CurrencyLive
	CurrencyList *CurrencyList
}

const (
	accessKey    = "8WNG6RGp8vKEJPFnzvX23JwGGDaGupyw"
	liveUrl      = "https://api.apilayer.com/currency_data/live"
	listUrl      = "https://api.apilayer.com/currency_data/list"
	timeFrameUrl = "https://api.apilayer.com/currency_data/timeframe"
	changeUrl    = "https://api.apilayer.com/currency_data/change"
)

type CurrencyLive struct {
	Success   bool               `json:"success"`
	Terms     string             `json:"terms"`
	Privacy   string             `json:"privacy"`
	Timestamp int64              `json:"timestamp"`
	Source    string             `json:"source"`
	Quotes    map[string]float64 `json:"quotes"`
}

type CurrencyTimeFrame struct {
	EndDate   string                 `json:"end_date"`
	Quotes    map[string]interface{} `json:"quotes"`
	Source    string                 `json:"source"`
	StartDate string                 `json:"start_date"`
	Success   bool                   `json:"success"`
	TimeFrame bool                   `json:"timeframe"`
}

type CurrencyTimeChange struct {
	Source    string                 `json:"source"`
	StartDate string                 `json:"start_date"`
	Change    bool                   `json:"change"`
	Quotes    map[string]interface{} `json:"quotes"`
	EndDate   string                 `json:"end_date"`
	Success   bool                   `json:"success"`
}

type CurrencyList struct {
	Success bool              `json:"success"`
	Symbols map[string]string `json:"symbols"`
}

func Init() {
	CurrencyService = &Service{}
	ctx := context.Background()
	CurrencyService.HttpGetCurrencyLive(ctx)
	CurrencyService.HttpGetCurrencyList(ctx)
	go CurrencyService.LoopGetCurrencyLive(ctx)
	go CurrencyService.LoopGetCurrencyList(ctx)
}

func (s *Service) LoopGetCurrencyLive(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Minute)
	for range ticker.C {
		go CurrencyService.HttpGetCurrencyLive(ctx)
	}
}

func (s *Service) LoopGetCurrencyList(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Minute)
	for range ticker.C {
		go CurrencyService.HttpGetCurrencyList(ctx)
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

func (s *Service) GetCurrencyLive(ctx context.Context) (res *model.Reply, timestamp int64) {
	data := make(map[string]interface{})
	currency, err := json.Marshal(s.CurrencyLive)
	if err != nil {
		return
	}
	timestamp = s.CurrencyLive.Timestamp
	iv := s.GetIV(timestamp)
	b, err := tool.GcmEncrypt(string(currency), "6143ec9acb9160154306ffb7d12ee141", []byte(iv))
	if err != nil {
		return &model.Reply{
			Errno:   model.CodeErrEncrypt,
			Message: model.CodeErrEncryptMessage,
		}, 0
	}
	data["ed"] = b
	res = &model.Reply{
		Data: data,
	}
	return
}

func (s *Service) GetCurrencyList(ctx context.Context) (res *model.Reply, timestamp int64) {
	data := make(map[string]interface{})
	currency, err := json.Marshal(s.CurrencyList)
	if err != nil {
		return
	}
	timestamp = time.Now().Unix()
	iv := s.GetIV(timestamp)
	b, err := tool.GcmEncrypt(string(currency), "6143ec9acb9160154306ffb7d12ee141", []byte(iv))
	if err != nil {
		return &model.Reply{
			Errno:   model.CodeErrEncrypt,
			Message: model.CodeErrEncryptMessage,
		}, 0
	}
	data["ed"] = b
	res = &model.Reply{
		Data: data,
	}
	return
}

func (s *Service) HttpGetCurrencyLive(ctx context.Context) (code int, res *CurrencyLive) {
	// 处理query
	//query := fmt.Sprintf("?access_key=%s", accessKey)
	req, _ := http.NewRequest("GET", liveUrl, nil)
	// 处理header
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "api.currencylayer.com")
	req.Header.Set("Accept", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("apikey", accessKey)
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	resp, err := (&http.Client{Timeout: time.Millisecond * 5000}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &s.CurrencyLive)
	fmt.Println(s.CurrencyLive)
	res = s.CurrencyLive
	return
}

func (s *Service) HttpGetCurrencyList(ctx context.Context) (code int, res *CurrencyList) {
	// 处理query
	//query := fmt.Sprintf("?access_key=%s", accessKey)
	req, _ := http.NewRequest("GET", listUrl, nil)
	// 处理header
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "api.currencylayer.com")
	req.Header.Set("Accept", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("apikey", accessKey)
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	resp, err := (&http.Client{Timeout: time.Millisecond * 5000}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &s.CurrencyList)
	fmt.Println(s.CurrencyList)
	res = s.CurrencyList
	return
}

func (s *Service) GetCurrencyTimeFrame(ctx context.Context, params map[string]string) (res *model.Reply, timestamp int64) {
	data := make(map[string]interface{})
	_, d := s.HttpGetCurrencyTimeFrame(ctx, params)
	currency, err := json.Marshal(d)
	if err != nil {
		return
	}
	timestamp = time.Now().Unix()
	iv := s.GetIV(timestamp)
	fmt.Println(currency, "currency")
	b, err := tool.GcmEncrypt(string(currency), "6143ec9acb9160154306ffb7d12ee141", []byte(iv))
	if err != nil {
		return &model.Reply{
			Errno:   model.CodeErrEncrypt,
			Message: model.CodeErrEncryptMessage,
		}, 0
	}
	data["ed"] = b
	res = &model.Reply{
		Data: data,
	}
	return
}

func (s *Service) GetCurrencyTimeChange(ctx context.Context, params map[string]string) (res *model.Reply, timestamp int64) {
	data := make(map[string]interface{})
	_, d := s.HttpGetCurrencyChange(ctx, params)
	currency, err := json.Marshal(d)
	if err != nil {
		return
	}
	timestamp = time.Now().Unix()
	iv := s.GetIV(timestamp)
	fmt.Println(currency, "currency")
	b, err := tool.GcmEncrypt(string(currency), "6143ec9acb9160154306ffb7d12ee141", []byte(iv))
	if err != nil {
		return &model.Reply{
			Errno:   model.CodeErrEncrypt,
			Message: model.CodeErrEncryptMessage,
		}, 0
	}
	data["ed"] = b
	res = &model.Reply{
		Data: data,
	}
	return
}

func (s *Service) HttpGetCurrencyTimeFrame(ctx context.Context, params map[string]string) (code int, res *CurrencyTimeFrame) {
	// 处理query
	var q string
	for k, v := range params {
		q += fmt.Sprintf("%s=%s&", k, v)
	}
	q = "?" + q
	req, _ := http.NewRequest("GET", timeFrameUrl+q, nil)
	// 处理header
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "api.currencylayer.com")
	req.Header.Set("Accept", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("apikey", accessKey)
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	resp, err := (&http.Client{Timeout: time.Millisecond * 1000}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res = &CurrencyTimeFrame{}
	json.Unmarshal(body, res)
	return
}

func (s *Service) HttpGetCurrencyChange(ctx context.Context, params map[string]string) (code int, res *CurrencyTimeChange) {
	// 处理query
	var q string
	for k, v := range params {
		q += fmt.Sprintf("%s=%s&", k, v)
	}
	q = "?" + q
	req, _ := http.NewRequest("GET", changeUrl+q, nil)
	// 处理header
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "api.currencylayer.com")
	req.Header.Set("Accept", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("apikey", accessKey)
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	resp, err := (&http.Client{Timeout: time.Millisecond * 1000}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res = &CurrencyTimeChange{}
	json.Unmarshal(body, res)
	return
}
