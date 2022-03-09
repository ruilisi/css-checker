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
  <a href="README-zh_CN.md">查看README in English</a>
</p>

## 目标：

`css_checker` 会检查 css 样式是否存在重复，并找出它们之间的差异 `css classes` 在几秒钟内具有高度相似性。它的设计目的是避免文件之间出现冗余或类似的 css，并在本地开发和 CI 等自动化方面都能很好地工作。

默认情况下，该软件包还支持颜色检查、长脚本和 未被使用的 css 警告，以帮助开发人员优化 css 文件。该项目由[协码啦团队](`https://xiemala.com`)开发，它该项目优化了数百个冗余的 css 类。

## 安装(以下二选一)：

#### 使用 go install:

```
go install github.com/ruilisi/css-checker@latest
```

(对于 go 1.17 之前的版本，请使用 `go get github.com/ruilisi/css-checker`)。或者从[releases](https://github.com/ruilisi/css-checker/releases)下载。

#### 使用 npm:

```
npm install -g css-checker-kit
```

## 用法:

#### 运行:

- `cd PROJECT_WITH_CSS_FILES` 并且直接运行:

```
css-checker
```

- （Alpha 功能：查找 js/jsx/ts/tsx/html 代码未引用的 class）: `css-checker -path=[YOUR_PROJECT_PATH] -unused`
- （要设置项目路径及忽略路径）: `css-checker -path=[YOUR_PROJECT_PATH] -ignores=node_modules,packages,others*`

![DEMO](https://assets.ruilisi.com/css-checker-demo.gif)

（这可以检查类之间的相似性，并显示相似类之间的差异（>=80%）。默认情况下，还会指出使用了多次的颜色、长脚本。查看“css checker-help”中的自定义选项。）

将带有“rgb/rgba/hsl/hsla/hex”的颜色将转换为 rbga 并一起比较。

#### 按路径运行:

- `css-checker -path=YOUR_PROJECT_PATH`

#### 按路径忽略:

- 用户可将要忽略的文件名填进`gitignore`，`css_check`就会自动忽略这些文件，可通过`-unrestricted=true`来取消忽略。
- 具体命令可参考下面的基本命令。
- `-ignores=node_modules,packages`: 用来忽略特定的文件夹。

#### 关于 yaml 文件：

- `css_Checker`参数默认读取`css-checker.example.yaml`用户可在其中直接添加参数并且不需要在参数前面加上`-`。
- `-config=css-checker. .yaml`: 可在`.`之间输入来重命名文件。

#### 基本命令:

- `-help`: 输出帮助并退出
- `-colors`: 是否检查颜色（默认为 true）
- `-ignores`: 输出被忽略的路径和文件(e.g. node_modules,\*.example.css)
- `-length-threshold`: 一个被认为是长脚本行（默认为 20）的单一样式值的最小整型长度（不包括键）
- `-long-line`: 是否检查重复的长脚本行（默认为 true）
- `-path`: 文件路径的字符串，默认为当前文件夹（默认为"."）
- `-sections`: 是否检查部分重复（默认为 true）
- `-sim`: 是否检查类似的 CSS 类（>=80% && <100%)(默认为 true)
- `-version`: 输出现在的版本和退出
- `-unused`: 检查未被使用的 CSS class (默认为 false, Beta 功能)

#### 输出:

![image.png](https://assets.ruilisi.com/t=yDNXWrmyg+V6mUzCAG7A==)

#### 我们如何获取类之间的相似性？:

0. hash 类中的每一行（比如代码中的`section`),生成 map ：`LineHash -> Section`.
1. 转换 map `LineHash -> Section` => `[SectionIndex1][SectionIndex2] -> Duplicated Hashes`, n 代表相同的哈希，section 代表 CSS 类。
2. 在 map 中: `[SectionIndex1][SectionIndex2]` -> `Duplicated Hashes`, 重复的哈希数表示类之间的重复行。

#### 相似度检查:

- 检查类之间的相似度(>=80% && < 100%)。这将在类之间打印相同的行。

![image.png](https://assets.ruilisi.com/bzljM=P4Mz+dmtHKNvdHtg==)

#### 相似阈值:

- `-sim-threshold=`: 用户可用来自定义（>=20% && <=60% ）的相似阀值。
- `yaml:"sections"`: 用户可用该参数设置为完全相似来查询。

#### 长脚本行检查：

长脚本可以保存为多种脚本来使您的生活更轻松。这仅当长脚本使用超过一次时才会发出警报。

![image.png](https://assets.ruilisi.com/5bdqZTuLTzJCaGSynA7+2w==)

#### 颜色检查:

- 检查 HEX/RGB/RGBA/HSL/HSLA 在代码中使用不止一次的颜色。支持不同的主题并可能在未来更新你的颜色集，你可以考虑把它们作为 CSS 的变量。

![image.png](https://assets.ruilisi.com/iqmnGQHwglb+pxE3kr3L1Q==)

#### 重复的 CSS 类:

- 与`相似度检查`相似，但是会把那些完全相同的类放在一起。

## 构建&释放：

- `make test-models`
- `make build`
- `make release`
