# AnimeAPI/classify
Zerobot-ACGImage插件的AI评分接口，也可单独引用。

# 接口说明

## func Init(dataPath string)
设置数据缓存路径

## func Flush()
手动刷新图片缓存及访问时间戳。

## func CanVisit(delay int64) bool
距上次返回`true`时间间隔大于`delay`秒则返回`true`并刷新时间戳，用以避免频繁访问。`Classify`函数并不会做验证，因此请务必在调用`Classify`前手动验证，以避免由于缓存时间戳不正确导致的无法加载图片。

请确保`delay`至少大于1，否则可能导致无法加载图片。

## func Classify(targeturl string, noimg bool) (int, int64, string, string)
用AI对`targeturl`指向的图片打分。返回值：class lastvisit dhash comment。

- 如果`noimg==true`，将不下载图片。
- 如果`noimg==false`，将下载图片到`dataPath/cache`+`lastvisit`。
- `dhash`为图片的dhash值的[base16384](https://github.com/fumiama/base16384)编码。
- `comment`为针对该`class`的评语，详见下方打分等级。

# 打分等级

> 普通图

- [0]这啥啊
- [1]普通欸
- [2]有点可爱
- [3]不错哦
- [4]很棒
- [5]我好啦!

> r18

由于模型并不精确，目前对于非`lolicon`图片有70%可能误判（`lolicon`图片可保证无误），请勿将其作为鉴黄模型使用，仅供娱乐。

- [6]影响不好啦!
- [7]太涩啦，🐛了!（非`lolicon`图片的最高等级）
- [8]已经🐛不动啦...（仅`lolicon`图片有此等级，表示误判但被纠正的图片）
