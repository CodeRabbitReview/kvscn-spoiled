//nolint
package controllers_test

import (
	"context"
	"github.com/miprokop/crd-kvd/api/v1beta1"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("KeyValueData controller", func() {
	Context("Inside of a new namespace", func() {
		ctx := context.Background()
		const (
			timeout               = time.Second * 10
			interval              = time.Millisecond * 250
			KeyValueDataName      = "test-keyvaluedata1"
			KeyValueDataNamespace = "default"
		)
		BeforeEach(func() {
			k8sClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
			reconciler.Client = k8sClient
		})
		AfterEach(func() {
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
			}
			k8sClient.Delete(ctx, keyValueData)
		})
		It("Should change KeyValueData Status to one successes action", func() {
			By("Creating new KeyValueData successfully")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "test-value-string",
					},
				},
			}
			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())

			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
			createdKeyValueData := v1beta1.KeyValueData{}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			if err != nil {
				panic(err)
			}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, keyValueDataLookupKey, &createdKeyValueData)
				return err != nil
			}, timeout, interval).Should(BeTrue())
			Expect(*createdKeyValueData.Status.SuccessSends).To(Equal(int32(1)))
			Expect(*createdKeyValueData.Status.FailedSends).To(Equal(int32(0)))
			Expect(len(createdKeyValueData.Status.Conditions)).To(Equal(len(keyValueData.Spec.Data)))
			Expect(*createdKeyValueData.Status.Conditions[0]).To(Equal(v1beta1.Condition{
				Key:            "test-key-string",
				Type:           v1beta1.AddedType,
				Status:         v1beta1.SuccessStatus,
				Reason:         "",
				Message:        "",
				LastInsertTime: createdKeyValueData.Status.Conditions[0].LastInsertTime,
			}))
		})

		It("Should change KeyValueData Status to one Unsuccessful action", func() {
			By("Creating new KeyValueData unsuccessfully")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "test-value-string",
					},
				},
			}
			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())

			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
			createdKeyValueData := v1beta1.KeyValueData{}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred(), "failed to callReconcile")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, keyValueDataLookupKey, &createdKeyValueData)
				return err != nil
			}, timeout, interval).Should(BeTrue())
			Expect(*createdKeyValueData.Status.SuccessSends).To(Equal(int32(0)))
			Expect(*createdKeyValueData.Status.FailedSends).To(Equal(int32(1)))
			Expect(len(createdKeyValueData.Status.Conditions)).To(Equal(len(keyValueData.Spec.Data)))
			Expect(*createdKeyValueData.Status.Conditions[0]).To(Equal(v1beta1.Condition{
				Key:            "test-key-string",
				Type:           v1beta1.FailedType,
				Status:         v1beta1.FailedStatus,
				Reason:         "",
				Message:        "500 Internal Server Error",
				LastInsertTime: nil,
			}))
		})

		It("Should change KeyValueData Status to one Unsuccessful action", func() {
			By("Empty key")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"": "test-value-string",
					},
				},
			}
			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())

			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
			createdKeyValueData := v1beta1.KeyValueData{}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred(), "failed to callReconcile")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, keyValueDataLookupKey, &createdKeyValueData)
				return err != nil
			}, timeout, interval).Should(BeTrue())
			Expect(*createdKeyValueData.Status.SuccessSends).To(Equal(int32(0)))
			Expect(*createdKeyValueData.Status.FailedSends).To(Equal(int32(1)))
			Expect(len(createdKeyValueData.Status.Conditions)).To(Equal(len(keyValueData.Spec.Data)))
			Expect(*createdKeyValueData.Status.Conditions[0]).To(Equal(v1beta1.Condition{
				Key:            "",
				Type:           v1beta1.FailedType,
				Status:         v1beta1.FailedStatus,
				Reason:         "empty key or value",
				Message:        "",
				LastInsertTime: nil,
			}))
		})

		It("Should change KeyValueData Status to one Unsuccessful action", func() {
			By("Empty value")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "",
					},
				},
			}
			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())

			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
			createdKeyValueData := v1beta1.KeyValueData{}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred(), "failed to callReconcile")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, keyValueDataLookupKey, &createdKeyValueData)
				return err != nil
			}, timeout, interval).Should(BeTrue())
			Expect(*createdKeyValueData.Status.SuccessSends).To(Equal(int32(0)))
			Expect(*createdKeyValueData.Status.FailedSends).To(Equal(int32(1)))
			Expect(len(createdKeyValueData.Status.Conditions)).To(Equal(len(keyValueData.Spec.Data)))
			Expect(*createdKeyValueData.Status.Conditions[0]).To(Equal(v1beta1.Condition{
				Key:            "test-key-string",
				Type:           v1beta1.FailedType,
				Status:         v1beta1.FailedStatus,
				Reason:         "empty key or value",
				Message:        "",
				LastInsertTime: nil,
			}))
		})

		It("Incorrect HTTP server URL", func() {
			By("Empty value")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string1": "123",
						"test-key-string2": "1234",
					},
				},
			}
			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())

			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
			createdKeyValueData := v1beta1.KeyValueData{}
			reconciler.ServerURL = ""
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred(), "failed to callReconcile")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, keyValueDataLookupKey, &createdKeyValueData)
				return err != nil
			}, timeout, interval).Should(BeTrue())
			Expect(*createdKeyValueData.Status.SuccessSends).To(Equal(int32(0)))
			Expect(*createdKeyValueData.Status.FailedSends).To(Equal(int32(2)))
			Expect(len(createdKeyValueData.Status.Conditions)).To(Equal(len(keyValueData.Spec.Data)))
			Expect(*createdKeyValueData.Status.Conditions[0]).To(Equal(v1beta1.Condition{
				Key:            "test-key-string1",
				Type:           v1beta1.FailedType,
				Status:         v1beta1.FailedStatus,
				Reason:         "can not send request data to server: Post \"\": unsupported protocol scheme \"\"",
				Message:        "",
				LastInsertTime: nil,
			}))
			Expect(*createdKeyValueData.Status.Conditions[1]).To(Equal(v1beta1.Condition{
				Key:            "test-key-string2",
				Type:           v1beta1.FailedType,
				Status:         v1beta1.FailedStatus,
				Reason:         "can not send request data to server: Post \"\": unsupported protocol scheme \"\"",
				Message:        "",
				LastInsertTime: nil,
			}))
		})

		It("Incorrect HTTP format server URL", func() {
			By("Empty value")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string1": "123",
						"test-key-string2": "1234",
					},
				},
			}
			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())

			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
			createdKeyValueData := v1beta1.KeyValueData{}
			reconciler.ServerURL = ""
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred(), "failed to callReconcile")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, keyValueDataLookupKey, &createdKeyValueData)
				return err != nil
			}, timeout, interval).Should(BeTrue())
			Expect(*createdKeyValueData.Status.SuccessSends).To(Equal(int32(0)))
			Expect(*createdKeyValueData.Status.FailedSends).To(Equal(int32(2)))
			Expect(len(createdKeyValueData.Status.Conditions)).To(Equal(len(keyValueData.Spec.Data)))
			Expect(*createdKeyValueData.Status.Conditions[0]).To(Equal(v1beta1.Condition{
				Key:            "test-key-string1",
				Type:           v1beta1.FailedType,
				Status:         v1beta1.FailedStatus,
				Reason:         "can not send request data to server: Post \"\": unsupported protocol scheme \"\"",
				Message:        "",
				LastInsertTime: nil,
			}))
			Expect(*createdKeyValueData.Status.Conditions[1]).To(Equal(v1beta1.Condition{
				Key:            "test-key-string2",
				Type:           v1beta1.FailedType,
				Status:         v1beta1.FailedStatus,
				Reason:         "can not send request data to server: Post \"\": unsupported protocol scheme \"\"",
				Message:        "",
				LastInsertTime: nil,
			}))
		})

		It("Incorrect HTTP server URL", func() {
			By("Empty value")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string1": "123",
					},
				},
			}
			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())

			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
			createdKeyValueData := v1beta1.KeyValueData{}
			reconciler.ServerURL = "postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require"
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred(), "failed to callReconcile")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, keyValueDataLookupKey, &createdKeyValueData)
				return err != nil
			}, timeout, interval).Should(BeTrue())
			Expect(*createdKeyValueData.Status.SuccessSends).To(Equal(int32(0)))
			Expect(*createdKeyValueData.Status.FailedSends).To(Equal(int32(1)))
			Expect(len(createdKeyValueData.Status.Conditions)).To(Equal(len(keyValueData.Spec.Data)))
			Expect(*createdKeyValueData.Status.Conditions[0]).To(Equal(v1beta1.Condition{
				Key:            "test-key-string1",
				Type:           v1beta1.FailedType,
				Status:         v1beta1.FailedStatus,
				Reason:         "can not create request to server: parse \"postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require\": net/url: invalid userinfo",
				Message:        "",
				LastInsertTime: nil,
			}))
		})
	})
})
