# chrodel

chronological file deletion tool

## 用法用量

```
Usage of ./chrodel:
  -dir string
    	文件目录
  -dry
    	调试开关，并不真的要删除文件
  -keep int
    	需要保留的日志天数
  -layout string
    	日期格式，参考 Go 'time' 包
  -match string
    	文件名匹配，其中必须包含 date 子匹配名，且需要和 --pattern 参数匹配
```

## 许可证

Guo Y.K., MIT License
