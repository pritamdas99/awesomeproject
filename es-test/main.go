package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	dbc "kubedb.dev/db-client-go/elasticsearch"
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
	//fmt.Println(config)
	kc, err := client.New(config, client.Options{
		Scheme: scm,
		Mapper: nil,
	})
	if err != nil {
		klog.Error(err)
	}
	fmt.Println("got kb client", kc)
	db := &api.Elasticsearch{}
	err = kc.Get(context.TODO(), types.NamespacedName{
		Name:      "es",
		Namespace: "demo",
	}, db)
	if err != nil {
		klog.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	esClient, err := dbc.NewKubeDBClientBuilder(kc, db).WithContext(ctx).WithURL("https://localhost:9200").GetElasticClient()
	if err != nil {
		klog.Error(err)
	}

	fmt.Println("got my client", esClient)

	response, err := esClient.ShardStats()
	if err != nil {
		klog.Error(err)
	}

	fmt.Println(response)
	for _, x := range response {
		fmt.Println(x.Index, x.Shard, x.State, x.UnassignedReason)
	}
}
