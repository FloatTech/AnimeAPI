package bilibili

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
)

func TestMemberCard(t *testing.T) {
	var card MemberCard
	data, err := web.GetData(fmt.Sprintf(MemberCardURL, 2))
	if err != nil {
		return
	}
	str := gjson.ParseBytes(data).Get("card").String()
	err = json.Unmarshal(binary.StringToBytes(str), &card)
	if err != nil {
		return
	}
	t.Logf("%+v\n", card)
}
