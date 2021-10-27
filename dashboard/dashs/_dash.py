import dash
from dash.dcc.Graph import Graph
from dash.exceptions import PreventUpdate
from dash import dcc, html
import dash_auth
import plotly.express as px
import plotly
import pandas as pd
import numpy as np
import json
from mongos import Stations
from influxs import Measurements

from configs import config

__all__ = ["app"]

DEFAULT_STATION = "tirTemp103143"

app = dash.Dash()
auth = dash_auth.BasicAuth(
    app,
    config["dash"]["user_pw"]
)


app.layout = html.Div([
    html.H1(
            id='header',
            children='Open Data visulaization Tyrol'
        ),
    dcc.Graph(id='graph-map'),
    dcc.Dropdown(id='dropdown-stations'),
    dcc.Graph(id='graph-time-series'),
    dcc.Store(
        id="data-stations", 
        storage_type='session'
    )
], className='container')

@app.callback(
    dash.dependencies.Output("data-stations", "data"),
    [dash.dependencies.Input("header", "children")],
    [dash.dependencies.State("data-stations", "data"),]
)
def get_data(_, d: str):
    if d is not None:
        raise PreventUpdate
    s = Stations()
    stats: pd.DataFrame = s.get_stations()
    return stats.to_json()

@app.callback(
    [
        dash.dependencies.Output('dropdown-stations', 'options'),
        dash.dependencies.Output('dropdown-stations', 'value')
    ],
    [
        dash.dependencies.Input("data-stations", "data")
    ]
)
def fill_stations_dropdown(stations: str):
    if stations is None:
        raise PreventUpdate
    df = pd.DataFrame.from_dict(json.loads(stations))
    options = [{"label": (row["name"] + ": " + row["measurement"]), "value": row["_id"]} for index, row in df.iterrows()]
    value = DEFAULT_STATION
    return options, value


@app.callback(
    [
        dash.dependencies.Output('graph-map', 'figure'),
        dash.dependencies.Output('graph-time-series', 'figure'),
    ],
    [
        dash.dependencies.Input('dropdown-stations', 'value')
    ],
    [
        dash.dependencies.State("data-stations", "data")
    ],
    )
def update_graph(station_id: str, stations: str):
    if stations is None or station_id is None:
        raise PreventUpdate
    df = pd.DataFrame.from_dict(json.loads(stations))
    map = px.scatter_mapbox(
        data_frame=df, 
        lat="lat", 
        lon="lon", 
        hover_name="name",
        hover_data=["measurement"],
        color="measurement",
        zoom=7, 
        height=600, 
        width=1000
    )
    map.update_layout(mapbox_style="open-street-map")
    map.update_layout(margin={"r":0,"t":0,"l":0,"b":0})

    m = Measurements()
    ts = m.for_station(
        station_id=station_id,
        date_from="2021-08-01", 
        date_to="2021-12-01", 
        conv_win_size=1
    )

    tsg = px.line(ts, x="_time", y="_value")


    return map, tsg