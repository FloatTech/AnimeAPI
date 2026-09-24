// Package pixiv (copied from plugin/pixiv/api/pixiv_api.go)
package pixiv

import (
	"errors"
	"fmt"
	"sort"

	"github.com/FloatTech/AnimeAPI/pixiv/model"
	log "github.com/sirupsen/logrus"
)

// PixivAPI ...
type PixivAPI struct {
	Client *Client
	Token  *TokenStore
}

// NewPixivAPI ...
func NewPixivAPI(refreshToken string) *PixivAPI {
	c := NewClient()
	return &PixivAPI{
		Client: c,
		Token:  NewTokenStore(refreshToken, c),
	}
}

// FetchPixivByPID ...
func (p *PixivAPI) FetchPixivByPID(pid int64) (*model.IllustCache, error) {
	url := fmt.Sprintf("https://app-api.pixiv.net/v1/illust/detail?illust_id=%d", pid)
	accessToken, err := p.Token.GetAccessToken()
	if err != nil {
		return nil, err
	}
	rawData, err := p.Client.SearchPixivIllustrations(accessToken, url)
	if err != nil {
		return nil, err
	}
	if rawData == nil || rawData.Illust == nil {
		return nil, errors.New("pixiv 返回数据为空或结构不匹配")
	}
	return convertToIllustCache(rawData.Illust)
}

// FetchPixivByUser ...
func (p *PixivAPI) FetchPixivByUser(uid int64, limit int, pids []int64) ([]model.IllustCache, error) {
	url := fmt.Sprintf("https://app-api.pixiv.net/v1/user/illusts?user_id=%d&type=illust", uid)

	excludeCache := make(map[int64]struct{})

	for _, pid := range pids {
		excludeCache[pid] = struct{}{}
	}
	return p.fetchPixivCommon(url, limit, nil, excludeCache)
}

// FetchPixivRecommend ...
func (p *PixivAPI) FetchPixivRecommend(limit int) ([]model.IllustCache, error) {
	firstURL := "https://app-api.pixiv.net/v1/illust/recommended?filter=for_ios"
	return p.fetchPixivCommon(firstURL, limit, nil, nil) // 不做R18过滤，不排缓存
}

// FetchPixivIllusts ...
func (p *PixivAPI) FetchPixivIllusts(keyword string, isR18Req bool, limit int, cachedIDs []int64) ([]model.IllustCache, error) {
	cachedMap := make(map[int64]struct{}, len(cachedIDs))
	if len(cachedIDs) > 0 {
		for _, id := range cachedIDs {
			cachedMap[id] = struct{}{}
		}
	}

	firstURL := buildPixivSearchURL(keyword)
	return p.fetchPixivCommon(firstURL, limit, &isR18Req, cachedMap, keyword)
}

// GetIllustsByKeyword ...
func (p *PixivAPI) GetIllustsByKeyword(keyword string, limit int, cachedIllust []model.IllustCache, cached []int64) ([]model.IllustCache, error) {

	r18Req := IsR18(keyword)
	keyword = RemoveR18Keywords(keyword)

	// 如果查到了，直接返回
	if len(cachedIllust) == limit {
		return cachedIllust, nil
	}

	// 设置一个保底的关键词
	if keyword == "" && r18Req {
		keyword = "R-18"
	}

	// 计算还需要几张图片
	needed := 0
	if len(cachedIllust) < limit {
		needed = limit - len(cachedIllust)
	}

	log.Printf("从数据库读到%d,还需要下载%d\n", len(cachedIllust), needed)
	// 缓存没数据 -> 调用Pixiv API拉取
	pixivResults, err := p.FetchPixivIllusts(keyword, r18Req, needed, cached)
	if err != nil && len(cachedIllust) == 0 {
		return nil, err
	}

	// 如果Pixiv也没查到直接返回空
	if len(pixivResults) == 0 && len(cachedIllust) == 0 {
		return nil, errors.New("这个关键词可能没有找到符合条件的图片或出现未知错误")
	}

	if len(cachedIllust) > 0 && len(pixivResults) == 0 {
		log.Println("http没有找到图片")
		return cachedIllust, nil
	}

	pixivResults = append(pixivResults, cachedIllust...)

	if len(pixivResults) >= limit {
		pixivResults = pixivResults[:limit]
	}

	log.Println("预计发送", len(pixivResults), "张图片")

	return pixivResults, nil
}

func (p *PixivAPI) fetchPixivCommon(
	firstURL string,
	limit int,
	isR18Req *bool,
	excludeCache map[int64]struct{},
	keywords ...string,
) ([]model.IllustCache, error) {

	accessToken, err := p.Token.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// 高质量图（≥1000）
	high := make([]model.IllustCache, 0, limit)
	// 低质量图（<1000）
	low := make([]model.IllustCache, 0, limit)

	seen := make(map[int64]struct{})
	url := firstURL

	for url != "" {
		rawData, err := p.Client.SearchPixivIllustrations(accessToken, url)
		if err != nil {
			return nil, err
		}

		for i := range rawData.Illusts {
			raw := &rawData.Illusts[i]

			// 去重
			if _, ok := seen[raw.ID]; ok {
				continue
			}
			if excludeCache != nil {
				if _, ok := excludeCache[raw.ID]; ok {
					continue
				}
			}
			seen[raw.ID] = struct{}{}

			tagNames := extractTagNames(raw.Tags)

			if isR18Req != nil && *isR18Req && !hasR18Tag(tagNames) {
				continue
			}

			// 转换
			ill, err := convertToIllustCache(raw)
			if err != nil {
				continue
			}
			if len(keywords) > 0 {
				ill.Keyword = keywords[0]
			}

			// R18过滤
			if isR18Req != nil && ill.R18 != *isR18Req {
				continue
			}

			// 判断高质量
			if ill.Bookmarks >= 1000 {
				high = append(high, *ill)
				// 高质量够了就直接返回
				if len(high) >= limit {
					return high[:limit], nil
				}
			}
			// 低质量池还没满 → 接受
			if len(low) <= limit {
				low = append(low, *ill)
			}
			// 低质量够 limit 就不再放入，避免爆炸增长
		}

		url = rawData.NextURL
	}

	// 低质量排序（按收藏数倒序）
	sort.Slice(low, func(i, j int) bool {
		return low[i].Bookmarks > low[j].Bookmarks
	})

	// 用低质量中的高质量去补齐不足的图
	high = append(high, low...)
	if len(high) > limit {
		high = high[:limit]
	}

	return high, nil
}
