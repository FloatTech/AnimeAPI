// Package wallet 货币系统
package wallet

import (
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/file"
	sql "github.com/FloatTech/sqlite"
)

// WalletSYS 货币系统
type WalletSYS struct {
	sync.RWMutex
	Db *sql.Sqlite
}

// Wallet 钱包
type Wallet struct {
	UID   int64
	Money int
}

var (
	sdb = &WalletSYS{
		Db: &sql.Sqlite{},
	}
)

func init() {
	if file.IsNotExist("data/Wallet/wallet.db") {
		_, err := os.Create("data/Wallet/wallet.db")
		if err != nil {
			panic(err)
		}
	}
	sdb.Db.DBPath = "data/Wallet/wallet.db"
	err := sdb.Db.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
}

// GetScoreInfo 获取钱包数据
func GetWalletInfo(uid int64) (money int, err error) {
	return sdb.getWallet(uid)
}

// GetWalletInfoGroup 获取多人钱包数据(sort = true,由高到低排序)
func GetWalletInfoGroup(uids []int64, sortable bool) (money []Wallet, err error) {
	return sdb.getWalletGroup(uids, sortable)
}

// InsertScoreInfo 更新钱包(money > 0 增加,money < 0 减少)
func InsertWalletInfo(uid int64, money int) error {
	lastMoney, err := sdb.getWallet(uid)
	if err == nil {
		err = sdb.setWallet(uid, lastMoney+money)
	}
	return err
}

// 获取钱包数据
func (sql *WalletSYS) getWallet(uid int64) (money int, err error) {
	sql.Lock()
	defer sql.Unlock()
	err = sql.Db.Create("WalletSYS", &Wallet{})
	if err != nil {
		return
	}
	info := Wallet{}
	uidstr := strconv.FormatInt(uid, 10)
	err = sql.Db.Find("WalletSYS", &info, "where uid is "+uidstr)
	if err != nil {
		err = sql.Db.Insert("WalletSYS", &Wallet{
			UID:   uid,
			Money: 0,
		})
		return
	}
	money = info.Money
	return
}

// 获取钱包数据组
func (sql *WalletSYS) getWalletGroup(uid []int64, sortable bool) (money []Wallet, err error) {
	sql.Lock()
	defer sql.Unlock()
	err = sql.Db.Create("WalletSYS", &Wallet{})
	if err != nil {
		return
	}
	money = make([]Wallet, len(uid))
	for _, info := range uid {
		var walletinfo Wallet
		uidstr := strconv.FormatInt(info, 10)
		err := sql.Db.Find("WalletSYS", &walletinfo, "where uid is "+uidstr)
		if err == nil {
			money = append(money, walletinfo)
		}
	}
	if sortable {
		sort.SliceStable(money, func(i, j int) bool {
			return money[i].Money > money[j].Money
		})
	}
	return
}

// 更新钱包
func (sql *WalletSYS) setWallet(uid int64, money int) (err error) {
	sql.Lock()
	defer sql.Unlock()
	err = sql.Db.Create("WalletSYS", &Wallet{})
	if err == nil {
		err = sql.Db.Insert("WalletSYS", &Wallet{
			UID:   uid,
			Money: money,
		})
	}
	return
}
