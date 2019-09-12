package golog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	DEFAULT_LOGFILE_SPLIT_SIZE = 1024 * 1024 * 100
	DEFAULT_LOGFILE_SPLIT_TYPE = "day"
	LOGFILE_SPLIT_BY_HOUR      = "hour"
	DEFAULT_CHECK_LOGFILE_TIME = time.Second * 1

	DEFAULT_LOGFILE_BACKUP = 10

	HOUR_TIME_LAYOUT = "2006-01-02_15"
	DAY_TIME_LAYOUT  = "2006-01-02"
)

type Handler interface {
	Work(*level, string, ...interface{})
}

type ConsoleHandler struct {
	writer      io.Writer
	errorWriter io.Writer
	Formatter   Formatter

	bufferPool sync.Pool
}

func NewConsleHandler() *ConsoleHandler {
	return &ConsoleHandler{writer: os.Stdout, errorWriter: os.Stderr, Formatter: NewDefaultFormatter()}
}

func (this *ConsoleHandler) SetFormatter(formatter Formatter) *ConsoleHandler {
	this.Formatter = formatter
	return this
}

func (this *ConsoleHandler) Work(level *level, format string, args ...interface{}) {
	b, ok := this.bufferPool.Get().(*bytes.Buffer)
	if !ok {
		b = &bytes.Buffer{}
	}
	msg := this.Formatter.Format(b, level.string(), format, args...)
	if level.isEnable(ERROR) {
		_, _ = this.errorWriter.Write(msg)
	} else {
		_, _ = this.writer.Write(msg)
	}
	b.Reset()
	this.bufferPool.Put(b)

}

type FileHandler struct {
	Writer           *os.File
	ErrorWriter      *os.File
	Formatter        Formatter
	CheckLogFileTime time.Duration
	splitBySize      bool
	splitByTime      bool
	splitSize        int64
	splitTimeType    string
	FileName         string

	logfileBackup int

	isStartSplit bool

	bufferPool sync.Pool

	lock sync.Mutex
}

func NewFileHandler(fileName string) *FileHandler {
	fileWriter, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return &FileHandler{
		FileName:         fileName,
		Writer:           fileWriter,
		ErrorWriter:      fileWriter,
		Formatter:        NewDefaultFormatter(),
		CheckLogFileTime: DEFAULT_CHECK_LOGFILE_TIME,
	}
}

