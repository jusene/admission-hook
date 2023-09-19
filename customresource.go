package main

import (
	"encoding/json"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

const (
	customResourcePatch1 string = `[
         { "op": "add", "path": "/data/mutation-stage-1", "value": "yes" }
     ]`
	customResourcePatch2 string = `[
         { "op": "add", "path": "/data/mutation-stage-2", "value": "yes" }
     ]`
)

func mutateCustomResource(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("mutating custom resource")
	cr := struct {
		metav1.ObjectMeta
		Data map[string]string
	}{}

	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &cr)
	if err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}

	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true

	if cr.Data["mutation-start"] == "yes" {
		reviewResponse.Patch = []byte(customResourcePatch1)
	}

	if cr.Data["mutation-stage-1"] == "yes" {
		reviewResponse.Patch = []byte(customResourcePatch2)
	}

	if len(reviewResponse.Patch) != 0 {
		pt := v1.PatchTypeJSONPatch
		reviewResponse.PatchType = &pt
	}

	return &reviewResponse
}

func admitCustomResource(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("admitting custom resource")
	cr := struct {
		metav1.ObjectMeta
		Data map[string]string
	}{}

	var raw []byte
	if ar.Request.Operation == v1.Delete {
		raw = ar.Request.OldObject.Raw
	} else {
		raw = ar.Request.Object.Raw
	}
	err := json.Unmarshal(raw, &cr)
	if err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}

	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true
	for k, v := range cr.Data {
		if k == "webhook-e2e-test" && v == "webhook-disallow" &&
			(ar.Request.Operation == v1.Create || ar.Request.Operation == v1.Update) {
			reviewResponse.Allowed = false
			reviewResponse.Result = &metav1.Status{
				Reason: "the custom resource contains unwanted data",
			}
		}
		if k == "webhook-e2e-test" && v == "webhook-nondeletable" && ar.Request.Operation == v1.Delete {
			reviewResponse.Allowed = false
			reviewResponse.Result = &metav1.Status{
				Reason: "the custom resource cannot be deleted because it contains unwanted key and value",
			}
		}
	}
	return &reviewResponse
}
