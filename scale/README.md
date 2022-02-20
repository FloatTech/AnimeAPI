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
> 注意：由于云函数内存较小，请将图片分辨率控制在`0.25MP`，即`500*500`之内
#### 必要参数
> GET https://bilibiliai.azurewebsites.net/api/scale?url=图片链接
#### 可选参数
> model=[conservative, no-denoise, denoise1x, denoise2x, denoise3x]


> scale=[2, 3, 4]


> tile=[0, 1, 2, 3, 4]

特别地，`scale=[3, 4]`时没有模型`[denoise1x, denoise2x]`

返回：`webp`格式的输出图片
