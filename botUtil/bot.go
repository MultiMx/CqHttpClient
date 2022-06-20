package botUtil

import (
	"context"
	"fmt"
	"github.com/Mmx233/tool"
	"github.com/MultiMx/CqHttpClient/client"
	botTransfer "github.com/MultiMx/CqHttpClient/transfer"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"unicode/utf8"
)

type bot struct{}

var Bot bot

func (s bot) Equal(c context.Context, key string) bool {
	return strings.TrimSpace(s.GetMsg(c)) == key
}

// Match prefix非正则，del只在匹配后执行
func (s bot) Match(c context.Context, prefix string, delPredix bool, allowEmpty bool) bool {
	if strings.HasPrefix(s.GetMsg(c), prefix) {
		if a := strings.TrimSpace(strings.TrimPrefix(s.GetMsg(c), prefix)); a == "" && !allowEmpty {
			return false
		} else if delPredix {
			s.WriteMsg(c, a)
		}
		return true
	}
	return false
}

func (s bot) MatchContain(c context.Context, p string) bool {
	return strings.Contains(s.GetMsg(c), p)
}

func (s bot) MatchReply(c context.Context, word string, ReplaceMsg bool) bool {
	if tool.Regexp.MatchExisting(`^\[CQ:reply,.*?\]\s?`+word, s.GetMsg(c)) {
		if ReplaceMsg {
			s.WriteMsg(c, Bot.GetRelyMsg(c))
		}
		return true
	}
	return false
}

// MatchAt
//
// delAt只在匹配命中后生效。word为正则表达式，i指是否区分大小写。
func (s bot) MatchAt(c context.Context, word string, delAt bool, i bool) bool {
	var pre string
	if i {
		pre = "(i?)"
	}
	if tool.Regexp.MatchExisting(pre+`^\[CQ:at,qq=`+fmt.Sprint(s.GetSelf(c))+`\]\s*`+word, s.GetMsg(c)) {
		if delAt {
			s.WriteMsg(c, tool.Regexp.Replace(`^\[CQ:at,qq=`+fmt.Sprint(s.GetSelf(c))+`\]\s*`, s.GetMsg(c), ""))
		}
		return true
	}
	return false
}

func (s bot) MatchRegexp(c context.Context, reg string) bool {
	return tool.Regexp.MatchExisting(reg, s.GetMsg(c))
}

func (bot) GetMap(c context.Context) map[string]interface{} {
	return c.Value("D").(map[string]interface{})
}

func (s bot) doString(c context.Context, index string) string {
	v, _ := s.GetMap(c)[index].(string)
	return v
}

func (s bot) doInt(c context.Context, index string) int {
	v, _ := s.GetMap(c)[index].(float64)
	return int(v)
}

func (s bot) GetGroupId(c context.Context) int {
	return s.doInt(c, "group_id")
}

func (s bot) GetMsg(c context.Context) string {
	return s.doString(c, "message")
}

func (s bot) GetUserId(c context.Context) int {
	return s.doInt(c, "user_id")
}

func (s bot) GetSelf(c context.Context) int {
	return s.doInt(c, "self_id")
}

func (s bot) GetRole(c context.Context) string {
	return s.GetMap(c)["sender"].(map[string]interface{})["role"].(string)
}

func (s bot) GetMessageId(c context.Context) int {
	return s.doInt(c, "message_id")
}

func (bot) GetImageId(c context.Context) []string {
	d := tool.Regexp.MatchValue(`\[CQ:image.*?file=(.*\.image).*?\]`, Bot.GetMsg(c))
	var t []string
	for _, v := range d {
		t = append(t, v[1])
	}
	return t
}

// GetPostType message notice等
func (s bot) GetPostType(c context.Context) string {
	return s.doString(c, "post_type")
}

// GetMsgType private或group
func (s bot) GetMsgType(c context.Context) string {
	return s.doString(c, "message_type")
}

func (s bot) GetRelyMsg(c context.Context) string {
	return Cq.DecodeType(c, "reply", false)[0].Data["text"]
}

/*func (s bot) GetTime(c context.Context) time.Time{
	m := *s.GetMap(c)
	t:=int64(m["time"].(float64))
	return time.Unix(t,0)
}*/

func (s bot) IsGroup(c context.Context) bool {
	if v, ok := s.GetMap(c)["message_type"]; ok && v == "group" {
		return true
	}
	return false
}

func (s bot) IsOpened(GroupId int) bool {
	return BotData.ReadGroups(GroupId)
}
func (s bot) IO(c context.Context) bool {
	return s.IsOpened(s.GetGroupId(c))
}

func (s bot) WriteMsg(c context.Context, a string) {
	s.GetMap(c)["message"] = a
}

func (s bot) GenShareByUrl(url string) (*botTransfer.Share, error) {
	res, e := client.Http.GetRequest(&tool.DoHttpReq{
		Url: url,
	})
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	t, e := goquery.NewDocumentFromReader(res.Body)
	if e != nil {
		return nil, e
	}
	Title := strings.Split(strings.TrimSpace(t.Find("title").Text()), " ")
	if len(Title) == 1 {
		Title = append(Title, "分享")
	}
	return &botTransfer.Share{
			Source:   Title[len(Title)-1],
			Title:    Title[0],
			Url:      url,
			ImageUrl: t.Find("meta[property*=image]").AttrOr("content", ""),
			Desc:     t.Find("meta[name=description]").AttrOr("content", " "),
		},
		e
}

func (s bot) GenShareCqCode(c *botTransfer.Share) string {
	var title, content string
	if utf8.RuneCountInString(c.Title) > 20 {
		title = string([]rune(c.Title)[:19]) + "…"
	} else {
		title = c.Title
	}
	if utf8.RuneCountInString(c.Desc) > 50 {
		content = string([]rune(c.Desc)[:49]) + "…"
	} else {
		content = c.Desc
	}
	return Cq.Make("share", map[string]string{
		"url":     c.Url,
		"title":   title,
		"content": content,
		"image":   c.ImageUrl,
	})
}

func (s bot) InGroup(GroupId int) bool {
	BotData.GroupsLock.RLock()
	_, ok := BotData.Groups[GroupId]
	BotData.GroupsLock.RUnlock()
	return ok
}
