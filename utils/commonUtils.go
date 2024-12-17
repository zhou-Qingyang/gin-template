package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"
)

func ParseDate(date string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", date, time.Local)
}

func HasContainInSlice(target string, slice []string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func GetPagination(page, pageSize int) (int, int, error) {
	if page <= 0 {
		return 0, 0, errors.New("page must be greater than 0")
	}
	if pageSize <= 0 {
		return 0, 0, errors.New("pageSize must be greater than 0")
	}
	offset := (page - 1) * pageSize
	return offset, pageSize, nil
}

func HasContainInSliceInt64(target int64, slice []int64) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func Md5(str string) string {
	// 计算 MD5 哈希值
	hash := md5.New()
	hash.Write([]byte(str))
	md5Hash := hash.Sum(nil)
	// 将 MD5 值转为十六进制字符串
	return hex.EncodeToString(md5Hash)
}
