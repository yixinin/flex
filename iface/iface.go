package iface

import (
	"errors"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

type Plugin interface {
	Name() string
	Handle(c *fasthttp.RequestCtx) error
	SetConfig(config string) error
}

type PluginError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e PluginError) Error() string {
	return fmt.Sprintf("plugin error: [code:%d, message:%s]", e.Code, e.Message)
}

func Error(code int, msg string) error {
	return &PluginError{
		Code:    code,
		Message: msg,
	}
}

func Wrap(err error, msg ...string) error {
	if len(msg) > 0 {
		return &PluginError{
			Code:    400,
			Message: fmt.Sprintf("%s %s", strings.Join(msg, ","), err.Error()),
		}
	}
	return &PluginError{
		Code:    400,
		Message: err.Error(),
	}
}

var (
	ErrorAbort = errors.New("abort")
)
