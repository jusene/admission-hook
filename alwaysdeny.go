package main

import (
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

// alwaysDeny all requests made to this function.
func alwaysDeny(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("calling always-deny")
	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = false
	reviewResponse.Result = &metav1.Status{Message: "this webhook denies all requests"}
	return &reviewResponse
}
