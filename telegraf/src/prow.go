package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/smallfish/simpleyaml"
	"github.com/theckman/yacspin"
)

// type error interface {
// 	Error() string
// }

// type errorString struct {
// 	s string
// }

// func (e *errorString) Error() string {
// 	return e.s
// }

type Job struct {
	id            string
	state         string
	state_int     int
	flake_type    string
	log_url       string
	log_yaml      string
	log_artifacts string
	start_time    string
	end_time      string
	name          string
	pull_request  string
	job_type      string // [ pr, periodic ]
	cloud_profile string // remove
	log_job_name  string // [ operator-e2e-aws-periodic-slack, operator-e2e-gcp-periodic-slack, operator-e2e, etc.]
	test_type     string // [ unit, index, e2e ]
}

// type Flake struct {
// 	e2eNotStartedLog string // https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs/origin-ci-test/logs/periodic-ci-openshift-oadp-operator-master-4.8-operator-e2e-aws-periodic-slack/1493435130520276992/artifacts/operator-e2e-aws-periodic-slack/e2e/build-log.txt

// }

var all_jobs = make(map[string]Job)
var ErrorLogger *log.Logger

func main() {
	//error logging
	file, err := os.OpenFile("/tmp/prow_error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// get cli args
	var url_to_scrape string
	var print_for_human bool

	flag.StringVar(&url_to_scrape, "url_to_scrape", "https://prow.ci.openshift.org/?job=*oadp*", "prow url to scrape, e.g. ")
	flag.BoolVar(&print_for_human, "print_for_human", false, "print for a human, not influxdb!!!!!")

	flag.Parse()

	// start dem spinners

	// required if stopSpinner is called
	//var spinner *yacspin.Spinner
	if print_for_human {
		spinner, err := start_spinner()
		spinner = spinner // get around compile error
		if err != nil {
			log.Printf("spinner failed")
		}
	}

	// start web scraping
	start_geziyor(url_to_scrape)

	// This may not be required
	// stop spinner
	// if print_for_human {
	// 	stopSpinnerOnSignal(spinner)
	// }

	// print output
	if print_for_human {
		print_human(all_jobs)
	} else {
		print_db(all_jobs)
	}
}

func start_geziyor(url_to_scrape string) {
	geziyor.NewGeziyor(&geziyor.Options{
		LogDisabled: true,
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url_to_scrape, g.Opt.ParseFunc)
		},
		ParseFunc: getProwJobs,
	}).Start()
}

func start_spinner() (*yacspin.Spinner, error) {
	// meh have some fun
	cfg := yacspin.Config{
		Frequency:         500 * time.Millisecond,
		Writer:            nil,
		ShowCursor:        false,
		HideCursor:        false,
		SpinnerAtEnd:      false,
		CharSet:           yacspin.CharSets[59],
		Prefix:            " ",
		Suffix:            " ",
		SuffixAutoColon:   true,
		Message:           " Getting your jobs",
		ColorAll:          true,
		Colors:            []string{"fgYellow"},
		StopCharacter:     "✓",
		StopColors:        []string{"fgGreen"},
		StopMessage:       "done",
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
		StopFailMessage:   "failed",
		NotTTY:            false,
	}

	spinner, err := yacspin.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to make spinner from struct: %s", err)
	}

	err = spinner.Start()
	time.Sleep(1 * time.Second)
	// end fun
	return spinner, err

}

func stopSpinnerOnSignal(spinner *yacspin.Spinner) {
	// ensure we stop the spinner before exiting, otherwise cursor will remain
	// hidden and terminal will require a `reset`
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh

		spinner.StopFailMessage("interrupted")

		// ignoring error intentionally
		_ = spinner.StopFail()

		os.Exit(0)
	}()
}

