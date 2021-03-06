package webhook

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/unravellingtechnologies/kgv/pkg/handlers"
	"io"
	admission "k8s.io/api/admission/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net/http"
)

// admissionHandler represents the HTTP handler for an admission webhook
type admissionHandler struct {
	decoder runtime.Decoder
}

// NewAdmissionHandler returns an instance of AdmissionHandler
func newAdmissionHandler() *admissionHandler {
	return &admissionHandler{
		decoder: serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer(),
	}
}

// Serve returns a webhook.HandlerFunc for an admission webhook
func (h *admissionHandler) Serve(hook handlers.Hook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		review := parseRequest(h, w, r)

		result, err := hook.Execute(review.Request)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		admissionResponse := admission.AdmissionReview{
			Response: &admission.AdmissionResponse{
				UID:     review.Request.UID,
				Allowed: result.Allowed,
				Result:  &meta.Status{Message: result.Msg},
			},
		}

		// set the patch operations for mutating admission
		if len(result.PatchOps) > 0 {
			patchBytes, err := json.Marshal(result.PatchOps)
			if err != nil {
				log.Error(err)
				http.Error(w, fmt.Sprintf("could not marshal JSON patch: %v", err), http.StatusInternalServerError)
			}
			admissionResponse.Response.Patch = patchBytes
		}

		res, err := json.Marshal(admissionResponse)
		if err != nil {
			log.Error(err)
			http.Error(w, fmt.Sprintf("could not marshal response: %v", err), http.StatusInternalServerError)
			return
		}

		log.Infof("Webhook [%s - %s] - Allowed: %t", r.URL.Path, review.Request.Operation, result.Allowed)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(res)
		if err != nil {
			log.Error("Error while writing the response", err)
			return
		}
	}
}

// parseRequest performs the validations to ensure the request is valid and well-formed
func parseRequest(h *admissionHandler, w http.ResponseWriter, r *http.Request) *admission.AdmissionReview {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprint("invalid method only POST requests are allowed"), http.StatusMethodNotAllowed)
		return nil
	}

	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, fmt.Sprint("only content type 'application/json' is supported"), http.StatusBadRequest)
		return nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not read request body: %v", err), http.StatusBadRequest)
		return nil
	}

	var review admission.AdmissionReview
	if _, _, err := h.decoder.Decode(body, nil, &review); err != nil {
		http.Error(w, fmt.Sprintf("could not deserialize request: %v", err), http.StatusBadRequest)
		return nil
	}

	if review.Request == nil {
		http.Error(w, "malformed admission review: request is nil", http.StatusBadRequest)
		return nil
	}

	return &review
}