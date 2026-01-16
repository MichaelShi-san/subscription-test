package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/MichaelShi-san/subscription-test/internal/model"
	"github.com/MichaelShi-san/subscription-test/internal/service"
)

type SubscriptionHandler struct {
	svc *service.SubscriptionService
	log *slog.Logger
}

func NewSubscriptionHandler(s *service.SubscriptionService, l *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{svc: s, log: l}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Service string `json:"service_name"`
		Price   int    `json:"price"`
		UserID  string `json:"user_id"`
		Start   string `json:"start_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start, err := time.Parse("01-2006", req.Start)
	if err != nil {
		http.Error(w, "invalid start_date format", http.StatusBadRequest)
		return
	}

	sub := &model.Subscription{
		Service:   req.Service,
		Price:     req.Price,
		UserID:    req.UserID,
		StartDate: start,
	}

	if err := h.svc.Create(r.Context(), sub); err != nil {
		h.log.Error("create subscription failed", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.log.Info("subscription created", "id", sub.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(items)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	_ = h.svc.Delete(r.Context(), id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var req struct {
		Service string     `json:"service_name"`
		Price   int        `json:"price"`
		EndDate *time.Time `json:"end_date"`
	}

	_ = json.NewDecoder(r.Body).Decode(&req)

	sub := &model.Subscription{
		ID:      id,
		Service: req.Service,
		Price:   req.Price,
		EndDate: req.EndDate,
	}

	_ = h.svc.Update(r.Context(), sub)
	w.WriteHeader(http.StatusOK)
}

func (h *SubscriptionHandler) TotalCost(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")

	from, _ := time.Parse("2006-01", r.URL.Query().Get("from"))
	to, _ := time.Parse("2006-01", r.URL.Query().Get("to"))

	sum, err := h.svc.TotalCost(r.Context(), userID, serviceName, from, to)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]int{"total": sum})
}
