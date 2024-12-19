import sys
import pandas as pd
import joblib
import requests
import os
import warnings
warnings.filterwarnings("ignore", category=UserWarning)


def download_from_ipfs(cid, output_path):
    url = f"https://gateway.pinata.cloud/ipfs/{cid}"
    response = requests.get(url)
    response.raise_for_status()
    with open(output_path, "wb") as f:
        f.write(response.content)

if __name__ == "__main__":
    # Arguments: datasetCID, modelCID, scalerCID
    # Example: python predict.py <datasetCID> <modelCID> <scalerCID>
    if len(sys.argv) < 4:
        print("Usage: python predict.py <datasetCID> <modelCID> <scalerCID>")
        sys.exit(1)

    dataset_cid = sys.argv[1]
    model_cid = sys.argv[2]
    scaler_cid = sys.argv[3]

    try:
        # Download dataset, model and scaler from IPFS
        download_from_ipfs(dataset_cid, "dataset.csv")
        download_from_ipfs(model_cid, "model.pkl")
        download_from_ipfs(scaler_cid, "scaler.pkl")

        # Load dataset
        data = pd.read_csv("dataset.csv")
        # For safety, ensure the dataset has at least 3 rows
        data = data.head(3)

        scaler = joblib.load("scaler.pkl")
        model = joblib.load("model.pkl")

        # Prepare input features
        X = data[['Hours Studied', 'Previous Scores', 'Extracurricular Activities', 'Sleep Hours', 'Sample Question Papers Practiced']]

        # Normalize features
        X_normalized = scaler.transform(X)

        # Make predictions
        predictions = model.predict(X_normalized)

        # Output only the 3 predictions
        for p in predictions:
            print(f"{p:.6f}")

    except Exception as e:
        print("Error occurred:", e)
        sys.exit(1)
