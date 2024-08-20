
[![Golang](https://img.shields.io/badge/Golang-fff.svg?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-purple?style=flat-square&logo=libreoffice)](LICENSE)
[![Latest Version](https://img.shields.io/github/v/tag/0x4f53/dnscovery?label=Version&style=flat-square&logo=semver)](https://github.com/0x4f53/dnscovery/releases)
[![Binaries](https://img.shields.io/badge/Binaries-Click%20Here-blue?style=flat-square&logo=dropbox)](.build/binaries/)

# 🌐 Dnscovery

<img src = preview.gif alt="dnscovery preview" width = "500dp">

A lightning-fast Golang tool to discover services embedded into DNS records


## 🚀 Features

- Takes just 1 second to resolve multiple domains**
- Queries multiple DNS servers concurrently
- More than 100 service signatures supported
- Easy to customize regexes and resolvers lists in YAML format
- Verbose JSON output for in-depth debugging

_** - depending on factors like internet speed, DNS server availability etc._


## 🖊️ Usage

```bash
Usage:
  dnscovery <domain1> <domain2>... [flags]

Flags:
  -o, --output string   Save output to file (in JSON format)
  -v, --verbose         Give extremely detailed information in output
```
### Examples
- Trying one domain
```bash
❯ ./dnscovery 0x4f.in

Reading resolvers...    [ 7 found! ]
Checking if online...   [ ✓ ONLINE ]

0x4f.in: OpenAI Domain, Cloudflare Mail, Google Workspace, Ethereum Name Service
```

- JSON output
```bash
❯ ./dnscovery 0x4f.in -o=output.json

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
❯ dnscovery nintendo.co.jp phase.dev huffpost.com redgear.com 0x4f.in lenovo.com apple.com microsoft.com netflix.com hackertyper.com tcl.com

Reading resolvers...    [ 7 found! ]
Checking if online...   [ ✓ ONLINE ]

apple.com: Apple, Facebook, Google Cloud Platform, Atlassian
0x4f.in: Ethereum Name Service, Google Workspace, Cloudflare Mail, OpenAI Domain
tcl.com: Google Cloud Platform
lenovo.com: Microsoft Office 365
microsoft.com: Microsoft Office 365, Microsoft Dynamics 365, Docusign
netflix.com: Dropbox, Apple, Docusign
nintendo.co.jp: Microsoft Office 365, Docusign, Adobe Creative Cloud, Google Cloud Platform, Apple
huffpost.com: Microsoft Office 365, Dropbox, Docusign, KnowBe4, Facebook, Google Cloud Platform
phase.dev: Google Cloud Platform, Gandi.net, Google Workspace
redgear.com: Google Workspace, Microsoft Office 365, Barracuda.com, Google Cloud Platform, Dropbox
hackertyper.com: Google Cloud Platform, Google Workspace
```

- Verbose mode
```bash
❯ ./dnscovery 0x4f.in -v

Reading resolvers...    [ 7 found! ]
Checking if online...   [ ✓ ONLINE ]

0x4f.in
  Resolved by: Control D (76.76.2.0) Cloudflare (1.1.1.1) Quad9 (9.9.9.9) OpenDNS (208.67.222.222) Google (8.8.4.4) Verisign (64.6.64.6)
  Services:
    OpenAI Domain
      openai-domain-verification=dv-ThXpvQCK0VDGRfFHh6hCP7cy
    Google Workspace
      v=spf1 include:_spf.mx.cloudflare.net include:_spf.google.com ~all
    Cloudflare Mail
      v=spf1 include:_spf.mx.cloudflare.net include:_spf.google.com ~all
    Ethereum Name Service
      ENS1 dnsname.ens.eth 0x6189345d91a667c4822A0afD7587a4994965a57C
```

## ⚙️ Building

To build this on your machine, you need to have Golang installed.
If you do, simply make build.sh executable and run it like so

```bash
chmod +x build.sh
./build.sh
```

## ⚙️ Installation
### Linux and macOS

Simply run the `./install.sh` script (don't 
have the time to put this on package managers)

```bash
chmod +x install.sh
sudo ./install.sh
```

And to uninstall

```bash
chmod +x uninstall.sh
sudo ./uninstall.sh
```

You can also find the binaries in [`.build/binaries`](.build/binaries/) if you want to directly run them
without installation

### Windows
You can find the exe files in [`.build/binaries`](.build/binaries/)


## ❓ Why I made this

I made this tool to check common services that multiple hosts use, by running it on a list of top 10,000 sites, 
for statistical purposes. This tool can also speed up a blue-teamer's inspection tasks or 
provide instant attack vectors for red-teamers to experiment with.


## 👍 Credits

- [NetSPI's Powershell scripts](https://github.com/NetSPI/PowerShell/blob/master/Resolve-DnsDomainValidationToken.ps1)

- [Google Dorks](https://www.freecodecamp.org/news/google-dorking-for-pentesters-a-practical-tutorial/) - good life skill to have

---

Copyright (c) 2024  Owais Shaikh

Licensed under the [MIT License](LICENSE)