package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/miekg/dns"
	"gopkg.in/yaml.v3"
)

var resolversFile = "resolvers.yaml"
var signaturesFile = "signatures.yaml"

var ErrorLog = log.New(os.Stderr, "Failed: ", 0)

var Signatures []Service
var Resolvers []Resolver

type Service struct {
	Name      string
	Signature []regexp.Regexp
}

type Record struct {
	Services []string
	Type     string
	Hostname string
	Value    string
	Class    string
	TTL      string
	Opcode   string
	Rcode    string
}

type Resolver struct {
	Name string
	IP   string
}

type Answers struct {
	Resolver Resolver
	Records  []Record
}

type Output struct {
	Host    string
	Answers []Answers
}

var recordTypes = []uint16{
	dns.TypeTXT,
	dns.TypeA,
	dns.TypeAAAA,
	dns.TypeMX,
	dns.TypeCNAME,
	dns.TypeSRV,
}

func queryDNS(domain string, recordType uint16, resolver Resolver) ([]Record, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), recordType)
	m.RecursionDesired = true

	var answers []Record

	r, _, err := c.Exchange(m, resolver.IP+":53")
	if err != nil {
		//ErrorLog = log.New(os.Stderr, fmt.Sprintf("Failed (%s): ", resolver.Name), 0)
		//ErrorLog.Println(err)
		return answers, err
	}

	for _, ans := range r.Answer {

		var value string
		if txt, ok := ans.(*dns.TXT); ok {
			value = txt.Txt[0]
		}

		var serviceNames []string

		for _, service := range Signatures {
			for _, signature := range service.Signature {
				if signature.MatchString(value) {
					serviceNames = append(serviceNames, service.Name)
					break
				}
			}
		}

		record := Record{
			Services: serviceNames,
			Hostname: ans.Header().Name,
			Type:     dns.TypeToString[ans.Header().Rrtype],
			Class:    dns.ClassToString[ans.Header().Class],
			TTL:      strconv.FormatUint(uint64(ans.Header().Ttl), 10),
			Opcode:   dns.OpcodeToString[r.Copy().Opcode],
			Rcode:    dns.RcodeToString[r.Copy().Rcode],
			Value:    value,
		}

		if len(serviceNames) > 0 {
			answers = append(answers, record)
		}

	}

	return answers, nil
}

func Dig(domain string) (Output, error) {
	var output Output
	output.Host = domain

	var answers []Answers
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, record := range recordTypes {
		for _, resolver := range Resolvers {
			wg.Add(1)
			go func(resolver Resolver, record uint16) {
				defer wg.Done()
				data, err := queryDNS(domain, record, resolver)
				if err != nil {
					//ErrorLog.Println(err)
				}
				mu.Lock()
				if len(data) > 0 {
					answers = append(answers,
						Answers{
							Resolver: resolver,
							Records:  data,
						},
					)
				}
				mu.Unlock()
			}(resolver, record)
		}
	}

	wg.Wait()
	output.Answers = answers

	return output, nil
}

func pwd() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func GetResolvers() ([]Resolver, error) {
	resolversFile = pwd() + "/" + resolversFile
	yamlFile, err := os.ReadFile(resolversFile)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	var data map[string][]string

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	var resolvers []Resolver

	for name, ips := range data {
		for _, ip := range ips {
			resolvers = append(resolvers, Resolver{Name: name, IP: ip})
		}
	}

	return resolvers, nil

}

func GetSignatures() ([]Service, error) {
	signaturesFile = pwd() + "/" + signaturesFile
	yamlFile, err := os.ReadFile(signaturesFile)
	if err != nil {
		log.Fatalf("Error reading signatures: %v", err)
	}

	var tempServices map[string][]string
	err = yaml.Unmarshal(yamlFile, &tempServices)
	if err != nil {
		log.Fatalf("Error unmarshalling signatures: %v", err)
	}

	var services []Service
	for name, patterns := range tempServices {
		var compiledPatterns []regexp.Regexp
		for _, pattern := range patterns {
			compiled, err := regexp.Compile(pattern)
			if err != nil {
				log.Fatalf("Error compiling regex: %v", err)
			}
			compiledPatterns = append(compiledPatterns, *compiled)
		}
		services = append(services, Service{
			Name:      name,
			Signature: compiledPatterns,
		})
	}

	return services, err

}

func RemoveDuplicatesAndEmptyStrings(input []string) []string {
	uniqueStrings := make(map[string]bool)
	var result []string

	for _, str := range input {
		if str != "" && !uniqueStrings[str] {
			uniqueStrings[str] = true
			result = append(result, str)
		}
	}

	return result
}

func CheckInternet() bool {
	c := new(dns.Client)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("."), dns.TypeA)
	m.RecursionDesired = true
	_, _, err := c.Exchange(m, Resolvers[0].IP+":53")
	return err == nil
}
