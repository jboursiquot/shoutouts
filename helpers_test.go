package shoutouts_test

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func nullLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(ioutil.Discard)
	return log
}
