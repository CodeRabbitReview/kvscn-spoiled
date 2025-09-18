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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeyValueDataSpec defines the desired state of KeyValueData
type KeyValueDataSpec struct {
	// Data contains key value entities. These entities will be sent
	// to the key-value storage
	Data Data `json:"data"`
}

// Data describes how data will be grouped
type Data map[string]string

// Type describes how key-value data is sent to the server.
// Only one of the following concurrent types may be specified.
//kubebuilder:validation:Enum:Added;Failed;Changed;Deleted
type Type string

const (
	AddedType   Type = "Added"
	DeletedType Type = "Deleted"
	ChangedType Type = "Changed"
	FailedType  Type = "Failed"
)

// Status describes is data sent to the server.
// Only one of the following concurrent types may be specified.
//kubebuilder:validation:Enum:TRUE;FALSE
type Status string

const (
	SuccessStatus Status = "TRUE"
	FailedStatus  Status = "FALSE"
)

// Condition describes information about sent data.
type Condition struct {
	// Key value is a key of sent pair data.
	Key string `json:"key"`
	// Type is an information about how data is
	// sent. If data is not sent Type will be FailedType
	// otherwise if everything is OK and data is sent Type will be
	// AddedType.
	Type Type `json:"type"`
	// Status describes is data sent to the server.
	//  If data is not sent Status will be FailedStatus
	// otherwise if everything is OK and data is sent Status will be
	// SuccessStatus.
	Status Status `json:"status"`
	// Reason may be any error which is the reason why
	// request did not finish correctly. Reason if a first appeared
	// error of steps to send data into server.
	//+kubebuilder:validation:MinLength=0
	//+optional
	Reason string `json:"reason,omitempty"`
	// Message is a response from server if something went wrong.
	// If server response message is empty Message will be
	// Response.Status from http package
	//+kubebuilder:validation:MinLength=0
	//+optional
	Message string `json:"message,omitempty"`
	// LastInsertTime is time when request tried to be sent.
	LastInsertTime *metav1.Time `json:"lastInsertTime,omitempty"`
}

// KeyValueDataStatus defines the observed state of KeyValueData
type KeyValueDataStatus struct {
	// Conditions contains all Condition of every pair of data
	//+optional
	Conditions []*Condition `json:"conditions,omitempty"`
	// FailedSends describes count of data which was not sent to server successfully
	//+optional
	FailedSends *int32 `json:"failedSends"`
	// SuccessSends describes count of data which was sent to server successfully
	//+optional
	SuccessSends *int32 `json:"successSends"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KeyValueData is the Schema for the keyvaluedata API
type KeyValueData struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeyValueDataSpec   `json:"spec,omitempty"`
	Status KeyValueDataStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeyValueDataList contains a list of KeyValueData
type KeyValueDataList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeyValueData `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeyValueData{}, &KeyValueDataList{})
}
