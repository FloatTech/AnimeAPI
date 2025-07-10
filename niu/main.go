// Package niu 牛牛大作战
package niu

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/FloatTech/floatbox/file"

	"github.com/FloatTech/AnimeAPI/wallet"
)

const (
	ur = "user"
	ac = "auction"
)

var (
	db         *gorm.DB
	globalLock sync.Mutex

	errCancelFail = errors.New("遇到不可抗力因素，注销失败！")
	// ErrNoBoys 表示当前没有男孩子可用的错误。
	ErrNoBoys = errors.New("暂时没有男孩子哦")

	// ErrNoGirls 表示当前没有女孩子可用的错误。
	ErrNoGirls = errors.New("暂时没有女孩子哦")

	// ErrNoNiuNiu 表示用户尚未拥有牛牛的错误。
	ErrNoNiuNiu = errors.New("你还没有牛牛呢,快去注册吧！")

	// ErrNoNiuNiuINAuction 表示拍卖行当前没有牛牛可用的错误。
	ErrNoNiuNiuINAuction = errors.New("拍卖行还没有牛牛呢")

	// ErrNoMoney 表示用户资金不足的错误。
	ErrNoMoney = errors.New("你的钱不够快去赚钱吧！")

	// ErrAdduserNoNiuNiu 表示对方尚未拥有牛牛，因此无法进行某些操作的错误。
	ErrAdduserNoNiuNiu = errors.New("对方还没有牛牛呢，不能🤺")

	// ErrCannotFight 表示无法进行战斗操作的错误。
	ErrCannotFight = errors.New("你要和谁🤺？你自己吗？")

	// ErrNoNiuNiuTwo 表示用户尚未拥有牛牛，无法执行特定操作的错误。
	ErrNoNiuNiuTwo = errors.New("你还没有牛牛呢，咋的你想凭空造一个啊")

	// ErrAlreadyRegistered 表示用户已经注册过的错误。
	ErrAlreadyRegistered = errors.New("你已经注册过了")

	// ErrInvalidPropType 表示传入的道具类别错误的错误。
	ErrInvalidPropType = errors.New("道具类别传入错误")

	// ErrInvalidPropUsageScope 表示道具使用域错误的错误。
	ErrInvalidPropUsageScope = errors.New("道具使用域错误")

	// ErrPropNotFound 表示找不到指定道具的错误。
	ErrPropNotFound = errors.New("道具不存在")
)

func init() {
	if file.IsNotExist("data/niuniu") {
		err := os.MkdirAll("data/niuniu", 0775)
		if err != nil {
			panic(err)
		}
	}

	sdb, err := gorm.Open("sqlite3", "data/niuniu/niuniu.db")
	if err != nil {
		panic(err)
	}

	if err = sdb.AutoMigrate(&niuNiuManager{}).Error; err != nil {
		panic(err)
	}

	db = sdb.LogMode(false)

}

// DeleteWordNiuNiu ...
func DeleteWordNiuNiu(gid, uid int64) error {
	globalLock.Lock()
	defer globalLock.Unlock()
	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return err
	}
	return deleteUserByID(gid, uid)
}

// SetWordNiuNiu length > 0 就增加 , length < 0 就减小
func SetWordNiuNiu(gid, uid int64, length float64) error {
	globalLock.Lock()
	defer globalLock.Unlock()
	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return err
	}
	m := map[string]interface{}{
		"length": length,
	}
	return updatesUserByID(gid, uid, m)
}

// GetWordNiuNiu ...
func GetWordNiuNiu(gid, uid int64) (float64, error) {
	globalLock.Lock()
	defer globalLock.Unlock()

	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return 0, err
	}

	info, err := getUserByID(gid, uid)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, ErrNoNiuNiu
	} else if err != nil {
		return 0, err
	}

	return info.Length, err
}

// GetRankingInfo 获取排行信息
func GetRankingInfo(gid int64, t bool) (BaseInfos, error) {
	globalLock.Lock()
	defer globalLock.Unlock()
	var (
		list users
	)

	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return nil, err
	}

	us, err := listUsers(gid)
	if err != nil {
		return nil, err
	}

	list = us.filter(t)
	list.sort(t)
	if len(list) == 0 {
		if t {
			return nil, ErrNoBoys
		}
		return nil, ErrNoGirls
	}
	f := make(BaseInfos, len(list))
	for i, info := range list {
		f[i] = BaseInfo{
			UID:    info.UserID,
			Length: info.Length,
		}
	}
	return f, nil
}

