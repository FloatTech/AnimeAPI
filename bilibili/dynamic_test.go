package bilibili

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
)

func TestSpaceHistory(t *testing.T) {
	data, err := web.GetData(fmt.Sprintf(SpaceHistoryURL, "667526012", "642279068898689029"))
	if err != nil {
		t.Fatal(err)
	}
	var desc Desc
	_ = json.Unmarshal([]byte(gjson.ParseBytes(data).Get("data.cards.0.desc").Raw), &desc)
	t.Logf("desc:%+v\n", desc)
	var card Card
	_ = json.Unmarshal([]byte(gjson.ParseBytes(data).Get("data.cards.0.card").Str), &card)
	t.Logf("card:%+v\n", card)
}

func TestCard2msg(t *testing.T) {
	data, err := web.GetData(fmt.Sprintf(SpaceHistoryURL, "667526012", "642279068898689029"))
	if err != nil {
		t.Fatal(err)
	}
	var dynamicCard DynamicCard
	_ = json.Unmarshal([]byte(gjson.ParseBytes(data).Get("data.cards.0").Raw), &dynamicCard)
	t.Logf("dynCard:%+v\n", dynamicCard)
}

func TestDynamicDetail(t *testing.T) {
	t.Log("cType = 1")
	t.Log(DynamicDetail("642279068898689029"))

	t.Log("cType = 2")
	t.Log(DynamicDetail("642470680290394121"))

	t.Log("cType = 2048")
	t.Log(DynamicDetail("642277677329285174"))

	t.Log("cType = 4")
	t.Log(DynamicDetail("642154347357011968"))

	t.Log("cType = 8")
	t.Log(DynamicDetail("675892999274627104"))

	t.Log("cType = 4308")
	t.Log(DynamicDetail("668598718656675844"))

	t.Log("cType = 64")
	t.Log(DynamicDetail("675966082178088963"))

	t.Log("cType = 256")
	t.Log(DynamicDetail("599253048535707632"))

	t.Log("cType = 4,投票类型")
	t.Log(DynamicDetail("677231070435868704"))
}
