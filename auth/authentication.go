package auth

import (
	"encoding/json"
	"net/http"
)

func CheckAuth(r *http.Request) (code int, s string) {
	if isInternalRequest() {
		return 200, ""
	} else {
		return GetSessionFromHeader(r)
	}
}

func isInternalRequest() (ok bool) {

	return false
}

func GetAuthenticationError(errCode int) (out []byte) {
	r := make(map[string]interface{})
	if errCode != 200 {
		e := make(map[string]interface{})
		e["error_code"] = errCode
		e["error_message"] = getErrorMessage(errCode)

		r["status"] = false
		r["response_code"] = 601 // error code for invalid authentication params
		r["response_message"] = getErrorMessage(601)
		r["error"] = e //specific error code

	}
	out, _ = json.Marshal(r)
	return out
}

func getErrorMessage(errCode int) (errMsg string) {
	switch errCode {
	case 601:
		errMsg = "Invalid Request"
	case 602:
		errMsg = "Invalid Timestamp"
	case 603:
		errMsg = "Invalid Hash"
	case 605:
		errMsg = "Invalid App Session"
	default:
		errMsg = ""
	}
	return errMsg
}
