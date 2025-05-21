from dash import Output, Input, dcc
import plotly.express as px
from services import data_service  # Importa el servicio de datos

def register_callbacks(app):
    @app.callback(
        Output('main-graph-output', 'children'),
        Input('tabs', 'value')
    )
    def update_main_graph(tab):
        df = data_service.get_data()
        fig = px.line(df, x='Fecha', y=tab.capitalize(), title=f'{tab.capitalize()} en el tiempo')
        return dcc.Graph(figure=fig)
