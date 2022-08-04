package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/test-infra/prow/apis/prowjobs/v1"
	yaml3 "sigs.k8s.io/yaml"
)

func Test_isFlake(t *testing.T) {
	type test struct {
		Name    string  `yaml:"name"`
		Mockjob TestJob `yaml:"arg"`
		Want    bool    `yaml:"want"`
	}

	type Tests struct {
		ATest []test `yaml:"Data_to_Test"`
	}

	yamlFile, err := ioutil.ReadFile("testData/isFlakeData.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var testdata Tests
	err = yaml.Unmarshal(yamlFile, &testdata)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, tt := range testdata.ATest {
		t.Run(tt.Name, func(t *testing.T) {
			fmt.Println(tt.Name)

			if strings.HasPrefix(tt.Name, "Negative") {
				unitTestfile, err := os.OpenFile("testData/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					fmt.Println("err in making logs.txt ")
				}

				ErrorLogger = log.New(unitTestfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
				job := TestJobTOJob(tt.Mockjob)
				isFlake(job)
				if _, err := os.Stat("testData/logs.txt"); err != nil {
					t.Errorf("Error log failed for flake ")
				} else {
					os.Remove("testData/logs.txt")
				}
			} else {
				job := TestJobTOJob(tt.Mockjob)

				if got := isFlake(job); got != tt.Want {
					t.Errorf("isFlake() = %v, want %v", got, tt.Want)
				}
			}
		})
	}
}

func Test_isDeprovisioningFlake(t *testing.T) {
	type test struct {
		Name    string  `yaml:"name"`
		Mockjob TestJob `yaml:"arg"`
		Want    bool    `yaml:"want"`
	}

	type Tests struct {
		ATest []test `yaml:"Data_to_Test"`
	}

	yamlFile, err := ioutil.ReadFile("testData/isDeprovisioningFlakeData.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var testdata Tests
	err = yaml.Unmarshal(yamlFile, &testdata)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, tt := range testdata.ATest {
		t.Run(tt.Name, func(t *testing.T) {
			fmt.Println(tt.Name)

			if strings.HasPrefix(tt.Name, "Invalid") {
				unitTestfile, err := os.OpenFile("testData/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					fmt.Println("err in making logs.txt ")
				}

				ErrorLogger = log.New(unitTestfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
				job := TestJobTOJob(tt.Mockjob)
				isDeprovisioningFlake(job)
				if _, err := os.Stat("testData/logs.txt"); err != nil {
					t.Errorf("Error log failed for flake ")
				} else {
					os.Remove("testData/logs.txt")
				}
			} else {
				job := TestJobTOJob(tt.Mockjob)

				if got := isDeprovisioningFlake(job); got != tt.Want {
					t.Errorf("isDeprovisioningFlake() = %v, want %v", got, tt.Want)
				}
			}
		})
	}
}

func Test_getJobType(t *testing.T) {
	type test struct {
		Name    string  `yaml:"name"`
		Mockjob TestJob `yaml:"arg"`
		Want    JobType `yaml:"want"`
	}

	type Tests struct {
		ATest []test `yaml:"Data_to_Test"`
	}

	yamlFile, err := ioutil.ReadFile("testData/getJobTypeData.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var testdata Tests
	err = yaml.Unmarshal(yamlFile, &testdata)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, tt := range testdata.ATest {
		t.Run(tt.Name, func(t *testing.T) {
			job := TestJobTOJob(tt.Mockjob)

			if got := getJobType(job); got != tt.Want {
				t.Errorf("getJobType() = %v, want %v", got, tt.Want)
			}
		})
	}
}

func Test_nameJob(t *testing.T) {
	type test struct {
		Name     string  `yaml:"name"`
		FileName string  `yaml:"filename"`
		Mockjob  TestJob `yaml:"arg"`
		Want     string  `yaml:"want"`
	}

	type Tests struct {
		ATest []test `yaml:"Data_to_Test"`
	}

	yamlFile, err := ioutil.ReadFile("testData/nameJobData.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var testdata Tests
	err = yaml.Unmarshal(yamlFile, &testdata)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, tt := range testdata.ATest {
		yaml_data, err := ioutil.ReadFile(tt.FileName)
		if err != nil {
			fmt.Println("err in file  ")
		}

		prow := v1.ProwJob{}
		unmarshalErr := yaml3.Unmarshal(yaml_data, &prow)
		if unmarshalErr != nil {
			fmt.Println("error in unmarshal")
		}

		t.Run(tt.Name, func(t *testing.T) {
			if strings.HasPrefix(tt.Name, "Negative") {
				unitTestfile, err := os.OpenFile("testData/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					fmt.Println("err in making logs.txt ")
				}

				ErrorLogger = log.New(unitTestfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
				job := TestJobTOJob(tt.Mockjob)
				nameJob(prow, job)
				if _, err := os.Stat("testData/logs.txt"); err != nil {
					t.Errorf("Error log failed for Empty name")
				} else {
					os.Remove("testData/logs.txt")
				}
			} else {
				job := TestJobTOJob(tt.Mockjob)

				if got := nameJob(prow, job); got != tt.Want {
					t.Errorf("nameJob() = %v, want %v", got, tt.Want)
				}
			}
		})
	}
}

func Test_startTime(t *testing.T) {
	type test struct {
		Name     string  `yaml:"name"`
		FileName string  `yaml:"filename"`
		Mockjob  TestJob `yaml:"arg"`
		Want     string  `yaml:"want"`
	}

	type Tests struct {
		ATest []test `yaml:"Data_to_Test"`
	}

	yamlFile, err := ioutil.ReadFile("testData/startTimeData.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var testdata Tests
	err = yaml.Unmarshal(yamlFile, &testdata)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, tt := range testdata.ATest {
		yaml_data, err := ioutil.ReadFile(tt.FileName)
		if err != nil {
			fmt.Println("err in file  ")
		}

		prow := v1.ProwJob{}
		unmarshalErr := yaml3.Unmarshal(yaml_data, &prow)
		if unmarshalErr != nil {
			fmt.Println("error in unmarshal")
		}

		t.Run(tt.Name, func(t *testing.T) {
			if strings.HasPrefix(tt.Name, "Negative") {
				unitTestfile, err := os.OpenFile("testData/badData_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					fmt.Println("err in making badData_logs.txt for startTime")
				}

				ErrorLogger = log.New(unitTestfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
				job := TestJobTOJob(tt.Mockjob)
				getStartTime(prow, job)
				if _, err := os.Stat("testData/badData_logs.txt"); err != nil {
					t.Errorf("Error log failed for start Time")
				} else {
					os.Remove("testData/badData_logs.txt")
				}
			} else {
				job := TestJobTOJob(tt.Mockjob)

				if got := getStartTime(prow, job); got != tt.Want {
					t.Errorf("getstartTime() = %v, want %v", got, tt.Want)
				}
			}
		})
	}
}

func Test_getStatus(t *testing.T) {
	type test struct {
		Name     string  `yaml:"name"`
		FileName string  `yaml:"filename"`
		Mockjob  TestJob `yaml:"arg"`
		Want     string  `yaml:"want"`
	}

	type Tests struct {
		ATest []test `yaml:"Data_to_Test"`
	}

	yamlFile, err := ioutil.ReadFile("testData/getStatusData.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var testdata Tests
	err = yaml.Unmarshal(yamlFile, &testdata)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, tt := range testdata.ATest {
		yaml_data, err := ioutil.ReadFile(tt.FileName)
		if err != nil {
			fmt.Println("err in file  ")
		}

		prow := v1.ProwJob{}
		unmarshalErr := yaml3.Unmarshal(yaml_data, &prow)
		if unmarshalErr != nil {
			fmt.Println("error in unmarshal")
		}

		t.Run(tt.Name, func(t *testing.T) {
			if strings.HasPrefix(tt.Name, "Negative") {
				unitTestfile, err := os.OpenFile("testData/badData_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					fmt.Println("err in making badData_logs.txt for getStatus")
				}

				ErrorLogger = log.New(unitTestfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
				job := TestJobTOJob(tt.Mockjob)
				getStatus(prow, job)
				if _, err := os.Stat("testData/badData_logs.txt"); err != nil {
					t.Errorf("Error log failed for geStatus")
				} else {
					os.Remove("testData/badData_logs.txt")
				}
			} else {
				job := TestJobTOJob(tt.Mockjob)

				if got := getStatus(prow, job); got != tt.Want {
					t.Errorf("getStatus() = %v, want %v", got, tt.Want)
				}
			}
		})
	}
}

func Test_endTime(t *testing.T) {
	type test struct {
		Name     string  `yaml:"name"`
		FileName string  `yaml:"filename"`
		Mockjob  TestJob `yaml:"arg"`
		Want     string  `yaml:"want"`
	}

	type Tests struct {
		ATest []test `yaml:"Data_to_Test"`
	}

	yamlFile, err := ioutil.ReadFile("testData/endTimeData.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var testdata Tests
	err = yaml.Unmarshal(yamlFile, &testdata)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, tt := range testdata.ATest {
		yaml_data, err := ioutil.ReadFile(tt.FileName)
		if err != nil {
			fmt.Println("err in file  ")
		}

		prow := v1.ProwJob{}
		unmarshalErr := yaml3.Unmarshal(yaml_data, &prow)
		if unmarshalErr != nil {
			fmt.Println("error in unmarshal")
		}

		t.Run(tt.Name, func(t *testing.T) {
			if strings.HasPrefix(tt.Name, "Negative") {
				unitTestfile, err := os.OpenFile("testData/badData_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					fmt.Println("err in making badData_logs.txt for endTime")
				}

				ErrorLogger = log.New(unitTestfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
				job := TestJobTOJob(tt.Mockjob)
				getEndTime(prow, job)
				if _, err := os.Stat("testData/badData_logs.txt"); err != nil {
					t.Errorf("Error log failed for end Time")
				} else {
					os.Remove("testData/badData_logs.txt")
				}
			} else {
				job := TestJobTOJob(tt.Mockjob)

				if got := getEndTime(prow, job); got != tt.Want {
					t.Errorf("getstartTime() = %v, want %v", got, tt.Want)
				}
			}
		})
	}
}

func Test_cluster_profile(t *testing.T) {
	type test struct {
		Name     string  `yaml:"name"`
		FileName string  `yaml:"filename"`
		Mockjob  TestJob `yaml:"arg"`
		Want     string  `yaml:"want"`
	}

	type Tests struct {
		ATest []test `yaml:"Data_to_Test"`
	}

	yamlFile, err := ioutil.ReadFile("testData/cluster_profileData.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var testdata Tests
	err = yaml.Unmarshal(yamlFile, &testdata)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, tt := range testdata.ATest {
		yaml_data, err := ioutil.ReadFile(tt.FileName)
		if err != nil {
			fmt.Println("err in file  ")
		}

		prow := v1.ProwJob{}
		unmarshalErr := yaml3.Unmarshal(yaml_data, &prow)
		if unmarshalErr != nil {
			fmt.Println("error in unmarshal")
		}

		t.Run(tt.Name, func(t *testing.T) {
			if strings.HasPrefix(tt.Name, "Negative") {
				unitTestfile, err := os.OpenFile("testData/badData_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					fmt.Println("err in making badData_logs.txt for clusterProfile  ")
				}

				ErrorLogger = log.New(unitTestfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
				job := TestJobTOJob(tt.Mockjob)
				getClusterProfile(prow, job)
				if _, err := os.Stat("testData/badData_logs.txt"); err != nil {
					t.Errorf("Error log failed for cluster Profile")
				} else {
					os.Remove("testData/badData_logs.txt")
				}
			} else {
				job := TestJobTOJob(tt.Mockjob)

				if got := getClusterProfile(prow, job); got != tt.Want {
					t.Errorf("getClusterProfile() = %v, want %v", got, tt.Want)
				}
			}
		})
	}
}

func Test_getStateInt(t *testing.T) {
	type test struct {
		Name   string `yaml:"name"`
		Status string `yaml:"status"`
		Want   int    `yaml:"want"`
	}

	type Tests struct {
		ATest []test `yaml:"Data_to_Test"`
	}

	yamlFile, err := ioutil.ReadFile("testData/getStateIntData.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var testdata Tests
	err = yaml.Unmarshal(yamlFile, &testdata)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, tt := range testdata.ATest {
		t.Run(tt.Name, func(t *testing.T) {
			if got := getStateInt(tt.Status); got != tt.Want {
				t.Errorf("getStateInt() = %v, want %v", got, tt.Want)
			}
		})
	}
}

func Test_getTargetTestName(t *testing.T) {
	type test struct {
		Name     string  `yaml:"name"`
		FileName string  `yaml:"filename"`
		Mockjob  TestJob `yaml:"arg"`
		Want     string  `yaml:"want"`
	}

	type Tests struct {
		ATest []test `yaml:"Data_to_Test"`
	}

	yamlFile, err := ioutil.ReadFile("testData/getTargetTestName.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var testdata Tests
	err = yaml.Unmarshal(yamlFile, &testdata)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, tt := range testdata.ATest {
		yaml_data, err := ioutil.ReadFile(tt.FileName)
		if err != nil {
			fmt.Println("err in file for yaml_data ")
		}

		prow := v1.ProwJob{}
		unmarshalErr := yaml3.Unmarshal(yaml_data, &prow)
		if unmarshalErr != nil {
			fmt.Println("error in unmarshal")
		}

		t.Run(tt.Name, func(t *testing.T) {
			if strings.HasPrefix(tt.Name, "Negative") {
				unitTestfile, err := os.OpenFile("testData/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					fmt.Println("err in making logs.txt ")
				}

				ErrorLogger = log.New(unitTestfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
				job := TestJobTOJob(tt.Mockjob)
				getTargetTestName(prow, job)
				if _, err := os.Stat("testData/logs.txt"); err != nil {
					t.Errorf("Error log failed for getTargetTestName ")
				} else {
					os.Remove("testData/logs.txt")
				}
			} else {
				job := TestJobTOJob(tt.Mockjob)

				if got := getTargetTestName(prow, job); got != tt.Want {
					t.Errorf("getTargetTestName() = %v, want %v", got, tt.Want)
				}
			}
		})
	}
}
