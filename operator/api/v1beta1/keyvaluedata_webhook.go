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

package v1beta1

import (
	"context"
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var keyvaluedatalog = logf.Log.WithName("keyvaluedata-resource")
var c client.Client
var ErrExistingKey = errors.New("this key already exists")
var ErrEmptyData = errors.New("empty data field")
var ErrEmptyKeyOrValue = errors.New("empty key or value")

func (r *KeyValueData) SetupWebhookWithManager(mgr ctrl.Manager) error {
	c = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-key-value-teamdev-com-v1beta1-keyvaluedata,mutating=true,failurePolicy=fail,sideEffects=None,groups=key-value.teamdev.com,resources=keyvaluedata,verbs=create;update,versions=v1beta1,name=mkeyvaluedata.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &KeyValueData{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *KeyValueData) Default() {
	keyvaluedatalog.Info("default", "name", r.Name)

}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-key-value-teamdev-com-v1beta1-keyvaluedata,mutating=false,failurePolicy=fail,sideEffects=None,groups=key-value.teamdev.com,resources=keyvaluedata,verbs=create;update,versions=v1beta1,name=vkeyvaluedata.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &KeyValueData{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *KeyValueData) ValidateCreate() error {
	keyvaluedatalog.Info("validate create", "name", r.Name)

	var createdResources KeyValueDataList
	if err := c.List(context.Background(), &createdResources); err != nil {
		keyvaluedatalog.Error(err, "can not get all KeyValueData resources")
		return err
	}
	for k, v := range r.Spec.Data {
		if k == "" || v == "" {
			return ErrEmptyKeyOrValue
		}
	}
	if len(r.Spec.Data) == 0 {
		return ErrEmptyData
	}
	for _, keyValueData := range createdResources.Items {
		if ok := r.containsAny(keyValueData.Spec.Data); ok {
			return ErrExistingKey
		}
	}
	return nil
}

func (r *KeyValueData) containsAny(d Data) bool {
	for k := range d {
		if _, ok := r.Spec.Data[k]; ok {
			return true
		}
	}
	return false
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *KeyValueData) ValidateUpdate(old runtime.Object) error {
	keyvaluedatalog.Info("validate create", "name", r.Name)

	var createdResources KeyValueDataList
	if err := c.List(context.Background(), &createdResources); err != nil {
		keyvaluedatalog.Error(err, "can not get all KeyValueData resources")
		return err
	}
	for k, v := range r.Spec.Data {
		if k == "" || v == "" {
			return ErrEmptyKeyOrValue
		}
	}
	if len(r.Spec.Data) == 0 {
		return ErrEmptyData
	}
	for _, keyValueData := range createdResources.Items {
		if ok := r.containsAny(keyValueData.Spec.Data); ok {
			return ErrExistingKey
		}
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *KeyValueData) ValidateDelete() error {
	keyvaluedatalog.Info("validate delete", "name", r.Name)

	return nil
}
