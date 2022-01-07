# AnimeAPI/classify
Zerobot-ACGImage插件的AI评分接口，也可单独引用。

# 接口说明
## func Classify(targetURL string, isNoNeedImg bool) (int, string, string, []byte)
用AI对`targeturl`指向的图片打分。返回值：class dhash comment data。

- 如果`noimg==true`，将不下载图片。
- 如果`noimg==false`，将下载图片到 data。
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
