package niu

import (
	"fmt"
	"github.com/RomiChan/syncx"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var (
	migratedGroups = syncx.Map[string, bool]{} // key: string, value: bool
)

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

func tableFor(gid int64, prefix string) *gorm.DB {

	switch prefix {
	case ur:
		err := ensureTable[userInfo](gid, ur)
		if err != nil {
			logrus.Errorf("ensureTable error: %v", err)
			return nil
		}
	case ac:
		err := ensureTable[AuctionInfo](gid, ac)
		if err != nil {
			logrus.Errorf("ensureTable error: %v", err)
			return nil
		}
	}

	tableName := fmt.Sprintf("group_%d_%s_info", gid, prefix)
	return db.Table(tableName)
}

func listUsers(gid int64) (users, error) {
	var us users
	err := tableFor(gid, ur).Find(&us).Error
	return us, err
}

func listAuction(gid int64) ([]AuctionInfo, error) {
	var as []AuctionInfo
	err := tableFor(gid, ac).Order("money DESC").Find(&as).Error
	return as, err
}

func createUser(gid int64, user *userInfo, fix string) error {
	return tableFor(gid, fix).Create(user).Error
}

func getUserByID(gid int64, uid int64) (*userInfo, error) {
	var user userInfo
	err := tableFor(gid, ur).Where("user_id = ?", uid).First(&user).Error
	return &user, err
}

func updatesUserByID(gid int64, id int64, fields map[string]interface{}) error {
	return tableFor(gid, ur).Where("user_id = ?", id).Updates(fields).Error
}

func deleteUserByID(gid int64, id int64) error {
	return tableFor(gid, ur).Where("user_id = ?", id).Delete(&userInfo{}).Error
}
