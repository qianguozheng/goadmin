package main

import (
	"net/http"
	"net/http/pprof"
)

// Pprof start http pprof.
func Pprof(addr string) {
	ps := http.NewServeMux()
	ps.HandleFunc("/debug/pprof/", pprof.Index)
	ps.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	ps.HandleFunc("/debug/pprof/profile", pprof.Profile)
	ps.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	go func() {
		if err := http.ListenAndServe(addr, ps); err != nil {
			panic(err)
		}
	}()
}
