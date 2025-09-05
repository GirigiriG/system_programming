package utils

import (
	"fmt"
	"syscall"
)

type FileStatInfo struct {
	NumberOfLines int
	NumberOfWords int
	NumberOfBytes int
}

func openFile(filePath string) (int, error) {
	fd, err := syscall.Open(filePath, syscall.O_RDWR|syscall.O_APPEND, 0644)
	if err != nil {
		fmt.Println("unable to open file", err.Error())
		return 0, err
	}
	return fd, nil
}

func Write(filePath string, p string) int {
	fd, err := openFile(filePath)
	fileOperationError(err)
	_, err = syscall.Write(fd, []byte(p))
	if err != nil {
		fmt.Println("error writing to file ", err.Error())
		return fd
	}
	return fd
}

func Read(filePath string) (string, error) {
	fd, err := openFile(filePath)
	fileOperationError(err)

	defer syscall.Close(fd)

	var stats syscall.Stat_t
	err = syscall.Fstat(fd, &stats)
	if err != nil {
		return "", err
	}

	buff := make([]byte, stats.Size)

	_, err = syscall.Read(fd, buff)
	if err != nil {
		fmt.Println("error reading file", err)
	}
	return string(buff), nil
}

func FileStats(filePath string) FileStatInfo {
	fd, err := openFile(filePath)
	fileOperationError(err)

	defer syscall.Close(fd)

	var stats syscall.Stat_t
	_ = syscall.Fstat(fd, &stats)

	sizOfBytes := int(stats.Size)
	var numberOfLines int

	buff := make([]byte, sizOfBytes)

	syscall.Read(fd, buff)

	var numberOfWords int
	var currWord string

	for i, v := range buff {
		w := string(v)
		currWord += w

		if w == " " || w == "\n" || i == sizOfBytes-1 && len(currWord) > 1 {
			numberOfWords++
			currWord = ""
		}
	}

	return FileStatInfo{
		NumberOfLines: numberOfLines,
		NumberOfBytes: sizOfBytes,
		NumberOfWords: numberOfWords,
	}
}

func fileOperationError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
