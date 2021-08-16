package pods

import (
	"encoding/json"
	"github.com/unravellingtechnologies/kgv/pkg/handlers"
	admission "k8s.io/api/admission/v1"
	"strings"

	apps "k8s.io/api/core/v1"
)

// NewValidationHook creates a new instance of pods validation hook
func NewValidationHook() handlers.Hook {
	return handlers.Hook{
		Create: validateCreate(),
	}
}

// NewMutationHook creates a new instance of pods mutation hook
//func NewMutationHook() kgv.Hook {
//	return kgv.Hook{
//		Create: mutateCreate(),
//	}
//}

// validateCreate function validates the pods
func validateCreate() handlers.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*handlers.Result, error) {
		pod, err := parsePod(r.Object.Raw)
		if err != nil {
			return &handlers.Result{Msg: err.Error()}, nil
		}

		for _, c := range pod.Spec.Containers {
			if strings.HasSuffix(c.Image, ":latest") {
				return &handlers.Result{Msg: "You cannot use the tag 'latest' in a container."}, nil
			}
		}

		return &handlers.Result{Allowed: true}, nil
	}
}

// parsePod function parses the pod object into the expected format
func parsePod(object []byte) (*apps.Pod, error) {
	var pod apps.Pod
	if err := json.Unmarshal(object, &pod); err != nil {
		return nil, err
	}

	return &pod, nil
}