// GetGroupUserRank 获取指定用户在群中的排名
func GetGroupUserRank(gid, uid int64) (int, error) {
	globalLock.Lock()
	defer globalLock.Unlock()

	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return 0, err
	}
	niu, err := getUserByID(gid, uid)
	if err != nil {
		return -1, err
	}

	group, err := listUsers(gid)
	if err != nil {
		return -1, err
	}

	return group.ranking(niu.Length, uid), nil
}

// View 查看牛牛
func View(gid, uid int64, name string) (string, error) {
	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return "", err
	}
	i, err := getUserByID(gid, uid)
	if err != nil {
		return "", ErrNoNiuNiu
	}
	niuniu := i.Length
	var result strings.Builder
	sexLong := "长"
	sex := "♂️"
	if niuniu < 0 {
		sexLong = "深"
		sex = "♀️"
	}
	niuniuList, err := listUsers(gid)
	if err != nil {
		return "", err
	}
	result.WriteString(fmt.Sprintf("\n📛%s<%s>的牛牛信息\n⭕性别:%s\n⭕%s度:%.2fcm\n⭕排行:%d\n⭕%s ",
		name, strconv.FormatInt(uid, 10),
		sex, sexLong, niuniu, niuniuList.ranking(niuniu, uid), generateRandomString(niuniu)))
	return result.String(), nil
}

// HitGlue 打胶
func HitGlue(gid, uid int64, prop string) (string, error) {
	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return "", err
	}
	niuniu, err := getUserByID(gid, uid)
	if err != nil {
		return "", ErrNoNiuNiuTwo
	}

	messages, err := niuniu.processDaJiao(prop)
	if err != nil {
		return "", err
	}

	if err = TableFor(gid, ur).Where("user_id = ?", uid).Save(niuniu).Error; err != nil {
		return "", err
	}

	return messages, nil
}

// Register 注册牛牛
func Register(gid, uid int64) (string, error) {
	globalLock.Lock()
	defer globalLock.Unlock()

	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return "", err
	}

	if _, err := getUserByID(gid, uid); err == nil {
		return "", ErrAlreadyRegistered
	}
	// 获取初始长度
	length := newLength()
	u := userInfo{
		UserID: uid,
		NiuID:  uuid.New(),
		Length: length,
	}

	if err := createUser(gid, &u, ur); err != nil {
		return "", err
	}

	if err := db.Model(&niuNiuManager{}).Create(&niuNiuManager{
		NiuID: u.NiuID,
	}).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("注册成功,你的牛牛现在有%.2fcm", u.Length), nil
}

// JJ ...
func JJ(gid, uid, adduser int64, prop string) (message string, adduserLength float64, err error) {
	globalLock.Lock()
	defer globalLock.Unlock()

	if err = ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return "", 0, ErrNoNiuNiu
	}

	myniuniu, err := getUserByID(gid, uid)
	if err != nil {
		return "", 0, ErrNoNiuNiu
	}

	adduserniuniu, err := getUserByID(gid, adduser)
	if err != nil {
		return "", 0, ErrAdduserNoNiuNiu
	}

	if uid == adduser {
		return "", 0, ErrCannotFight
	}

	message, err = myniuniu.processJJ(adduserniuniu, prop)
	if err != nil {
		return "", 0, err
	}

	if err = TableFor(gid, ur).Where("user_id =?", uid).Update("length", myniuniu.Length).Error; err != nil {
		return "", 0, err
	}

	if err = TableFor(gid, ur).Where("user_id =?", adduser).Update("length", adduserniuniu.Length).Error; err != nil {
		return "", 0, err
	}

	adduserLength = adduserniuniu.Length

	return
}

// Cancel 注销牛牛
func Cancel(gid, uid int64) (string, error) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return "", err
	}
	_, err := getUserByID(gid, uid)
	if err != nil {
		return "", ErrNoNiuNiuTwo
	}
	err = deleteUserByID(gid, uid)
	if err != nil {
		return "", errCancelFail
	}
	err = db.Model(&niuNiuManager{}).Where("niu_id = ?", uid).Update("status", 2).Error
	return "注销成功,你已经没有牛牛了", err
}

