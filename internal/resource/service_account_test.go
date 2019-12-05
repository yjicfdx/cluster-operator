package resource_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	rabbitmqv1beta1 "github.com/pivotal/rabbitmq-for-kubernetes/api/v1beta1"
	"github.com/pivotal/rabbitmq-for-kubernetes/internal/resource"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("ServiceAccount", func() {
	var (
		serviceAccount        *corev1.ServiceAccount
		instance              rabbitmqv1beta1.RabbitmqCluster
		serviceAccountBuilder *resource.ServiceAccountBuilder
		builder               *resource.RabbitmqResourceBuilder
	)

	BeforeEach(func() {
		instance = rabbitmqv1beta1.RabbitmqCluster{
			ObjectMeta: v1.ObjectMeta{
				Name:      "a name",
				Namespace: "a namespace",
			},
		}
		builder = &resource.RabbitmqResourceBuilder{
			Instance: &instance,
		}
		serviceAccountBuilder = builder.ServiceAccount()
	})

	Context("Build", func() {
		BeforeEach(func() {
			obj, err := serviceAccountBuilder.Build()
			serviceAccount = obj.(*corev1.ServiceAccount)
			Expect(err).NotTo(HaveOccurred())
		})

		It("generates a ServiceAccount with the correct name and namespace", func() {
			Expect(serviceAccount.Name).To(Equal(builder.Instance.ChildResourceName("server")))
			Expect(serviceAccount.Namespace).To(Equal(builder.Instance.Namespace))
		})

		It("only creates the required labels", func() {
			labels := serviceAccount.Labels
			Expect(len(labels)).To(Equal(3))
			Expect(labels["app.kubernetes.io/name"]).To(Equal(instance.Name))
			Expect(labels["app.kubernetes.io/component"]).To(Equal("rabbitmq"))
			Expect(labels["app.kubernetes.io/part-of"]).To(Equal("pivotal-rabbitmq"))
		})
	})

	Context("Build with instance that has labels", func() {
		BeforeEach(func() {
			instance.Labels = map[string]string{
				"app.kubernetes.io/foo": "bar",
				"foo":                   "bar",
				"rabbitmq":              "is-great",
				"foo/app.kubernetes.io": "edgecase",
			}

			obj, err := serviceAccountBuilder.Build()
			serviceAccount = obj.(*corev1.ServiceAccount)
			Expect(err).NotTo(HaveOccurred())
		})

		It("has the labels from the CRD on the serviceAccount", func() {
			testLabels(serviceAccount.Labels)
		})

		It("also has the required labels", func() {
			labels := serviceAccount.Labels
			Expect(labels["app.kubernetes.io/name"]).To(Equal(instance.Name))
			Expect(labels["app.kubernetes.io/component"]).To(Equal("rabbitmq"))
			Expect(labels["app.kubernetes.io/part-of"]).To(Equal("pivotal-rabbitmq"))
		})
	})

	Context("Update", func() {
		BeforeEach(func() {
			instance.Labels = map[string]string{
				"app.kubernetes.io/foo": "bar",
				"foo":                   "bar",
				"rabbitmq":              "is-great",
				"foo/app.kubernetes.io": "edgecase",
			}

			serviceAccount = &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name": "rabbit-labelled",
					},
				},
			}
			Expect(serviceAccountBuilder.Update(serviceAccount)).To(Succeed())
		})

		It("adds labels from the CRD on the service account", func() {
			testLabels(serviceAccount.Labels)
		})

		It("persists the labels it had before Update", func() {
			Expect(serviceAccount.Labels).To(HaveKeyWithValue("app.kubernetes.io/name", "rabbit-labelled"))
		})
	})
})