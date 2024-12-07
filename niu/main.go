package niu

import (
	"errors"
	"fmt"
	"github.com/FloatTech/AnimeAPI/wallet"
	"github.com/FloatTech/floatbox/file"
	sql "github.com/FloatTech/sqlite"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	db                   = &model{}
	ErrNoBoys            = errors.New("暂时没有男孩子哦")
	ErrNoGirls           = errors.New("暂时没有女孩子哦")
	ErrNoNiuNiu          = errors.New("你还没有牛牛呢,快去注册吧！")
	ErrNoNiuNiuINAuction = errors.New("拍卖行还没有牛牛呢")
	ErrNoMoney           = errors.New("你的钱不够快去赚钱吧！")
	ErrAdduserNoNiuNiu   = errors.New("对方还没有牛牛呢，不能🤺")
	ErrCannotFight       = errors.New("你要和谁🤺？你自己吗？")
	ErrNoNiuNiuTwo       = errors.New("你还没有牛牛呢，咋的你想凭空造一个啊")
	ErrAlreadyRegistered = errors.New("你已经注册过了")
)

func init() {
	if file.IsNotExist("data/niuniu") {
		err := os.MkdirAll("data/niuniu", 775)
		if err != nil {
			panic(err)
		}
	}
	db.sql = sql.New("data/niuniu/niuniu.db")
	err := db.sql.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
}

// SetWordNiuNiu length > 0 就增加 , length < 0 就减小
func SetWordNiuNiu(gid, uid int64, length float64) error {
	db.Lock()
	defer db.Unlock()
	niu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return err
	}
	niu.Length += length
	return db.setWordNiuNiu(gid, niu)
}

func GetWordNiuNiu(gid, uid int64) (float64, error) {
	db.RLock()
	defer db.RUnlock()

	niu, err := db.getWordNiuNiu(gid, uid)
	return niu.Length, err
}

func GetRankingInfo(gid int64, t bool) (BaseInfos, error) {
	var (
		list users
		err  error
	)
	niuOfGroup, err := db.getAllNiuNiuOfGroup(gid)
	if err != nil {
		if t {
			return nil, ErrNoBoys
		}
		return nil, ErrNoGirls
	}
	list = niuOfGroup.filter(t)
	f := make(BaseInfos, len(list))
	for i, info := range list {
		f[i] = BaseInfo{
			UID:    info.UID,
			Length: info.Length,
		}
	}
	return f, nil
}

// GetGroupUserRank 获取指定用户在群中的排名
func GetGroupUserRank(gid, uid int64) (int, error) {
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

// View 查看牛牛
func View(gid, uid int64, name string) (*strings.Builder, error) {
	i, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return nil, ErrNoNiuNiu
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

// HitGlue 打胶
func HitGlue(gid, uid int64, prop string) (string, error) {
	niuniu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return "", ErrNoNiuNiuTwo
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

// Register 注册牛牛
func Register(gid, uid int64) (string, error) {
	if _, err := db.getWordNiuNiu(gid, uid); err == nil {
		return "", ErrAlreadyRegistered
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

// JJ ...
func JJ(gid, uid, adduser int64, prop string) (message string, adduserLength float64, err error) {

	myniuniu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return "", 0, ErrNoNiuNiu
	}

	adduserniuniu, err := db.getWordNiuNiu(gid, adduser)
	if err != nil {
		return "", 0, ErrAdduserNoNiuNiu
	}

	if uid == adduser {
		return "", 0, ErrCannotFight
	}

	message, err = myniuniu.processJJuAction(adduserniuniu, prop)
	if err != nil {
		return "", 0, err
	}

	if err = db.setWordNiuNiu(gid, myniuniu); err != nil {
		return "", 0, err
	}

	if err = db.setWordNiuNiu(gid, adduserniuniu); err != nil {
		return "", 0, err
	}

	adduserLength = adduserniuniu.Length
	return
}

// Cancel 注销牛牛
func Cancel(gid, uid int64) (string, error) {
	_, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return "", ErrNoNiuNiuTwo
	}
	err = db.deleteWordNiuNiu(gid, uid)
	if err != nil {
		err = errors.New("遇到不可抗力因素，注销失败！")
	}
	return "注销成功,你已经没有牛牛了", err
}

// Redeem 赎牛牛
func Redeem(gid, uid int64, lastLength float64) error {
	money := wallet.GetWalletOf(uid)
	if money < 150 {
		var builder strings.Builder
		walletName := wallet.GetWalletName()
		builder.WriteString("赎牛牛需要150")
		builder.WriteString(walletName)
		builder.WriteString("，快去赚钱吧，目前仅有:")
		builder.WriteString(strconv.Itoa(money))
		builder.WriteString("个")
		builder.WriteString(walletName)
		return errors.New(builder.String())
	}

	if err := wallet.InsertWalletOf(uid, -150); err != nil {
		return err
	}

	niu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return ErrNoNiuNiu
	}

	niu.Length = lastLength

	return db.setWordNiuNiu(gid, niu)
}

// Store 牛牛商店
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
		return ErrNoMoney
	}

	if err = wallet.InsertWalletOf(uid, -money); err != nil {
		return err
	}

	return db.setWordNiuNiu(uid, info)
}

