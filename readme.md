
[![Golang](https://img.shields.io/badge/Golang-fff.svg?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-purple?style=flat-square&logo=libreoffice)](LICENSE)
[![Latest Version](https://img.shields.io/github/v/tag/0x4f53/dnsservices?label=Version&style=flat-square&logo=semver)](https://github.com/0x4f53/subs/releases)
[![Binaries](https://img.shields.io/badge/Binaries-Click%20Here-blue?style=flat-square&logo=dropbox)](.build/binaries/)

# DNSservices

![DNSservices GIF](terminal.gif)

A lightning-fast Golang tool to discover services embedded into DNS records

## Features
- Takes just 2 seconds to resolve a domain**
- Queries multiple DNS servers concurrently
- More than 100 service signatures supported!
- Easy to customize regexes and resolvers lists in YAML format
- Verbose JSON output for in-depth debugging

_**- depending on factors like internet speed, DNS server availability etc._

## Usage

```bash
Usage:
  dnsservices <domain1> <domain2>... [flags]

Flags:
  -h, --help            Help
  -o, --output string   Save output to file (in JSON format) (default "output.json")
  -v, --verbose         Give extremely detailed information in output
```
### Examples
- Trying one domain
```bash
❯ ./dnsservices 0x4f.in
Checking if online...   [ ✓ ONLINE ]
Looking up '0x4f.in'... [ 7 resolvers found! ]
Found services: OpenAI Domain, Ethereum Name Service, Cloudflare Mail, Google Workspace
```

- JSON output
```bash
❯ ./dnsservices 0x4f.in -o=output.json
Checking if online...   [ ✓ ONLINE ]
Looking up '0x4f.in'... [ 7 resolvers found! ]
Output saved to 'output.json'
❯ cat output.json
{
  "Host": "0x4f.in",
  "Answers": [
    {
      "Resolver": {
        "Name": "Google",
        "IP": "8.8.4.4"
      },
      "Records": [
        {
          "Services": [
            "Ethereum Name Service"
          ],
          "Type": "TXT",
          "Hostname": "0x4f.in.",
          "Value": "ENS1 dnsname.ens.eth 0x6189345d91a667c4822A0afD7587a4994965a57C",
    ...
```

- Trying multiple domains
```bash
❯ ./dnsservices 0x4f.in blackhat.com zoho.com
Checking if online...   [ ✓ ONLINE ]

Looking up '0x4f.in'... [ 7 resolvers found! ]
Found services: Ethereum Name Service, Cloudflare Mail, Google Workspace, OpenAI Domain

Looking up 'blackhat.com'...    [ 7 resolvers found! ]
Found services: Google Search Console, Microsoft Office 365, Twilio SendGrid, Google Workspace

Looking up 'zeron.one'...       [ 7 resolvers found! ]
Found services: Microsoft Office 365, Pinterest, Google Search Console, MongoDB, Atlassian
```

- Verbose mode
```bash
❯ ./dnsservices 0x4f.in -v
Checking if online...   [ ✓ ONLINE ]
Looking up '0x4f.in'... [ 7 resolvers found! ]

DNS providers that responded:
  1. Quad9 (9.9.9.9)
  2. Google (8.8.4.4)
  4. Cloudflare (1.1.1.1)

Found services: 
  1. Ethereum Name Service (in TXT record):
        ENS1 dnsname.ens.eth 0x6189345d91a667c4822A0afD7587a4994965a57C
  2. OpenAI Domain (in TXT record):
        openai-domain-verification=dv-ThXpvQCK0VDGRfFHh6hCP7cy
  3. Cloudflare Mail (in TXT record):
        v=spf1 include:_spf.mx.cloudflare.net include:_spf.google.com ~all
  4. Google Workspace (in TXT record):
        v=spf1 include:_spf.mx.cloudflare.net include:_spf.google.com ~all
```

## Building

To build this on your machine, you need to have Golang installed.
If you do, simply make build.sh executable and run it like so

```bash
❯ chmod +x build.sh
❯ ./build.sh
```

## Installation
### Linux and macOS

Simply run the `./install.sh` script (don't 
have the time to put this on package managers)

```bash
chmod +x install.sh
./install.sh
```

You can also find the binaries in [`.build/binaries`](.build/binaries/) if you want to directly run them
without installation

### Windows
You can find the exe files in [`.build/binaries`](.build/binaries/)


---

Copyright (c) 2024  Owais Shaikh

Licensed under the [MIT License](LICENSE)