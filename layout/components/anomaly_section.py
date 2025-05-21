from dash import html, dcc, Input, Output, State, callback
from services.prediction_service import predict_anomaly
import pandas as pd

def render():
    return html.Div([
        html.H3('Anomalías', className='section-title'),
        html.Div([
            html.Label('Presión:', className='input-label'),
            dcc.Input(id='input-presion', type='number', placeholder='Ingrese presión', className='input-field'),
            html.Label('Temperatura:', className='input-label'),
            dcc.Input(id='input-temperatura', type='number', placeholder='Ingrese temperatura', className='input-field'),
            html.Label('Volumen:', className='input-label'),
            dcc.Input(id='input-volumen', type='number', placeholder='Ingrese volumen', className='input-field'),
        ], className='input-container'),
        html.Button('Realizar prediccion', id='predict-button', className='action-button'),
        dcc.ConfirmDialog(id='popup-dialog', message=""),  # Popup para mostrar el mensaje
        html.Div(id='prediction-result', className='result-message')  # Contenedor para el resultado
    ])

@callback(
    [Output('popup-dialog', 'displayed'),  # Controla si el popup se muestra
     Output('popup-dialog', 'message')],  # Mensaje del popup
    [Input('predict-button', 'n_clicks')],
    [State('selected-client', 'data'),  # Obtener el cliente seleccionado
     State('client-data', 'data'),  # Obtener los datos del cliente
     State('input-presion', 'value'),
     State('input-temperatura', 'value'),
     State('input-volumen', 'value')]
)
def handle_prediction(n_clicks, client_name, client_data, presion, temperatura, volumen):
    if n_clicks is None:
        return False, ""

    if not client_name:
        return True, "Error: No se ha seleccionado un cliente."

    if not client_data:
        return True, "Error: No hay datos disponibles para el cliente seleccionado."

    # Convertir los datos del cliente a un DataFrame
    df = pd.DataFrame(client_data)

    # Obtener la última fecha y hora registrada
    if 'Fecha' not in df.columns:
        return True, "Error: Los datos del cliente no contienen una columna de fecha."
    df['Fecha'] = pd.to_datetime(df['Fecha'])  # Asegurarse de que la columna 'Fecha' esté en formato datetime
    ultima_fecha = df['Fecha'].max()  # Obtener la última fecha
    fecha_mas_una_hora = ultima_fecha + pd.Timedelta(hours=1)  # Sumar una hora

    # Formatear la fecha como un número flotante (YYYYMMDD.HHMMSS)
    fecha_formateada = float(fecha_mas_una_hora.strftime('%Y%m%d.%H%M%S'))

    # Llamar al servicio de predicción
    result = predict_anomaly(client_name, fecha_formateada, presion, temperatura, volumen)

    # Manejar la respuesta
    if "error" in result:
        return True, f"Error al realizar la predicción: {result['error']}"
    elif result.get("anomalia"):
        return True, "⚠️ Resultado: ¡Es una anomalía!"
    else:
        return True, "✅ Resultado: No es una anomalía."