// Sell 出售牛牛
func Sell(gid, uid int64) (string, error) {
	niu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return "", ErrNoNiuNiu
	}
	money, t, message := profit(niu.Length)
	if !t {
		return "", errors.New(message)
	}
	err = wallet.InsertWalletOf(uid, money)
	if err != nil {
		return message, err
	}
	u := AuctionInfo{
		UserId: niu.UID,
		Length: niu.Length,
		Money:  money * 2,
	}
	err = db.setNiuNiuAuction(gid, &u)
	return message, err
}

// ShowAuction 展示牛牛拍卖行
func ShowAuction(gid int64) ([]AuctionInfo, error) {
	db.RLock()
	defer db.RUnlock()
	return db.getAllNiuNiuAuction(gid)
}

// Auction 购买牛牛
func Auction(gid, uid int64, i int) (string, error) {
	auction, err := db.getAllNiuNiuAuction(gid)
	if err != nil {
		return "", ErrNoNiuNiuINAuction
	}
	err = wallet.InsertWalletOf(uid, -auction[i].Money)
	if err != nil {
		return "", ErrNoMoney
	}

	niu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		niu = &userInfo{
			UID: uid,
		}
	}
	niu.Length = auction[i].Length

	if auction[i].Money > 500 {
		niu.WeiGe = 2
		niu.Artifact = 2
	}

	if err = db.setWordNiuNiu(gid, niu); err != nil {
		return "", err
	}
	err = db.deleteNiuNiuAuction(gid, auction[i].ID)
	if err != nil {
		return "", err
	}
	if auction[i].Money > 500 {
		return fmt.Sprintf("恭喜你购买成功,当前长度为%.2fcm,此次购买将赠送你%d个伟哥,%d个媚药",
			niu.Length, niu.WeiGe, niu.Artifact), nil
	}
	return fmt.Sprintf("恭喜你购买成功,当前长度为%.2fcm", niu.Length), nil
}

// Bag 牛牛背包
func Bag(gid, uid int64) (string, error) {
	niu, err := db.getWordNiuNiu(gid, uid)
	if err != nil {
		return "", ErrNoNiuNiu
	}

	var result strings.Builder
	result.Grow(100)

	result.WriteString("当前牛牛背包如下\n")
	result.WriteString(fmt.Sprintf("伟哥: %v\n", niu.WeiGe))
	result.WriteString(fmt.Sprintf("媚药: %v\n", niu.Philter))
	result.WriteString(fmt.Sprintf("击剑神器: %v\n", niu.Artifact))
	result.WriteString(fmt.Sprintf("击剑神稽: %v\n", niu.ShenJi))

	return result.String(), nil
}
