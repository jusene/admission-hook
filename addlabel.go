package main

import (
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

const (
	addFirstLabelPatch string = `[
         { "op": "add", "path": "/metadata/labels", "value": {"added-label": "yes"}},
         { "op": "add", "path": "/spec/selector/matchLabels","value": {"added-label": "yes"}},
         { "op": "add", "path": "/spec/template/metadata/labels","value": {"added-label": "yes"}}
     ]`
	addAdditionalLabelPatch string = `[
         { "op": "add", "path": "/metadata/labels/added-label", "value": "yes" },
         { "op": "add", "path": "/spec/template/metadata/labels/added-label", "value": "yes"}
     ]`
	updateLabelPatch string = `[
         { "op": "replace", "path": "/metadata/labels/added-label", "value": "yes" },
         { "op": "replace", "path": "/spec/template/metadata/labels/added-label", "value": "yes" }
     ]`
)

// Add a label {"added-label": "yes"} to the object
func addLabel(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("calling add-label")
	obj := struct {
		metav1.ObjectMeta `json:"metadata,omitempty"`
	}{}
	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &obj)
	if err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}

	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true

	pt := v1.PatchTypeJSONPatch
	labelValue, hasLabel := obj.ObjectMeta.Labels["added-label"]
	switch {
	case len(obj.ObjectMeta.Labels) == 0:
		reviewResponse.Patch = []byte(addFirstLabelPatch)
		reviewResponse.PatchType = &pt
	case !hasLabel:
		reviewResponse.Patch = []byte(addAdditionalLabelPatch)
		reviewResponse.PatchType = &pt
	case labelValue != "yes":
		reviewResponse.Patch = []byte(updateLabelPatch)
		reviewResponse.PatchType = &pt
	default:
		// already set
	}
	fmt.Println(string(reviewResponse.Patch))
	return &reviewResponse
}
