package niu

import (
	"errors"
	"fmt"
	"github.com/FloatTech/AnimeAPI/wallet"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/rendercard"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	db = &model{}
)

func init() {
	if file.IsNotExist("data/niuniu") {
		err := os.MkdirAll("data/niuniu", 0755)
		if err != nil {
			panic(err)
		}
	}
	err := db.sql.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
}

// SetWordNiuNiu length > 0 就增加 , length < 0 就减小
func SetWordNiuNiu(gid, uid int64, length float64) error {
	niu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return err
	}
	niu.Length += length
	return db.setWordNiuNiu(gid, niu)
}

func GetWordNiuNiu(gid, uid int64) (float64, error) {
	niu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return 0, err
	}
	return niu.Length, nil
}

func DeleteWordNiuNiu(gid, uid int64) error {
	return db.deleteWordNiuNiu(gid, uid)
}

func GetRankingInfo(gid int64, t bool) ([]*rendercard.RankInfo, error) {
	var (
		s    = "牛牛深度"
		f    []*rendercard.RankInfo
		list users
		err  error
	)
	if t {
		s = "牛牛长度"
	}
	niuOfGroup, err := db.getAllNiuNiuOfGroup(gid)
	if err != nil {
		if t {
			err = errors.New("暂时没有男孩子哦")
		} else {
			err = errors.New("暂时没有女孩子哦")
		}
		return nil, err
	}

	if t {
		list = niuOfGroup.positive()
		niuOfGroup.sort(t)
	} else {
		list = niuOfGroup.negative()
		niuOfGroup.sort(!t)
	}
	for i, info := range list {
		f[i] = &rendercard.RankInfo{
			BottomLeftText: fmt.Sprintf("QQ:%d", info.UID),
			RightText:      fmt.Sprintf("%s:%.2fcm", s, info.Length),
		}
	}

	return f, nil
}

// GetRankingOfSpecifiedUser 获取指定用户在群中的排名
func GetRankingOfSpecifiedUser(gid, uid int64) (int, error) {
	niu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return -1, err
	}
	group, err := db.getAllNiuNiuOfGroup(gid)
	if err != nil {
		return -1, err
	}
	return group.ranking(niu.Length, uid), nil
}

func View(gid, uid int64, name string) (*strings.Builder, error) {
	i, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return nil, errors.New("你还没有牛牛呢不能查看")
	}
	niuniu := i.Length
	var result strings.Builder
	sexLong := "长"
	sex := "♂️"
	if niuniu < 0 {
		sexLong = "深"
		sex = "♀️"
	}
	niuniuList, err := db.getAllNiuNiuOfGroup(gid)
	if err != nil {
		return nil, err
	}
	result.WriteString(fmt.Sprintf("\n📛%s<%s>的牛牛信息\n⭕性别:%s\n⭕%s度:%.2fcm\n⭕排行:%d\n⭕%s ",
		name, strconv.FormatInt(uid, 10),
		sex, sexLong, niuniu, niuniuList.ranking(niuniu, uid), generateRandomString(niuniu)))
	return &result, nil
}

func ProcessHitGlue(gid, uid int64, prop string) (string, error) {
	niuniu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return "", errors.New("请先注册牛牛！")
	}

	messages, err := niuniu.processNiuNiuAction(prop)
	if err != nil {
		return "", err
	}
	if err = db.setWordNiuNiu(gid, niuniu); err != nil {
		return "", err
	}
	return messages, nil
}

func Register(gid, uid int64) (string, error) {
	if _, err := db.getWordNiuNiu(gid, uid); err == nil {
		return "", errors.New("你已经注册过了")
	}
	// 获取初始长度
	length := db.newLength()
	u := userInfo{
		UID:    uid,
		Length: length,
	}
	if err := db.setWordNiuNiu(gid, &u); err != nil {
		return "", err
	}
	return fmt.Sprintf("注册成功,你的牛牛现在有%.2fcm", u.Length), nil
}

func JJ(gid, uid, adduser int64, prop string) (message string, err error) {
	myniuniu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return "", errors.New("你还没有牛牛快去注册一个吧！")
	}
	adduserniuniu, err := db.getWordNiuNiu(gid, adduser)
	if err != nil {
		return "", errors.New("对方还没有牛牛呢，不能🤺")
	}

	if uid == adduser {
		return "", errors.New("你要和谁🤺？你自己吗？")
	}

	message, err = myniuniu.processJJuAction(adduserniuniu, prop)
	if err != nil {
		return "", err
	}

	if err = db.setWordNiuNiu(gid, myniuniu); err != nil {
		return "", err
	}

	if err = db.setWordNiuNiu(gid, adduserniuniu); err != nil {
		return "", err
	}
	return
}

func Cancel(gid, uid int64) (string, error) {
	_, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return "", errors.New("你还没有牛牛呢，咋的你想凭空造一个啊")
	}
	err = db.deleteWordNiuNiu(gid, uid)
	if err != nil {
		err = errors.New("遇到不可抗力因素，注销失败！")
	}
	return "注销成功,你已经没有牛牛了", err
}

func Redeem(gid, uid int64, lastLength float64) error {
	niu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return err
	}
	niu.Length = lastLength
	return db.setWordNiuNiu(gid, niu)
}

func Store(gid, uid int64, n int) error {
	info, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return err
	}

	money, err := info.purchaseItem(n)
	if err != nil {
		return err
	}

	if wallet.GetWalletOf(uid) < money {
		return errors.New("你还没有足够的ATRI币呢,不能购买")
	}

	if err = wallet.InsertWalletOf(uid, -money); err != nil {
		return err
	}

	return db.setWordNiuNiu(uid, info)
}
