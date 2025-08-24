// Package niu 牛牛大作战
package niu

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"math"
	"math/rand"
	"sort"
	"time"
)

var (
	daJiaoProps = []string{"伟哥", "媚药"}
	jjPorps     = []string{"击剑神器", "击剑神稽"}
)

// Rm ...
type Rm struct {
	NiuID     uuid.UUID
	TimeLimit time.Time
	Count     int
	Length    float64
}

type users []*userInfo

type niuNiuManager struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	NiuID     uuid.UUID `gorm:"type:varchar(36);unique_index"`
	Status    int       `gorm:"default:0"` // 0正常 1拍卖 2注销
}

// UserInfo 结构体
type userInfo struct {
	gorm.Model
	UserID   int64     `gorm:"column:user_id;index"`
	NiuID    uuid.UUID `gorm:"type:char(36);index"`
	Length   float64   `gorm:"default:1"`
	WeiGe    int       `gorm:"default:0"`
	MeiYao   int       `gorm:"default:0"`
	Artifact int       `gorm:"default:0"`
	ShenJi   int       `gorm:"default:0"`
	Buff2    int       `gorm:"default:0"`
	Buff3    int       `gorm:"default:0"`
	Buff4    int       `gorm:"default:0"`
	Buff5    int       `gorm:"default:0"`
}

// AuctionInfo 结构体
type AuctionInfo struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UserID    int64     `gorm:"column:user_id;index"`
	NiuID     uuid.UUID `gorm:"type:varchar(36);unique_index"`
	Length    float64   `gorm:"default:0.01"`
	Money     int
}

// BaseInfo ...
type BaseInfo struct {
	UID    int64
	Length float64
}

// BaseInfos ...
type BaseInfos []BaseInfo

func (m users) filter(pos bool) users {
	if pos {
		return m.positive()
	}
	return m.negative()
}

func (m users) positive() users {
	var m1 []*userInfo
	for _, i2 := range m {
		if i2.Length > 0 {
			m1 = append(m1, i2)
		}
	}
	return m1
}

func (m users) negative() users {
	var m1 []*userInfo
	for _, i2 := range m {
		if i2.Length <= 0 {
			m1 = append(m1, i2)
		}
	}
	return m1
}

func (m users) sort(isDesc bool) {
	t := func(i, j int) bool {
		return m[i].Length < m[j].Length
	}
	if isDesc {
		t = func(i, j int) bool {
			return m[i].Length > m[j].Length
		}
	}
	sort.Slice(m, t)
}

func (m users) ranking(niuniu float64, uid int64) int {
	m.sort(niuniu > 0)
	for i, user := range m {
		if user.UserID == uid {
			return i + 1
		}
	}
	return -1
}

func (u *userInfo) useWeiGe() (string, float64) {
	niuniu := u.Length
	reduce := math.Abs(hitGlue(niuniu))
	niuniu += reduce
	return randomChoice([]string{
		fmt.Sprintf("哈哈，你这一用伟哥，牛牛就像是被激发了潜能，增加了%.2fcm！看来今天是个大日子呢！", reduce),
		fmt.Sprintf("你这是用了什么神奇的道具？牛牛竟然增加了%.2fcm，简直是牛气冲天！", reduce),
		fmt.Sprintf("使用伟哥后，你的牛牛就像是开启了加速模式，一下增加了%.2fcm，这成长速度让人惊叹！", reduce),
	}), niuniu
}

func (u *userInfo) useMeiYao() (string, float64) {
	niuniu := u.Length
	reduce := math.Abs(hitGlue(niuniu))
	niuniu -= reduce
	return randomChoice([]string{
		fmt.Sprintf("你使用媚药,咿呀咿呀一下使当前长度发生了一些变化，当前长度%.2fcm", niuniu),
		fmt.Sprintf("看来你追求的是‘微观之美’，故意使用媚药让牛牛凹进去了%.2fcm！", reduce),
		fmt.Sprintf("缩小奇迹’在你身上发生了，牛牛凹进去了%.2fcm，你的选择真是独特！", reduce),
	}), niuniu
}

func (u *userInfo) useArtifact(adduserniuniu float64) (string, float64, float64) {
	myLength := u.Length
	difference := myLength - adduserniuniu
	var (
		change float64
	)
	if difference > 0 {
		change = hitGlue(myLength + adduserniuniu)
	} else {
		change = hitGlue((myLength + adduserniuniu) / 2)
	}
	myLength += change
	return randomChoice([]string{
		fmt.Sprintf("凭借神秘道具的力量，你让对方在你的长度面前俯首称臣！你的长度增加了%.2fcm，当前长度达到了%.2fcm", change, myLength),
		fmt.Sprintf("神器在手，天下我有！你使用道具后，长度猛增%.2fcm，现在的总长度是%.2fcm，无人能敌！", change, myLength),
		fmt.Sprintf("这就是道具的魔力！你轻松增加了%.2fcm，让对手望尘莫及，当前长度为%.2fcm！", change, myLength),
		fmt.Sprintf("道具一出，谁与争锋！你的长度因道具而增长%.2fcm，现在的长度是%.2fcm，霸气尽显！", change, myLength),
		fmt.Sprintf("使用道具的你，如同获得神助！你的长度增长了%.2fcm，达到%.2fcm的惊人长度，胜利自然到手！", change, myLength),
	}), myLength, adduserniuniu - change/1.3
}

