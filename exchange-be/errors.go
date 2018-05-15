package main

import (
	"net/http"
)

type errorc struct {
	Code int
	Text string
	Log bool
}

var (
errLogin = errorc{http.StatusBadRequest, "You are not logged in.", false}
errWithdraw = errorc{http.StatusBadRequest, "Can't withdraw more money than you have.", false}
errExistingEmail = errorc{http.StatusBadRequest, "The email is asready exist.", false}
errWrongEmailPassword = errorc{http.StatusBadRequest, "Wrong email or password.", false}
errNoCookie = errorc{http.StatusBadRequest, "No required cookie. You are not logged in.", false}
errParseFloat = errorc{http.StatusBadRequest, "Can't parse float number.", false}
errNoToken = errorc{http.StatusInternalServerError, "No free tokens in getUniqueId.", true}
errLength = errorc{http.StatusInternalServerError, "Mismatched length in getRandoms.", true}
errGetRandom = errorc{http.StatusInternalServerError, "Get randoms error.", true}
errHashGen = errorc{http.StatusInternalServerError, "Generate hash from password error.", true}

errNo = errorc{0, "", false}
)
