package httpapi

import (
    "encoding/json"
    "errors"
    "net/http"
    "time"

    "pricemon/internal/domain"
)

type createSubReq struct {
    UserID   string                      `json:"userId"`
    Ref      domain.ProductRef           `json:"ref"`
    Rule     domain.Rule                 `json:"rule"`
    Targets  []domain.NotificationTarget `json:"targets"`
    CronSpec string                      `json:"cronSpec,omitempty"`
}

func (s *Server) createSubscription(w http.ResponseWriter, r *http.Request) {
    var req createSubReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeErr(w, http.StatusBadRequest, err)
        return
    }
    if req.UserID == "" || req.Ref.Site == "" || req.Ref.SKU == "" || len(req.Targets) == 0 {
        writeErr(w, http.StatusBadRequest, errors.New("userId, ref, and targets required"))
        return
    }

    sub := domain.Subscription{
        ID:        newID(),
        UserID:    req.UserID,
        Ref:       req.Ref,
        Rule:      req.Rule,
        Targets:   req.Targets,
        CronSpec:  req.CronSpec,
        Active:    true,
        CreatedAt: time.Now().UTC(),
    }
    if err := s.subs.Put(r.Context(), sub); err != nil {
        writeErr(w, http.StatusInternalServerError, err)
        return
    }
    if err := s.sched.ScheduleSubscription(sub); err != nil {
        writeErr(w, http.StatusInternalServerError, err)
        return
    }
    writeJSON(w, http.StatusCreated, sub)
}

func (s *Server) getSubscription(w http.ResponseWriter, r *http.Request) {
    sub, err := s.subs.Get(r.Context(), r.PathValue("id"))
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeErr(w, http.StatusNotFound, err)
            return
        }
        writeErr(w, http.StatusInternalServerError, err)
        return
    }
    writeJSON(w, http.StatusOK, sub)
}

func (s *Server) deleteSubscription(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if err := s.subs.Delete(r.Context(), id); err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeErr(w, http.StatusNotFound, err)
            return
        }
        writeErr(w, http.StatusInternalServerError, err)
        return
    }
    s.sched.UnscheduleSubscription(id)
    w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, err error) {
    writeJSON(w, status, map[string]string{"error": err.Error()})
}

// TODO: replace with github.com/google/uuid → uuid.NewString()
func newID() string { return time.Now().UTC().Format("20060102T150405.000000000") }
