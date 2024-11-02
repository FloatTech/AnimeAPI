package niuniu

import (
	"github.com/FloatTech/floatbox/file"
	sql "github.com/FloatTech/sqlite"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type model struct {
	sql sql.Sqlite
	sync.RWMutex
}

type UserInfo struct {
	UID       int64
	Length    float64
	UserCount int
	WeiGe     int // 伟哥
	Philter   int // 媚药
	Artifact  int // 击剑神器
	ShenJi    int // 击剑神稽
	Buff1     int // 暂定
	Buff2     int // 暂定
	Buff3     int // 暂定
	Buff4     int // 暂定
	Buff5     int // 暂定
}

var (
	db = &model{}
)

func init() {
	if file.IsNotExist("data/niuniu") {
		err := os.MkdirAll("data/niuniu", 0755)
		if err != nil {
			panic(err)
		}
	}
	err := db.sql.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
}

// GetInitNiuNiu ...
func GetInitNiuNiu() float64 {
	return db.randLength()
}

// CreateGIDTable 创建一个表格
func CreateGIDTable(gid int64) error {
	return db.createGIDTable(gid)
}

// FindNiuNiuInfo 查询一个NiuNiu数据
func FindNiuNiuInfo(gid, uid int64) (*UserInfo, error) {
	return db.findNiuNiu(gid, uid)
}

// InsertNiuNiuInfo 修改NiuNiu数据
func InsertNiuNiuInfo(gid int64, u *UserInfo) error {
	return db.insertNiuNiu(u, gid)
}

// DeleteNiuNiu 删除一个NiuNiu
func DeleteNiuNiu(gid, uid int64) error {
	return db.deleteNiuNiu(gid, uid)
}

// GetAllNiuNiuInfo 获取当前群组所有NiuNiu数据
func GetAllNiuNiuInfo(gid int64) ([]*UserInfo, error) {
	return db.readAllTable(gid)
}

func (db *model) randLength() float64 {
	return float64(rand.Intn(9)+1) + (float64(rand.Intn(100)) / 100)
}

func (db *model) createGIDTable(gid int64) error {
	db.Lock()
	defer db.Unlock()
	return db.sql.Create(strconv.FormatInt(gid, 10), &UserInfo{})
}

func (db *model) findNiuNiu(gid, uid int64) (*UserInfo, error) {
	db.RLock()
	defer db.RUnlock()
	u := UserInfo{}
	err := db.sql.Find(strconv.FormatInt(gid, 10), &u, "where UID = "+strconv.FormatInt(uid, 10))
	return &u, err
}

func (db *model) insertNiuNiu(u *UserInfo, gid int64) error {
	db.Lock()
	defer db.Unlock()
	return db.sql.Insert(strconv.FormatInt(gid, 10), u)
}

func (db *model) deleteNiuNiu(gid, uid int64) error {
	db.Lock()
	defer db.Unlock()
	return db.sql.Del(strconv.FormatInt(gid, 10), "where UID = "+strconv.FormatInt(uid, 10))
}

func (db *model) readAllTable(gid int64) ([]*UserInfo, error) {
	db.Lock()
	defer db.Unlock()
	a, err := sql.FindAll[UserInfo](&db.sql, strconv.FormatInt(gid, 10), "where UserCount  = 0")
	return a, err
}
