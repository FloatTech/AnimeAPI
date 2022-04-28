# scale
采用[叔叔的模型](https://github.com/bilibili/ailab)放大二次元图片

<table>
	<tr>
		<td align="center"><img src="test/0.jpg"></td>
		<td align="center"><img src="test/out.webp"></td>
	</tr>
    <tr>
		<td align="center">输入</td>
		<td align="center">输出</td>
	</tr>
</table>

## API
> 注意：由于云函数内存较小, 请将图片分辨率控制在`0.25MP`, 即`500*500`之内
```go
func Get(u string, model, scale, tile int) ([]byte, error)
func Post(body io.Reader, model, scale, tile int) ([]byte, error)
```
返回：`webp`格式的输出图片
