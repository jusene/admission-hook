package main

import (
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

const (
	configMapPatch1 string = `[
         { "op": "add", "path": "/data/mutation-stage-1", "value": "yes" }
     ]`
	configMapPatch2 string = `[
         { "op": "add", "path": "/data/mutation-stage-2", "value": "yes" }
     ]`
)

// deny configmaps with specific key-value pair.
func admitConfigMaps(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("admitting configmaps")
	configMapResource := metav1.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}

	if ar.Request.Resource != configMapResource {
		klog.Errorf("expect resource to be %s", configMapResource)
		return nil
	}

	var raw []byte
	if ar.Request.Operation == v1.Delete {
		raw = ar.Request.OldObject.Raw
	} else {
		raw = ar.Request.Object.Raw
	}

	configmap := corev1.ConfigMap{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &configmap); err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}

	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true
	for k, v := range configmap.Data {
		if k == "webhook-e2e-test" && v == "webhook-disallow" &&
			(ar.Request.Operation == v1.Create || ar.Request.Operation == v1.Update) {
			reviewResponse.Allowed = false
			reviewResponse.Result = &metav1.Status{
				Reason: "the configmap contains unwanted key and value",
			}
		}

		if k == "webhook-e2e-test" && v == "webhook-nondeletable" && ar.Request.Operation == v1.Delete {
			reviewResponse.Allowed = false
			reviewResponse.Result = &metav1.Status{
				Reason: "the configmap cannot be deleted because it contains unwanted key and value",
			}
		}
	}
	return &reviewResponse
}

func mutateConfigmaps(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("mutating configmaps")
	configMapResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}
	if ar.Request.Resource != configMapResource {
		klog.Errorf("expect resource to be %s", configMapResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	configmap := corev1.ConfigMap{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &configmap); err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}
	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true
	if configmap.Data["mutation-start"] == "yes" {
		reviewResponse.Patch = []byte(configMapPatch1)
	}
	if configmap.Data["mutation-stage-1"] == "yes" {
		reviewResponse.Patch = []byte(configMapPatch2)
	}

	if len(reviewResponse.Patch) != 0 {
		pt := v1.PatchTypeJSONPatch
		reviewResponse.PatchType = &pt
	}

	return &reviewResponse
}
