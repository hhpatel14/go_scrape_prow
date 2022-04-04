# go_scrape_prow
###### tags: `Documentation`

## Commands to run unit tests 
```
cd go_scrape_prow/telegraf/src
go build .
go test -v -coverprofile cover.out ./
go tool cover -html=cover.out -o cover.html
open cover.html
```

