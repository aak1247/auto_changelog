package utils

import (
	"bufio"
	"github.com/aak1247/gchangelog/configs"
	"io/ioutil"
	"os"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func InsertToFile(path string, content string, skipRows int) error {
	if !FileExists(path) {
		// 创建文件
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		// 写入文件头
		file.WriteString(configs.DefaultHead)
		file.Close()
	}
	// 打开文件以读取
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取文件的所有内容
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// 确保 skipRows 不超过文件的行数
	if skipRows > len(lines) {
		skipRows = len(lines)
	}

	// 将内容插入到指定位置
	newLines := append(lines[:skipRows], append(strings.Split(content, "\n"), lines[skipRows:]...)...)

	// 将新的内容写回到文件中
	output := strings.Join(newLines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}
