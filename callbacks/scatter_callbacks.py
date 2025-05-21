from dash import Output, Input, dcc
import plotly.express as px
from services import data_service  # Importa el servicio de datos

def register_callbacks(app):
    @app.callback(
        Output('scatter-graph', 'figure'),
        [Input('x-axis', 'value'),
         Input('y-axis', 'value')]
    )
    def update_scatter(x_axis, y_axis):
        df = data_service.get_data()
        fig = px.scatter(df, x=x_axis, y=y_axis, trendline="ols")
        return fig
