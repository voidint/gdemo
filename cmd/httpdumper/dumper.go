package main

import (
	"errors"
	"time"

	"github.com/astaxie/beego/httplib"
)

// HTTPDumper HTTP协议文件转储器
type HTTPDumper struct {
	url      string
	method   string
	header   map[string]string
	checksum ChecksumFunc
	hexHash  string
}

// NewHTTPDumper 实例化http转储器
func NewHTTPDumper(url, method string, header map[string]string) *HTTPDumper {
	return &HTTPDumper{
		url:    url,
		method: method,
		header: header,
	}
}

// ChecksumFunc 目标文件摘要校验函数定义。
// dstFilename 待计算摘要的目标文件。
// hexHash 目标文件的预期摘要值。
// 若实际计算的摘要值与预期一致，返回true，否则为false。当摘要计算过程中发生错误，将返回非nil的error值。
type ChecksumFunc func(dstFilename, hexHash string) (ok bool, err error)

// Checksum 设置checksum的函数实现
func (d *HTTPDumper) Checksum(hexHash string, checksum ChecksumFunc) *HTTPDumper {
	d.hexHash = hexHash
	d.checksum = checksum
	return d
}

var (
	// ErrChecksum 无效的checksum
	ErrChecksum = errors.New("invalid checksum")
)

// Dump 下载远程文件并存储到本地指定位置
func (d *HTTPDumper) Dump(dstFilename string) (err error) {
	req := httplib.NewBeegoRequest(d.url, d.method)
	for k, v := range d.header {
		req.Header(k, v)
	}
	if err = req.SetTimeout(time.Second, time.Second).Retries(3).ToFile(dstFilename); err != nil {
		return err
	}

	if d.checksum != nil {
		ok, err := d.checksum(dstFilename, d.hexHash)
		if err != nil {
			return err
		}
		if !ok {
			return ErrChecksum
		}
	}
	return nil
}
