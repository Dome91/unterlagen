package httpx

import "net/http"

func Redirect(writer http.ResponseWriter, url string, code int) {
	writer.Header().Set("HX-Location", url)
	writer.WriteHeader(code)
}
