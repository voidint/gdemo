package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

// Options 选项
type Options struct {
	All  bool
	List bool
}

var opts Options

func init() {
	flag.BoolVar(&opts.All, "a", false, "Include directory entries whose names begin with a dot (.).")
	flag.BoolVar(&opts.List, "l", false, "List in long format. If the output is to a terminal, a total sum for all the file sizes is output on a line before the long listing.")
	flag.Parse()
}

func main() {
	var target string
	if length := len(os.Args); length > 1 && !strings.HasPrefix(os.Args[length-1], "-") {
		target = os.Args[length-1]
	} else {
		target, _ = os.Getwd()
	}

	f, err := os.Open(target)
	handleError(err)

	defer f.Close()

	items, err := getFileInfos(f)
	handleError(err)

	sort.Slice(items, func(i, j int) bool {
		for n := 0; n <= i && n <= j; n++ {
			return items[i].Name()[n] < items[j].Name()[n]
		}
		return i < j
	})

	if opts.List {
		handleError(NewListRender(target, items).Render(opts.All, os.Stdout))
		return
	}
	handleError(NewRawRender(items).Render(opts.All, os.Stdout))
}

// Renderer 文本渲染器
type Renderer interface {
	Render(all bool, out io.Writer) error
}

// ListRender 列表渲染器
type ListRender struct {
	dir   string
	items []os.FileInfo
}

// NewListRender 实例化列表渲染器
func NewListRender(dir string, items []os.FileInfo) Renderer {
	return &ListRender{
		dir:   dir,
		items: items,
	}
}

// Render 渲染
func (rd *ListRender) Render(all bool, out io.Writer) error {
	var total int64
	var buf strings.Builder
	for _, item := range rd.items {
		if !all && strings.HasPrefix(item.Name(), ".") {
			continue
		}

		var uname, gname string
		var nlink uint16

		if stat, ok := item.Sys().(*syscall.Stat_t); ok && stat != nil {
			if user, _ := user.LookupId(strconv.Itoa(int(stat.Uid))); user != nil {
				uname = user.Name
			}
			if group, _ := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid)); group != nil {
				gname = group.Name
			}
			nlink = stat.Nlink
			total += stat.Blocks
		}

		fname := item.Name()
		if item.Mode()&os.ModeSymlink != 0 {
			dest, _ := os.Readlink(filepath.Join(rd.dir, fname))
			fname = fmt.Sprintf("%s -> %s", fname, dest)
		}

		buf.WriteString(fmt.Sprintf("%s\t%d\t%s\t%s\t%d\t%d\t%d\t%s\t%s\n",
			item.Mode().String(),
			nlink,
			uname,
			gname,
			item.Size(),
			item.ModTime().Month(),
			item.ModTime().Day(),
			item.ModTime().Format("15:04"),
			fname,
		))
	}
	fmt.Fprintln(out, fmt.Sprintf("total %d", total))
	fmt.Fprint(out, buf.String())
	return nil
}

// RawRender 原样渲染器
type RawRender struct {
	items []os.FileInfo
}

// NewRawRender 实例化原样渲染器
func NewRawRender(items []os.FileInfo) Renderer {
	return &RawRender{
		items: items,
	}
}

// Render 渲染
func (rd *RawRender) Render(all bool, out io.Writer) error {
	for _, item := range rd.items {
		if !all && strings.HasPrefix(item.Name(), ".") {
			continue
		}
		fmt.Fprintln(out, item.Name())
	}
	return nil
}

func getFileInfos(f *os.File) (items []os.FileInfo, err error) {
	fInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !fInfo.IsDir() {
		return []os.FileInfo{fInfo}, nil
	}

	return f.Readdir(-1)
}

func handleError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
