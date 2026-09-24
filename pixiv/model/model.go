// Package model ...
package model

import (
	"github.com/jinzhu/gorm"
	"gorm.io/datatypes"
)

// IllustCache 插画缓存表
type IllustCache struct {
	gorm.Model

	PID         int64          `gorm:"unique_index:idx_keyword_pid;not null;column:pid"` // Pixiv 作品 ID
	UID         int64          `gorm:"default:0;not null;column:uid"`                    // 插画作者的id
	Keyword     string         `gorm:"unique_index:idx_keyword_pid;type:varchar(255)"`   // 搜索关键词
	Title       string         `gorm:"type:varchar(255)"`                                // 标题
	AuthorName  string         `gorm:"type:varchar(255)"`                                // 用户名
	ImageURL    string         `gorm:"type:varchar(500)"`                                // 大图地址
	OriginalURL string         `gorm:"type:varchar(500)"`                                // 原图地址
	R18         bool           `gorm:"not null;default:false"`                           // 是否为 R-18 作品
	Bookmarks   int64          // 收藏数
	TotalView   int64          // 总浏览数
	CreateDate  string         // 创建日期
	PageCount   int64          `gorm:"default:1"` // 页数
	Tags        datatypes.JSON `gorm:"type:json"` // 插画的所有标签 方便后续查找
}

// SentImage 已发送记录表
type SentImage struct {
	gorm.Model

	GroupID int64 `gorm:"index:idx_group_pid;not null"`            // 群组 ID
	PID     int64 `gorm:"index:idx_group_pid;not null;column:pid"` // 插画 PID
}

// GroupR18Permission R18权限表
type GroupR18Permission struct {
	gorm.Model
	GroupID int64 `gorm:"unique_index"`
}

// RefreshToken ...
type RefreshToken struct {
	gorm.Model

	User int64 `gorm:"unique"`

	Token string
}

// PixivProxyConfig 保存 Pixiv API 客户端的代理配置
type PixivProxyConfig struct {
	gorm.Model

	Name  string `gorm:"unique_index"`
	Proxy string
}

// RootEntity ...
type RootEntity struct {
	Illusts         []IllustsEntity `json:"illusts"`
	Illust          *IllustsEntity  `json:"illust"`
	NextURL         string          `json:"next_url"`
	SearchSpanLimit int64           `json:"search_span_limit"`
	ShowAi          bool            `json:"show_ai"`
}

// IllustsEntity ...
type IllustsEntity struct {
	ID              int64                `json:"id"`
	Title           string               `json:"title"`
	Type            string               `json:"type"`
	ImageURLs       ImageUrlsEntity      `json:"image_urls"`
	Caption         string               `json:"caption"`
	Restrict        int64                `json:"restrict"`
	User            UserEntity           `json:"user"`
	Tags            []TagsEntity         `json:"tags"`
	CreateDate      string               `json:"create_date"`
	PageCount       int64                `json:"page_count"`
	Width           int64                `json:"width"`
	Height          int64                `json:"height"`
	SanityLevel     int64                `json:"sanity_level"`
	XRestrict       int64                `json:"x_restrict"`
	MetaSinglePage  MetaSinglePageEntity `json:"meta_single_page"`
	MetaPages       []MetaPage           `json:"meta_pages"`
	TotalView       int64                `json:"total_view"`
	TotalBookmarks  int64                `json:"total_bookmarks"`
	IsBookmarked    bool                 `json:"is_bookmarked"`
	Visible         bool                 `json:"visible"`
	IsMuted         bool                 `json:"is_muted"`
	IllustAiType    int64                `json:"illust_ai_type"`
	IllustBookStyle int64                `json:"illust_book_style"`
}

// MetaPage ...
type MetaPage struct {
	ImageURLs struct {
		SquareMedium string `json:"square_medium"`
		Medium       string `json:"medium"`
		Large        string `json:"large"`
		Original     string `json:"original"`
	} `json:"image_urls"`
}

// ImageUrlsEntity ...
type ImageUrlsEntity struct {
	SquareMedium string `json:"square_medium"`
	Medium       string `json:"medium"`
	Large        string `json:"large"`
}

// UserEntity ...
type UserEntity struct {
	ID               int64                  `json:"id"`
	Name             string                 `json:"name"`
	Account          string                 `json:"account"`
	ProfileImageUrls ProfileImageUrlsEntity `json:"profile_image_urls"`
	IsFollowed       bool                   `json:"is_followed"`
	IsAcceptRequest  bool                   `json:"is_accept_request"`
}

// ProfileImageUrlsEntity ...
type ProfileImageUrlsEntity struct {
	Medium string `json:"medium"`
}

// TagsEntity ...
type TagsEntity struct {
	Name           string      `json:"name"`
	TranslatedName interface{} `json:"translated_name"`
}

// MetaSinglePageEntity ...
type MetaSinglePageEntity struct {
	OriginalImageURL string `json:"original_image_url"`
}
