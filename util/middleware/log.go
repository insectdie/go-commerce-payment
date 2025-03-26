package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"time"
)

func LoggerMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			recorder := httptest.NewRecorder()
			next.ServeHTTP(recorder, r)

			for k, v := range recorder.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(recorder.Code)
			recorder.Body.WriteTo(w)

			responseTime := time.Since(start).Seconds()
			logMessage := fmt.Sprintf("%s - [%s] - \"%s %s %s\" %d %s - [%s]\n",
				r.RemoteAddr,
				time.Now().Format(time.RFC1123),
				r.Method,
				r.URL.Path,
				r.Proto,
				recorder.Code,
				r.UserAgent(),
				fmt.Sprintf("%.9fµs", responseTime),
			)
			log.Print(logMessage)
		})
	}
}

func ApplyMiddleware(h http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	handler := http.Handler(h)
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler.ServeHTTP
}

func IsSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

func DebugOutput(i interface{}) {
	var mapRes map[string]interface{}
	switch iData := i.(type) {
	case string:
		log.Println(iData)
		return
	case []byte:
		json.Unmarshal(iData, &mapRes)
		s, _ := json.MarshalIndent(mapRes, "", "\t")
		fmt.Print(string(s))
	default:
		switch IsSlice(iData) {
		case true:
			s, _ := json.MarshalIndent(i, "", "\t")
			fmt.Print(string(s))
		default:
			s, _ := json.MarshalIndent(i, "", "\t")
			fmt.Print(string(s))
		}
	}
}
