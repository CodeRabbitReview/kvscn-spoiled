/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	kvdv1beta1 "github.com/miprokop/crd-kvd/api/v1beta1"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sort"
	"time"
)

// KeyValueDataReconciler reconciles a KeyValueData object
type KeyValueDataReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	HTTPClient    *http.Client
	ServerURL     string
	FinalizerName string
}

// DataRequest pair or key and value
type DataRequest struct {
	Key    string `json:"key"`
	Entity string `json:"entity"`
}

//+kubebuilder:rbac:groups=key-value.teamdev.com,resources=keyvaluedata,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=key-value.teamdev.com,resources=keyvaluedata/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=key-value.teamdev.com,resources=keyvaluedata/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
// Reconcile gets created resource and write condition of request
// Reconcile goes throw Spec.Data map and gets key and value data from it.
// Try to send request to the key value storage - if an error appears Reconcile writes error condition
// if everything is okay with request and status code does not equal http.StatusCreated
// tries to get message from server. If there are any messages will be used Response.Status otherwise this message.
// Also writes error into condition.
// Sort conditions of all requests in alphabet order.
// Update the status of corresponding resource.
func (r *KeyValueDataReconciler) Reconcile(ctx context.Context,
	req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var keyValueData kvdv1beta1.KeyValueData
	if err := r.Get(ctx, req.NamespacedName, &keyValueData); err != nil {
		logger.Error(err, "unable to fetch KeyValueData")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info("got key value data resource", "name: ", keyValueData.Name,
		"time: ", time.Now().String())

	if keyValueData.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("add default finalizer", "resource name: ", keyValueData.Name)
		if !controllerutil.ContainsFinalizer(&keyValueData, r.FinalizerName) {
			controllerutil.AddFinalizer(&keyValueData, r.FinalizerName)

			if err := r.Update(ctx, &keyValueData); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		logger.Info("delete corresponding data before resource deletion",
			"resource name: ", keyValueData.Name, "time: ", time.Now().String())
		if controllerutil.ContainsFinalizer(&keyValueData, r.FinalizerName) {
			for k := range keyValueData.Spec.Data {
				_, _, _ = r.createRequests(ctx, DataRequest{
					Key: k,
				}, http.MethodDelete,
					http.StatusNoContent, kvdv1beta1.DeletedType)
			}
			controllerutil.RemoveFinalizer(&keyValueData, r.FinalizerName)
			if err := r.Update(ctx, &keyValueData); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}
	var conditions []*kvdv1beta1.Condition
	var successSends int32
	var failedSends int32
	if len(keyValueData.Status.Conditions) != 0 {
		for i, c := range keyValueData.Status.Conditions {
			if v, ok := keyValueData.Spec.Data[c.Key]; !ok {
				_, _, _ = r.createRequests(ctx, DataRequest{
					Key: c.Key,
				},
					http.MethodDelete, http.StatusNoContent, kvdv1beta1.DeletedType)
				logger.Info("key to delete in update", "key: ", c.Key)
				keyValueData.Status.Conditions[i] = keyValueData.Status.Conditions[len(keyValueData.Status.Conditions)-1]
				keyValueData.Status.Conditions[len(keyValueData.Status.Conditions)-1] = nil
				keyValueData.Status.Conditions = keyValueData.Status.Conditions[:len(keyValueData.Status.Conditions)-1]
			} else {
				condition, success, failed := r.createRequests(ctx, DataRequest{
					Key:    c.Key,
					Entity: v,
				},
					http.MethodPut, http.StatusCreated, kvdv1beta1.ChangedType)
				conditions = append(conditions, condition)
				successSends += success
				failedSends += failed
			}
		}
	}

	contains := func(conditions []*kvdv1beta1.Condition, key string) bool {
		for i := 0; i < len(conditions); i++ {
			if conditions[i].Key == key {
				return true
			}
		}
		return false
	}

	for k, v := range keyValueData.Spec.Data {
		if !contains(conditions, k) {
			condition, success, failed := r.createRequests(ctx, DataRequest{
				Key:    k,
				Entity: v,
			}, http.MethodPost, http.StatusCreated, kvdv1beta1.AddedType)
			conditions = append(conditions, condition)
			successSends += success
			failedSends += failed
		}
	}

	sort.Slice(conditions, func(i, j int) bool {
		return conditions[i].Key < conditions[j].Key
	})

	keyValueData.Status.Conditions = conditions
	keyValueData.Status.FailedSends = &failedSends
	keyValueData.Status.SuccessSends = &successSends
	err := r.Client.Status().Update(ctx, &keyValueData)
	logger.Info("updated status", "resource name: ", keyValueData.Name)
	return ctrl.Result{}, err
}

func (r *KeyValueDataReconciler) createRequests(ctx context.Context, entity DataRequest,
	method string, expectedStatusCode int, neededType kvdv1beta1.Type) (*kvdv1beta1.Condition, int32, int32) {
	logger := log.FromContext(ctx)

	logger.Info("request key", "key: ", entity.Key)
	marshaledRequestData, err := json.Marshal(entity)
	if err != nil {
		logger.Error(err, "can not marshall req data")
		return &kvdv1beta1.Condition{
			Key:    entity.Key,
			Type:   kvdv1beta1.FailedType,
			Status: kvdv1beta1.FailedStatus,
			Reason: fmt.Sprintf("%s: %s",
				"can not marshall req data", err.Error()),
		}, 0, 1
	}
	logger.Info("marshaled resource", "key: ", entity.Key, "time: ", time.Now().String())

	request, err := http.NewRequest(method, r.ServerURL,
		bytes.NewBuffer(marshaledRequestData))
	if err != nil && err != io.EOF {
		logger.Error(err, "can not create http request")
		logger.Error(err, "can not read req data to http request")
		return &kvdv1beta1.Condition{
			Key:    entity.Key,
			Type:   kvdv1beta1.FailedType,
			Status: kvdv1beta1.FailedStatus,
			Reason: fmt.Sprintf("%s: %s",
				"can not create request to server", err.Error()),
		}, 0, 1
	}
	response, err := r.HTTPClient.Do(request)
	if err != nil {
		logger.Error(err, "can not send request")
		return &kvdv1beta1.Condition{
			Key:    entity.Key,
			Type:   kvdv1beta1.FailedType,
			Status: kvdv1beta1.FailedStatus,
			Reason: fmt.Sprintf("%s: %s",
				"can not send request data to server", err.Error()),
		}, 0, 1
	}
	logger.Info("sent request to server", "key: ", entity.Key, "time: ", time.Now().String())
	request.Body.Close()

	if response.StatusCode != expectedStatusCode {
		var b []byte
		var reason string
		_, err = response.Body.Read(b)
		if err != nil && err != io.EOF {
			reason = fmt.Sprintf("data is not created: %s", err)
		}
		logger.Error(fmt.Errorf("incorrect status code"), "server sent unexpected response",
			"expected status code", http.StatusCreated,
			"got", response.StatusCode, "body data", string(b), "key: ", entity.Key)

		m := string(b)
		if len(b) == 0 {
			m = response.Status
		}

		return &kvdv1beta1.Condition{
			Key:     entity.Key,
			Type:    kvdv1beta1.FailedType,
			Status:  kvdv1beta1.FailedStatus,
			Reason:  reason,
			Message: m,
		}, 0, 1
	}
	response.Body.Close()

	logger.Info("send value to server successfully",
		"key: ", entity.Key, "time: ", time.Now().String())

	return &kvdv1beta1.Condition{
		Key:            entity.Key,
		Type:           neededType,
		Status:         kvdv1beta1.SuccessStatus,
		LastInsertTime: &metav1.Time{Time: time.Now()},
	}, 1, 0
}

// SetupWithManager sets up the controller with the Manager. With GenerationChanged predicate into this manages.
func (r *KeyValueDataReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kvdv1beta1.KeyValueData{}, builder.WithPredicates(&predicate.GenerationChangedPredicate{})).
		Complete(r)
}
