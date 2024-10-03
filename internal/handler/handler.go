package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ankodd/todo-server/internal/storage"
	"github.com/ankodd/todo-server/pkg/models/todo"
	"github.com/ankodd/todo-server/pkg/utils/write"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Storage     *storage.Storage
	IdleTimeout time.Duration
}

func (h *Handler) Create(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handler.Create"
		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		var t todo.Todo

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			write.Write(err.Error(), w, http.StatusBadRequest)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		err = h.Storage.Insert(ctx, &t)
		if err != nil {
			write.Write(err.Error(), w, http.StatusInternalServerError)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		select {
		case <-ctx.Done():
			write.Write(fmt.Sprintf("Error: %s", ctx.Err()), w, http.StatusRequestTimeout)
			log.Printf("%s: %s\n", fn, ctx.Err())
			return
		default:
			w.WriteHeader(http.StatusCreated)
			log.Printf("%s\n: Created", fn)
		}
	}
}

func (h *Handler) FetchAll(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handler.FetchAll"
		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		todos, err := h.Storage.FetchAll(ctx)
		if err != nil {
			write.Write(err.Error(), w, http.StatusInternalServerError)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		select {
		case <-ctx.Done():
			write.Write(fmt.Sprintf("Error: %s", ctx.Err()), w, http.StatusRequestTimeout)
			log.Printf("%s: %s\n", fn, ctx.Err())
			return
		default:
			write.Write(todos, w, http.StatusOK)
			log.Printf("%s\n: Fetched", fn)
		}
	}
}

func (h *Handler) Update(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handler.Update"
		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		ID := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(ID, 10, 64)
		if err != nil {
			write.Write(err.Error(), w, http.StatusBadRequest)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		var t todo.Todo

		err = json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			write.Write(err.Error(), w, http.StatusBadRequest)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		err = h.Storage.Update(ctx, &t, id)
		if err != nil {
			write.Write(err.Error(), w, http.StatusInternalServerError)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		select {
		case <-ctx.Done():
			write.Write(fmt.Sprintf("Error: %s", ctx.Err()), w, http.StatusRequestTimeout)
			log.Printf("%s: %s\n", fn, ctx.Err())
			return
		default:
			w.WriteHeader(http.StatusOK)
			log.Printf("%s: Updated", fn)
		}
	}
}

func (h *Handler) Delete(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handler.Delete"
		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		ID := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(ID, 10, 64)
		if err != nil {
			write.Write(err.Error(), w, http.StatusBadRequest)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		err = h.Storage.Delete(ctx, id)
		if err != nil {
			write.Write(err.Error(), w, http.StatusInternalServerError)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		select {
		case <-ctx.Done():
			write.Write(fmt.Sprintf("Error: %s", ctx.Err()), w, http.StatusRequestTimeout)
			log.Printf("%s: %s\n", fn, ctx.Err())
			return
		default:
			w.WriteHeader(http.StatusOK)
			log.Printf("%s\n: Deleted", fn)
		}
	}
}

func (h *Handler) CountEntries(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handler.CountEntries"
		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		cnt, err := h.Storage.CountEntries(ctx)
		if err != nil {
			write.Write(err.Error(), w, http.StatusInternalServerError)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		select {
		case <-ctx.Done():
			write.Write(fmt.Sprintf("Error: %s", ctx.Err()), w, http.StatusRequestTimeout)
			log.Printf("%s: %s\n", fn, ctx.Err())
			return
		default:
			write.Write(fmt.Sprintf("Count: %d", cnt), w, http.StatusOK)
			log.Printf("%s: Counted", fn)
		}
	}
}
