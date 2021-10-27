import influxdb_client
import pandas as pd

from configs import config

class Measurements:

    def __init__(self) -> None:
        self.client = influxdb_client.InfluxDBClient(
            url=config["influx"]["host"], 
            token=config["influx"]["token"], 
            org=config["influx"]["org"]
        )
        self.query_api = self.client.query_api()
    

    def for_station(self, station_id: str, date_from: str, date_to: str, conv_win_size=1) -> pd.DataFrame:
        return self.query_api.query_data_frame(f"""
            from(bucket:"{config["influx"]["bucket"]}") 
                |> range(start: {date_from}, stop: {date_to}) 
                |> filter(fn: (r) => r["stationId"] == "{station_id}")
                |> aggregateWindow(every: {conv_win_size}h, fn: mean, createEmpty: false)
                |> yield(name: "mean")
        """)

