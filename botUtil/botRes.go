package botUtil

import "encoding/json"

type botRes struct{}

var BotRes botRes

func (*botRes) decode(content []byte) map[string]interface{} {
	var t map[string]interface{}
	_ = json.Unmarshal(content, &t)
	return t
}

func (*botRes) checkStatus(s map[string]interface{}) bool {
	if s == nil {
		return false
	}
	if s["status"] == "ok" {
		return true
	}
	return false
}

func (*botRes) doInt(d map[string]interface{}, a string) int {
	return int(d["data"].(map[string]interface{})[a].(float64))
}

func (a *botRes) GetMsgId(res map[string]interface{}) int {
	if !a.checkStatus(res) {
		return 0
	}
	return a.doInt(res, "message_id")
}
