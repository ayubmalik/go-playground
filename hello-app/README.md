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
gcloud services enable secretmanager.googleapis.com

```

### IAM Config

Allow build service access
```
gcloud iam service-accounts create cloud-build-go  --description="Build and test Go applications" --display-name="Cloud Build Go (TDSSCHEDULES)"
gcloud iam service-accounts list --filter="email:-compute@developer.gserviceaccount.com"
gcloud iam service-accounts add-iam-policy-binding 843958366025-compute@developer.gserviceaccount.com --member="serviceAccount:cloud-build-go@tdsschedules.iam.gserviceaccount.com" --role="roles/iam.serviceAccountUser"
```

Service account permissions
```
gcloud projects add-iam-policy-binding tdsschedules --member='serviceAccount:cloud-build-go@tdsschedules.iam.gserviceaccount.com' --role='roles/cloudbuild.builds.builder'
gcloud projects add-iam-policy-binding tdsschedules --member='serviceAccount:cloud-build-go@tdsschedules.iam.gserviceaccount.com' --role='roles/artifactregistry.writer'
gcloud projects add-iam-policy-binding tdsschedules --member='serviceAccount:cloud-build-go@tdsschedules.iam.gserviceaccount.com' --role='roles/storage.objectCreator'
gcloud projects add-iam-policy-binding tdsschedules --member='serviceAccount:cloud-build-go@tdsschedules.iam.gserviceaccount.com' --role='roles/run.developer'
gcloud projects add-iam-policy-binding tdsschedules --member='serviceAccount:cloud-build-go@tdsschedules.iam.gserviceaccount.com' --role='roles/secretmanager.admin'
gcloud projects add-iam-policy-binding tdsschedules --member='serviceAccount:service-843958366025@gcp-sa-cloudbuild.iam.gserviceaccount.com' --role='roles/secretmanager.admin'
```


### Create Docker Repo

```
gcloud artifacts repositories create hello-app --repository-format=docker --location=us-central1 --description="Hello App Repo"

```

### Github connection

```
 gcloud builds connections create github github-connection-goplayground --region=us-central1
```
