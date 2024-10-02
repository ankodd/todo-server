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
		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		var t todo.Todo

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			write.Write(err.Error(), w, http.StatusBadRequest)
			log.Printf("handler.Create: %s\n", err.Error())
			return
		}

		err = h.Storage.Insert(ctx, &t)
		if err != nil {
			write.Write(err.Error(), w, http.StatusInternalServerError)
			log.Printf("handler.Create: %s\n", err.Error())
			return
		}

		select {
		case <-ctx.Done():
			write.Write(fmt.Sprintf("Error: %s", ctx.Err()), w, http.StatusRequestTimeout)
			log.Printf("handler.Create: %s\n", ctx.Err())
			return
		default:
			w.WriteHeader(http.StatusCreated)
			log.Println("handler.Create: Created")
		}
	}
}

func (h *Handler) FetchAll(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		todos, err := h.Storage.FetchAll(ctx)
		if err != nil {
			write.Write(err.Error(), w, http.StatusInternalServerError)
			log.Printf("handler.FetchAll: %s\n", err.Error())
			return
		}

		select {
		case <-ctx.Done():
			write.Write(fmt.Sprintf("Error: %s", ctx.Err()), w, http.StatusRequestTimeout)
			log.Printf("handler.FetchAll: %s\n", ctx.Err())
			return
		default:
			write.Write(todos, w, http.StatusOK)
			log.Println("handler.FetchAll: Fetched")
		}
	}
}

func (h *Handler) Update(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		ID := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(ID, 10, 64)
		if err != nil {
			write.Write(err.Error(), w, http.StatusBadRequest)
			log.Printf("handler.Update: %s\n", err.Error())
			return
		}

		var t todo.Todo

		err = json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			write.Write(err.Error(), w, http.StatusBadRequest)
			log.Printf("handler.Update: %s\n", err.Error())
			return
		}

		err = h.Storage.Update(ctx, &t, id)
		if err != nil {
			write.Write(err.Error(), w, http.StatusInternalServerError)
			log.Printf("handler.Update: %s\n", err.Error())
			return
		}

		select {
		case <-ctx.Done():
			write.Write(fmt.Sprintf("Error: %s", ctx.Err()), w, http.StatusRequestTimeout)
			log.Printf("handler.Update: %s\n", ctx.Err())
			return
		default:
			w.WriteHeader(http.StatusOK)
			log.Println("handler.Update: Updated")
		}
	}
}

func (h *Handler) Delete(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		ID := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(ID, 10, 64)
		if err != nil {
			write.Write(err.Error(), w, http.StatusBadRequest)
			log.Printf("handler.Delete: %s\n", err.Error())
			return
		}

		err = h.Storage.Delete(ctx, id)
		if err != nil {
			write.Write(err.Error(), w, http.StatusInternalServerError)
			log.Printf("handler.Delete: %s\n", err.Error())
			return
		}

		select {
		case <-ctx.Done():
			write.Write(fmt.Sprintf("Error: %s", ctx.Err()), w, http.StatusRequestTimeout)
			log.Printf("handler.Delete: %s\n", ctx.Err())
			return
		default:
			w.WriteHeader(http.StatusOK)
			log.Println("handler.Delete: Deleted")
		}
	}
}

func (h *Handler) CountEntries(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		cnt, err := h.Storage.CountEntries(ctx)
		if err != nil {
			write.Write(err.Error(), w, http.StatusInternalServerError)
			log.Printf("handler.CountEntries: %s\n", err.Error())
			return
		}

		select {
		case <-ctx.Done():
			write.Write(fmt.Sprintf("Error: %s", ctx.Err()), w, http.StatusRequestTimeout)
			log.Printf("handler.CountEntries: %s\n", ctx.Err())
			return
		default:
			write.Write(fmt.Sprintf("Count: %d", cnt), w, http.StatusOK)
			log.Println("handler.CountEntries: Counted")
		}

	}
}
