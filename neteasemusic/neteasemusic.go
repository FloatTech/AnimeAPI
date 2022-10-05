// Package netease 网易云音乐一些简单的API:搜歌、下歌、搜歌词、下歌词
package netease

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	"github.com/pkg/errors"
)

// wyy搜歌结果
type searchResult struct {
	Result struct {
		Songs []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				ID        int           `json:"id"`
				Name      string        `json:"name"`
				PicURL    interface{}   `json:"picUrl"`
				Alias     []interface{} `json:"alias"`
				AlbumSize int           `json:"albumSize"`
				PicID     int           `json:"picId"`
				FansGroup interface{}   `json:"fansGroup"`
				Img1V1URL string        `json:"img1v1Url"`
				Img1V1    int           `json:"img1v1"`
				Trans     interface{}   `json:"trans"`
			} `json:"artists"`
			Album struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Artist struct {
					ID        int           `json:"id"`
					Name      string        `json:"name"`
					PicURL    interface{}   `json:"picUrl"`
					Alias     []interface{} `json:"alias"`
					AlbumSize int           `json:"albumSize"`
					PicID     int           `json:"picId"`
					FansGroup interface{}   `json:"fansGroup"`
					Img1V1URL string        `json:"img1v1Url"`
					Img1V1    int           `json:"img1v1"`
					Trans     interface{}   `json:"trans"`
				} `json:"artist"`
				PublishTime int64 `json:"publishTime"`
				Size        int   `json:"size"`
				CopyrightID int   `json:"copyrightId"`
				Status      int   `json:"status"`
				PicID       int64 `json:"picId"`
				Mark        int   `json:"mark"`
			} `json:"album"`
			Duration    int           `json:"duration"`
			CopyrightID int           `json:"copyrightId"`
			Status      int           `json:"status"`
			Alias       []interface{} `json:"alias"`
			Rtype       int           `json:"rtype"`
			Ftype       int           `json:"ftype"`
			Mvid        int           `json:"mvid"`
			Fee         int           `json:"fee"`
			RURL        interface{}   `json:"rUrl"`
			Mark        int           `json:"mark"`
		} `json:"songs"`
		SongCount int `json:"songCount"`
	} `json:"result"`
	Code int `json:"code"`
}

// 歌词内容
type musicLrc struct {
	LyricVersion int    `json:"lyricVersion"`
	Lyric        string `json:"lyric"`
	Code         int    `json:"code"`
}

// 搜索网易云音乐歌曲
//
// keyword:搜索内容 n:输出数量
//
// list:map[歌曲名称]歌曲ID
func SearchMusic(keyword string, n int) (list map[string]int, err error) {
	list = make(map[string]int, n)
	requestURL := "http://music.163.com/api/search/get/web?type=1&limit=" + strconv.Itoa(n) + "&s=" + url.QueryEscape(keyword)
	data, err := web.GetData(requestURL)
	if err != nil {
		return
	}
	var searchResult searchResult
	err = json.Unmarshal(data, &searchResult)
	if err != nil {
		return
	}
	if searchResult.Code != 200 {
		err = errors.Errorf("Status Code: %d", searchResult.Code)
		return
	}
	for _, musicinfo := range searchResult.Result.Songs {
		musicName := musicinfo.Name
		// 歌手信息
		artistsmun := len(musicinfo.Artists)
		if artistsmun != 0 {
			musicName += " - "
			for i, artistsinfo := range musicinfo.Artists {
				if artistsinfo.Name != "" {
					musicName += artistsinfo.Name
				}
				if i != 0 && i < artistsmun-1 {
					musicName += "&"
				}
			}
		}
		// 出自信息
		if len(musicinfo.Alias) != 0 {
			musicName += " - " + musicinfo.Alias[0].(string)
		}
		// 记录歌曲信息
		list[musicName] = musicinfo.ID
	}
	return
}

// 下载网易云音乐(歌曲ID，歌曲名称，下载路径)
func DownloadMusic(musicID int, musicName, pathOfMusic string) error {
	downMusic := pathOfMusic + "/" + musicName + ".mp3"
	musicURL := "http://music.163.com/song/media/outer/url?id=" + strconv.Itoa(musicID)
	if file.IsNotExist(downMusic) {
		// 检查歌曲是否存在
		response, err := http.Head(musicURL)
		if err != nil {
			return err
		}
		_ = response.Body.Close()
		if response.StatusCode != 200 {
			err = errors.Errorf("Status Code: %d", response.StatusCode)
			return err
		}
		// 下载歌曲
		err = file.DownloadTo(musicURL, downMusic, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// 搜索网易云音乐歌词(歌曲ID)
func SreachLrc(musicID int) (lrc string, err error) {
	musicURL := "http://music.163.com/api/song/media?id=" + strconv.Itoa(musicID)
	data, err := web.GetData(musicURL)
	if err != nil {
		return
	}
	var lrcinfo musicLrc
	err = json.Unmarshal(data, &lrcinfo)
	if err != nil {
		return
	}
	lrc = lrcinfo.Lyric
	return
}

// 下载网易云音乐歌词(歌曲ID，歌曲名称，下载路径)
func DownloadLrc(musicID int, musicName, pathOfMusic string) error {
	err := os.MkdirAll(pathOfMusic, 0777)
	if err != nil {
		return err
	}
	downfile := pathOfMusic + "/" + musicName + ".lrc"
	musicURL := "http://music.163.com/api/song/media?id=" + strconv.Itoa(musicID)
	if file.IsNotExist(downfile) {
		data, err := web.GetData(musicURL)
		if err != nil {
			return err
		}
		var lrcinfo musicLrc
		err = json.Unmarshal(data, &lrcinfo)
		if err != nil {
			return err
		}
		if lrcinfo.Code != 200 {
			err = errors.Errorf("Status Code: %d", lrcinfo.Code)
			return err
		}
		if lrcinfo.Lyric != "" {
			err = os.WriteFile(downfile, binary.StringToBytes(lrcinfo.Lyric), 0666)
			if err != nil {
				return err
			}
		} else {
			err = errors.Errorf("该歌曲无歌词")
			return err
		}
	}
	return nil
}
