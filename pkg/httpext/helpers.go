package httpext

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
)

func Abort(w http.ResponseWriter, err error, code int) {
	response := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(response)
}

func AbortWithoutContent(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func AbortPlain(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	String(w, err.Error())
}

func AbortJSON(w http.ResponseWriter, body interface{}, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func JSON(w http.ResponseWriter, body interface{}) {
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}

func XML(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(xml.Header))
	_ = xml.NewEncoder(w).Encode(body)
}

func Data(w http.ResponseWriter, data []byte) {
	w.Header().Add("Content-Type", http.DetectContentType(data))
	_, _ = w.Write(data)
}

func String(w http.ResponseWriter, str string) {
	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write([]byte(str))
}

func SanitizeURL(uri string) string {
	cleanRegexp, _ := regexp.Compile("[\\d\\w\\-/]+")
	return cleanRegexp.FindString(uri)
}

func MakeAbsoluteURL(r *http.Request, target string, isSecured bool) string {
	if isSecured {
		return fmt.Sprintf("https://%s%s", r.Host, target)
	} else {
		return fmt.Sprintf("http://%s%s", r.Host, target)
	}
}
