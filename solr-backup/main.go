package main

import (
	"fmt"
	"kmodules.xyz/client-go/tools/clientcmd"
	"kubedb.dev/apimachinery/pkg/factory"
	"os"
	"path/filepath"

	"k8s.io/client-go/util/homedir"
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

func main() {
	fmt.Println("start........................")
	config, err := clientcmd.BuildConfigFromContext(kubeconfigPath, kubeContext)
	if err != nil {
		fmt.Println("Failed to get config")
	}
	fmt.Println("got congig", config)
	config.Burst = 100
	config.QPS = 100
	KBClient, err := factory.NewUncachedClient(config)
	if err != nil {
		fmt.Println("Failed to get clinet")
	}
	fmt.Println("got kb client", KBClient)
	//db := &api.Solr{}
	//err = KBClient.Get(context.TODO(), types.NamespacedName{
	//	Name:      "solr-combined",
	//	Namespace: "demo",
	//}, db)
	//if err != nil {
	//	klog.Error(err)
	//}
	//slClient, err := dbc.NewKubeDBClientBuilder(KBClient, db).WithContext(context.TODO()).WithLog(klog.Background()).WithURL("http://localhost:8983").GetSolrClient()
	//if err != nil {
	//	klog.Error(err)
	//}
	//
	//response, err := slClient.BackupRestoreCollection("BACKUP", "book", "book-backup", "s3:/pritam", "kubedb-linode")
	//if err != nil {
	//	klog.Error(err)
	//}
	//
	//responseBody, err := slClient.DecodeResponse(response)
	//if err != nil {
	//	klog.Error(err)
	//}
	//
	//status, err := slClient.GetResponseStatus(responseBody)
	//if err != nil {
	//	klog.Error(err)
	//}
	//
	//if status != 0 {
	//	klog.Errorf(fmt.Sprintf("status is not 0"))
	//}
	//
	//fmt.Println(responseBody)
}
