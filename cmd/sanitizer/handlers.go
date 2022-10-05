package main

import (
	"io"
	"net/http"

	goaway "github.com/TwiN/go-away"
	"github.com/sirupsen/logrus"
)

type healthHandler struct {
	log *logrus.Logger
}

func (hh healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hh.log.Traceln("healthHandler.ServeHTTP")
	w.WriteHeader(http.StatusOK)
}

type sanitizeHandler struct {
	log *logrus.Logger
}

func (ch sanitizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ch.log.Infoln("checkHandler.ServeHTTP")
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		ch.log.WithError(err).Errorln("failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if goaway.IsProfane(string(bs)) {
		ch.log.WithField("text", string(bs)).Warnln("profanity detected")
	}

	body := goaway.Censor(string(bs))

	if _, err := w.Write([]byte(body)); err != nil {
		ch.log.WithError(err).Errorln("failed to write response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
