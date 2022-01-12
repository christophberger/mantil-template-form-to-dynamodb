package test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/mantil-io/go-mantil-template/api/form"
)

func TestForm(t *testing.T) {
	api := httpexpect.New(t, apiURL)

	req := form.DefaultRequest{}
	api.POST("/form").
		WithJSON(req).
		Expect().
		Status(http.StatusNoContent)

	saveReq := form.Form{
		Name:         "mantil",
		CanYouAttend: "Yes,  I'll be there",
		Count:        "20",
		Items: []string{
			"Drinks",
			"Sides/Appetizers",
		},
		Restrictions: "Nope",
		Email:        "email@emailnotfound.com",
	}
	api.POST("/form/save").
		WithJSON(saveReq).
		Expect().
		Status(http.StatusOK)

	api.POST("/form/list").
		Expect().
		ContentType("application/json").
		Status(http.StatusOK)
}
