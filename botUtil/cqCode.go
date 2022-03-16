package botUtil

import (
	"context"
	"github.com/Mmx233/tool"
	"html"
	"strings"
)

type cq struct{}

var Cq cq

type CqCode struct {
	Type string
	Data map[string]string
}

func (*cq) RMType(c context.Context, Type string) {
	Bot.WriteMsg(c, tool.Regexp.Replace(`\[CQ:`+Type+`,.*?\]`, Bot.GetMsg(c), ""))
}

func (a *cq) RMAll(c context.Context) {
	a.RMType(c, ".*?")
}

func (a *cq) DecodeType(c context.Context, Type string, del bool) []*CqCode {
	r := `\[CQ:` + Type + `,(.*?)\]`
	var d []*CqCode
	t := tool.Regexp.MatchValue(r, Bot.GetMsg(c))
	for _, v := range t {
		var c = CqCode{
			Type: Type,
			Data: make(map[string]string),
		}
		for _, vv := range strings.Split(v[1], ",") {
			vvv := strings.Split(vv, "=")
			c.Data[vvv[0]] = vvv[1]
		}
		d = append(d, &c)
	}
	if del {
		a.RMType(c, Type)
	}
	return d
}

func (a *cq) Decode(c context.Context, del bool) []*CqCode {
	var d []*CqCode
	t := tool.Regexp.MatchValue(`\[CQ:(.*?),(.*?)\]`, Bot.GetMsg(c))
	for _, v := range t {
		var c = CqCode{
			Type: v[1],
			Data: make(map[string]string),
		}
		for _, vv := range strings.Split(v[2], ",") {
			vvv := strings.Split(vv, "=")
			c.Data[vvv[0]] = vvv[1]
		}
		d = append(d, &c)
	}
	if del {
		a.RMAll(c)
	}
	return d
}

func (*cq) Make(Type string, data map[string]string) string {
	Type = "[CQ:" + Type
	for k, v := range data {
		Type = Type + "," + k + "=" + v
	}
	Type = Type + "]"
	return Type
}

func (*cq) DataConvert(a string) string {
	t := map[string]string{
		",": "，",
	}
	for k, v := range t {
		a = strings.Replace(a, k, v, -1)
	}
	return html.EscapeString(a)
}
