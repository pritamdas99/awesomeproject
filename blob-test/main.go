package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
	_ "gocloud.dev/blob/s3blob"
	"gomodules.xyz/pointer"
	"io"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/tools/clientcmd"
	"kmodules.xyz/client-go/tools/portforward"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	accessKeyID, secretAccessKey string
)

func init() {
	accessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
}

func main() {
	ctx := context.TODO()
	_, sess, err := connect(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to S3 bucket: %v\n", err)
	}
	//iter := bucket.List(nil)
	//for {
	//	obj, err := iter.Next(ctx)
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(obj.Key)
	//
	//}
	//defer func() {
	//	_ = bucket.Close()
	//}()

	//fmt.Println("-----------List done--------------")
	//fmt.Println("connected..................", bucket)

	key := "/solr/bar.txt" // "/solr/pritam/car-backup1/"
	//_ = key
	err = put(sess, key)
	if err != nil {
		fmt.Println("FATAL LOGE GERJEEJDJDD")
		log.Fatal(err)
	}
	//err = List(sess, "/")
	//if err != nil {
	//	log.Fatal(err)
	//}
	err = check(sess, key)
	if err != nil {
		log.Fatal(err)
	}

	//err = write(ctx, bucket, key)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = read(ctx, bucket, key)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
func List(sess *session.Session, prefix string) error {
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("appscode-testing"),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Name)
	return nil
}

var (
	storageClass   = "standard"
	kubeconfigPath = func() string {
		kubecfg := os.Getenv("KUBECONFIG")
		if kubecfg != "" {
			return kubecfg
		}
		return filepath.Join(homedir.HomeDir(), ".kube", "config")
	}()
	kubeContext = ""
)

func connect(ctx context.Context) (*blob.Bucket, *session.Session, error) {
	config, err := clientcmd.BuildConfigFromContext(kubeconfigPath, kubeContext)
	kubeClient := kubernetes.NewForConfigOrDie(config)
	fmt.Println(accessKeyID, secretAccessKey)

	tunnel := portforward.NewTunnel(portforward.TunnelOptions{
		Client:    kubeClient.CoreV1().RESTClient(),
		Config:    config,
		Resource:  string(core.ResourceServices),
		Namespace: "demo",
		Name:      "s3proxy-s3",
		Remote:    80,
	})
	if err := tunnel.ForwardPort(); err != nil {
		return nil, nil, err
	}
	port := tunnel.Local
	klog.Info(fmt.Sprintf("http://localhost:%v", port))
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String(fmt.Sprintf("http://localhost:%v", port)),
		//Endpoint:    aws.String("https://us-east-1.linodeobjects.com"),
		S3ForcePathStyle: pointer.BoolP(true),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	fmt.Println("WE HAVE GOT THE BUCKET", err)

	bucket, err := s3blob.OpenBucket(ctx, sess, "arnob-test123", nil)
	if err != nil {
		return nil, nil, err
	} else {
		fmt.Println("WHADHDUNKXD********************************")
	}
	return bucket, sess, nil
}

func write(ctx context.Context, bucket *blob.Bucket, key string) error {
	writeCtx, cancelWrite := context.WithCancel(ctx)
	defer cancelWrite()

	w, err := bucket.NewWriter(writeCtx, key, nil)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = fmt.Fprintln(w, "Hello, World!")
	if err != nil {
		return err
	}
	return nil
}

func read(ctx context.Context, bucket *blob.Bucket, key string) error {
	// Open the key "foo.txt" for reading with the default options.
	r, err := bucket.NewReader(ctx, key, nil)
	if err != nil {
		return err
	}
	defer r.Close()
	// Readers also have a limited view of the blob's metadata.
	fmt.Println("Content-Type:", r.ContentType())
	fmt.Println()
	// Copy from the reader to stdout.
	if _, err := io.Copy(os.Stdout, r); err != nil {
		return err
	}
	return nil
}

func put(sess *session.Session, key string) error {
	svc := s3.New(sess)
	fmt.Println("SOMETHING HERE 1")

	// Define the POSIX metadata
	metadata := map[string]*string{
		"x-amz-meta-mode": aws.String("0774"), // POSIX permissions
		"x-amz-meta-uid":  aws.String("1001"), // User ID
		"x-amz-meta-gid":  aws.String("1001"), // Group ID
	}

	directoryPaths := []string{
		"sunny-gui",
	}
	_ = directoryPaths
	//
	//for _, dir := range directoryPaths {
	//	_, err := svc.PutObject(&s3.PutObjectInput{
	//		Bucket: aws.String("appscode-testing"),
	//		Key:    aws.String(dir),
	//		Body:   strings.NewReader(""),
	//		//Metadata: metadata,
	//	})
	//	if err != nil {
	//		fmt.Println("ehehe")
	//		log.Fatal(err)
	//	}
	//	//err := putObjectWithMetadata(sess, "kubedb", dir, "", metadata)
	//	//if err != nil {
	//	//	log.Fatal(err)
	//	//}
	//}
	fmt.Println("WHAT THE HELL")

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:   aws.String("arnob-test123"),
		Key:      aws.String(key),
		Body:     strings.NewReader("Hello, World!"),
		Metadata: metadata,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Object uploaded successfully with POSIX metadata")
	return nil
}

func check(sess *session.Session, key string) error {
	svc := s3.New(sess)
	headObjectOutput, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String("arnob-test123"),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	// Print the metadata
	log.Println("Metadata:")
	for k, v := range headObjectOutput.Metadata {
		log.Printf("%s: %s\n", k, *v)
	}

	headObjectOutput, err = svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String("arnob-test123"),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	log.Println("Metadata for directory:")
	for k, v := range headObjectOutput.Metadata {
		log.Printf("%s: %s\n", k, *v)
	}
	return nil
}
