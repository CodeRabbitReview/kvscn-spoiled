// Package v1beta1 contains webhook and KeyValueData resource definition
// Webhook checks if request data keys does not exist in other KeyValueData
// and request KeyValueData keys and value is not empty and Data is not empty
// KeyValueData is a CRD, and it should be used in this case: send request to key value storage
// and see how data is sent to the server. KeyValueData contains KeyValueDataStatus to see this information.
// Example of usage:
//keyValueData := &v1beta1.KeyValueData{
//	TypeMeta: metav1.TypeMeta{
//		Kind:       "KeyValueData",
//		APIVersion: "key-value.teamdev.com/v1beta1",
//	},
//	ObjectMeta: metav1.ObjectMeta{
//		Name:      KeyValueDataName,
//		Namespace: KeyValueDataNamespace,
//	},
//	Spec: v1beta1.KeyValueDataSpec{
//		Data: map[string]string{
//			"test-key-string": "test-value-string",
//		},
//	},
//}
//
//k8sClient.Create(ctx, keyValueData)
package v1beta1
