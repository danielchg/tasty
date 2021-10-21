/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

var result map[string]interface{}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List operators",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var operators []string
		kubeconfig, _ := os.LookupEnv("KUBECONFIG")
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Panicln("failed to create K8s config")
		}
		client, err := dynamic.NewForConfig(config)
		if err != nil {
			log.Panicln("Failed to create K8s clientset")
		}
		packagemanifests := schema.GroupVersionResource{Group: "packages.operators.coreos.com", Version: "v1", Resource: "packagemanifests"}
		list, err := client.Resource(packagemanifests).Namespace("openshift-marketplace").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for _, d := range list.Items {
			// namespace, _, err := unstructured.NestedString(d.Object, "metadata", "namespace")
			// if err != nil {
			//	log.Printf("Error getting namespace %v", err)
			//	continue
			//}
			operators = append(operators, d.GetName())
		}
		sort.Strings(operators)
		for _, operator := range operators {
			fmt.Println(operator)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}