func (u *userInfo) useShenJi(adduserniuniu float64) (string, float64, float64) {
	myLength := u.Length
	difference := myLength - adduserniuniu
	var (
		change float64
	)
	if difference > 0 {
		change = hitGlue(myLength + adduserniuniu)
	} else {
		change = hitGlue((myLength + adduserniuniu) / 2)
	}
	myLength -= change
	var r string
	if myLength > 0 {
		r = randomChoice([]string{
			fmt.Sprintf("哦吼！？看来你的牛牛因为使用了神秘道具而缩水了呢🤣🤣🤣！缩小了%.2fcm！", change),
			fmt.Sprintf("哈哈，看来这个道具有点儿调皮，让你的长度缩水了%.2fcm！现在你的长度是%.2fcm，下次可得小心使用哦！", change, myLength),
			fmt.Sprintf("使用道具后，你的牛牛似乎有点儿害羞，缩水了%.2fcm！现在的长度是%.2fcm，希望下次它能挺直腰板！", change, myLength),
			fmt.Sprintf("哎呀，这个道具的效果有点儿意外，你的长度减少了%.2fcm，现在只有%.2fcm了！下次选道具可得睁大眼睛！", change, myLength),
		})
	} else {
		r = randomChoice([]string{
			fmt.Sprintf("哦哟，小姐姐真是玩得一手好游戏，使用道具后数值又降低了%.2fcm，小巧得更显魅力！", change),
			fmt.Sprintf("看来小姐姐喜欢更加精致的风格，使用道具后，数值减少了%.2fcm，更加迷人了！", change),
			fmt.Sprintf("小姐姐的每一次变化都让人惊喜，使用道具后，数值减少了%.2fcm，更加优雅动人！", change),
			fmt.Sprintf("小姐姐这是在展示什么是真正的精致小巧，使用道具后，数值减少了%.2fcm，美得不可方物！", change),
		})
	}
	return r, myLength, adduserniuniu + 0.7*change
}

func (u *userInfo) applyProp(props string) error {
	propsMap := map[string]struct {
		itemCount *int
		errMsg    string
	}{
		"伟哥":   {&u.WeiGe, "你还没有伟哥呢,不能使用"},
		"媚药":   {&u.MeiYao, "你还没有媚药呢,不能使用"},
		"击剑神器": {&u.Artifact, "你还没有击剑神器呢,不能使用"},
		"击剑神稽": {&u.ShenJi, "你还没有击剑神稽呢,不能使用"},
	}

	propInfo, ok := propsMap[props]
	if !ok {
		return ErrPropNotFound
	}

	return u.useItem(propInfo.itemCount, propInfo.errMsg)
}

func (u *userInfo) useItem(itemCount *int, errMsg string) error {
	if *itemCount > 0 {
		*itemCount--
		return nil
	}
	fmt.Println(*u)
	return errors.New(errMsg)
}

func (u *userInfo) checkProps(props, propSort string) error {
	validProps := map[string][]string{
		"dajiao": daJiaoProps,
		"jj":     jjPorps,
	}

	// 检查是否是有效道具类别
	validPropsList, ok := validProps[propSort]
	if !ok {
		return ErrInvalidPropType
	}

	validPropsMap := make(map[string]struct{})
	for _, prop := range validPropsList {
		validPropsMap[prop] = struct{}{}
	}

	// 如果道具属于有效道具，返回 nil
	if _, exists := validPropsMap[props]; exists {
		return nil
	}

	// 检查是否相反
	conflictingProps := daJiaoProps
	if propSort == "dajiao" {
		conflictingProps = jjPorps
	}

	// 如果道具属于冲突集合,返回
	for _, conflictProp := range conflictingProps {
		if props == conflictProp {
			return ErrInvalidPropUsageScope
		}
	}

	return ErrPropNotFound
}

func (u *userInfo) purchaseItem(productID, quantity int) (money int, err error) {

	m := map[int]int{
		1: 300, // 伟哥
		2: 300, // 媚药
		3: 500, // 击剑神器
		4: 500, // 击剑神稽
	}

	var ok bool
	if money, ok = m[productID]; !ok {
		return 0, ErrInvalidProductID
	}

	money *= quantity

	switch productID {
	case 1:
		u.WeiGe += quantity
	case 2:
		u.MeiYao += quantity
	case 3:
		u.Artifact += quantity
	case 4:
		u.ShenJi += quantity
	}

	return
}

func (u *userInfo) processDaJiao(props string) (messages string, err error) {
	var f float64

	info := *u
	if props != "" {
		err = u.checkProps(props, "dajiao")
		if err != nil {
			return "", err
		}
		if err = u.applyProp(props); err != nil {
			return "", err
		}
	}
	switch {
	case u.WeiGe-info.WeiGe != 0:
		messages, f = u.useWeiGe()
		u.Length = f
	case u.MeiYao-info.MeiYao != 0:
		messages, f = u.useMeiYao()
		u.Length = f
	default:
		messages, f = hitGlueNiuNiu(u.Length)
		u.Length = f
	}
	return
}

func (u *userInfo) processJJ(adduserniuniu *userInfo, props string) (string, error) {
	var (
		fencingResult string
		f             float64
		f1            float64
		info          userInfo
		err           error
	)
	info = *u
	if props != "" {
		err = u.checkProps(props, "jj")
		if err != nil {
			return "", err
		}
		if err = u.applyProp(props); err != nil {
			return "", err
		}
	}
	switch {
	case u.ShenJi-info.ShenJi != 0:
		fencingResult, f, f1 = u.useShenJi(adduserniuniu.Length)
		u.Length = f
		adduserniuniu.Length = f1

	case u.Artifact-info.Artifact != 0:
		fencingResult, f, f1 = u.useArtifact(adduserniuniu.Length)
		u.Length = f
		adduserniuniu.Length = f1

	default:
		fencingResult, f, f1 = fencing(u.Length, adduserniuniu.Length)
		u.Length = f
		adduserniuniu.Length = f1
	}
	return fencingResult, err
}

func newLength() float64 {
	return float64(rand.Intn(9)+1) + (float64(rand.Intn(100)) / 100)
}
