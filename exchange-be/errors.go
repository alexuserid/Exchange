package main

import (
	"errors"
	"net/http"
)

type errorc struct {
	Code int
	Text string
	Err error
}

func fullError(errc *errorc, err error) *errorc {
	e := errc
	e.Err = err
	return &e
}

var (
	errLogin              = &errorc{http.StatusBadRequest, "You are not logged in.", nil}
	errWithdraw           = &errorc{http.StatusBadRequest, "Can't withdraw more money than you have.", nil}
	errExistingEmail      = &errorc{http.StatusBadRequest, "The email is asready exist.", nil}
	errWrongEmailPassword = &errorc{http.StatusBadRequest, "Wrong email or password.", nil}
	errNoCookie           = &errorc{http.StatusBadRequest, "You are not logged in.", nil}
	errParseFloat         = &errorc{http.StatusBadRequest, "Can't parse float number.", nil}
	errNoToken            = &errorc{http.StatusInternalServerError, "No free tokens in getUniqueId.", errors.New("get token: no free tokens after 100 iteration")}
	errLength             = &errorc{http.StatusInternalServerError, "Mismatched length in getRandoms.", errors.New("getRandom32: rand.Read error")}
	errGetRandom          = &errorc{http.StatusInternalServerError, "Get randoms error.", nil}
	errHashGen            = &errorc{http.StatusInternalServerError, "Generate hash from password error.", nil}
)
