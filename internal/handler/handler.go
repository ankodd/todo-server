package handler

import (
	"context"
	"encoding/json"
	"github.com/ankodd/todo-server/internal/metrics"
	"github.com/ankodd/todo-server/internal/storage/sqlite"
	"github.com/ankodd/todo-server/pkg/models/http/response"
	"github.com/ankodd/todo-server/pkg/models/todo"
	"github.com/ankodd/todo-server/pkg/utils/write"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Storage     *sqlite.Storage
	IdleTimeout time.Duration
	Response    response.Response
	Metrics     *metrics.Metrics
}

func (h *Handler) Create(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handler.Create"

		start := time.Now()

		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		h.setResponse(response.New(
			nil,
			http.StatusCreated,
			nil, w),
		)

		defer func() {
			h.Metrics.ObserveRequest(time.Since(start), h.Response.Status)
			if h.Response.Err != nil {
				h.Metrics.IncError()
			}

			h.Metrics.IncRequest()

			log.Printf("%s: Response: %+v\n", fn, h.Response.Log())
		}()

		var out todo.Todo

		if err := json.NewDecoder(r.Body).Decode(&out); err != nil {
			h.Response.MakeErr(err, http.StatusBadRequest)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		data := out

		if err := h.Storage.Insert(ctx, &data); err != nil {
			h.Response.MakeErr(err, http.StatusInternalServerError)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		select {
		case <-ctx.Done():
			h.Response.MakeErr(ctx.Err(), http.StatusRequestTimeout)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, ctx.Err())
			return
		default:
			write.Write(&h.Response)
			log.Printf("%s: Created\n", fn)
		}
	}
}

func (h *Handler) FetchAll(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handler.FetchAll"

		start := time.Now()

		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		h.setResponse(response.New(
			nil,
			http.StatusOK,
			nil, w),
		)

		defer func() {
			h.Metrics.ObserveRequest(time.Since(start), h.Response.Status)
			if h.Response.Err != nil {
				h.Metrics.IncError()
			}

			h.Metrics.IncRequest()

			log.Printf("%s: Response: %+v\n", fn, h.Response.Log())
		}()

		todos, err := h.Storage.FetchAll(ctx)
		h.Response.Data = &todos

		if err != nil {
			h.Response.MakeErr(err, http.StatusInternalServerError)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		select {
		case <-ctx.Done():
			h.Response.MakeErr(ctx.Err(), http.StatusRequestTimeout)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, ctx.Err())
			return
		default:
			write.Write(&h.Response)
			log.Printf("%s: Fetched\n", fn)
		}
	}
}

func (h *Handler) Update(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handler.Update"

		start := time.Now()

		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		h.setResponse(response.New(
			todo.Todo{},
			http.StatusOK,
			nil, w),
		)

		defer func() {
			h.Metrics.ObserveRequest(time.Since(start), h.Response.Status)
			if h.Response.Err != nil {
				h.Metrics.IncError()
			}

			h.Metrics.IncRequest()

			log.Printf("%s: Response: %+v\n", fn, h.Response.Log())
		}()

		ID := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(ID, 10, 64)
		if err != nil {
			h.Response.MakeErr(err, http.StatusBadRequest)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		err = json.NewDecoder(r.Body).Decode(&h.Response.Data)
		if err != nil {
			h.Response.MakeErr(err, http.StatusBadRequest)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		data := h.Response.Data.(todo.Todo)
		err = h.Storage.Update(ctx, &data, id)
		if err != nil {
			h.Response.MakeErr(err, http.StatusInternalServerError)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		select {
		case <-ctx.Done():
			h.Response.MakeErr(ctx.Err(), http.StatusRequestTimeout)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, ctx.Err())
			return
		default:
			w.WriteHeader(http.StatusOK)
			log.Printf("%s: Updated\n", fn)
		}
	}
}

func (h *Handler) Delete(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handler.Delete"

		start := time.Now()

		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		h.setResponse(response.New(
			nil,
			http.StatusOK,
			nil, w),
		)

		defer func() {
			h.Metrics.ObserveRequest(time.Since(start), h.Response.Status)
			if h.Response.Err != nil {
				h.Metrics.IncError()
			}

			h.Metrics.IncRequest()

			log.Printf("%s: Response: %+v\n", fn, h.Response.Log())
		}()

		ID := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(ID, 10, 64)
		if err != nil {
			h.Response.MakeErr(err, http.StatusBadRequest)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		err = h.Storage.Delete(ctx, id)
		if err != nil {
			h.Response.MakeErr(err, http.StatusInternalServerError)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		select {
		case <-ctx.Done():
			h.Response.MakeErr(ctx.Err(), http.StatusRequestTimeout)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, ctx.Err())
			return
		default:
			write.Write(&h.Response)
			log.Printf("%s: Deleted\n", fn)
		}
	}
}

func (h *Handler) CountEntries(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handler.CountEntries"

		start := time.Now()

		ctx, cancel := context.WithTimeout(ctx, h.IdleTimeout)
		defer cancel()

		h.setResponse(response.New(
			nil,
			http.StatusOK,
			nil, w),
		)

		defer func() {
			h.Metrics.ObserveRequest(time.Since(start), h.Response.Status)
			if h.Response.Err != nil {
				h.Metrics.IncError()
			}

			h.Metrics.IncRequest()

			log.Printf("%s: Response: %+v\n", fn, h.Response.Log())
		}()

		cnt, err := h.Storage.CountEntries(ctx)
		if err != nil {
			h.Response.MakeErr(err, http.StatusInternalServerError)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, err.Error())
			return
		}

		h.Response.Data = map[string]int64{"count": cnt}

		select {
		case <-ctx.Done():
			h.Response.MakeErr(ctx.Err(), http.StatusRequestTimeout)
			write.Write(&h.Response)
			log.Printf("%s: %s\n", fn, ctx.Err())
			return
		default:
			write.Write(&h.Response)
			log.Printf("%s: Counted\n", fn)
		}
	}
}

func (h *Handler) setResponse(r response.Response) {
	h.Response = r
}
