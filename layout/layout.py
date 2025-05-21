from layout.components import filter_section, anomaly_section, graph_tabs
from services import data_service
from dash import dcc, html, Input, Output, callback
import pandas as pd
from dash import callback

app_layout = html.Div([
    dcc.Store(id='selected-client', storage_type='memory'),
    dcc.Store(id='client-data', storage_type='memory'),

    html.Div([
        html.Div([
            html.H1('Contugas - Sistema de Detección de Anomalías', className='header-title'),
            html.Button('CONTUGAS', className='header-button')
        ], className='header-content')
    ], className='header'),

    html.Div(className='main-body', children=[
        html.Div(className='left-column', children=[
            filter_section.render(),
            anomaly_section.render()
        ]),

        html.Div(className='right-column', children=[
            html.Div(className='main-graph-section', children=[
                html.H3(id='main-title', children='Comportamiento y Predicción - Cliente'),
                graph_tabs.render(),
                dcc.Graph(id='scatter-graph')
            ]),

            html.Div(className='bottom-section', children=[
                html.Div(className='variables-section', children=[
                    html.H3('Relación entre Variables'),
                    html.Div([
                        html.Div([
                            html.Label('Eje X:'),
                            dcc.Dropdown(id='x-axis', options=[], value=None)
                        ], className='axis-selector'),
                        html.Div([
                            html.Label('Eje Y:'),
                            dcc.Dropdown(id='y-axis', options=[], value=None)
                        ], className='axis-selector')
                    ], className='axis-selectors'),  # Contenedor para los selectores
                    dcc.Graph(id='variables-graph')  # Gráfica para "Relación entre Variables"
                ]),
                html.Div(className='summary-section', children=[
                    html.H3('Resumen Estadístico'),
                    html.Div(id='summary-output'),
                    html.Div([
                        html.P('Promedio: ', id='mean-output'),
                        html.P('Máximo: ', id='max-output'),
                        html.P('Mínimo: ', id='min-output'),
                        html.P('Desviación Estándar: ', id='std-output')
                    ], className='statistics-values'),

                ])
            ])
        ])
    ])
])

@callback(
    [Output('selected-client', 'data'),
     Output('client-data', 'data'),
     Output('main-title', 'children')],
    Input('client-dropdown', 'value')
)
def update_client_and_data(client_id):
    if not client_id:
        return None, None, 'Comportamiento y Predicción'

    # Obtener los datos del cliente como un DataFrame
    data = data_service.get_data_columns(client_id)

    # Convertir el DataFrame a una lista de diccionarios (JSON serializable)
    data_as_dict = data.to_dict('records') if not data.empty else None

    return client_id, data_as_dict, f'Comportamiento y Predicción - Cliente {client_id}'


@callback(
    [Output('x-axis', 'options'),
     Output('y-axis', 'options'),
     Output('scatter-graph', 'figure')],
    [Input('client-data', 'data'),
     Input('graph-tabs', 'value'),
     Input('date-picker-range', 'start_date'),
     Input('date-picker-range', 'end_date')]
)
def update_graph(client_data, selected_tab, start_date, end_date):
    if not client_data or not selected_tab:
        return [], [], {
            'data': [],
            'layout': {
                'title': 'Seleccione un cliente, una pestaña y un rango de fechas para visualizar la gráfica'
            }
        }

    # Convertir la lista de diccionarios de nuevo a un DataFrame
    df = pd.DataFrame(client_data)

    # Validar que las fechas seleccionadas existan y filtrar el DataFrame
    if start_date and end_date:
        df['Fecha'] = pd.to_datetime(df['Fecha'])  # Asegurarse de que 'Fecha' esté en formato datetime
        df = df[(df['Fecha'] >= start_date) & (df['Fecha'] <= end_date)]
    else:
        # Si no hay rango de fechas seleccionado, usar todos los datos
        df['Fecha'] = pd.to_datetime(df['Fecha'])

    # Generar las opciones para los dropdowns
    dropdown_options = [{'label': col, 'value': col} for col in df.columns]

    # Validar los valores seleccionados para los ejes
    x_axis_value = 'Fecha'
    y_axis_value = selected_tab
    print(f"Valores únicos en 'Anomalia': {df['Anomalia'].unique()}")
    # Dividir los datos en normales y anómalos
    normal_data = df[df['Anomalia'] == 0]
    # Imprimir la cantidad de datos normales y anómalos
    print(f"Datos normales: {len(normal_data)}")
    anomalous_data = df[df['Anomalia'] == 1]
    print(f"Datos anómalos: {len(anomalous_data)}")
    q1 = df[y_axis_value].quantile(0.02)  # Primer cuartil
    q3 = df[y_axis_value].quantile(0.98)  # Tercer cuartil

    # Crear la figura de la gráfica
    fig = {
        'data': [
            # Puntos normales
            {
                'x': normal_data[x_axis_value],
                'y': normal_data[y_axis_value],
                'type': 'scatter',
                'mode': 'markers',
                'name': 'Puntos normales',
                'marker': {'color': 'blue'}
            },
            # Puntos anómalos
            {
                'x': anomalous_data[x_axis_value],
                'y': anomalous_data[y_axis_value],
                'type': 'scatter',
                'mode': 'markers',
                'name': 'Puntos anómalos',
                'marker': {'color': 'red'}
            },
            # Línea para Q1 (primer cuartil)
            {
                'x': df[x_axis_value],
                'y': [q1] * len(df),
                'type': 'scatter',
                'mode': 'lines',
                'name': 'Límite inferior',
                'line': {'color': 'red', 'dash': 'dash'}
            },
            # Línea para Q3 (tercer cuartil)
            {
                'x': df[x_axis_value],
                'y': [q3] * len(df),
                'type': 'scatter',
                'mode': 'lines',
                'name': 'Límite superior',
                'line': {'color': 'red', 'dash': 'dash'}
            }
        ],
        'layout': {
            'title': f'{y_axis_value} vs {x_axis_value} (Límites Intercuartiles)',
            'xaxis': {'title': x_axis_value},
            'yaxis': {'title': y_axis_value}
        }
    }

    return dropdown_options, dropdown_options, fig


