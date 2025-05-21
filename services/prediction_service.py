import requests

def predict_anomaly(client_name, fecha, presion, temperatura, volumen):
    url = 'http://localhost:8080/predict'
    headers = {'Content-Type': 'application/json'}
    payload = {
        "client_name": client_name,
        "fecha": fecha,
        "presion": presion,
        "temperatura": temperatura,
        "volumen": volumen
    }

    try:
        print(payload)
        response = requests.post(url, json=payload, headers=headers)
        response.raise_for_status()  # Lanza una excepción si el código de estado no es 2xx
        print(response.json())
        return response.json()  # Devuelve la respuesta como un diccionario
    except requests.exceptions.RequestException as e:
        return {"error": str(e)}