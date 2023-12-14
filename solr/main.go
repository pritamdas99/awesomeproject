package main

import (
	"fmt"
	xmlpatch "github.com/maxyurk/go-xml-patch"
	"os"
	"reflect"
)

var keys = map[string]string{
	"maxBooleanClauses": "solr.max.booleanClauses",
	"sharedLib":         "solr.sharedLib",
	"hostPort":          "solr.port.advertise",
}

var shardHandlerFactory = map[string]interface{}{
	"socketTimeout": 600000,
	"connTimeout":   60000,
}

var solrCloud = map[string]interface{}{
	"host":                     "",
	"hostPort":                 80,
	"hostContext":              "solr",
	"genericCoreNodeNames":     true,
	"zkClientTimeout":          30000,
	"distribUpdateSoTimeout":   600000,
	"distribUpdateConnTimeout": 60000,
	"zkCredentialsProvider":    "org.apache.solr.common.cloud.DefaultZkCredentialsProvider",
	"zkACLProvider":            "org.apache.solr.common.cloud.DefaultZkACLProvider",
}

var solr map[string]interface{} = map[string]interface{}{
	"maxBooleanClauses":   1024,
	"sharedLib":           "",
	"allowPaths":          "",
	"solrCloud":           solrCloud,
	"shardHandlerFactory": shardHandlerFactory,
}

func getXMLConfigElement(name string, key string, value interface{}, kind string) string {
	if key == "" {
		key = name
	}
	pp := reflect.TypeOf(value).Kind()
	if pp == reflect.Int {
		return fmt.Sprintf("<%s name=\"%s\">${%s:%d}</%s>\n", kind, name, key, value, kind)
	} else if pp == reflect.String {
		return fmt.Sprintf("<%s name=\"%s\">${%s:%s}</%s>\n", kind, name, key, value, kind)
	} else if pp == reflect.Bool {
		return fmt.Sprintf("<%s name=\"%s\">${%s:%t}</%s>\n", kind, name, key, value, kind)
	}
	return ""
}

func intend(ss string, level int) string {
	//fmt.Println(level)
	level *= 2
	for level > 0 {
		level--
		ss += " "
	}
	return ss
}

func rec(mp map[string]interface{}, level int) string {
	//fmt.Println("startit", mp)
	ss := ""
	for x, y := range mp {
		//fmt.Println("eleents", x)
		kind := reflect.TypeOf(y).Kind()
		if kind == reflect.Int {
			val := y.(int)
			ss = intend(ss, level)
			ss = ss + getXMLConfigElement(x, keys[x], val, "int")
		} else if kind == reflect.String {
			val := y.(string)
			ss = intend(ss, level)
			ss = ss + getXMLConfigElement(x, keys[x], val, "str")
		} else if kind == reflect.Bool {
			val := y.(bool)
			ss = intend(ss, level)
			ss = ss + getXMLConfigElement(x, keys[x], val, "bool")
		} else {
			ss = intend(ss, level)
			if x == "shardHandlerFactory" {
				ss = ss + fmt.Sprintf("<%s name=\"shardHandlerFactory\" class=\"HttpShardHandlerFactory\">\n", x)
			} else {
				ss = ss + fmt.Sprintf("<%s>\n", x)
			}
			//	fmt.Println(x, y)
			v, ok := y.(map[string]interface{})
			//fmt.Println("get it", v)
			if !ok {
				fmt.Println("failed to decode")
			}
			ss += rec(v, level+1)
			ss = intend(ss, level)
			ss = ss + fmt.Sprintf("</%s>\n", x)
		}
	}
	return ss
}

func main() {

	//keys = map[string]string{
	//	"maxBooleanClauses": "solr.max.booleanClauses",
	//	"sharedLib":         "solr.sharedLib",
	//	"hostPort":          "solr.port.advertise",
	//}

	//shardHandlerFactory = map[string]interface{}{
	//	"socketTimeout": 600000,
	//	"connTimeout":   60000,
	//}
	//solrCloud = map[string]interface{}{
	//	"host":                     "",
	//	"hostPort":                 80,
	//	"hostContext":              "solr",
	//	"genericCoreNodeNames":     true,
	//	"zkClientTimeout":          30000,
	//	"distribUpdateSoTimeout":   600000,
	//	"distribUpdateConnTimeout": 60000,
	//	"zkCredentialsProvider":    "org.apache.solr.common.cloud.DefaultZkCredentialsProvider",
	//	"zkACLProvider":            "org.apache.solr.common.cloud.DefaultZkACLProvider",
	//}
	//
	//solr := map[string]interface{}{
	//	"maxBooleanClauses":   1024,
	//	"sharedLib":           "",
	//	"allowPaths":          "",
	//	"solrCloud":           solrCloud,
	//	"shardHandlerFactory": shardHandlerFactory,
	//}

	ss := "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n<solr>\n" + "  <str name=\"coreRootDirectory\">var/solr/data</str>\n" + rec(solr, 1)
	ss += "  <metrics enabled=\"${metricsEnabled:true}\"/>\n"
	ss += "</solr>\n"
	fmt.Println(ss)
	sd := "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n<solr>\n  <str name=\"key\">val</str>\n  <solrCloud>\n     <str name=\"host\">${solr.host:}</str>\n  </solrCloud>\n</solr>"
	fmt.Println(sd)

	target, err := os.ReadFile("target.xml")
	if err != nil {
		fmt.Println("error 1", err)
	}
	diff, err := os.ReadFile("diff.xml")
	if err != nil {
		fmt.Println("error 2", err)
	}
	patch, err := xmlpatch.Patch(target, diff)
	if err != nil {
		panic(err)
	}
	fmt.Println(patch)
}
