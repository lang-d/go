### 默认日志输出格式

```text
[2019-09-16 11:16:34.638][info] 一些信息
[2019-09-16 11:16:34.638][warn] 一些警告信息
[2019-09-16 11:16:34.638][error] 一些错误信息
```
### 使用方法
#### 使用默认的设置打印日志

```go
    import log "github.com/lang-d/golog" 
	log.Infof("%s,%d", "445", 4)
    	log.Debugf("%s,%d", "445", 4)
    	log.Warnf("%s,%d", "445", 4)
    	log.Errorf("%s,%d", "445", 4)
    	log.Fatalf("%s,%d", "445", 4)
```

#### 更改日志的输出格式
这里可以设置一个新的handler用来替换默认的handler,默认的handler有两个,一个FileHandler,一个ConsoHandler,区别在于FileHandler将日志输出到文件,ConsoleHandler将日志输出到控制台

```go
log.SetHandler(log.NewConsleHandler().
		SetFormatter(log.NewDefaultFormatter().
			SetTemplate(log.ASCTIME_MARK + " -- " + log.LEVELNAME_MARK + " -- " + log.MSG_MARK).
			SetTimeLayout("2006-01-02")))
```
可以通过添加一个FileHandler的方式将日志记录到文件,FileHandler 支持对日志进行按时间分割和按大小分割,默认是不会对日志进行分割的

```go
log.AddHandler(log.NewFileHandler("log/run.log").
		SetFormatter(log.NewDefaultFormatter()).SetSplitByTime(true, log.LOGFILE_SPLIT_BY_HOUR))
```

如果提供的handler不能满足需求,你可以自己实现一个