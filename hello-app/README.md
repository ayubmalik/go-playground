# hello-app

To test Google Cloud deployment

## GCP Instructions

```
gcloud projects list

gcloud config set project tdsschedules

gcloud services list --available

```

### Enable Cloud Build Services

```
gcloud services enable run.googleapis.com

gcloud services enable cloudbuild.googleapis.com

gcloud services enable artifactregistry.googleapis.com

```

### IAM Config

```
gcloud iam service-accounts create cloud-build-go  --description="Build and test Go applications" --display-name="Cloud Build Go (TDSSCHEDULES)"

```
