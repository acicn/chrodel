package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func exit(err *error) {
	if *err != nil {
		log.Println("exited with error:", (*err).Error())
		os.Exit(1)
	} else {
		log.Println("exited")
	}
}

var (
	optDir    string
	optMatch  string
	optLayout string
	optKeep   int
	optDry    bool
)

const (
	SubExpNameDate = "date"
)

func main() {
	var err error
	defer exit(&err)

	flag.StringVar(&optDir, "dir", "", "日志文件目录")
	flag.StringVar(&optMatch, "match", "", "日志文件名匹配，其中必须包含 date 子匹配名，且需要和 --pattern 参数匹配")
	flag.StringVar(&optLayout, "layout", "", "日期格式，参考 Go 'time' 包")
	flag.IntVar(&optKeep, "keep", 0, "需要保留的日志天数")
	flag.BoolVar(&optDry, "dry", false, "调试开关，并不真的要删除日志文件")
	flag.Parse()

	optDir = strings.TrimSpace(optDir)

	optMatch = strings.TrimSpace(optMatch)

	if optMatch == "" {
		err = errors.New("缺失参数 --match")
		return
	}

	var match *regexp.Regexp
	if match, err = regexp.Compile(optMatch); err != nil {
		return
	}

	matchDateIndex := -1

	for i, name := range match.SubexpNames() {
		if name == SubExpNameDate {
			matchDateIndex = i
		}
	}

	if matchDateIndex < 0 {
		err = errors.New("--match 参数缺乏 date 子匹配，请参考 Go 'regexp/syntax' 包，使用 (?P<date>XXXXX) 定义 date 子匹配")
		return
	}

	if optKeep <= 0 {
		err = errors.New("缺失参数 --keep")
		return
	}

	now := time.Now()

	if err = filepath.Walk(optDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		subMatches := match.FindStringSubmatch(filepath.Base(path))
		if len(subMatches) == 0 {
			log.Println(path, "未匹配")
			return nil
		}

		var date time.Time
		if date, err = time.Parse(optLayout, subMatches[matchDateIndex]); err != nil {
			log.Println(path, "无法解析日期:", err.Error())
			return nil
		}

		if now.Sub(date) < time.Hour*24*time.Duration(optKeep) {
			log.Println(path, "无需删除")
			return nil
		}

		if optDry {
			log.Println(path, "即将删除")
		} else {
			if err = os.Remove(path); err != nil {
				log.Println(path, "删除失败:", err.Error())
			} else {
				log.Println(path, "已删除")
			}
		}

		return nil
	}); err != nil {
		return
	}
}
