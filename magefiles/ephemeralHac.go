package main

import (
	"context"
	"fmt"
	"time"

	"github.com/RedHatInsights/ephemeral-namespace-operator/apis/cloud.redhat.com/v1alpha1"
	insightsUtils "github.com/RedHatInsights/rhc-osdk-utils/utils"
	client "github.com/redhat-appstudio/e2e-tests/pkg/apis/kubernetes"
	"github.com/redhat-appstudio/e2e-tests/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (hac EphemeralHac) TestTarget() error {
	kubeClient, err := client.NewK8SClient()
	if err != nil {
		return err
	}
	namespaceReservation := &v1alpha1.NamespaceReservation{
		ObjectMeta: v1.ObjectMeta{
			Name: "rhopp-reservation",
		},
		Spec: v1alpha1.NamespaceReservationSpec{
			Requester: "rhopp",
			Duration:  insightsUtils.StringPtr("1h"),
			Pool:      "default",
		},
	}

	// userSignup := &toolchainApi.UserSignup{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      "user1",
	// 		Namespace: "toolchain-host-operator",
	// 		Annotations: map[string]string{
	// 			"toolchain.dev.openshift.com/user-email": "user1@user.us",
	// 		},
	// 		Labels: map[string]string{
	// 			"toolchain.dev.openshift.com/email-hash": md5.CalcMd5("user1@user.us"),
	// 		},
	// 	},
	// 	Spec: toolchainApi.UserSignupSpec{
	// 		Userid:   "user1",
	// 		Username: "user1",
	// 		States:   []toolchainApi.UserSignupState{"approved"},
	// 	},
	// }
	fmt.Printf("Creating: %+v\n", namespaceReservation)
	if err := kubeClient.KubeRest().Create(context.TODO(), namespaceReservation); err != nil {
		return err
	}
	err = utils.WaitUntil(func() (done bool, err error) {
		err = kubeClient.KubeRest().Get(context.TODO(), types.NamespacedName{Name: "rhopp-reservation"}, namespaceReservation)
		if err != nil {
			return false, err
		}
		fmt.Printf("Waiting. %+v\n", namespaceReservation)
		if namespaceReservation.Status.State == "active" {
			return true, nil
		}
		return false, nil
	}, 2*time.Minute)

	if err != nil {
		return err
	}
	return nil
}
