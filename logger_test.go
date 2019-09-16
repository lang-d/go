package golog_test

import (
	log "golog"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	log.Infof("%s,%d", "445", 4)
	log.Debugf("%s,%d", "445", 4)
	log.Warnf("%s,%d", "445", 4)
	log.Errorf("%s,%d", "445", 4)
	log.Fatalf("%s,%d", "445", 4)
}

func TestFormatter(t *testing.T) {
	log.SetHandler(log.NewConsleHandler().
		SetFormatter(log.NewDefaultFormatter().
			SetTemplate(log.ASCTIME_MARK + " -- " + log.LEVELNAME_MARK + " -- " + log.MSG_MARK).
			SetTimeLayout("2006-01-02")))
	log.Infof("%s,%d", "445", 4)
	log.Debugf("%s,%d", "445", 4)
	log.Warnf("%s,%d", "445", 4)
	log.Errorf("%s,%d", "445", 4)
	log.Fatalf("%s,%d", "445", 4)
}

func TestFileHanlder(t *testing.T) {
	log.SetHandler(log.NewFileHandler("log/run.log").
		SetFormatter(log.NewDefaultFormatter()))
	for i := 0; i < 10000; i++ {
		log.Infof("%s,%d", "445", 4)
		log.Debugf("%s,%d", "445", 4)
		log.Warnf("%s,%d", "445", 4)
		log.Errorf("%s,%d", "445", 4)
	}

}

func TestFileHanlderSplit(t *testing.T) {
	log.SetHandler(log.NewFileHandler("log/run.log").
		SetFormatter(log.NewDefaultFormatter()).SetSplitBySize(true, 10*1024*1024))
	for i := 0; i < 1000000; i++ {
		log.Infof("%s,%d", "445", 4)
		log.Debugf("%s,%d", "445", 4)
		log.Warnf("%s,%d", "445", 4)
		log.Errorf("%s,%d", "445", 4)
	}

}

func TestFileHanlderSplitByTime(t *testing.T) {
	log.AddHandler(log.NewFileHandler("log/run.log").
		SetFormatter(log.NewDefaultFormatter()).SetSplitByTime(true, log.LOGFILE_SPLIT_BY_HOUR))
	for i := 0; i < 4000; i++ {
		log.Infof("%s,%d", "445", 4)
		log.Debugf("%s,%d", "445", 4)
		log.Warnf("%s,%d", "445", 4)
		log.Errorf("%s,%d", "445", 4)
		time.Sleep(time.Second * 1)
	}

}
