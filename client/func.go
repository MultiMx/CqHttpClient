package client

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Mmx233/tool"
	"github.com/MultiMx/CqHttpClient/botUtil"
	botTransfer "github.com/MultiMx/CqHttpClient/transfer"
	"github.com/MultiMx/CqHttpClient/util"
	"html"
	"time"
)

func genMsgWithReply(msg string, MsgId int) string {
	if MsgId != 0 {
		return botUtil.Cq.Make("reply", map[string]string{
			"id": fmt.Sprint(MsgId),
		}) + msg
	} else {
		return msg
	}
}

func action(NeedID bool, action string, query map[string]interface{}) (int, error) {
	res, e := do(!NeedID, action, query)
	return botUtil.BotRes.GetMsgId(res), e
}

func SendGroupMsg(groupId int, msg string, MsgId int, NeedID bool) (int, error) {
	return action(NeedID, "send_group_msg", map[string]interface{}{
		"group_id": groupId,
		"message":  genMsgWithReply(msg, MsgId),
	})
}
func SGMR(c context.Context, msg string, NeedID bool) (int, error) {
	return SendGroupMsg(botUtil.Bot.GetGroupId(c), msg, botUtil.Bot.GetMessageId(c), NeedID)
}
func SGM(c context.Context, msg string, NeedID bool) (int, error) {
	return SendGroupMsg(botUtil.Bot.GetGroupId(c), msg, 0, NeedID)
}

func SendPrivateMsg(UserId int, msg string, MsgId int, NeedID bool) (int, error) {
	return action(NeedID, "send_private_msg", map[string]interface{}{
		"user_id": UserId,
		"message": genMsgWithReply(msg, MsgId),
	})
}
func SPMR(c context.Context, msg string, NeedID bool) (int, error) {
	return SendPrivateMsg(botUtil.Bot.GetUserId(c), msg, botUtil.Bot.GetMessageId(c), NeedID)
}
func SPM(c context.Context, msg string, NeedID bool) (int, error) {
	return SendPrivateMsg(botUtil.Bot.GetUserId(c), msg, 0, NeedID)
}

func AutoReply(c context.Context, msg string, NeedId bool) (int, error) {
	if botUtil.Bot.IsGroup(c) {
		return SGMR(c, msg, NeedId)
	} else {
		return SPMR(c, msg, NeedId)
	}
}
func Auto(c context.Context, msg string, NeedID bool) (int, error) {
	if botUtil.Bot.IsGroup(c) {
		return SGM(c, msg, NeedID)
	} else {
		return SPM(c, msg, NeedID)
	}
}

func DelMsg(msgId int) error {
	_, e := do(true, "delete_msg", map[string]interface{}{
		"message_id": msgId,
	})
	return e
}
func DM(c context.Context) error {
	return DelMsg(botUtil.Bot.GetMessageId(c))
}

func Ocr(image string) (string, error) {
	m, e := do(false, "ocr_image", map[string]interface{}{
		"image": image,
	})
	if e != nil {
		return "", e
	}
	if m["status"].(string) == "failed" {
		return "", errors.New("failed")
	}
	t2 := m["data"].(map[string]interface{})["texts"].([]interface{})
	var s string
	for _, v := range t2 {
		s += v.(map[string]interface{})["text"].(string)
	}
	return s, nil
}

func Ban(GroupId int, UserId int, d time.Duration) error {
	_, e := do(true, "set_group_ban", map[string]interface{}{
		"group_id": GroupId,
		"user_id":  UserId,
		"duration": uint(d.Seconds()),
	})
	return e
}

func KickMember(GroupId int, UserId int, reject bool) error {
	_, e := do(false, "set_group_kick", map[string]interface{}{
		"group_id":           GroupId,
		"user_id":            UserId,
		"reject_add_request": reject,
	})
	return e
}

func ShareGroup(GroupId int, share *botTransfer.Share, NeedId bool) (int, error) {
	return SendGroupMsg(GroupId, botUtil.Bot.GenShareCqCode(share), 0, NeedId)
}
func Share(c context.Context, share *botTransfer.Share, NeedId bool) (int, error) {
	return Auto(c, botUtil.Bot.GenShareCqCode(share), NeedId)
}

func TecentTts(c context.Context, word string, needId bool) (int, error) {
	return Auto(c, botUtil.Cq.Make("tts", map[string]string{
		"text": html.EscapeString(word),
	}), needId)
}

func YoudaoTts(c context.Context, word string, needId bool) (int, error) {
	_, n, e := util.Http.GetBytes(&tool.DoHttpReq{
		Url: "https://tts.youdao.com/fanyivoice",
		Query: map[string]interface{}{
			"le":      "auto",
			"keyfrom": "speaker-target",
			"word":    word,
		},
	})
	if e != nil {
		return 0, e
	}
	return Auto(c, botUtil.Cq.Make("record", map[string]string{
		"file": "base64://" + base64.StdEncoding.EncodeToString(n),
	}), needId)
}

func GetGroups() ([]int, error) {
	d, e := do(false, "get_group_list", nil)
	if e != nil {
		return nil, e
	}
	t := d["data"].([]interface{})
	var m []int
	for _, v := range t {
		m = append(m, int(v.(map[string]interface{})["group_id"].(float64)))
	}
	return m, nil
}

func RenewGroups(groupOpened []int) error {
	var t = make(map[int]bool)
	d, e := GetGroups()
	if e != nil || d == nil {
		return e
	}
	for _, v := range d {
		t[v] = false
	}
	for _, v := range groupOpened {
		if _, ok := t[v]; ok {
			t[v] = true
		}
	}
	botUtil.BotData.GroupsLock.Lock()
	botUtil.BotData.Groups = t
	botUtil.BotData.GroupsLock.Unlock()
	return nil
}

func GetMembersLastMsgTime(GroupId int) (map[int]time.Time, error) {
	var m = make(map[int]time.Time)
	d, e := do(false, "get_group_member_list", map[string]interface{}{
		"group_id": GroupId,
		"no_cache": true,
	})
	if e != nil {
		return nil, e
	}
	t := d["data"].([]interface{})
	for _, v := range t {
		m[int(v.(map[string]interface{})["user_id"].(float64))] = time.Unix(int64(v.(map[string]interface{})["last_sent_time"].(float64)), 0)
	}
	return m, nil
}
