# Zerobot-AnimeAPI-Classify
Zerobot-ACGImage插件的AI评分接口，也可单独引用。

# 接口说明

## func Init(dataPath string)
设置数据缓存路径

## func Flush()
手动刷新图片缓存及访问时间戳。

## func CanVisit(delay int64) bool
距上次返回`true`时间间隔大于`delay`秒则返回`true`并刷新时间戳，用以避免频繁访问。`Classify`函数并不会做验证，因此请务必在调用`Classify`前手动验证，以避免由于缓存时间戳不正确导致的无法加载图片。

请确保`delay`至少大于1，否则可能导致无法加载图片。

## func Classify(ctx *zero.Ctx, targeturl string, noimg bool)
用AI对`targeturl`指向的图片打分。

- 如果`noimg==true`，将用打分回复`ctx.Event.MessageID`指示的消息。
- 如果`noimg==false`，将发送该图片并针对该图片回复打分。

# 打分等级

> 普通图

- [0]这啥啊
- [1]普通欸
- [2]有点可爱
- [3]不错哦
- [4]很棒
- [5]我好啦!

> r18

由于模型并不精确，目前有较大可能误判，请勿将其作为鉴黄模型使用，仅供娱乐。

- [6]影响不好啦!
- [7]太涩啦，🐛了!
- [8]已经🐛不动啦...