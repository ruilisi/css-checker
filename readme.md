# CSS CHECKER

### Basic commands

- `-help`: prints help and exits
- `-colors`: whether to check colors (default true)
- `-ignores`: string paths and files to be ignored (e.g. node_modules,\*.example.css)
- `-length-threshold`: int Min length of a single style value (no including the key) that to be considered as long script line (default 20)
- `-long-line`: whether to check duplicated long script lines (default true)
- `-path`: string set path to files, default to be current folder (default ".")
- `-sections`: whether to check sections duplications (default true)
- `-sim`: whether to check similar css classes (>=80% && < 100%) (default true)
- `-version`: prints current version and exits

### Output:

![image.png](https://assets.ruilisi.com/t=yDNXWrmyg+V6mUzCAG7A==)

### How we get similarities between classes?:

0. Hash each line of class (aka. `section` in our code), Generate map: `LineHash -> Section`.
1. Convert map `LineHash -> Section` => `[SectionIndex1][SectionIndex2] -> Duplicated Hashes`, n for identical hash, section stands for css class.
2. In map: `[SectionIndex1][SectionIndex2]` -> `Duplicated Hashes`, number of the duplicated hashes stands for duplicated lines between classes.

### Similarity Check:

Check the similarity (>=80% && < 100%) between classes. This will print the same line in between classes.

![image.png](https://assets.ruilisi.com/bzljM=P4Mz+dmtHKNvdHtg==)

### Long Script Line Check:

Long scripts can be saved as varirables to make your life easiler. This will only alert when long scriptes are used for more then once.

![image.png](https://assets.ruilisi.com/5bdqZTuLTzJCaGSynA7+2w==)

### Colors Check:

Check colors in HEX/RGB/RGBA/HSL/HSLA that used more then once in your code. As for supporting of diffrent themes and possible future updates of you color set, you may consider to put them as css variables.

![image.png](https://assets.ruilisi.com/iqmnGQHwglb+pxE3kr3L1Q==)

### Duplicated CSS Classes:

Similar to `Similarity Check` but put those classes that are total identical to each other.
