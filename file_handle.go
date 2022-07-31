
package main

import (
	"bufio" //缓存IO
	"fmt"
	"os"
	// "log"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 创建一个输出文件
func createOutuputFile(fileName string) {

	var f = "css-check.html"
	if (len(fileName) != 0) {
		f = fileName
	}

	// 初始化这个html文件
	file, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE, 0666)

	// 关闭文件句柄
	defer file.Close()
	// 创建一个writer
	writer := bufio.NewWriter(file)
	content := "<html><head><style>p {margin-bottom:6px;margin-top:6px;}.t{font-size:24px;font-weight:bold;margin-bottom: 15px; border-top: 1px solid; padding-top: 14px;}</style></head><body style='background: black; color: white'>\n"
	_, err = writer.WriteString(content)
	if err != nil {
		fmt.Printf("write file err, %v", err)
	}
	// 最后真正写入文件
	writer.Flush()
	fmt.Println("write file success")


	
}


// 获取输出文件引用
func getHtmlFile(fileName string) *os.File{
	file,_ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return file
}

// 写字符串到文件
func writeToFile(file *os.File,txt string , hasParameter bool) {
	if (!hasParameter) {
		return 
	}
	file.WriteString(txt)
}


