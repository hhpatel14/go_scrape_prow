### Clone “go-scrape-prow” and cd go-scrape-prow/openshift
### Create new project  
	oc new-project go-scrape-prow

### BuildConfig/ImageStream of telegraf and grafanaSidecar
	oc apply -f buildConfig.yaml
	
### Start telegraf and grafanaSidecar build
	oc start-build telegraf
	oc start-build grafana-sidecar
	
Optional: Add webhook to your repo. Using buildConfig details copy webhook URL with secret 
### Deploy elements of go-scrape-prow 
	oc apply -f go-scrape-prow.yaml 
Influx, grafana, and telegraf pods should be up and running. 


#### Telegraf 
Logs should show “Successfully connected to output.influxdb”

To Debug: Access telegraf container terminal, and navigate to `/etc/telegraf/telegraf.conf` for the main configuration file and `/etc/telegraf/telegraf.d` for the directory of configuration files.

Example: 
```
	cd /etc/telegraf
	telegraf  --config telegraf.conf  --test
```
[List of commands and flags](https://docs.influxdata.com/telegraf/v1.21/administration/commands/)

#### Influx
To Debug: Access influx container terminal. 

Example: 
```
	$ influx
	  Connected to http://localhost:8086 version 1.7.1
	  InfluxDB shell version: 1.7.1
	  Enter an InfluxQL query
	  > 
```

Helpful commands: 
```
	show databases
	use telegraf
	show measurments
	select * from build . 
```
[InfluxDB command line interface](https://docs.influxdata.com/influxdb/v1.8/tools/shell/)

#### Grafana
Tools to export/import front-end changes:

Access grafanaSidecar contianer terminal
1. Create an api key using create-api-key.py  For example: `./create-api-key.py --url http://admin:admin@grafana:3000 —key-name NameofYourKey`
Use --help flag for more infomation 
2. Import/Export changes using import-grafana.py/export-grafana.py. For example: `./import-grafana.py --host http://grafana:3000  --key NameOfYourKey`

### Cron-job to renew Telegraf pod
Telegraf hits the system's max process limit after running for a while. This will result in exec-resource being temporarily unavailable. [Filed issues on telegraf repo](https://github.com/influxdata/telegraf/issues/3657) suggest increasing ulimit by updating `/etc/security/limits.conf` to solve ulimit exhaustion. However, openshift container host limits default pid_limit = 1024 set in crio.conf. [(More information on the issue)](https://bugzilla.redhat.com/show_bug.cgi?id=1869832). This can be overridden by [creating a `ContainerRuntimeConfig(ctrcfg)` to patch the maximum PIDs to avoid the exhaustion limit.](https://access.redhat.com/solutions/5597061)

This solution may break at some point depending on configuration, and it is not reliable when the usage of resources is unknown.  Therefore, renew the telegraf pod by creating a cron job to avoid resource exhaustion. 
	
	oc appy -f telegraf-cronJob.yaml
	
Note: Change `spec.schedule as required. [crontab.guru](https://crontab.guru/)
