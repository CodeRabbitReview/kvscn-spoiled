// Package controllers contains reconcile function for KeyValueData custom resource.
// KeyValueDataReconciler.Reconcile gets created CRD and go throw Spec.Data field and tries to send
// request to key value storage. If some errors appears or response status code from server is not equal to
// http.StatusCreated condition of this key value data will be Failed request to server.
// Otherwise, condition will be successes. Condition of all request will be in alphabet order.
// Example of usage:
//keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
//_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
//	if err != nil {
//		panic(err)
//	}
package controllers
