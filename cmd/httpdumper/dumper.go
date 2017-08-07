package main

import (
	"errors"
	"time"

	"github.com/astaxie/beego/httplib"
)

// HTTPDumper HTTP协议文件转储器
type HTTPDumper struct {
	url         string
	method      string
	header      map[string]string
	checksum    ChecksumFunc
	hexHash     string
	connTimeout time.Duration
	rwTimeout   time.Duration
	retries     int
}

// NewHTTPDumper 实例化http转储器
func NewHTTPDumper(url, method string, header map[string]string) *HTTPDumper {
	return &HTTPDumper{
		url:    url,
		method: method,
		header: header,
	}
}

// Checksum 设置checksum的函数实现
func (d *HTTPDumper) Checksum(hexHash string, checksum ChecksumFunc) *HTTPDumper {
	d.hexHash = hexHash
	d.checksum = checksum
	return d
}

// Timeout 设置连接超时时间和读写超时时间
func (d *HTTPDumper) Timeout(connTimeout, readWriteTimeout time.Duration) *HTTPDumper {
	d.connTimeout = connTimeout
	d.rwTimeout = readWriteTimeout
	return d
}

// Retries 设置重试次数
func (d *HTTPDumper) Retries(retries int) *HTTPDumper {
	d.retries = retries
	return d
}

var (
	// ErrChecksum 无效的checksum错误
	ErrChecksum = errors.New("invalid checksum")
)

// Dump 下载远程文件并存储到本地指定位置
func (d *HTTPDumper) Dump(dstFilename string) (err error) {
	req := httplib.NewBeegoRequest(d.url, d.method)
	for k, v := range d.header {
		req.Header(k, v)
	}

	if d.connTimeout > 0 && d.rwTimeout > 0 {
		req.SetTimeout(d.connTimeout, d.rwTimeout)
	}

	if err = req.Retries(d.retries).ToFile(dstFilename); err != nil {
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
