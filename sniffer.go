package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	outputFile    string
	outputFlagSet bool
	verbose       bool
	input         []string
)

var rootCmd = &cobra.Command{
	Use:   "dnsservices <domain1> <domain2>...",
	Short: "dnsservices",
	Long:  "dnsservices - discover service-related information from DNS records.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input = args
		if cmd.Flags().Changed("output") {
			outputFlagSet = true
		} else {
			outputFlagSet = false
		}
	},
}

func printServices(output Output) {

	if verbose {
		printProviders(output)
	}

	var services []string
	for _, item := range output.Answers {
		for _, record := range item.Records {
			for _, service := range record.Services {
				if verbose {
					service = fmt.Sprintf("%s (in %s record):\n\t%s", service, record.Type, record.Value)
				}
				services = append(services, service)
			}
		}
	}
	services = RemoveDuplicatesAndEmptyStrings(services)

	if verbose {
		fmt.Println()
	}
	fmt.Printf("Found services: ")
	if verbose {
		fmt.Println()
	}

	for index, record := range services {
		if verbose {
			fmt.Printf("  %s. %s\n", strconv.Itoa(index+1), record)
			continue
		}

		fmt.Printf("%s", record)

		if index != len(services)-1 {
			fmt.Printf(", ")
		}

	}

	fmt.Println()

}

func saveAsJson(output Output) {
	outputBytes, err := json.MarshalIndent(output, "", "  ")

	if err != nil {
		ErrorLog.Println(err)
		os.Exit(-1)
	}

	os.WriteFile(outputFile, outputBytes, 0644)
}

func printProviders(output Output) {

	var providers []string

	for _, item := range output.Answers {
		for _, record := range item.Records {
			if record.Value != "" {
				providers = append(providers, item.Resolver.Name+" ("+item.Resolver.IP+")")
			}
		}
	}

	providers = RemoveDuplicatesAndEmptyStrings(providers)

	fmt.Println("\nDNS providers with this host:")

	for index, item := range providers {
		fmt.Printf("  %s. %s\n", strconv.Itoa(index+1), item)
	}
}

func main() {

	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "output.json", "Save output to file (in JSON format)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Give extremely detailed information in output")
	rootCmd.Flags().BoolP("help", "h", false, "Help")

	if err := rootCmd.Execute(); err != nil {
		if len(os.Args) < 2 {
			ErrorLog.Println("Please provide at least one domain name as an input")
		}
		os.Exit(-2)
	}

	resolvers, err := GetResolvers()
	if err != nil {
		ErrorLog.Println(err)
		os.Exit(-1)
	}
	Resolvers = resolvers

	fmt.Println("Checking if online...")

	if !CheckInternet() {
		ErrorLog.Println("Couldn't connect to the internet. Please check your connection and try again!")
		os.Exit(-2)
	}

	fmt.Print("\t[ ✓ ONLINE ]\n")

	signatures, err := GetSignatures()
	if err != nil {
		ErrorLog.Println(err)
		os.Exit(-1)
	}
	Signatures = signatures

	for _, domain := range input {

		fmt.Println("Looking up '" + domain + "'...\t[ " + strconv.Itoa(len(resolvers)) + " resolvers found! ]")

		output, err := Dig(domain)

		if err != nil {
			ErrorLog.Println(err)
			os.Exit(-1)
		}

		if outputFlagSet {
			fmt.Println("Output saved to '" + outputFile + "'")
			saveAsJson(output)
			return
		}

		printServices(output)

	}

}
