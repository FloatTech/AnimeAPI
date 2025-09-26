package bilibili

import (
	"testing"
)

func TestGetAllGuard(t *testing.T) {
	guardUser, err := GetAllGuard("628537")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", guardUser)
}

func TestGetDynamicDetail(t *testing.T) {
	cfg := NewCookieConfig("config.json")
	detail, err := cfg.GetDynamicDetail("851252197280710664")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", detail)
}

func TestArticleInfo(t *testing.T) {
	card, err := GetArticleInfo("17279244")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(card.ToArticleMessage("17279244"))
}

func TestMemberCard(t *testing.T) {
	card, err := GetMemberCard(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", card)
}

func TestVideoInfo(t *testing.T) {
	card, err := GetVideoInfo("10007")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(card.ToVideoMessage())
	card, err = GetVideoInfo("BV1xx411c7mD")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(card.ToVideoMessage())
	card, err = GetVideoInfo("bv1xx411c7mD")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(card.ToVideoMessage())
	card, err = GetVideoInfo("BV1mF411j7iU")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(card.ToVideoMessage())
}

func TestLiveRoomInfo(t *testing.T) {
	card, err := GetLiveRoomInfo("83171", "b_ut=7;buvid3=0;i-wanna-go-back=-1;innersign=0;")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(card.ToMessage())
}
