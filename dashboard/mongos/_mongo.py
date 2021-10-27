import pymongo
from configs import config
import pandas as pd
import pyproj

__all__ = ["Stations"]

class Stations():

    def __init__(self) -> None:
        self.client = pymongo.MongoClient(config["mongo"]["host"])
        self.db = self.client[config["mongo"]["database"]]
        self.collection = self.db[config["mongo"]["collection"]]

    @staticmethod
    def translate_measurement(s: pd.Series) -> pd.Series:
        translator = {
            "pegel": "water level",
            "prec": "precipitation",
            "temp": "temperature"
        }
        return s.map(lambda x: translator[x])


    @classmethod
    def epsg2gps(cls, epsg: pd.DataFrame) -> pd.DataFrame:
        return_df = epsg.copy()
        df = epsg.apply(
            lambda x: pyproj.transform(
                pyproj.Proj(x["format"]),
                pyproj.Proj("epsg:4326"), # this is gps coordinates
                x["high"],
                x["right"]
                ),
            axis=1
        )
        gps =  pd.DataFrame(
            df.to_list(),
            columns=["lat", "lon"]
        )
        return_df["measurement"] = cls.translate_measurement(return_df["measurement"])
        return return_df.drop(columns=["right", "high", "format"]).join([gps]).sort_values(by=["measurement", "name"])

    def get_stations(self) -> pd.DataFrame:
        l = []
        for i in self.collection.find():
            l.append(i)
        df = pd.DataFrame(l)
        self.epsg2gps(df)
        return self.epsg2gps(df)