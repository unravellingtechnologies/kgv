package deployments

import (
	"encoding/json"
	"github.com/unravellingtechnologies/kgv/pkg/handlers"
	admission "k8s.io/api/admission/v1"
	apps "k8s.io/api/apps/v1"
)

// NewValidationHook creates a new instance of deployment validation hook
func NewValidationHook() handlers.Hook {
	return handlers.Hook{
		Create: validateCreate(),
		Update: validateUpdate(),
	}
}

// validateCreate function validates deployments being newly created
func validateCreate() handlers.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*handlers.Result, error) {
		deployment, err := parseDeployment(r.Object.Raw)
		if err != nil {
			return &handlers.Result{Msg: err.Error()}, nil
		}

		if deployment.Namespace == "special" {
			return &handlers.Result{Msg: "You cannot create a deployment in `special` namespace."}, nil
		}

		return &handlers.Result{Allowed: true}, nil
	}
}

// validateUpdate validates deployments being updated
func validateUpdate() handlers.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*handlers.Result, error) {
		deployment, err := parseDeployment(r.OldObject.Raw)
		if err != nil {
			return &handlers.Result{Msg: err.Error()}, nil
		}

		if deployment.Namespace == "special" {
			return &handlers.Result{Msg: "You cannot create a deployment in `special` namespace."}, nil
		}

		return &handlers.Result{Allowed: true}, nil
	}
}

// parseDeployment function parses a Deployment manifest into an usable object
func parseDeployment(object []byte) (*apps.Deployment, error) {
	var deployment apps.Deployment
	if err := json.Unmarshal(object, &deployment); err != nil {
		return nil, err
	}

	return &deployment, nil
}