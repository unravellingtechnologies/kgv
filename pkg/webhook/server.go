package webhook

import (
	log "github.com/sirupsen/logrus"
	"github.com/unravellingtechnologies/kgv/pkg/handlers/deployments"
	"github.com/unravellingtechnologies/kgv/pkg/handlers/pods"
	"net/http"
)

// SetupListeners function sets up the needed listeners for the webhook to function
func SetupListeners() *http.ServeMux {

	admissionHandler := newAdmissionHandler()

	deploymentValidation := deployments.NewValidationHook()
	podValidation := pods.NewValidationHook()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz())
	mux.HandleFunc("/v1/validate/pods", admissionHandler.Serve(podValidation))
	mux.HandleFunc("/v1/validate/deployments", admissionHandler.Serve(deploymentValidation))

	return mux
}

// healthz function provides a healthcheck endpoint for the webhook
func healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Ok!"))
		if err != nil {
			log.Error("Not able to answer request", err)
			return
		}
	}
}