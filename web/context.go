package web

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/steambap/captcha"
	"image/color"
	"micro-libs/utils/dtype"
	"micro-libs/utils/tool"
	"net/http"
)

type Context struct {
	ctx echo.Context
}

func (c *Context) Echo() echo.Context {
	return c.ctx
}

func (c *Context) Error(err error) error {
	return CtxResult(c.ctx, ParseError(err))
}

func (c *Context) JsonError(code int, msg ...interface{}) error {
	return CtxError(c.ctx, code, msg...)
}

func (c *Context) JsonSuccess(data interface{}, msg ...string) error {
	return CtxSuccess(c.ctx, data, msg...)
}

func (c *Context) Bind(req interface{}) error {
	return c.ctx.Bind(req)
}

func (c *Context) BindValid(req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.ctx.Validate(req); err != nil {
		return err
	}
	return nil
}

func (c *Context) Session() *sessions.Session {
	s, _ := session.Get("SESSION", c.Echo())
	return s
}

func (c *Context) SessionDo(closure func(*sessions.Session) interface{}) interface{} {
	s := c.Session()
	res := closure(s)
	_ = s.Save(c.ctx.Request(), c.ctx.Response())
	return res
}

func (c *Context) SessId() string {
	res := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.ID
	})
	return res.(string)
}

func (c *Context) SessOpts() *sessions.Options {
	res := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.Options
	})
	return res.(*sessions.Options)
}

func (c *Context) SessOptsSet(opts *sessions.Options) {
	_ = c.SessionDo(func(ss *sessions.Session) interface{} {
		ss.Options = opts
		return nil
	})
}

func (c *Context) SessFlashAdd(val interface{}, key ...string) {
	_ = c.SessionDo(func(ss *sessions.Session) interface{} {
		ss.AddFlash(val, key...)
		return nil
	})
}

func (c *Context) SessFlash(key ...string) []interface{} {
	res := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.Flashes(key...)
	})
	return res.([]interface{})
}

func (c *Context) SessGetValues() map[interface{}]interface{} {
	res := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.Values
	})
	return res.(map[interface{}]interface{})
}

func (c *Context) SessSetValues(values map[interface{}]interface{}) {
	_ = c.SessionDo(func(ss *sessions.Session) interface{} {
		for k, v := range values {
			ss.Values[k] = v
		}
		return nil
	})
}

func (c *Context) SessHas(key interface{}) bool {
	_, b := c.SessGetValues()[key]
	return b
}

func (c *Context) SessGet(key interface{}) interface{} {
	return c.SessGetValues()[key]
}

func (c *Context) SessSet(key, val interface{}) {
	c.SessSetValues(map[interface{}]interface{}{
		key: val,
	})
}

func (c *Context) SessDel(key ...interface{}) {
	c.SessionDo(func(s *sessions.Session) interface{} {
		for _, k := range key {
			if _, ok := s.Values[k]; ok {
				delete(s.Values, k)
			}
		}
		return nil
	})
}

func (c *Context) SessClean() {
	c.SessionDo(func(s *sessions.Session) interface{} {
		s.Values = make(map[interface{}]interface{})
		return nil
	})
}

// 创建验证码
func (c *Context) CaptchaNew(key string, width, height int, setOpt ...captcha.SetOption) error {
	var opt captcha.SetOption
	if len(setOpt) > 0 && setOpt[0] != nil {
		opt = setOpt[0]
	} else {
		opt = func(opt *captcha.Options) {
			opt.BackgroundColor = color.White
			opt.CharPreset = "0123456789"
			opt.CurveNumber = 1
			opt.FontDPI = 80
		}
	}

	data, err := captcha.New(width, height, opt)
	if err != nil {
		return c.Error(err)
	}

	c.SessSet(key, data.Text)

	return data.WriteImage(c.ctx.Response().Writer)
}

// 验证验证码
func (c *Context) CaptchaCheck(key string, captcha string) bool {
	codeText := c.SessionDo(func(s *sessions.Session) interface{} {
		if c, ok := s.Values[key]; ok {
			return c
		}
		return ""
	})
	return dtype.ParseStr(codeText) == captcha
}

func NewContext(ctx echo.Context) *Context {
	c := &Context{
		ctx: ctx,
	}
	return c
}

// BindValid 绑定验证
func BindValid(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	return nil
}

// GetRealIP 获取真实IP
func GetRealIP(c echo.Context) string {
	return tool.GetHttpRealIP(c.Request())
}

// CtxError 返回错误信息
func CtxError(ctx echo.Context, code int, err ...interface{}) error {
	var msg string

	if len(err) > 0 && err[0] != nil {
		switch ee := err[0].(type) {
		case error:
			msg = ee.Error()
		case string:
			msg = ee
		}
	}

	if msg == "" {
		msg = http.StatusText(code)
		if msg == "" {
			msg = "unknown error"
		}
	}

	return CtxResult(ctx, NewResult(code, msg, ""))
}

// CtxSuccess 返回成功消息
func CtxSuccess(ctx echo.Context, data interface{}, msg ...string) error {
	if len(msg) > 0 && msg[0] != "" {
		return CtxResult(ctx, NewResult(0, msg[0], data))
	}
	return CtxResult(ctx, NewResult(0, "success", data))
}

// CtxResult 自定义返回内容
func CtxResult(ctx echo.Context, res *Result) error {
	return ctx.JSON(http.StatusOK, res)
}

func GetStr(c echo.Context, key string) string {
	return dtype.ParseStr(c.Get(key))
}

func GetInt(c echo.Context, key string) int {
	return dtype.ParseInt(c.Get(key))
}

func GetInt32(c echo.Context, key string) int32 {
	return dtype.ParseInt32(c.Get(key))
}

func GetInt64(c echo.Context, key string) int64 {
	return dtype.ParseInt64(c.Get(key))
}

func GetFloat32(c echo.Context, key string) float32 {
	return dtype.ParseFloat32(c.Get(key))
}

func GetFloat64(c echo.Context, key string) float64 {
	return dtype.ParseFloat64(c.Get(key))
}

func GetBool(c echo.Context, key string) bool {
	return dtype.ParseBool(c.Get(key))
}

func GetParamStr(c echo.Context, key string) string {
	return dtype.ParseStr(c.Param(key))
}

func GetParamInt(c echo.Context, key string) int {
	return dtype.ParseInt(c.Param(key))
}

func GetParamInt32(c echo.Context, key string) int32 {
	return dtype.ParseInt32(c.Param(key))
}

func GetParamInt64(c echo.Context, key string) int64 {
	return dtype.ParseInt64(c.Param(key))
}

func GetParamFloat32(c echo.Context, key string) float32 {
	return dtype.ParseFloat32(c.Param(key))
}

func GetParamFloat64(c echo.Context, key string) float64 {
	return dtype.ParseFloat64(c.Param(key))
}

func GetParamBool(c echo.Context, key string) bool {
	return dtype.ParseBool(c.Param(key))
}