// Redeem 赎牛牛
func Redeem(gid, uid int64, lastLength float64) error {
	globalLock.Lock()
	defer globalLock.Unlock()
	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return err
	}

	_, err := getUserByID(gid, uid)
	if err != nil {
		return ErrNoNiuNiu
	}

	/*	var n niuNiuManager
		if err = db.Where("niu_id = ?", niuID).First(&n).Error; err != nil {
			return err
		}

		switch n.Status {
		case 1:
			return errors.New("你的牛牛已经被拍卖无法赎回")
		case 2:
			return errors.New("你的牛牛已经被注销无法赎回")
		}*/

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

	if err = wallet.InsertWalletOf(uid, -150); err != nil {
		return err
	}

	return TableFor(gid, ur).Where("user_id = ?", uid).Update("length", lastLength).Error
}

// Store 牛牛商店
func Store(gid, uid int64, n int) error {
	globalLock.Lock()
	defer globalLock.Unlock()
	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return err
	}

	info, err := getUserByID(gid, uid)
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

	return TableFor(gid, ur).Save(info).Error
}

// Sell 出售牛牛
func Sell(gid, uid int64) (string, error) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return "", err
	}
	if err := ensureUserInfoTable[AuctionInfo](gid, ac); err != nil {
		return "", err
	}
	niu, err := getUserByID(gid, uid)
	if err != nil {
		return "", ErrNoNiuNiu
	}
	money, t, message := profit(niu.Length)
	if !t {
		return "", errors.New(message)
	}

	if err = deleteUserByID(gid, uid); err != nil {
		return "", err
	}

	err = wallet.InsertWalletOf(uid, money)
	if err != nil {
		return message, err
	}

	if err = db.Model(&niu).Where("niu_id = ?", niu.NiuID).Update("status", 1).Error; err != nil {
		return message, err
	}

	u := AuctionInfo{
		UserID: uid,
		NiuID:  niu.NiuID,
		Length: niu.Length,
		Money:  money * 2,
	}

	if err = TableFor(gid, ac).Create(&u).Error; err != nil {
		return "", err
	}

	db.Model(&niuNiuManager{}).Where("niu_id = ?", niu.NiuID).Update("status", 1)

	return message, err
}

// ShowAuction 展示牛牛拍卖行
func ShowAuction(gid int64) ([]AuctionInfo, error) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if err := ensureUserInfoTable[AuctionInfo](gid, ac); err != nil {
		return nil, err
	}
	return listAuction(gid)
}

// Auction 购买牛牛
func Auction(gid, uid int64, index int) (string, error) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if err := ensureUserInfoTable[AuctionInfo](gid, ac); err != nil {
		return "", err
	}

	infos, err := listAuction(gid)
	if len(infos) == 0 || err != nil {
		return "", ErrNoNiuNiuINAuction
	}

	var info AuctionInfo
	if err = TableFor(gid, ac).Where("id = ?", index).First(&info).Error; err != nil {
		return "", err
	}

	if err = wallet.InsertWalletOf(uid, -info.Money); err != nil {
		return "", ErrNoMoney
	}

	niu, err := getUserByID(gid, uid)

	if err != nil {
		niu.UserID = uid
	}

	niu.Length = info.Length
	niu.NiuID = info.NiuID

	if info.Money >= 500 {
		niu.WeiGe += 2
		niu.MeiYao += 2
	}

	if err = TableFor(gid, ac).Delete(&info).Error; err != nil {
		return "", err
	}

	if err = TableFor(gid, ur).Save(&niu).Error; err != nil {
		return "", err
	}

	if err = db.Model(&niuNiuManager{}).Where("niu_id = ?", niu.NiuID).Update("status", 0).Error; err != nil {
		return "", err
	}

	if info.Money >= 500 {
		return fmt.Sprintf("恭喜你购买成功,当前长度为%.2fcm,此次购买将赠送你2个伟哥,2个媚药", niu.Length), nil
	}

	return fmt.Sprintf("恭喜你购买成功,当前长度为%.2fcm", niu.Length), nil
}

// Bag 牛牛背包
func Bag(gid, uid int64) (string, error) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if err := ensureUserInfoTable[userInfo](gid, ur); err != nil {
		return "", err
	}
	niu, err := getUserByID(gid, uid)
	if err != nil {
		return "", ErrNoNiuNiu
	}

	var result strings.Builder
	result.Grow(100)

	result.WriteString("当前牛牛背包如下\n")
	result.WriteString(fmt.Sprintf("伟哥: %v\n", niu.WeiGe))
	result.WriteString(fmt.Sprintf("媚药: %v\n", niu.MeiYao))
	result.WriteString(fmt.Sprintf("击剑神器: %v\n", niu.Artifact))
	result.WriteString(fmt.Sprintf("击剑神稽: %v\n", niu.ShenJi))

	return result.String(), nil
}
