package niu

import (
	"fmt"

	"github.com/RomiChan/syncx"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var (
	migratedGroups = syncx.Map[string, struct{}]{} // key: string, value: struct{}
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
	migratedGroups.Store(table, struct{}{})
	return nil
}

func tableFor(gid int64, prefix string) *gorm.DB {
	switch prefix {
	case usr:
		err := ensureTable[userInfo](gid, usr)
		if err != nil {
			logrus.Errorf("ensureTable error: %v", err)
			return nil
		}
	case auct:
		err := ensureTable[AuctionInfo](gid, auct)
		if err != nil {
			logrus.Errorf("ensureTable error: %v", err)
			return nil
		}
	}

	return db.Table(fmt.Sprintf("group_%d_%s_info", gid, prefix))
}

func listUsers(gid int64) (users, error) {
	var users users
	err := tableFor(gid, usr).Find(&users).Error
	return users, err
}

func listAuction(gid int64) ([]AuctionInfo, error) {
	var as []AuctionInfo
	err := tableFor(gid, auct).Order("money DESC").Find(&as).Error
	return as, err
}

func createUser(gid int64, user *userInfo, fix string) error {
	return tableFor(gid, fix).Create(user).Error
}

func getUserByID(gid int64, uid int64) (*userInfo, error) {
	var user userInfo
	err := tableFor(gid, usr).Where("user_id = ?", uid).First(&user).Error
	return &user, err
}

func updatesUserByID(gid int64, id int64, fields map[string]interface{}) error {
	return tableFor(gid, usr).Where("user_id = ?", id).Updates(fields).Error
}

func deleteUserByID(gid int64, id int64) error {
	return tableFor(gid, usr).Where("user_id = ?", id).Delete(&userInfo{}).Error
}
