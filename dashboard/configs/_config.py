import dotenv
import os

__all__ = ["config"]

dotenv.load_dotenv("../.env")

config = {
    "influx": dict(
        token=os.getenv("INFLUX_DB_TOKEN"),
        host=os.getenv("INFLUX_DB_HOST"),
        org=os.getenv("INFLUX_DB_ORG"),
        bucket=os.getenv("INFLUX_DB_BUCKET")
    ),
    "mongo": dict(
        host= os.getenv("MONGO_DB_HOST"),
        database=os.getenv("MONGO_DB_DB"),
        collection=os.getenv("MONGO_DB_COLLECTION")
    ),
    "dash": dict(
        user_pw = {'dashuser': 'pw'}
    )
}