func getProwJobs(g *geziyor.Geziyor, r *client.Response) {

	// to dump entire html body
	//fmt.Println(string(r.Body))
	rows := r.HTMLDoc.Find("#builds > tbody > tr")
	//log.Printf("length: %d", rows.Length())
	rows.Each(func(i int, s *goquery.Selection) {
		link := s.Find("td:nth-child(8) > a")
		my_url, _ := link.Attr("href")
		//log.Printf(my_url)

		u, _ := url.Parse(my_url)
		id := u.Path[strings.LastIndex(u.Path, "/")+1:]
		//log.Printf(id)

		this_job := Job{id, "", 4, u.String(), "", "", "", "", "", "not_found", "", "", "", "", ""}
		all_jobs[id] = this_job

	})
	getJobDetails(all_jobs)
	getYAMLDetails(all_jobs)
}

func getJobDetails(all_jobs map[string]Job) {
	log_yaml_base := "https://prow.ci.openshift.org"
	for id, job := range all_jobs {
		//log.Printf("%+v\n", job)
		//log.Printf(id)
		response, err := http.Get(job.log_url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		// Create a goquery document from the HTTP response
		document, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatal("Error loading HTTP response body. ", err)
		}

		//TODO
		// Count the number children in #links-card
		// if 3 periodic
		// if 5 pr
		fmt.Println("CHILDREN  $$$$$$$$$$")

		document.Find("#links-card > a").Map(func(i int, sel *goquery.Selection) string {
			return fmt.Sprintf("%d: %s", i+1, sel.Text())
		})

		// Get the Prow job YAML link
		document.Find("#links-card > a:nth-child(2)").Each(func(i int, s *goquery.Selection) {
			yaml_link, ok := s.Attr("href")
			if ok {
				job.log_yaml = log_yaml_base + yaml_link
				//log.Printf(job.log_yaml)
			}
		})

		// Get the Prow job Artifacts link
		document.Find("#links-card > a:nth-child(3)").Each(func(i int, s *goquery.Selection) {
			artifact_link, ok := s.Attr("href")
			if ok {
				job.log_artifacts = artifact_link
				//log.Printf(job.log_artifacts)
			}
		})

		// Get pull request
		document.Find("#links-card > a:nth-child(4)").Each(func(i int, s *goquery.Selection) {
			pull_request, ok := s.Attr("href")
			if ok {
				job.pull_request = pull_request
				//log.Printf(job.log_artifacts)
			}
		})

		all_jobs[id] = job
	}
}

func getYAMLDetails(all_jobs map[string]Job) {
	for id, job := range all_jobs {
		//log.Printf(id)
		response, err := http.Get(job.log_yaml)
		if err != nil {
			print_human_row(job)
		}
		defer response.Body.Close()
		yaml_data, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			print_human_row(job)
		}
		yaml, err := simpleyaml.NewYaml(yaml_data)
		if err != nil {
			print_human_row(job)
		}

		clusterProfile, err := yaml.Get("metadata").Get("labels").Get("ci-operator.openshift.io/cloud").String()
		if err != nil {
			print_human_row(job)
		}

		name, _ := yaml.Get("metadata").Get("annotations").Get("prow.k8s.io/job").String()
		if len(name) < 1 {
			print_human_row(job)
			name = "name_not_found"
			job.name = name
		}

		if strings.HasPrefix(name, "pull") {
			job.job_type = "pull"
		} else {
			job.job_type = "periodic"
		}

		//log.Printf("name:" + name)

		// Get Start / Stop time
		start, err := yaml.Get("status").Get("startTime").String()
		if err != nil {
			print_human_row(job)
		}
		end, err := yaml.Get("status").Get("completionTime").String()
		if err != nil {
			print_human_row(job)
		}

		if clusterProfile == "azure4" {
			clusterProfile = "azure"
		}
		job.cloud_profile = clusterProfile
		job.start_time = start
		job.end_time = end
		job.name = name

		status, err := yaml.Get("status").Get("state").String()
		if err != nil {
			print_human_row(job)
		}

		if status == "failure" {
			if job.job_type == "periodic" {

				// factor out this bit of code to a function
				// getLogFile(log string){

				//}
				buid_log_url := job.log_artifacts + "artifacts/operator-e2e-" + job.cloud_profile + "-periodic-slack/e2e"
				fmt.Println("@@@@@@@@@@@@" + buid_log_url)
				buildlog_response, err := http.Get(buid_log_url)
				if err != nil {
					print_human_row(job)
				}
				defer buildlog_response.Body.Close()

				buidlog, err := io.ReadAll(buildlog_response.Body)
				if err != nil {
					print_human_row(job)
				}
				fmt.Println(string(buidlog))

				if !strings.Contains(string(buidlog), "build_log.txt") {
					status = "flake"
				}

			}
			build_log_url := strings.TrimPrefix(job.log_url, "https://prow.ci.openshift.org/view/gs/")
			fmt.Println("https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs/" + build_log_url + "/artifacts/")

		}
		// get state
		//  0 = success, 1 = pending, 2 = failure 3 = aborted, 4 = unknown
		state_int := 10
		state := ""
		switch status {

		case "success":
			state_int = 0
			state = "success"
		case "pending":
			state_int = 1
			state = "pending"
		case "failure":
			state_int = 2
			state = "failure"
		case "aborted":
			state_int = 3
			state = "aborted"
		case "flake":
			state_int = 4
			state = "flake"
		default:
			state_int = 10
			state = "unknown"
		}

		job.state = state
		job.state_int = state_int

		// update object w/ success, failure status
		all_jobs[id] = job
	}
}

