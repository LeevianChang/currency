package http

import (
	"encoding/json"
	"fmt"
	"go-currency/model"
	"go-currency/service"
	"go-currency/tool"
	"strconv"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
)

func Init(eng *gin.Engine) {
	eng.POST("/v1/currencies/live", getCurrencyLive)
	eng.POST("/v1/currencies/list", getCurrencyList)
	eng.POST("/v1/currencies/timeframe", getCurrencyFrame)
	eng.POST("/v1/currencies/change", getCurrencyChange)

}

type Params struct {
	TimeStamp int64             `json:"timestamp"`
	Params    map[string]string `json:"req_params"`
}

type FrameParams struct {
	TimeStamp int64             `json:"timestamp"`
	Params    map[string]string `json:"req_params"`
}

type ChangeParams struct {
	TimeStamp int64             `json:"timestamp"`
	Params    map[string]string `json:"req_params"`
}

const (
	XTimeStamp = "X-Timestamp"
)

func getCurrencyLive(c *gin.Context) {
	v := new(struct {
		Params    string `form:"params" validate:"required"`
		TimeStamp int64  `form:"timestamp"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	header := c.Request.Header
	if header != nil {
		times, ok := header[XTimeStamp]
		fmt.Println(times, "times")
		if !ok {
			c.JSON(200, &model.Reply{
				Message: model.CodeErrTimeMessage,
				Errno:   model.CodeErrTime,
			})
			return
		}
		if len(times) == 1 {
			v.TimeStamp, _ = strconv.ParseInt(times[0], 10, 64)
		}
	}
	now := time.Now().Unix()
	if v.TimeStamp < now-120 || v.TimeStamp > now+120 {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrTimeMessage,
			Errno:   model.CodeErrTime,
		})
		return
	}
	iv := service.CurrencyService.GetIV(v.TimeStamp)

	params, err := tool.GcmDecrypt(v.Params, "6143ec9acb9160154306ffb7d12ee141", []byte(iv))
	if err != nil {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrDecryptMessage,
			Errno:   model.CodeErrDecrypt,
		})
		return
	}
	p := &Params{}
	err = json.Unmarshal([]byte(params), p)
	fmt.Println(p.TimeStamp, v.TimeStamp)
	if p.TimeStamp != v.TimeStamp {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrTimeMessage,
			Errno:   model.CodeErrTime,
		})
		return
	}
	data, t := service.CurrencyService.GetCurrencyLive(c)
	c.Header(XTimeStamp, fmt.Sprintf("%d", t))
	c.JSON(200, data)
}

func getCurrencyList(c *gin.Context) {
	v := new(struct {
		Params    string `form:"params" validate:"required"`
		TimeStamp int64  `form:"timestamp"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	header := c.Request.Header
	if header != nil {
		times, ok := header[XTimeStamp]
		if !ok {
			c.JSON(200, &model.Reply{
				Message: model.CodeErrTimeMessage,
				Errno:   model.CodeErrTime,
			})
			return
		}
		if len(times) == 1 {
			v.TimeStamp, _ = strconv.ParseInt(times[0], 10, 64)
		}
	}
	now := time.Now().Unix()
	if v.TimeStamp < now-120 || v.TimeStamp > now+120 {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrTimeMessage,
			Errno:   model.CodeErrTime,
		})
		return
	}
	iv := service.CurrencyService.GetIV(v.TimeStamp)

	params, err := tool.GcmDecrypt(v.Params, "6143ec9acb9160154306ffb7d12ee141", []byte(iv))
	if err != nil {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrDecryptMessage,
			Errno:   model.CodeErrDecrypt,
		})
		return
	}
	p := &Params{}
	err = json.Unmarshal([]byte(params), p)
	fmt.Println(p.TimeStamp, v.TimeStamp)
	if p.TimeStamp != v.TimeStamp {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrTimeMessage,
			Errno:   model.CodeErrTime,
		})
		return
	}
	data, t := service.CurrencyService.GetCurrencyList(c)
	c.Header(XTimeStamp, fmt.Sprintf("%d", t))
	c.JSON(200, data)
}

func getCurrencyFrame(c *gin.Context) {
	v := new(struct {
		Params    string `form:"params" validate:"required"`
		TimeStamp int64  `form:"timestamp"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	header := c.Request.Header
	if header != nil {
		times, ok := header[XTimeStamp]
		if !ok {
			c.JSON(200, &model.Reply{
				Message: model.CodeErrTimeMessage,
				Errno:   model.CodeErrTime,
			})
			return
		}
		if len(times) == 1 {
			v.TimeStamp, _ = strconv.ParseInt(times[0], 10, 64)
		}
	}
	now := time.Now().Unix()
	if v.TimeStamp < now-120 || v.TimeStamp > now+120 {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrTimeMessage,
			Errno:   model.CodeErrTime,
		})
		return
	}
	iv := service.CurrencyService.GetIV(v.TimeStamp)

	params, err := tool.GcmDecrypt(v.Params, "6143ec9acb9160154306ffb7d12ee141", []byte(iv))
	if err != nil {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrDecryptMessage,
			Errno:   model.CodeErrDecrypt,
		})
		return
	}
	p := &FrameParams{}
	err = json.Unmarshal([]byte(params), p)
	fmt.Println(p.TimeStamp, v.TimeStamp)
	if p.TimeStamp != v.TimeStamp {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrTimeMessage,
			Errno:   model.CodeErrTime,
		})
		return
	}
	data, t := service.CurrencyService.GetCurrencyTimeFrame(c, p.Params)
	c.Header(XTimeStamp, fmt.Sprintf("%d", t))
	c.JSON(200, data)
}

func Limiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.JSON(httpError.StatusCode, &model.Reply{
				Message: model.CodeErrLimitMessage,
				Errno:   model.CodeLimit,
			})
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func getCurrencyChange(c *gin.Context) {
	v := new(struct {
		Params    string `form:"params" validate:"required"`
		TimeStamp int64  `form:"timestamp"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	header := c.Request.Header
	if header != nil {
		times, ok := header[XTimeStamp]
		if !ok {
			c.JSON(200, &model.Reply{
				Message: model.CodeErrTimeMessage,
				Errno:   model.CodeErrTime,
			})
			return
		}
		if len(times) == 1 {
			v.TimeStamp, _ = strconv.ParseInt(times[0], 10, 64)
		}
	}
	now := time.Now().Unix()
	if v.TimeStamp < now-120 || v.TimeStamp > now+120 {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrTimeMessage,
			Errno:   model.CodeErrTime,
		})
		return
	}
	iv := service.CurrencyService.GetIV(v.TimeStamp)

	params, err := tool.GcmDecrypt(v.Params, "6143ec9acb9160154306ffb7d12ee141", []byte(iv))
	if err != nil {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrDecryptMessage,
			Errno:   model.CodeErrDecrypt,
		})
		return
	}
	p := &ChangeParams{}
	err = json.Unmarshal([]byte(params), p)
	fmt.Println(p.TimeStamp, v.TimeStamp)
	if p.TimeStamp != v.TimeStamp {
		c.JSON(200, &model.Reply{
			Message: model.CodeErrTimeMessage,
			Errno:   model.CodeErrTime,
		})
		return
	}
	data, t := service.CurrencyService.GetCurrencyTimeChange(c, p.Params)
	c.Header(XTimeStamp, fmt.Sprintf("%d", t))
	c.JSON(200, data)
}
