package values

import (
	"strings"

	"github.com/go-logr/logr"
	"github.com/loft-sh/utils/pkg/helm"
)

var EKSAPIVersionMap = map[string]string{
	"1.27": "public.ecr.aws/eks-distro/kubernetes/kube-apiserver:v1.27.3-eks-1-27-7",
	"1.26": "public.ecr.aws/eks-distro/kubernetes/kube-apiserver:v1.26.6-eks-1-26-13",
	"1.25": "public.ecr.aws/eks-distro/kubernetes/kube-apiserver:v1.25.11-eks-1-25-17",
	"1.24": "public.ecr.aws/eks-distro/kubernetes/kube-apiserver:v1.24.15-eks-1-24-21",
	"1.23": "public.ecr.aws/eks-distro/kubernetes/kube-apiserver:v1.23.17-eks-1-23-26",
}

var EKSControllerVersionMap = map[string]string{
	"1.27": "public.ecr.aws/eks-distro/kubernetes/kube-controller-manager:v1.27.3-eks-1-27-7",
	"1.26": "public.ecr.aws/eks-distro/kubernetes/kube-controller-manager:v1.26.6-eks-1-26-13",
	"1.25": "public.ecr.aws/eks-distro/kubernetes/kube-controller-manager:v1.25.11-eks-1-25-17",
	"1.24": "public.ecr.aws/eks-distro/kubernetes/kube-controller-manager:v1.24.15-eks-1-24-21",
	"1.23": "public.ecr.aws/eks-distro/kubernetes/kube-controller-manager:v1.23.17-eks-1-23-26",
}

var EKSEtcdVersionMap = map[string]string{
	"1.27": "public.ecr.aws/eks-distro/etcd-io/etcd:v3.5.7-eks-1-27-7",
	"1.26": "public.ecr.aws/eks-distro/etcd-io/etcd:v3.5.7-eks-1-26-13",
	"1.25": "public.ecr.aws/eks-distro/etcd-io/etcd:v3.5.7-eks-1-25-17",
	"1.24": "public.ecr.aws/eks-distro/etcd-io/etcd:v3.5.7-eks-1-24-21",
	"1.23": "public.ecr.aws/eks-distro/etcd-io/etcd:v3.5.7-eks-1-23-26",
}

var EKSCoreDNSVersionMap = map[string]string{
	"1.27": "public.ecr.aws/eks-distro/coredns/coredns:v1.10.1-eks-1-27-7",
	"1.26": "public.ecr.aws/eks-distro/coredns/coredns:v1.9.3-eks-1-26-13",
	"1.25": "public.ecr.aws/eks-distro/coredns/coredns:v1.9.3-eks-1-25-17",
	"1.24": "public.ecr.aws/eks-distro/coredns/coredns:v1.9.3-eks-1-24-21",
	"1.23": "public.ecr.aws/eks-distro/coredns/coredns:v1.8.7-eks-1-23-26",
}

func getDefaultEKSReleaseValues(chartOptions *helm.ChartOptions, log logr.Logger) (string, error) {
	serverVersionString := GetKubernetesVersion(chartOptions.KubernetesVersion)
	serverMinorInt, err := GetKubernetesMinorVersion(chartOptions.KubernetesVersion)
	if err != nil {
		return "", err
	}

	apiImage := EKSAPIVersionMap[serverVersionString]
	controllerImage := EKSControllerVersionMap[serverVersionString]
	etcdImage := EKSEtcdVersionMap[serverVersionString]
	corednsImage, ok := EKSCoreDNSVersionMap[serverVersionString]
	if !ok {
		if serverMinorInt > 27 {
			log.Info("officially unsupported host server version, will fallback to virtual cluster version v1.27", "serverVersion", serverVersionString)
			apiImage = EKSAPIVersionMap["1.27"]
			controllerImage = EKSControllerVersionMap["1.27"]
			etcdImage = EKSEtcdVersionMap["1.27"]
			corednsImage = EKSCoreDNSVersionMap["1.27"]
		} else {
			log.Info("officially unsupported host server version, will fallback to virtual cluster version v1.23", "serverVersion", serverVersionString)
			apiImage = EKSAPIVersionMap["1.23"]
			controllerImage = EKSControllerVersionMap["1.23"]
			etcdImage = EKSEtcdVersionMap["1.23"]
			corednsImage = EKSCoreDNSVersionMap["1.23"]
		}
	}

	// build values
	values := `api:
  image: ##API_IMAGE##
controller:
  image: ##CONTROLLER_IMAGE##
etcd:
  image: ##ETCD_IMAGE##
coredns:
  image: ##COREDNS_IMAGE##
`
	values = strings.ReplaceAll(values, "##API_IMAGE##", apiImage)
	values = strings.ReplaceAll(values, "##CONTROLLER_IMAGE##", controllerImage)
	values = strings.ReplaceAll(values, "##ETCD_IMAGE##", etcdImage)
	values = strings.ReplaceAll(values, "##COREDNS_IMAGE##", corednsImage)
	return addCommonReleaseValues(values, chartOptions)
}
