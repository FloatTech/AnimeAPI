package niu

import (
	"fmt"
	"github.com/RomiChan/syncx"
	"github.com/jinzhu/gorm"
	"sync"
)

var (
	migratedGroups = syncx.Map[string, bool]{} // key: string, value: bool
	tableHooks     []tableHook
	hooksMtx       sync.RWMutex
)

type tableHook func(gid int64) error

type model struct {
	*gorm.DB
}

func ensureTable[T userInfo | AuctionInfo](gid int64, prefix string) error {
	table := fmt.Sprintf("group_%d_%s_info", gid, prefix)
	if _, ok := migratedGroups.Load(table); ok {
		return nil
	}
	err := db.Table(table).AutoMigrate(new(T)).Error
	if err != nil {
		return err
	}

	// 设置为已迁移
	migratedGroups.Store(table, true)
	return nil
}

func ensureUserInfo(gid int64) error {
	return ensureTable[userInfo](gid, ur)
}

func ensureAuctionInfo(gid int64) error {
	return ensureTable[AuctionInfo](gid, ac)
}

// registerTableHook 注册钩子
func registerTableHook(h ...tableHook) {
	hooksMtx.Lock()
	defer hooksMtx.Unlock()
	tableHooks = append(tableHooks, h...)
}

// TableFor 大写是为了防止数据操作哪里有问题留个保底可以在zbp的项目里直接改
func TableFor(gid int64, prefix string) *model {
	// 先执行钩子
	hooksMtx.RLock()
	for _, h := range tableHooks {
		if err := h(gid); err != nil {
			panic(fmt.Sprintf("执行表钩子失败: %v", err))
		}
	}
	hooksMtx.RUnlock()

	tableName := fmt.Sprintf("group_%d_%s_info", gid, prefix)
	return &model{db.Table(tableName)}
}

func listUsers(gid int64) (users, error) {
	var us users
	err := TableFor(gid, ur).Find(&us).Error
	return us, err
}

func listAuction(gid int64) ([]AuctionInfo, error) {
	var as []AuctionInfo
	err := TableFor(gid, ac).Order("money DESC").Find(&as).Error
	return as, err
}

func createUser(gid int64, user *userInfo, fix string) error {
	return TableFor(gid, fix).Create(user).Error
}

func getUserByID(gid int64, uid int64) (*userInfo, error) {
	var user userInfo
	err := TableFor(gid, ur).Where("user_id = ?", uid).First(&user).Error
	return &user, err
}

func updatesUserByID(gid int64, id int64, fields map[string]interface{}) error {
	return TableFor(gid, ur).Where("user_id = ?", id).Updates(fields).Error
}

func deleteUserByID(gid int64, id int64) error {
	return TableFor(gid, ur).Where("user_id = ?", id).Delete(&userInfo{}).Error
}
