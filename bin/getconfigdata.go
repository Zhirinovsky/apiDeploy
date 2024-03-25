package bin

import (
	"bufio"
	"os"
)

func GetConfigData() map[string]string {
	data := make(map[string]string)
	var key, value string
	file, err := os.Open("config.txt")
	CheckErr(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		CheckErr(scanner.Err())
		subdata := scanner.Text()
		for index, char := range subdata {
			if string(char) == ":" {
				value = subdata[index+2:]
				break
			} else {
				key += string(char)
			}
		}
		data[key] = value
		key = ""
	}
	return data
}
