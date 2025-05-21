import pandas as pd
import requests

# Datos de ejemplo (simulando una fuente de datos externa)
data = {
    'Fecha': [],
    'Volumen': [],
    'Presion': [],
    'Temperatura': [],
    'Clientes': []
}
df = pd.DataFrame(data)


def get_clients():
    """Obtiene los datos de clientes desde un endpoint remoto."""
    try:
        # Realizamos la solicitud GET al endpoint
        response = requests.get("http://localhost:8080/clients")
        response.raise_for_status()  # Lanza una excepción si la respuesta tiene un error HTTP
        clients = response.json()  # Parseamos la respuesta JSON

        # Retornamos los clientes como un DataFrame
        return pd.DataFrame({'Clientes': clients})
    except requests.RequestException as e:
        print(f"Error al obtener los datos de clientes: {e}")
        # En caso de error, devolvemos un DataFrame vacío
        return pd.DataFrame({'Clientes': []})


def get_data_columns(client_name):
    """
    Obtiene los datos de un cliente específico desde un endpoint remoto.
    """
    try:
        # Definimos la URL del endpoint y el payload
        url = "http://localhost:8080/clients/data"
        payload = {
            "client_name": client_name
        }
        headers = {"Content-Type": "application/json"}

        # Realizamos la solicitud POST
        response = requests.post(url, json=payload, headers=headers)
        response.raise_for_status()  # Lanza una excepción si la respuesta tiene un error HTTP

        # Parseamos la respuesta JSON
        data = response.json()

        # Verificamos si los datos son válidos

        if data and isinstance(data, list):
            midataframe = pd.DataFrame(data)
            print(f"Valores únicos en 'Anomalia io': {midataframe['Anomalia'].unique()}")
            return  midataframe
        else:
            print(f"No se encontraron datos para {client_name}")
            return pd.DataFrame()  # Retornamos un DataFrame vacío si no hay datos
    except requests.RequestException as e:
        print(f"Error al obtener los datos para {client_name}: {e}")
        return pd.DataFrame()  # En caso de error, retornamos un DataFrame vacío