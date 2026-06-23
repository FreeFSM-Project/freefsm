package handlers

import (
	"net/http"
	"strconv"

	"github.com/MartialM1nd/freefsm/internal/ent"
	"github.com/MartialM1nd/freefsm/internal/ent/asset"
	"github.com/MartialM1nd/freefsm/internal/ent/customer"
	"github.com/MartialM1nd/freefsm/internal/ent/estimate"
	"github.com/MartialM1nd/freefsm/internal/ent/invoice"
	"github.com/MartialM1nd/freefsm/internal/ent/item"
	"github.com/MartialM1nd/freefsm/internal/ent/job"
	"github.com/MartialM1nd/freefsm/internal/ent/project"
	"github.com/go-chi/chi/v5"
)

func requireActiveObject(client *ent.Client, objectType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil || id <= 0 {
				http.NotFound(w, r)
				return
			}
			if !activeObjectExists(r, client, objectType, id) {
				http.Error(w, "archived records are read-only", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func activeObjectExists(r *http.Request, client *ent.Client, objectType string, id int64) bool {
	switch objectType {
	case "customer":
		exists, err := client.Customer.Query().Where(customer.IDEQ(id), customer.DeletedAtIsNil()).Exist(r.Context())
		return err == nil && exists
	case "job":
		exists, err := client.Job.Query().Where(job.IDEQ(id), job.DeletedAtIsNil()).Exist(r.Context())
		return err == nil && exists
	case "project":
		exists, err := client.Project.Query().Where(project.IDEQ(id), project.DeletedAtIsNil()).Exist(r.Context())
		return err == nil && exists
	case "estimate":
		exists, err := client.Estimate.Query().Where(estimate.IDEQ(id), estimate.DeletedAtIsNil()).Exist(r.Context())
		return err == nil && exists
	case "invoice":
		exists, err := client.Invoice.Query().Where(invoice.IDEQ(id), invoice.DeletedAtIsNil()).Exist(r.Context())
		return err == nil && exists
	case "asset":
		exists, err := client.Asset.Query().Where(asset.IDEQ(id), asset.DeletedAtIsNil()).Exist(r.Context())
		return err == nil && exists
	case "item":
		exists, err := client.Item.Query().Where(item.IDEQ(id), item.DeletedAtIsNil()).Exist(r.Context())
		return err == nil && exists
	default:
		return false
	}
}
