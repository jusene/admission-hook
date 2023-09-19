package main

import (
	"fmt"
	"k8s.io/api/admission/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func admitCRD(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("admitting crd")

	resource := "customresourcedefinitions"
	v1beta1GVR := metav1.GroupVersionResource{
		Group:    apiextensionsv1beta1.GroupName,
		Version:  "v1beta1",
		Resource: resource,
	}
	v1GVR := metav1.GroupVersionResource{
		Group:    apiextensionsv1.GroupName,
		Version:  "v1",
		Resource: resource,
	}

	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true

	raw := ar.Request.Object.Raw
	var labels map[string]string

	switch ar.Request.Resource {
	case v1beta1GVR:
		crd := apiextensionsv1beta1.CustomResourceDefinition{}
		deserializer := codecs.UniversalDeserializer()
		if _, _, err := deserializer.Decode(raw, nil, &crd); err != nil {
			klog.Error(err)
			return toV1AdmissionResponse(err)
		}

		labels = crd.Labels
	case v1GVR:
		crd := apiextensionsv1.CustomResourceDefinition{}
		deserializer := codecs.UniversalDeserializer()
		if _, _, err := deserializer.Decode(raw, nil, &crd); err != nil {
			klog.Error(err)
			return toV1AdmissionResponse(err)
		}
		labels = crd.Labels
	default:
		err := fmt.Errorf("expect resource to be one of [%v, %v] but got %v", v1beta1GVR, v1GVR, ar.Request.Resource)
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}

	if v, ok := labels["webhook-e2e-test"]; ok {
		if v == "webhook-disallow" {
			reviewResponse.Allowed = false
			reviewResponse.Result = &metav1.Status{Message: "the crd contains unwanted label"}
		}
	}
	return &reviewResponse
}