func (this *FileHandler) createLogFile() *FileHandler {
	fileWriter, err := os.OpenFile(this.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	this.Writer = fileWriter
	return this
}

func (this *FileHandler) SetLogFileBackup(backup int) *FileHandler {
	this.logfileBackup = backup
	return this
}

func (this *FileHandler) SetFormatter(formatter Formatter) *FileHandler {
	this.Formatter = formatter
	return this
}

func (this *FileHandler) SetFileName(fileName string) *FileHandler {
	this.FileName = fileName
	return this
}
func (this *FileHandler) SetCheckLogFileTime(checkTime time.Duration) *FileHandler {
	this.CheckLogFileTime = checkTime
	return this
}

func (this *FileHandler) SetSplitBySize(splitBySize bool, size int64) *FileHandler {
	this.splitByTime = false
	this.splitBySize = splitBySize
	this.splitSize = size
	if this.splitSize <= 0 {
		panic(errors.New(fmt.Sprintf("bad split size: %d", size)))
	}

	if this.logfileBackup == 0 {
		this.logfileBackup = DEFAULT_LOGFILE_BACKUP
	}

	if this.splitBySize && !this.isStartSplit {
		go this.split()
	}
	return this
}

func (this *FileHandler) SetSplitByTime(splitByTime bool, splitTimeType string) *FileHandler {
	this.splitByTime = splitByTime
	this.splitTimeType = splitTimeType

	if this.splitTimeType != DEFAULT_LOGFILE_SPLIT_TYPE && this.splitTimeType != LOGFILE_SPLIT_BY_HOUR {
		panic(errors.New(fmt.Sprintf("bad split type: %s", splitTimeType)))
	}

	if this.logfileBackup == 0 {
		this.logfileBackup = DEFAULT_LOGFILE_BACKUP
	}

	if this.splitByTime && !this.isStartSplit {
		go this.split()
	}
	return this
}

func (this *FileHandler) Work(level *level, format string, args ...interface{}) {
	b, ok := this.bufferPool.Get().(*bytes.Buffer)
	if !ok {
		b = &bytes.Buffer{}
	}
	msg := this.Formatter.Format(b, level.string(), format, args...)
	if level.isEnable(ERROR) {
		_, _ = this.ErrorWriter.Write(msg)
	} else {
		_, _ = this.Writer.Write(msg)
	}
	b.Reset()
	this.bufferPool.Put(b)

}

func (this *FileHandler) doSplitBySize() {
	fileInfo, err := os.Stat(this.FileName)
	if err != nil {
		panic(err)
	}
	fmt.Println(this.FileName)
	fmt.Println(fileInfo.Size())
	if fileInfo.Size() >= this.splitSize {
		this.lock.Lock()
		defer this.lock.Unlock()

		files, err := filepath.Glob(this.FileName + ".*")
		if err != nil {
			panic(err)
		}

		_ = this.Writer.Close()

		if len(files) == 0 {
			err = os.Rename(this.FileName, this.FileName+".1")
			if err != nil {
				panic(err)
			}
		} else {
			sort.Strings(files)
			filesNum := len(files)

			for i, _ := range files {

				file := files[filesNum-i-1]

				_ = os.Rename(file, strings.Replace(file, fmt.Sprintf(".%d", i+1), fmt.Sprintf(".%d", i+2), -1))
			}

			err = os.Rename(this.FileName, this.FileName+".1")
			if err != nil {
				panic(err)
			}

			files = append(files, this.FileName+".1")
			sort.Strings(files)

			if len(files) > this.logfileBackup {
				delFiles := files[this.logfileBackup:]
				for _, f := range delFiles {
					_ = os.Remove(f)
				}
			}

		}
		this.createLogFile()
	}
}

func (this *FileHandler) doSplitBytime() {
	now := time.Now()
	if this.splitTimeType == LOGFILE_SPLIT_BY_HOUR {
		lastHour := now.Add(-time.Hour)
		fileName := fmt.Sprintf("%s.%s", this.FileName, lastHour.Format(HOUR_TIME_LAYOUT))
		_, err := os.Stat(fileName)
		if os.IsNotExist(err) {
			this.lock.Lock()
			defer this.lock.Unlock()

			_ = this.Writer.Close()

			_ = os.Rename(this.FileName, fileName)
			// todo
			this.createLogFile()
		}
	} else if this.splitTimeType == DEFAULT_LOGFILE_SPLIT_TYPE {
		if now.Hour() == 0 {
			lastDay := now.Add(-time.Hour * 24)
			fileName := fmt.Sprintf("%s.%s", this.FileName, lastDay.Format(DAY_TIME_LAYOUT))
			_, err := os.Stat(fileName)
			if os.IsNotExist(err) {
				this.lock.Lock()
				defer this.lock.Unlock()

				_ = this.Writer.Close()

				_ = os.Rename(this.FileName, fileName)

				this.createLogFile()
			}
		}
	}

	files, err := filepath.Glob(this.FileName + ".*")
	if err != nil {
		panic(err)
	}

	if this.logfileBackup > len(files) {
		sort.Strings(files)
		_files := files[0 : this.logfileBackup-len(files)]
		for _, file := range _files {
			_ = os.Remove(file)
		}
	}

}

func (this *FileHandler) split() {
	this.isStartSplit = true
	for {
		if this.splitBySize {
			this.doSplitBySize()
		} else if this.splitByTime {
			this.doSplitBytime()
		} else {
			this.isStartSplit = false
			break
		}
		time.Sleep(this.CheckLogFileTime)
	}
}
