package niu

import (
	"fmt"
	"github.com/RomiChan/syncx"
	"github.com/jinzhu/gorm"
)

var (
	migratedGroups = syncx.Map[string, bool]{} // key: string, value: bool
)

func ensureUserInfoTable[T userInfo | AuctionInfo](gid int64, prefix string) error {
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

// TableFor 大写是为了防止数据操作哪里有问题留个保底可以在zbp的项目里直接改
func TableFor(gid int64, prefix string) *gorm.DB {
	return db.Table(fmt.Sprintf("group_%d_%s_info", gid, prefix))
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