@callback(
    Output('variables-graph', 'figure'),
    [Input('client-data', 'data'),
     Input('x-axis', 'value'),
     Input('y-axis', 'value'),
     Input('date-picker-range', 'start_date'),
     Input('date-picker-range', 'end_date')]

)
def update_variables_graph(client_data, x_axis_value, y_axis_value, start_date, end_date):
    if not client_data or not x_axis_value or not y_axis_value:
        return {
            'data': [],
            'layout': {
                'title': 'Seleccione variables para los ejes X e Y y un rango de fechas'
            }
        }

    # Convertir la lista de diccionarios de nuevo a un DataFrame
    df = pd.DataFrame(client_data)

    # Validar que las fechas seleccionadas existan y filtrar el DataFrame
    if start_date and end_date:
        df['Fecha'] = pd.to_datetime(df['Fecha'])  # Asegurarse de que 'Fecha' esté en formato datetime
        df = df[(df['Fecha'] >= start_date) & (df['Fecha'] <= end_date)]
    else:
        # Si no hay rango de fechas seleccionado, usar todos los datos
        df['Fecha'] = pd.to_datetime(df['Fecha'])

    # Validar que las columnas seleccionadas existan en el DataFrame
    if x_axis_value not in df.columns or y_axis_value not in df.columns:
        return {
            'data': [],
            'layout': {
                'title': 'Seleccione variables válidas para los ejes X e Y'
            }
        }

    # Crear la figura de la gráfica
    fig = {
        'data': [
            {
                'x': df[x_axis_value],
                'y': df[y_axis_value],
                'type': 'scatter',
                'mode': 'markers',  # Mostrar puntos en lugar de líneas
                'name': f'{y_axis_value} vs {x_axis_value}'
            }
        ],
        'layout': {
            'title': f'Relación entre {y_axis_value} y {x_axis_value}',
            'xaxis': {'title': x_axis_value},
            'yaxis': {'title': y_axis_value}
        }
    }

    return fig


@callback(
    [Output('mean-output', 'children'),
     Output('max-output', 'children'),
     Output('min-output', 'children'),
     Output('std-output', 'children')],
    [Input('client-data', 'data'),
     Input('graph-tabs', 'value'),
     Input('date-picker-range', 'start_date'),
     Input('date-picker-range', 'end_date')]
)
def update_statistics(client_data, selected_tab, start_date, end_date):
    if not client_data or not selected_tab:
        return 'Promedio: N/A', 'Máximo: N/A', 'Mínimo: N/A', 'Desviación Estándar: N/A'

    # Convertir la lista de diccionarios de nuevo a un DataFrame
    df = pd.DataFrame(client_data)

    # Validar que las fechas seleccionadas existan y filtrar el DataFrame
    if start_date and end_date:
        df['Fecha'] = pd.to_datetime(df['Fecha'])  # Asegurarse de que 'Fecha' esté en formato datetime
        df = df[(df['Fecha'] >= start_date) & (df['Fecha'] <= end_date)]

    # Validar que la columna seleccionada en la pestaña exista en el DataFrame
    if selected_tab not in df.columns:
        return 'Promedio: N/A', 'Máximo: N/A', 'Mínimo: N/A', 'Desviación Estándar: N/A'

    # Calcular estadísticas
    mean_value = df[selected_tab].mean()
    max_value = df[selected_tab].max()
    min_value = df[selected_tab].min()
    std_value = df[selected_tab].std()

    # Retornar los valores formateados
    return (
        f'Promedio: {mean_value:.2f}',
        f'Máximo: {max_value:.2f}',
        f'Mínimo: {min_value:.2f}',
        f'Desviación Estándar: {std_value:.2f}'
    )