func print_human(all_jobs map[string]Job) {
	for _, my_job := range all_jobs {
		fmt.Printf("%+v\n\n\n", my_job)
	}
}

func print_human_row(my_job Job) {
	ErrorLogger.Printf("%+v\n", my_job)

}

func print_db(all_jobs map[string]Job) {
	for _, my_job := range all_jobs {

		// ensure all the rows have the required data for entry
		if my_job.end_time == "" {
			print_human_row(my_job)
			break
		}

		if my_job.name == "" {
			print_human_row(my_job)
			break
		}

		if my_job.state_int == 10 {
			// job state not known.
			// this also causes duplicate entries for build_id
			// the state is eventually written by prow and the job in question
			// will be captured, most likely as aborted.
			print_human_row(my_job)
			break
		}

		// datestamps
		st, _ := time.Parse(time.RFC3339, my_job.start_time)
		et, _ := time.Parse(time.RFC3339, my_job.end_time)
		duration := fmt.Sprint(et.Sub(st).Seconds())
		timestamp := fmt.Sprint(st.Unix() * 1000000000)
		timestamp_int, _ := strconv.Atoi(timestamp)
		if timestamp_int < 1 {
			fmt.Println("timestamp is wrong, break out")
			break
		}

		// log.Printf(my_job.start_time)
		// log.Printf(my_job.end_time)
		// log.Printf(st.String())
		// log.Printf(et.String())
		// log.Printf("%f", st.Unix())

		// influxdb line format
		// https://docs.influxdata.com/influxdb/v2.1/reference/syntax/line-protocol/

		build_string := "build," +
			"job_name=" + my_job.name +
			",build_id=" + strconv.Quote(my_job.id) +
			",pull_request=" + strconv.Quote(my_job.pull_request) +
			",start_time=" + my_job.start_time +
			",end_time=" + my_job.end_time +
			",duration=" + duration + //seconds
			",state_int=" + strconv.Itoa(my_job.state_int) +
			",state=" + my_job.state +
			" " + //space required for influxdb format
			"job_name=" + strconv.Quote(my_job.name) +
			",build_id=" + strconv.Quote(my_job.id) +
			",pull_request=" + strconv.Quote(my_job.pull_request) +
			",start_time=" + strconv.Quote(my_job.start_time) +
			",end_time=" + strconv.Quote(my_job.end_time) +
			",duration=" + duration +
			",state_int=" + strconv.Itoa(my_job.state_int) +
			",state=" + strconv.Quote(my_job.state) +
			",log=" + strconv.Quote(my_job.log_url) +
			" " +
			timestamp // this timestap is the job recorded timestamp in influx

		fmt.Println(build_string)
	}
}
