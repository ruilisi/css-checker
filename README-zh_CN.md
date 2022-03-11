<p align="center">
  <a href="https://ruilisi.com/">
    <img alt="CSS-CHECKER" src="https://assets.ruilisi.com/cgULF9oHro3e1kSHXTfZYA==" width="211"/>
  </a>
</p>
<h1 align="center">Css Checker - 让CSS小而美</h1>
<p align="center">
  <a href="https://drone.ruilisi.com/ruilisi/css-checker" title="Build Status">
    <img src="https://drone.ruilisi.com/api/badges/ruilisi/css-checker/status.svg?ref=refs/heads/master">
  </a>
</p>
<p align="center">
  <a href="README.md">English Version</a>
</p>

## 目标

`css-checker` 会检查您的 css 样式是否有重复，并在几秒钟内找到具有高度相似性的 `css classes` 之间的差异。它可以避免文件之间出现冗余或类似的 css，并适用于本地开发和 CI 等自动化。

默认情况下它还支持颜色检查、长脚本以及未使用的 CSS classes 的警告。该项目由 [Xiemala Team](`https://xiemala.com`) 提供，它有助于为该项目中的开发人员删除数百个类似的 css classes。

## 安装

#### 使用 go:

```
go install github.com/ruilisi/css-checker@latest
```

(对于 go 1.17 之前的版本，请使用 `go get github.com/ruilisi/css-checker`)。或者从[releases](https://github.com/ruilisi/css-checker/releases)下载。

#### 使用 npm:

```
npm install -g css-checker-kit
```

## 用法

#### 运行

- `cd PROJECT_WITH_CSS_FILES` 并且直接运行:

```
css-checker
```

- （Alpha 功能：查找 js/jsx/ts/tsx/html 代码中未引用的 class）: `css-checker -path=[YOUR_PROJECT_PATH] -unused`

![DEMO](https://assets.ruilisi.com/css-checker-demo.gif)

（它可以检查 classes 之间的相似性，并显示（>=80%）相似度的 classes 之间的差异。默认情况下，他还能找出使用了多次的颜色、长脚本。可以用“css checker-help”来查看自定义选项。）

它能将带有“rgb/rgba/hsl/hsla/hex”的颜色转换为 rbga 并一起比较。

#### 按路径运行

- `css-checker -path=YOUR_PROJECT_PATH`

#### 忽略文件

- CSS-Checker 默认忽略 `.gitignore` 中的路径（可以使用 `-unrestricted=true` 来禁用此功能以读取所有文件）。
- 您可以使用：`-ignores=node_modules,packages`来添加要忽略的额外路径。

#### 配置文件

- `css-checker.yaml`：CSS-Checker 会在您的项目路径中读取此 yaml 文件进行设置，您可以使用 `Basic Commands` 中的部分参数来设置此文件（不用带前导“-”）。
- 此项目中还提供了一个名为“css-checker.example.yaml”的示例 yaml 文件，您可以将其命名为“css-checker.yaml”使用。
- 您可以使用 `-config=YOUR_CONFIG_FILE_PATH`来指定您的配置文件。

#### 基本命令

- `colors`: 是否检查颜色（默认为 true）
- `config`:设置配置文件路径 (string, default './css-checker.yaml') (string, default '')
- `ignores`: 输出被忽略的路径和文件(e.g. node_modules,\*.example.css)
- `length-threshold`: 被视为长脚本行的单个样式值（不包括键）的最小长度（默认 20）
- `long-line`: 是否检查重复的长脚本行（默认为 true）
- `path`: 文件路径的字符串，默认为当前文件夹（默认为"."）
- `sections`: 是否检查部分重复（默认为 true）
- `sim`: 是否检查相似的 css classes（默认 true）
- `sim-threshold`：相似性检查的阈值（$\geq20%$ && $\lt100%$）（int 类型，如80表示80%，请注意此为相似性检查控制，完全相同的css classes检查由 `sections`控制）（默认为 80）
- `unrestricted`：搜索所有文件（gitignore）
- `unused`：是否检查未使用的 classes（Beta）
- `version`：打印当前版本并退出

#### 输出:

![image.png](https://assets.ruilisi.com/t=yDNXWrmyg+V6mUzCAG7A==)

#### 我们是如何获取 classes 之间的相似性的？

0. hash classes 中的每一行（比如代码中的`section`),生成 map ：`LineHash -> Section`.
1. 转换 map `LineHash -> Section` => `[SectionIndex1][SectionIndex2] -> Duplicated Hashes`, n 代表相同的哈希，section 代表 CSS classes 。
2. 在 map 中: `[SectionIndex1][SectionIndex2]` -> `Duplicated Hashes`, 重复的哈希数表示在 classes 之间的重复行。

#### 相似性检查

检查 classes 之间的相似性 ($\geq(sim-threshold)$ && $\lt100$)，该功能会标注相似`classes`之前的diff (默认开启)。

- $sim-threshold$：使用 `-sim-threshold=` 参数或在配置 yaml 文件中设置 `sim-threshold:`，默认 80，最少 20。

![image.png](https://assets.ruilisi.com/bzljM=P4Mz+dmtHKNvdHtg==)

#### 重复的 CSS Classes

类似于 `Similarity Check` 但仅检查完全相同的 classes 较相似性检查效率更高 (默认开启)。

#### 长脚本行检查

将长脚本保存为多种脚本能让您的代码更简洁。只有当长脚本使用超过一次时才会发出警报 (默认开启)。

![image.png](https://assets.ruilisi.com/5bdqZTuLTzJCaGSynA7+2w==)

#### 颜色检查

- 检查 HEX/RGB/RGBA/HSL/HSLA 在代码中使用不止一次的颜色。您可以将其存为`CSS`变量, 以方便在未来可能的颜色及主题更新 (默认开启)。

![image.png](https://assets.ruilisi.com/iqmnGQHwglb+pxE3kr3L1Q==)

## 构建&释放

- `make test-models`
- `make build`
- `make release`
