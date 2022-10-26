package http

import (
	"context"
	"io"
	"net/http"

	"github.com/utkarsh-pro/use/pkg/config"
	"github.com/utkarsh-pro/use/pkg/storage"
	"github.com/utkarsh-pro/use/pkg/storage/errors"
	"github.com/utkarsh-pro/use/pkg/utils"
)

type Transport struct {
	srv     *http.Server
	storage storage.Storage
}

// New returns a new HTTP transport
func New(storage storage.Storage) *Transport {
	return &Transport{
		storage: storage,
	}
}

func (t *Transport) Setup(addr string) error {
	t.srv = &http.Server{
		Addr: addr,
	}

	// Setup all the handlers.
	t.setupHandlers()

	if err := t.srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}

		return err
	}

	return nil
}

func (t *Transport) Shutdown() error {
	return t.srv.Shutdown(context.TODO())
}

func (t *Transport) setupHandlers() {
	http.HandleFunc("/version", createHTTPMethodsHandler([]string{http.MethodGet}, t.versionHandler))
	http.HandleFunc("/api/set", createHTTPMethodsHandler([]string{http.MethodGet}, t.setHandler))
	http.HandleFunc("/api/get", createHTTPMethodsHandler([]string{http.MethodGet}, t.getHandler))
	http.HandleFunc("/api/delete", createHTTPMethodsHandler([]string{http.MethodGet}, t.deleteHandler))
	http.HandleFunc("/api/len", createHTTPMethodsHandler([]string{http.MethodGet}, t.lenHandler))
	http.HandleFunc("/api/exists", createHTTPMethodsHandler([]string{http.MethodGet}, t.existsHandler))
	http.HandleFunc("/api/snapshot", createHTTPMethodsHandler([]string{http.MethodGet}, t.snapshotHandler))
}

func (t *Transport) setHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	val := r.URL.Query().Get("val")

	if err := t.storage.Set(key, []byte(val)); err != nil {
		if err == errors.ErrReadOnlyStorage {
			w.WriteHeader(http.StatusTeapot)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *Transport) getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	val, err := t.storage.Get(key)
	if err != nil {
		if err == errors.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(val))
}

func (t *Transport) deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if err := t.storage.Delete(key); err != nil {
		if err == errors.ErrReadOnlyStorage {
			w.WriteHeader(http.StatusTeapot)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *Transport) lenHandler(w http.ResponseWriter, r *http.Request) {
	len, err := t.storage.Len()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, utils.IntToString(len))
}

func (t *Transport) existsHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	exists, err := t.storage.Exists(key)
	if err != nil {
		if err == errors.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, utils.BoolToString(exists))
}

func (t *Transport) snapshotHandler(w http.ResponseWriter, r *http.Request) {
	err := t.storage.PhysicalSnapshot(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}
}

func (t *Transport) versionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, config.Version)
}

func createHTTPMethodsHandler(method []string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !utils.SliceContains(
			method,
			r.Method,
			func(v1, v2 string) bool { return v1 == v2 },
		) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		handler(w, r)
	}
}
