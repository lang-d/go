package golog_test

import (
	log "golog"
	"testing"
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
			SetTemplate("%(asctime)s -- [%(levelname)s] -- %(message)s").
			SetTimeLayout("2006")))
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
		SetFormatter(log.NewDefaultFormatter()).SetSplitBySize(true, 10*1024))
	for i := 0; i < 10000000; i++ {
		log.Infof("%s,%d", "445", 4)
		log.Debugf("%s,%d", "445", 4)
		log.Warnf("%s,%d", "445", 4)
		log.Errorf("%s,%d", "445", 4)
	}

}
