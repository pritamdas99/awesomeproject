package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	dbc "kubedb.dev/db-client-go/solr"
	"os"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientSetScheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubedbscheme "kubedb.dev/apimachinery/client/clientset/versioned/scheme"
)

var (
	kubeconfigPath = func() string {
		kubecfg := os.Getenv("KUBECONFIG")
		if kubecfg != "" {
			return kubecfg
		}
		return filepath.Join(homedir.HomeDir(), ".kube", "config")
	}()
	kubeContext = ""
)

var (
	scm = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientSetScheme.AddToScheme(scm))
	utilruntime.Must(kubedbscheme.AddToScheme(scm))
}

func main() {
	fmt.Println("start........................")
	config, err := clientcmd.BuildConfigFromContext(kubeconfigPath, kubeContext)
	if err != nil {
		fmt.Println("Failed to get config")
	}
	kc, err := client.New(config, client.Options{
		Scheme: scm,
		Mapper: nil,
	})
	if err != nil {
		klog.Error(err)
	}
	fmt.Println("got kb client", kc)
	db := &api.Solr{}
	err = kc.Get(context.TODO(), types.NamespacedName{
		Name:      "solr-combined",
		Namespace: "demo",
	}, db)
	if err != nil {
		klog.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	slClient, err := dbc.NewKubeDBClientBuilder(kc, db).WithContext(ctx).WithLog(klog.Background()).WithURL("http://localhost:8983").GetSolrClient()
	if err != nil {
		klog.Error(err)
	}

	fmt.Println("got my client", slClient)

	response, err := slClient.BackupCollection(context.TODO(), "book", "book-backup", "s3:/hello", "kubedb-linode-s3")
	if err != nil {
		klog.Error(err)
	}

	responseBody, err := slClient.DecodeResponse(response)
	if err != nil {
		klog.Error(err)
	}

	status, err := slClient.GetResponseStatus(responseBody)
	if err != nil {
		klog.Error(err)
	}

	if status != 0 {
		klog.Errorf(fmt.Sprintf("status is not 0"))
	} else {
		fmt.Println("***************************************status is 0")
	}

	fmt.Println(responseBody)
}
