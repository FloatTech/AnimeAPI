// Package wallet 货币系统
package wallet

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/file"
	sql "github.com/FloatTech/sqlite"
)

// Storage 货币系统
type Storage struct {
	sync.RWMutex
	db sql.Sqlite
}

// Wallet 钱包
type Wallet struct {
	UID   int64
	Money int
}

var (
	sdb = &Storage{
		db: sql.New("data/wallet/wallet.db"),
	}
	walletName = "Atri币"
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
	err = sdb.db.Create("storage", &Wallet{})
	if err != nil {
		panic(err)
	}
}

// GetWalletName 获取货币名称
func GetWalletName() string {
	return walletName
}

// SetWalletName 设置货币名称
func SetWalletName(name string) {
	walletName = name
}

// GetWalletOf 获取钱包数据
func GetWalletOf(uid int64) (money int) {
	return sdb.getWalletOf(uid).Money
}

// GetGroupWalletOf 获取多人钱包数据
//
// if sort == true,由高到低排序; if sort == false,由低到高排序
func GetGroupWalletOf(sortable bool, uids ...int64) (wallets []Wallet, err error) {
	return sdb.getGroupWalletOf(sortable, uids...)
}

// InsertWalletOf 更新钱包(money > 0 增加,money < 0 减少)
func InsertWalletOf(uid int64, money int) error {
	sdb.Lock()
	defer sdb.Unlock()
	lastMoney := sdb.getWalletOf(uid)
	newMoney := lastMoney.Money + money
	if newMoney < 0 {
		newMoney = 0
	}
	return sdb.updateWalletOf(uid, newMoney)
}

// 获取钱包数据 no lock
func (s *Storage) getWalletOf(uid int64) (wallet Wallet) {
	uidstr := strconv.FormatInt(uid, 10)
	_ = s.db.Find("storage", &wallet, "WHERE uid = ?", uidstr)
	return
}

// 获取钱包数据组
func (s *Storage) getGroupWalletOf(sortable bool, uids ...int64) (wallets []Wallet, err error) {
	s.RLock()
	defer s.RUnlock()
	wallets = make([]Wallet, 0, len(uids))
	sort := "ASC"
	if sortable {
		sort = "DESC"
	}
	info := Wallet{}
	q, sl := sql.QuerySet("WHERE uid", "IN", uids)
	err = s.db.FindFor("storage", &info, q+" ORDER BY money "+sort, func() error {
		wallets = append(wallets, info)
		return nil
	}, sl...)
	return
}

// 更新钱包 no lock
func (s *Storage) updateWalletOf(uid int64, money int) (err error) {
	return s.db.Insert("storage", &Wallet{
		UID:   uid,
		Money: money,
	})
}
