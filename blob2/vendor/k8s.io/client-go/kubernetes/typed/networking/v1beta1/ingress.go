/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	context "context"

	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	applyconfigurationsnetworkingv1beta1 "k8s.io/client-go/applyconfigurations/networking/v1beta1"
	gentype "k8s.io/client-go/gentype"
	scheme "k8s.io/client-go/kubernetes/scheme"
)

// IngressesGetter has a method to return a IngressInterface.
// A group's client should implement this interface.
type IngressesGetter interface {
	Ingresses(namespace string) IngressInterface
}

// IngressInterface has methods to work with Ingress resources.
type IngressInterface interface {
	Create(ctx context.Context, ingress *networkingv1beta1.Ingress, opts v1.CreateOptions) (*networkingv1beta1.Ingress, error)
	Update(ctx context.Context, ingress *networkingv1beta1.Ingress, opts v1.UpdateOptions) (*networkingv1beta1.Ingress, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, ingress *networkingv1beta1.Ingress, opts v1.UpdateOptions) (*networkingv1beta1.Ingress, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*networkingv1beta1.Ingress, error)
	List(ctx context.Context, opts v1.ListOptions) (*networkingv1beta1.IngressList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *networkingv1beta1.Ingress, err error)
	Apply(ctx context.Context, ingress *applyconfigurationsnetworkingv1beta1.IngressApplyConfiguration, opts v1.ApplyOptions) (result *networkingv1beta1.Ingress, err error)
	// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
	ApplyStatus(ctx context.Context, ingress *applyconfigurationsnetworkingv1beta1.IngressApplyConfiguration, opts v1.ApplyOptions) (result *networkingv1beta1.Ingress, err error)
	IngressExpansion
}

// ingresses implements IngressInterface
type ingresses struct {
	*gentype.ClientWithListAndApply[*networkingv1beta1.Ingress, *networkingv1beta1.IngressList, *applyconfigurationsnetworkingv1beta1.IngressApplyConfiguration]
}

// newIngresses returns a Ingresses
func newIngresses(c *NetworkingV1beta1Client, namespace string) *ingresses {
	return &ingresses{
		gentype.NewClientWithListAndApply[*networkingv1beta1.Ingress, *networkingv1beta1.IngressList, *applyconfigurationsnetworkingv1beta1.IngressApplyConfiguration](
			"ingresses",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *networkingv1beta1.Ingress { return &networkingv1beta1.Ingress{} },
			func() *networkingv1beta1.IngressList { return &networkingv1beta1.IngressList{} },
			gentype.PrefersProtobuf[*networkingv1beta1.Ingress](),
		),
	}
}