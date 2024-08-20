package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
)

var (
	outputFile    string
	outputFlagSet bool
	verbose       bool
	input         []string
)

var rootCmd = &cobra.Command{
	Use:   "dnscovery <domain1> <domain2>...",
	Short: "dnscovery",
	Long:  "dnscovery - discover service-related information from DNS records.",
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

func printBasic(output Output) {

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

	fmt.Printf(output.Host + ": ")

	for index, record := range services {

		fmt.Printf("%s", record)

		if index != len(services)-1 {
			fmt.Printf(", ")
		}

	}

	fmt.Println()

}

func printVerbose(outputList []Output) {
	for _, output := range outputList {
		fmt.Println(output.Host)
		fmt.Print("  Resolved by:")
		for _, answer := range output.Answers {
			fmt.Print(" " + answer.Resolver.Name + " (" + answer.Resolver.IP + ")")
		}
		var services []string
		fmt.Print("\n  Services:\n")
		for _, answer := range output.Answers {
			for _, record := range answer.Records {
				for _, service := range record.Services {
					serviceString := "    " + service + "\n      " + record.Value
					services = append(services, serviceString)
				}
			}
		}
		services = RemoveDuplicatesAndEmptyStrings(services)
		for _, service := range services {
			fmt.Println(service)
		}
		fmt.Println()
	}

}

func saveAsJson(output []Output) {
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

	fmt.Print("\nReading resolvers...\t[ " + strconv.Itoa(len(resolvers)) + " found! ]")

	fmt.Print("\nChecking if online...")

	if !CheckInternet() {
		fmt.Println()
		ErrorLog.Println("Couldn't connect to the internet. Please check your connection and try again!")
		os.Exit(-2)
	}

	fmt.Print("\t[ ✓ ONLINE ]\n\n")

	signatures, err := GetSignatures()
	if err != nil {
		ErrorLog.Println(err)
		os.Exit(-1)
	}

	Signatures = signatures

	var finalOutput []Output

	var wg sync.WaitGroup

	for _, domain := range input {
		shuffleResolvers(resolvers)
		Resolvers = resolvers

		wg.Add(1)
		go func(domain string) {
			defer wg.Done()

			output, err := Dig(domain)
			if err != nil {
				ErrorLog.Println(err)
				return
			}

			finalOutput = append(finalOutput, output)

			if !verbose {
				printBasic(output)
			}

		}(domain)
	}

	wg.Wait()

	if verbose {
		printVerbose(finalOutput)
	}

	if outputFlagSet {
		fmt.Println("Output saved to '" + outputFile + "'")
		saveAsJson(finalOutput)
	}

}
