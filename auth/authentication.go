package auth

import (
	"net/http"
)

func AuthenticateRequest(r *http.Request, authMethod string) (bool, string) {

	switch authMethod {
	case "header": //authenticate request from header parameters
		return AuthenticateFromHeader(r)

	default: //invalid request
		return false, ""
	}
}
