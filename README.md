# CSS CHECKER

`css-checker` checks your css styles for duplications and find the diff among `css classes` with high similarity in seconds. It is designed to avoid redundant or similar css between files and to work well for both local development, and for automation like CI.

Colors check, long scripts, unused CSS classes warning of css are also supported by default. This project is provided by [Xiemala Team](`https://xiemala.com`), it helps in remove hundreds of similar css classes for developers in this project.

## Install

```
go install github.com/ruilisi/css-checker@latest
```

(With go version before 1.17, use `go get github.com/ruilisi/css-checker`). Or download from [releases](https://github.com/ruilisi/css-checker/releases)

## Usage

#### Run

- `cd PROJECT_WITH_CSS_FILES` and just run:

```
css-checker
```

- (Alpha Feature: Find classes that not referred by your js/jsx/ts/tsx/html code): `css-checker -path=[YOUR_PROJECT_PATH] -unused`

- (To Set your project path and ignore paths): `css-checker -path=[YOUR_PROJECT_PATH] -ignores=node_modules,packages,others*`

![DEMO](https://assets.ruilisi.com/css-checker-demo.gif)

(Check and show the diff among simullar classes (>=80%). Colors, long scripts that used more then once will also be pointed out by default. Check `css-checker -help` for customized options.)

Colors with `rgb/rgba/hsl/hsla/hex` will be converted to rbga and compared together.

#### Run with path and ignores

- `css-checker -path=YOUR_PROJECT_PATH -ignores=node_modules,packages`

#### Basic commands

- `-help`: prints help and exits
- `-colors`: whether to check colors (default true)
- `-ignores`: string paths and files to be ignored (e.g. node_modules,\*.example.css)
- `-length-threshold`: int Min length of a single style value (no including the key) that to be considered as long script line (default 20)
- `-long-line`: whether to check duplicated long script lines (default true)
- `-path`: string set path to files, default to be current folder (default ".")
- `-sections`: whether to check sections duplications (default true)
- `-sim`: whether to check similar css classes (>=80% && < 100%) (default true)
- `-version`: prints current version and exits

#### Output:

![image.png](https://assets.ruilisi.com/t=yDNXWrmyg+V6mUzCAG7A==)

#### How we get similarities between classes?:

0. Hash each line of class (aka. `section` in our code), Generate map: `LineHash -> Section`.
1. Convert map `LineHash -> Section` => `[SectionIndex1][SectionIndex2] -> Duplicated Hashes`, n for identical hash, section stands for css class.
2. In map: `[SectionIndex1][SectionIndex2]` -> `Duplicated Hashes`, number of the duplicated hashes stands for duplicated lines between classes.

#### Similarity Check:

Check the similarity (>=80% && < 100%) between classes. This will print the same line in between classes.

![image.png](https://assets.ruilisi.com/bzljM=P4Mz+dmtHKNvdHtg==)

#### Long Script Line Check:

Long scripts can be saved as varirables to make your life easiler. This will only alert when long scriptes are used for more then once.

![image.png](https://assets.ruilisi.com/5bdqZTuLTzJCaGSynA7+2w==)

#### Colors Check:

Check colors in HEX/RGB/RGBA/HSL/HSLA that used more then once in your code. As for supporting of diffrent themes and possible future updates of you color set, you may consider to put them as css variables.

![image.png](https://assets.ruilisi.com/iqmnGQHwglb+pxE3kr3L1Q==)

#### Duplicated CSS Classes:

Similar to `Similarity Check` but put those classes that are total identical to each other.

## Build & Release

- `make test-models`
- `make build`
- `make release`
