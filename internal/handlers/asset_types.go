package handlers

import (
	"net/http"
	"strconv"

	"github.com/MartialM1nd/freefsm/internal/ent"
	"github.com/MartialM1nd/freefsm/internal/services"
	"github.com/MartialM1nd/freefsm/internal/templates"
	"github.com/go-chi/chi/v5"
)

type AssetTypeHandler struct {
	svc           *services.AssetTypeService
	assetStatusSvc *services.AssetStatusService
}

func NewAssetTypeHandler(svc *services.AssetTypeService, assetStatusSvc *services.AssetStatusService) *AssetTypeHandler {
	return &AssetTypeHandler{svc: svc, assetStatusSvc: assetStatusSvc}
}

func (h *AssetTypeHandler) Show(w http.ResponseWriter, r *http.Request) {
	types, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	statuses, err := h.assetStatusSvc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	typeRows := make([]templates.AssetTypeRow, len(types))
	for i, t := range types {
		typeRows[i] = templates.AssetTypeRow{
			ID:        t.ID,
			Name:      t.Name,
			SortOrder: t.SortOrder,
		}
	}
	
	statusRows := make([]templates.AssetStatusRow, len(statuses))
	for i, s := range statuses {
		statusRows[i] = templates.AssetStatusRow{
			ID:        s.ID,
			Name:      s.Name,
			Color:     s.Color,
			SortOrder: s.SortOrder,
		}
	}
	
	templates.AssetSettingsPage(
		templates.AssetTypeListPageData{Types: typeRows},
		templates.AssetStatusListPageData{Statuses: statusRows},
	).Render(r.Context(), w)
}

func (h *AssetTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", 400)
		return
	}
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", 400)
		return
	}
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))
	if _, err := h.svc.Create(r.Context(), name, sortOrder); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/asset-types?flash=Type+created", http.StatusSeeOther)
}

func (h *AssetTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", 400)
		return
	}
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", 400)
		return
	}
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))
	if _, err := h.svc.Update(r.Context(), id, name, sortOrder); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/asset-types?flash=Type+updated", http.StatusSeeOther)
}

func (h *AssetTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/asset-types?flash=Type+deleted", http.StatusSeeOther)
}

type AssetStatusHandler struct {
	svc *services.AssetStatusService
}

func NewAssetStatusHandler(svc *services.AssetStatusService) *AssetStatusHandler {
	return &AssetStatusHandler{svc: svc}
}

func (h *AssetStatusHandler) List(w http.ResponseWriter, r *http.Request) {
	statuses, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	rows := make([]templates.AssetStatusRow, len(statuses))
	for i, s := range statuses {
		rows[i] = templates.AssetStatusRow{
			ID:        s.ID,
			Name:      s.Name,
			Color:     s.Color,
			SortOrder: s.SortOrder,
		}
	}
	templates.AssetStatusList(templates.AssetStatusListPageData{Statuses: rows}).Render(r.Context(), w)
}

func (h *AssetStatusHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", 400)
		return
	}
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", 400)
		return
	}
	color := r.FormValue("color")
	if color == "" {
		color = "#6B7280"
	}
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))
	if _, err := h.svc.Create(r.Context(), name, color, sortOrder); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/asset-statuses?flash=Status+created", http.StatusSeeOther)
}

func (h *AssetStatusHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", 400)
		return
	}
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", 400)
		return
	}
	color := r.FormValue("color")
	if color == "" {
		color = "#6B7280"
	}
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))
	if _, err := h.svc.Update(r.Context(), id, name, color, sortOrder); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/asset-statuses?flash=Status+updated", http.StatusSeeOther)
}

func (h *AssetStatusHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/asset-statuses?flash=Status+deleted", http.StatusSeeOther)
}

func assetTypeRow(t *ent.AssetType) templates.AssetTypeRow {
	return templates.AssetTypeRow{
		ID:        t.ID,
		Name:      t.Name,
		SortOrder: t.SortOrder,
	}
}

func assetStatusRow(s *ent.AssetStatus) templates.AssetStatusRow {
	return templates.AssetStatusRow{
		ID:        s.ID,
		Name:      s.Name,
		Color:     s.Color,
		SortOrder: s.SortOrder,
	}
}
