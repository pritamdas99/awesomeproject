package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	dbc "kubedb.dev/db-client-go/solr"
	"os"
	"path/filepath"
	"sort"
	"sync"
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

type CoreList struct {
	coreName   string
	collection string
}

type UpdateList struct {
	target     string
	replica    string
	collection string
}

func CheckupStatus(async string, slClient *dbc.Client) error {
	var wg sync.WaitGroup
	wg.Add(1)
	var errr error
	go func() {
		defer wg.Done()
		asyncId := async
		for {
			resp, err := slClient.RequestStatus(asyncId)
			if err != nil {
				klog.Error(fmt.Sprintf("Failed to get response for asyncId %s. Error: %v", asyncId, err))
				errr = err
				return
			}

			responseBody, err := slClient.DecodeResponse(resp)
			if err != nil {
				klog.Error(fmt.Sprintf("Failed to decode response for asyncId %s. Error: %v", asyncId, err))
				errr = err
				return
			}

			_, err = slClient.GetResponseStatus(responseBody)
			if err != nil {
				klog.Error(fmt.Sprintf("status is non zero while checking status for asyncId %s. Error %v", asyncId, err))
				errr = err
				return
			}

			state, err := slClient.GetAsyncStatus(responseBody)
			if err != nil {
				klog.Error(fmt.Sprintf("status is non zero while checking state of async for asyncId %s. Error %v", asyncId, err))
				errr = err
				return
			}
			klog.Info(fmt.Sprintf("State for asyncid %v is %v\n", asyncId, state))
			if state == "completed" {
				klog.Info("Status is completed for ", asyncId)
				err = flushStatus(asyncId, slClient)
				if err != nil {
					errr = err
					return
				}
				return
			} else if state == "failed" {
				klog.Info(fmt.Sprintf("API call for asyncId %s failed", asyncId))
				err = flushStatus(asyncId, slClient)
				if err != nil {
					errr = err
					return
				}
				errr = fmt.Errorf("response for asyncid %v. failed with response %v", asyncId, responseBody)
				break
			} else if state == "notfound" {
				klog.Info(fmt.Sprintf("API call for asyncid %s not found", asyncId))
				break
			}
			time.Sleep(10 * time.Second)
		}
	}()
	wg.Wait()
	return errr
}

func flushStatus(asyncId string, slClient *dbc.Client) error {
	resp, err := slClient.FlushStatus(asyncId)
	if err != nil {
		return err
	}

	responseBody, err := slClient.DecodeResponse(resp)
	if err != nil {
		return err
	}

	_, err = slClient.GetResponseStatus(responseBody)
	if err != nil {
		klog.Error("status is non zero while flushing status", err)
		return err
	}

	return nil
}

func CleanupAsync(slClient *dbc.Client, async string) error {
	var wg sync.WaitGroup
	wg.Add(1)
	var errr error
	go func() {
		defer wg.Done()
		asyncId := async
		for {
			resp, err := slClient.RequestStatus(asyncId)
			if err != nil {
				klog.Error(fmt.Sprintf("Failed to get response for asyncId %s. Error: %v", asyncId, err))
				errr = err
				break
			}

			responseBody, err := slClient.DecodeResponse(resp)
			if err != nil {
				klog.Error(fmt.Sprintf("Failed to decode response for asyncId %s. Error: %v", asyncId, err))
				errr = err
				break
			}

			_, err = slClient.GetResponseStatus(responseBody)
			if err != nil {
				klog.Error(fmt.Sprintf("status is non zero while checking status for asyncId %s. Error %v", asyncId, err))
				errr = err
				break
			}

			state, err := slClient.GetAsyncStatus(responseBody)
			if err != nil {
				klog.Error(fmt.Sprintf("status is non zero while checking state of async for asyncId %s. Error %v", asyncId, err))
				errr = err
				break
			}
			klog.Info(fmt.Sprintf("State for asyncid %v is %v\n", asyncId, state))
			if state == "completed" || state == "notfound" || state == "failed" {
				err := flushStatus(asyncId, slClient)
				if err != nil {
					klog.Error(fmt.Sprintf("Failed to flush api call for asyncId %s. Error %v", asyncId, err))
					time.Sleep(20 * time.Second)
					continue
				}
				break
			}
			time.Sleep(10 * time.Second)
		}
	}()
	wg.Wait()
	return errr
}

func balance(slClient *dbc.Client) {
	async := "balance-replica"
	err := CleanupAsync(slClient, async)
	if err != nil {
		klog.Error(fmt.Sprintf("Failed to clean asyncid******************************* %v\n", async))
		time.Sleep(30 * time.Second)
		return
	} else {
		klog.Info(fmt.Sprintf("Cleanup async successful for %v", async))
		time.Sleep(10 * time.Second)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := slClient.BalanceReplica(async)
		if err != nil {
			klog.Error(fmt.Errorf("failed to do balance request. err %v", err))
			return
		}
		responseBody, err := slClient.DecodeResponse(resp)
		if err != nil {
			klog.Error(fmt.Errorf("failed to decode response for async %s, err %v", async, err))
			return
		}
		_, err = slClient.GetResponseStatus(responseBody)
		if err != nil {
			klog.Error(fmt.Errorf("failed to decode response for async %s, err %v", async, err))
			return
		}

		err = CheckupStatus(async, slClient)
		if err != nil {
			klog.Error("Error while checking status************ ", err)
		}
	}()
	wg.Wait()
}

