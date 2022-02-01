# scale
采用[叔叔的模型](https://github.com/bilibili/ailab)放大二次元图片

## API
#### 必要参数
> GET https://sayuri.fumiama.top/scale/?url=图片链接

#### 可选参数
> model=[conservative, no-denoise, denoise1x, denoise2x, denoise3x]


> scale=[2, 3, 4]


> tile=[0, 1, 2, 3, 4]

特别地，`scale=[3, 4]`时没有模型`[denoise1x, denoise2x]`

返回：`webp`格式的输出图片

## 效果

<table>
	<tr>
		<td align="center"><img src="test/in.png"></td>
		<td align="center"><img src="test/out.webp"></td>
	</tr>
    <tr>
		<td align="center">输入</td>
		<td align="center">输出</td>
	</tr>
</table>