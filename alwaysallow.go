package main

import (
	"fmt"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"time"
)

// alwaysAllowDelayFiveSeconds sleeps for five seconds and allows all requests made to this function.
func alwaysAllowDelayFiveSeconds(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("always-allow-with-delay sleeping for 5 seconds")
	time.Sleep(5 * time.Second)
	klog.V(2).Info("calling always allow")
	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true
	reviewResponse.Result = &metav1.Status{Message: "this webhook allows all requests"}
	fmt.Println(reviewResponse)
	return &reviewResponse
}
