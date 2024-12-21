import sys
import pandas as pd
import joblib
import requests
import warnings
import random

warnings.filterwarnings("ignore", category=UserWarning)

def download_from_ipfs(cid, output_path):
    url = f"https://gateway.pinata.cloud/ipfs/{cid}"
    response = requests.get(url)
    response.raise_for_status()
    with open(output_path, "wb") as f:
        f.write(response.content)

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print("Usage: python predict.py <datasetCID> <modelCID> <scalerCID>")
        sys.exit(1)

    dataset_cid = sys.argv[1]
    model_cid = sys.argv[2]
    scaler_cid = sys.argv[3]

    # Download files
    download_from_ipfs(dataset_cid, "dataset.csv")
    download_from_ipfs(model_cid, "model.pkl")
    download_from_ipfs(scaler_cid, "scaler.pkl")

    data = pd.read_csv("dataset.csv")
    # If data is small, just re-sample from it, or generate random values
    # Randomize input between 1 and 10000:
    # Ensure the columns exist in dataset, else adapt accordingly.
    if len(data) < 3:
        # If dataset too small, generate random rows
        rows = []
        for _ in range(3):
            row = {
                'Hours Studied': random.randint(1, 10000),
                'Previous Scores': random.randint(1, 10000),
                'Extracurricular Activities': random.randint(1, 10000),
                'Sleep Hours': random.randint(1, 10000),
                'Sample Question Papers Practiced': random.randint(1, 10000)
            }
            rows.append(row)
        data = pd.DataFrame(rows)
    else:
        # Randomly sample 3 rows from dataset
        data = data.sample(n=3, replace=True)

        # Optionally modify the sampled rows to have random values 1-10000
        # to match the user's requirement more literally:
        data['Hours Studied'] = [random.randint(1, 10000) for _ in range(3)]
        data['Previous Scores'] = [random.randint(1, 10000) for _ in range(3)]
        data['Extracurricular Activities'] = [random.randint(1, 10000) for _ in range(3)]
        data['Sleep Hours'] = [random.randint(1, 10000) for _ in range(3)]
        data['Sample Question Papers Practiced'] = [random.randint(1, 10000) for _ in range(3)]

    scaler = joblib.load("scaler.pkl")
    model = joblib.load("model.pkl")

    X = data[['Hours Studied', 'Previous Scores', 'Extracurricular Activities', 'Sleep Hours', 'Sample Question Papers Practiced']]
    X_normalized = scaler.transform(X)
    predictions = model.predict(X_normalized)

    for p in predictions:
        print(f"{p:.6f}")
