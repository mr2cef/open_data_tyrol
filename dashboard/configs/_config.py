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
        host= f"""mongodb://{os.getenv("MONGO_DB_USER")}:{os.getenv("MONGO_DB_PASSWORD")}@{os.getenv("MONGO_DB_HOST").split('//')[-1]}""",
        database=os.getenv("MONGO_DB_DB"),
        collection=os.getenv("MONGO_DB_COLLECTION")
    ),
    "dash": dict(
        user_pw = {'dashuser': 'pw'}
    )
}