func run(lst []UpdateList, slClient *dbc.Client) {
	var wg sync.WaitGroup
	for _, x := range lst {
		target := x.target
		replica := x.replica
		collection := x.collection
		async := fmt.Sprintf("%s-%s-%s", replica, collection, target)
		err := CleanupAsync(slClient, async)
		if err != nil {
			klog.Error(fmt.Sprintf("Failed to clean asyncid******************************* %v\n", async))
			time.Sleep(30 * time.Second)
			return
		} else {
			klog.Info(fmt.Sprintf("Cleanup async successful for %v", async))
			time.Sleep(10 * time.Second)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := slClient.MoveReplica(target, replica, collection, async)
			if err != nil {
				klog.Error(fmt.Errorf("failed to do request for target %s, replica %s, collection %s, err %v", target, replica, collection, err))
				return
			}
			responseBody, err := slClient.DecodeResponse(resp)
			if err != nil {
				klog.Error(fmt.Errorf("failed to decode response for target %s, replica %s, collection %s, err %v", target, replica, collection, err))
				return
			}
			_, err = slClient.GetResponseStatus(responseBody)
			if err != nil {
				klog.Error(fmt.Errorf("failed to decode response for target %s, replica %s, collection %s, err %v", target, replica, collection, err))
				return
			}

			err = CheckupStatus(async, slClient)
			if err != nil {
				klog.Error("Error while checking status************ ", err)
			}
		}()
	}
	wg.Wait()
}

func down(nodeList []string, x int, mp map[string][]CoreList, slClient *dbc.Client) {
	n := len(nodeList)
	ls2 := nodeList[n-x:]
	ls1 := nodeList[:n-x]
	fmt.Println("ls1 ", ls1)
	fmt.Println("ls2 ", ls2)
	ar := make([]UpdateList, 0)
	for _, node := range ls2 {
		for _, core := range mp[node] {
			id := -1
			mx := 1000000000
			for j, l1 := range ls1 {
				if len(mp[l1]) < mx {
					mx = len(mp[l1])
					id = j
				}
			}
			ar = append(ar, UpdateList{
				target:     ls1[id],
				replica:    core.coreName,
				collection: core.collection,
			})
			mp[ls1[id]] = append(mp[ls1[id]], core)
			fmt.Println(core.coreName, core.collection, ls1[id])
		}
	}
	run(ar, slClient)
}
func up(nodeList []string, mp map[string][]CoreList, slClient *dbc.Client) {
	for _, x := range nodeList {
		if _, ok := mp[x]; !ok {
			mp[x] = make([]CoreList, 0)
		}
	}
	ar := make([]UpdateList, 0)
	for {
		mn := 10000000000
		minNode := ""
		mx := -1
		maxNode := ""
		for x, y := range mp {
			n := len(y)
			if mx < n {
				mx = n
				maxNode = x
			}

			if mn > n {
				mn = n
				minNode = x
			}
		}
		if maxNode == minNode || mx-mn <= 1 {
			break
		}
		target := minNode
		core := mp[maxNode][0].coreName
		collection := mp[maxNode][0].collection
		mp[minNode] = append(mp[minNode], mp[maxNode][0])
		mp[maxNode] = mp[maxNode][1:]
		ar = append(ar, UpdateList{
			target:     target,
			replica:    core,
			collection: collection,
		})
		fmt.Println(target, core, collection)
	}
	run(ar, slClient)
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
		Name:      "solr-cluster",
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

	//response, err := slClient.BackupCollection(context.TODO(), "book", "book-backup", "s3:/hello", "kubedb-linode-s3")
	//if err != nil {
	//	klog.Error(err)
	//}

	response, err := slClient.GetClusterStatus()
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

	clusterInfo, ok := responseBody["cluster"].(map[string]interface{})
	if !ok {
		klog.Error(fmt.Errorf("did not find cluster %v\n", responseBody))
	}
	collections, ok := clusterInfo["collections"].(map[string]interface{})
	if !ok {
		klog.Error("didn't find collections")
	}
	mp := make(map[string][]CoreList)
	for collection, info := range collections {
		collectionInfo := info.(map[string]interface{})
		shardInfo := collectionInfo["shards"].(map[string]interface{})
		for _, info := range shardInfo {
			shardInfo := info.(map[string]interface{})
			replicaInfo := shardInfo["replicas"].(map[string]interface{})
			for core, info := range replicaInfo {
				coreInfo := info.(map[string]interface{})
				nodeName := coreInfo["node_name"].(string)
				if _, ok := mp[nodeName]; !ok {
					mp[nodeName] = make([]CoreList, 0)
				}
				mp[nodeName] = append(mp[nodeName], CoreList{
					coreName:   core,
					collection: collection,
				})
			}
		}
	}

	nodeList := make([]string, 0)

	liveNodes, ok := clusterInfo["live_nodes"]
	if ok {
		fmt.Println("got livenodes")
	} else {
		fmt.Println("failed to get that")
	}
	xx := liveNodes.([]interface{})
	for _, node := range xx {
		x := node.(string)
		nodeList = append(nodeList, x)
	}
	sort.Strings(nodeList)
	fmt.Println(nodeList)

	//if db.Spec.Topology != nil {
	//	r := *db.Spec.Topology.Data.Replicas
	//	c := *db.Spec.Topology.Coordinator.Replicas
	//	nodeList = nodeList[c : r+c]
	//	fmt.Println("nodes ", r, nodeList)
	//}

	//fmt.Println(responseBody)

	//balance(slClient)
	//down(nodeList, 1, mp, slClient)
	up(nodeList, mp, slClient)
}
