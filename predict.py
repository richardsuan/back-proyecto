from flask import Flask, request, jsonify
import joblib
import numpy as np

app = Flask(__name__)

@app.route('/predict', methods=['POST'])
def predict():
    data = request.get_json()

    # Cargar el modelo desde el archivo
    model_path = data['model_path']
    model = joblib.load(model_path)

    # Preparar los datos para la predicción
    features = np.array([[data['presion'], data['temperatura'], data['volumen']]])

    # Realizar la predicción
    prediction = model.predict(features)[0]  # Cambia según el método del modelo
    is_anomaly = int(prediction == -1)  # -1 indica anomalía en IsolationForest

    return jsonify({"is_anomaly": is_anomaly})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)