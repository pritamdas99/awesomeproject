package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/types"
	"log"
	"net/http"
	"time"
)

//func main() {
//	config := hazelcast.Config{}
//	cc := &config.Cluster
//	cc.Network.SetAddresses("hzkk.svc.default:5701")
//	cc.Discovery.UsePublicIP = true
//	ctx := context.TODO()
//	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
//	if err != nil {
//		fmt.Println("Failed to get the client ready**************************************************")
//		panic(err)
//	}
//	fmt.Println("Successful connection!")
//	fmt.Println("Starting to fill the map with random entries.")
//	m, err := client.GetMap(ctx, "map")
//	if err != nil {
//		fmt.Println("Failed to get map &&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")
//		panic(err)
//	}
//	for {
//		num := rand.Intn(100_000)
//		key := fmt.Sprintf("key-%d", num)
//		value := fmt.Sprintf("value-%d", num)
//		if _, err = m.Put(ctx, key, value); err != nil {
//			fmt.Println("ERR:", err.Error())
//		} else {
//			if mapSize, err := m.Size(ctx); err != nil {
//				fmt.Println("ERR:", err.Error())
//			} else {
//				fmt.Println("Current map size:", mapSize)
//			}
//		}
//	}
//}

//func main() {
//	config := hazelcast.Config{}
//	cc := &config.Cluster
//	cc.Network.SetAddresses("10.2.0.96:5701")
//	cc.Discovery.UsePublicIP = true
//	cc.Security.Credentials.Password = "password"
//	cc.Security.Credentials.Username = "user"
//	ctx := context.TODO()
//	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("what : ", client, err)
//	fmt.Println("adders : ", client.Name(), err)
//	fmt.Println("Successful connection!")
//	fmt.Println("Starting to fill the map with random entries.")
//	m, err := client.GetMap(ctx, "mapp")
//	if err != nil {
//		fmt.Println("Failed to get map &&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")
//		panic(err)
//	}
//	for {
//		//time.Sleep(time.Second * 30)
//		num := rand.Intn(100_000)
//		key := fmt.Sprintf("key-%d", num)
//		value := fmt.Sprintf("value-%d", num)
//		if _, err = m.Put(ctx, key, value); err != nil {
//			fmt.Println("ERR A:", err.Error())
//		} else {
//			if mapSize, err := m.Size(ctx); err != nil {
//				fmt.Println("ERR B:", err.Error())
//			} else {
//				fmt.Println("Current map size:", mapSize, m.Name())
//			}
//		}
//	}
//}

func main() {
	ctx := context.TODO()

	// create the default configuration
	config := hazelcast.Config{}
	// optionally set member addresses manually
	config.Cluster.Network.SetAddresses("hazelcast-hazelcast-enterprise:5701")
	config.Cluster.Network.ConnectionTimeout = types.Duration(time.Second * 10)
	config.Cluster.Security.Credentials.Password = "password"
	config.Cluster.Security.Credentials.Username = "user"
	// create and start the client with the configuration provider
	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	// handle client start error

	if err != nil {
		log.Fatal(err)
	}
	// get a map
	people, err := client.GetMap(ctx, "people")
	if err != nil {
		log.Fatal(err)
	}
	personName := "Jane Doe"
	// set a value in the map
	if err = people.Set(ctx, personName, 30); err != nil {
		log.Fatal(err)
	}
	// get a value from the map
	age, err := people.Get(ctx, personName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is %d years old.\n", personName, age)
	client.Shutdown(ctx)

	url := "http://hazelcast-hazelcast-enterprise.default.svc:5701/hazelcast/health"

	// Create the data to send in the request

	// Create a new HTTP POST request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	auth := "user:passowrd"
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	// Send the request
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding response body:", err)
		return
	}

	// Print the result as a map
	fmt.Println("Response as Map:", result)
	defer resp.Body.Close()
	// stop the client to release resources
}
