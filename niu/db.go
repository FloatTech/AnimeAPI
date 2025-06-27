package niu

import (
	"fmt"
	"github.com/RomiChan/syncx"
	"gorm.io/gorm"
)

var (
	migratedGroups = syncx.Map[string, bool]{} // key: string, value: bool
)

type t interface {
	userInfo | AuctionInfo
}

func ensureUserInfoTable[T t](gid int64, prefix string) error {
	table := fmt.Sprintf("group_%d_%s_info", gid, prefix)
	if _, ok := migratedGroups.Load(table); ok {
		return nil
	}
	err := db.Table(table).AutoMigrate(new(T))
	if err != nil {
		return err
	}

	// 设置为已迁移
	migratedGroups.Store(table, true)
	return nil
}

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

func getUserByID(gid int64, uid int64, fix string) (*userInfo, error) {
	var user userInfo
	err := TableFor(gid, fix).Where("user_id = ?", uid).First(&user).Error
	return &user, err
}

func updatesUserByID(gid int64, id int64, fields map[string]interface{}, fix string) error {
	return TableFor(gid, fix).Where("user_id = ?", id).Updates(fields).Error
}

func deleteUserByID(gid int64, id int64, fix string) error {
	return TableFor(gid, fix).Where("user_id = ?", id).Delete(&userInfo{}).Error
}
