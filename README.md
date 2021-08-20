# Open Data Tyrol Collector
## Welcome

Open Data Tyrol collects all available data from the Hydrologischer Dienst Tirol such as water level of various rivers, temperature and precipitation. These data are avialable at `https://data.gv.at`. 

I created a go webserver written in `GO` which gets the data when called with a `/GET` to `http://127.0.0.1:8080/collect`. The data is stored in a `influxdb v2.0` hosted by docker (port `8086`).


## Installation
This is a Linux only install guide. 
### Prerequesites
Install `docker` and `docker-compose` on your system. Start your docker daemon.

### Set up InfluxDB
Create a directory for persisting your influxdb data. In the sequel I will refer to it as `<YOUR_INFLUX_DB_DIR>`. 

```
$ docker run -p 8086:8086 \
      -v <YOUR_INFLUX_DB_DIR>:/var/lib/influxdb2 \
      influxdb:2.0
```
Then go to http://127.0.0.1:8086. Set up user, password, organisation and initial bucket (We will not use this. Use a dummy name). When you are in the main menue. Create a new bucket (`Data>Buckets`) which we will then use to store the data in. Furthermore, under `Data>Tockens` you will find the access tocken for `I/O`. 

For the next setup you need the `Tocken` (`<YOUR_INFLUX_DB_TOCKEN>`), `Bucket` (`<YOUR_INFLUX_DB_BUCKET>`) and the `Organisation` (`<YOUR_INFLUX_DB_ORG>`) you use.


### Get Going
Get the repo at
```
git clone https://github.com/mr2cef/open_data_tyrol
```
Open `docker-compose.yaml` and fill in your `<YOUR_INFLUX_DB_*>` data. Start the services by
``` 
docker-compose up
```


## Getting the Data
You can now request to collect the data by 
```
curl http://127.0.0.1:8080/collect
```


## How to continue? 
Set up a cronjob which executes the `curl` from above in a regular time interval.