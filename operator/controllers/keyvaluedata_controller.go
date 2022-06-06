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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sort"
	"time"
)

// KeyValueDataReconciler reconciles a KeyValueData object
type KeyValueDataReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	HttpClient *http.Client
	ServerURL  string
}

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
func (r *KeyValueDataReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var keyValueData kvdv1beta1.KeyValueData

	if err := r.Get(ctx, req.NamespacedName, &keyValueData); err != nil {
		logger.Error(err, "unable to fetch KeyValueData")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	entities := keyValueData.Spec.Data

	var requestStatuses []*kvdv1beta1.Condition
	var successSends int32
	var failedSends int32
	for k, e := range entities {
		reqData := &DataRequest{
			Key:    k,
			Entity: e,
		}

		marshalRequestData, err := json.Marshal(reqData)
		if err != nil {
			logger.Error(err, "can not marshall req data")
			requestStatuses = append(requestStatuses, &kvdv1beta1.Condition{
				Key:    k,
				Type:   kvdv1beta1.FailedType,
				Status: kvdv1beta1.FailedStatus,
				Reason: fmt.Sprintf("%s: %s",
					"can not marshall req data", err.Error()),
			})
			failedSends++
			continue
		}

		postRequest, err := http.NewRequest(http.MethodPost, r.ServerURL,
			bytes.NewBuffer(marshalRequestData))
		if err != nil && err != io.EOF {
			logger.Error(err, "can not create http request")
			logger.Error(err, "can not read req data to http request")
			requestStatuses = append(requestStatuses, &kvdv1beta1.Condition{
				Key:    k,
				Type:   kvdv1beta1.FailedType,
				Status: kvdv1beta1.FailedStatus,
				Reason: fmt.Sprintf("%s: %s",
					"can not marshall req data", err.Error()),
			})
			failedSends++
			continue
		}
		response, err := r.HttpClient.Do(postRequest)
		if err != nil {
			logger.Error(err, "can not send request")
			requestStatuses = append(requestStatuses, &kvdv1beta1.Condition{
				Key:    k,
				Type:   kvdv1beta1.FailedType,
				Status: kvdv1beta1.FailedStatus,
				Reason: fmt.Sprintf("%s: %s",
					"can not marshall req data", err.Error()),
				LastInsertTime: &metav1.Time{Time: time.Now()},
			})
			failedSends++
			continue
		}
		postRequest.Body.Close()

		if response.StatusCode != http.StatusCreated {
			var b []byte
			_, err = response.Body.Read(b)
			if err != nil {
				return ctrl.Result{}, err
			}
			logger.Error(fmt.Errorf("incorrect status code"), "server sent unexpected response",
				"expected status code", http.StatusCreated,
				"got", response.StatusCode, "body data", string(b))

			m := string(b)
			if len(b) == 0 {
				m = fmt.Sprintf("%s", response.Status)
			}

			requestStatuses = append(requestStatuses, &kvdv1beta1.Condition{
				Key:    k,
				Type:   kvdv1beta1.FailedType,
				Status: kvdv1beta1.FailedStatus,
				Reason: fmt.Sprintf("%s: %s",
					"data is not created", string(b)),
				Message: m,
			})
			failedSends++
			continue
		}
		response.Body.Close()

		requestStatuses = append(requestStatuses, &kvdv1beta1.Condition{
			Key:            k,
			Type:           kvdv1beta1.SuccessStatus,
			Status:         kvdv1beta1.AddedType,
			LastInsertTime: &metav1.Time{Time: time.Now()},
		})
		successSends++
	}

	sort.Slice(requestStatuses, func(i, j int) bool {
		if requestStatuses[i].Key < requestStatuses[j].Key {
			return true
		}
		return false
	})

	keyValueData.Status.Conditions = requestStatuses
	keyValueData.Status.FailedSends = &failedSends
	keyValueData.Status.SuccessSends = &successSends
	err := r.Client.Status().Update(ctx, &keyValueData)
	if err != nil {
		logger.Error(err, "can not update status", "resource name:", keyValueData.Name)
	}
	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeyValueDataReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kvdv1beta1.KeyValueData{}).
		Complete(r)
}
