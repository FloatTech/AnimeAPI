# AnimeAPI/picture
Zerobot 的索取图片接口
# 使用
引用时，一般使用如下形式
```go
import (
  "github.com/FloatTech/AnimeAPI/picture"
)

zero.OnKeywordGroup([]string{"欲匹配短语1", "欲匹配短语2" ...}, picture.CmdMatch(), picture.MustGiven())
```
# 函数定义
## func CmdMatch() zero.Rule
命令匹配
## func Exists() zero.Rule
消息含有图片返回`true`
## func MustGiven() zero.Rule
消息不存在图片阻塞60秒至有图片，超时返回`false`
