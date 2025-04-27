package handler

import (
	"broker/models"
	"net/http"
)

func AddProperty(w http.ResponseWriter, r *http.Request) {
	handleCreateRequest[models.CreateProperty, models.Property](w, r, "configurator.property.add")
}
