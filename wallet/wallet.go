// Package wallet 货币系统
package wallet

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/file"
	sql "github.com/FloatTech/sqlite"
)

// Storage 货币系统
type Storage struct {
	sync.RWMutex
	db *sql.Sqlite
}

// Wallet 钱包
type Wallet struct {
	UID   int64
	Money int
}

var (
	sdb = &Storage{
		db: &sql.Sqlite{
			DBPath: "data/wallet/wallet.db",
		},
	}
)

func init() {
	if file.IsNotExist("data/wallet") {
		err := os.MkdirAll("data/wallet", 0755)
		if err != nil {
			panic(err)
		}
	}
	err := sdb.db.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
}

// GetWalletOf 获取钱包数据
func GetWalletOf(uid int64) (money int, err error) {
	return sdb.getWalletOf(uid)
}

// GetWalletInfoGroup 获取多人钱包数据
//
// if sort == true,由高到低排序; if sort == false,由低到高排序
func GetGroupWalletOf(uids []int64, sortable bool) (money []Wallet, err error) {
	return sdb.getGroupWalletOf(uids, sortable)
}

// InsertWalletOf 更新钱包(money > 0 增加,money < 0 减少)
func InsertWalletOf(uid int64, money int) error {
	lastMoney, err := sdb.getWalletOf(uid)
	if err == nil {
		err = sdb.updateWalletOf(uid, lastMoney+money)
	}
	return err
}

// 获取钱包数据
func (sql *Storage) getWalletOf(uid int64) (money int, err error) {
	sql.Lock()
	defer sql.Unlock()
	err = sql.db.Create("storage", &Wallet{})
	if err != nil {
		return
	}
	info := Wallet{}
	uidstr := strconv.FormatInt(uid, 10)
	_ = sql.db.Find("storage", &info, "where uid is "+uidstr)
	money = info.Money
	return
}

// 获取钱包数据组
func (sql *Storage) getGroupWalletOf(uids []int64, issorted bool) (money []Wallet, err error) {
	uidstr := make([]string, 0, len(uids))
	for _, uid := range uids {
		uidstr = append(uidstr, strconv.FormatInt(uid, 10))
	}
	sql.Lock()
	defer sql.Unlock()
	err = sql.db.Create("storage", &Wallet{})
	if err != nil {
		return
	}
	money = make([]Wallet, 0, len(uids))
	sort := "ASC"
	if issorted {
		sort = "DESC"
	}
	info := Wallet{}
	err = sql.db.FindFor("storage", &info, "where uid IN ("+strings.Join(uidstr, ", ")+") ORDER BY money "+sort, func() error {
		money = append(money, info)
		return nil
	})
	return
}

// 更新钱包
func (sql *Storage) updateWalletOf(uid int64, money int) (err error) {
	sql.Lock()
	defer sql.Unlock()
	err = sql.db.Create("storage", &Wallet{})
	if err == nil {
		err = sql.db.Insert("storage", &Wallet{
			UID:   uid,
			Money: money,
		})
	}
	return
}
