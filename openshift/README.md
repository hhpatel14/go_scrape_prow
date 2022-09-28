Clone “go-scrape-prow” and cd go-scrape-prow/openshift
Create new project  oc new-project go-scrape-prow
Create Source Crecret 
Add a new SSH key to your Github. Adding a new SSH key to your GitHub account
Create a new source secret in your cluster 
Click on workloads → Secrets
Navigate to project “go-scrape-prow”  and click on Create (button on the upper right corner) 
Paste SSH and click Create
 

Build telegraf and grafanaSidecar
oc apply -f buildConfig.yaml
Optional: Add webhook to your repo. Use buildConfig details and, copy URL with details 
Deploy elements of go-scrape-prow 
oc apply -f go-scrape-prow.yaml 
Influx, grafana, and telegraf pod should be up and running. 
	Telegraf log should show “Successfully connected to output.influxdb”

