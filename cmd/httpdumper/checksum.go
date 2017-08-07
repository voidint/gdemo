package main

import (
	stdmd5 "crypto/md5"
	"fmt"
	"io"
	"os"
)

// ChecksumFunc 目标文件摘要校验函数定义。
// dstFilename 待计算摘要的目标文件。
// hexHash 目标文件的预期摘要值。
// 若实际计算的摘要值与预期一致，返回true，否则为false。当摘要计算过程中发生错误，将返回非nil的error值。
type ChecksumFunc func(dstFilename, hexHash string) (ok bool, err error)

// MD5Checksum 使用MD5摘要算法计算的目标文件的摘要值是否与预期摘要值一致。
func MD5Checksum(dstFilename, hexHash string) (ok bool, err error) {
	md5, err := FileMD5(dstFilename)
	if err != nil {
		return false, err
	}
	return md5 == hexHash, nil
}

// FileMD5 计算文件的md5值
func FileMD5(filename string) (md5 string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := stdmd5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
