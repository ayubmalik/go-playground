TDS Schedules

# cloudbuild

see [cloudbuild docs](https://cloud.google.com/build/docs/building/build-go)

create service account

```shell
gcloud iam service-accounts create cloud-build-go \
--description="Build and test Go applications" \
--display-name="TDS Cloud Build Go" \
--project="tdsschedules"
```

allow computer service account to switch to role

```shell
gcloud iam service-accounts list --filter="email:-compute@developer.gserviceaccount.com"

gcloud iam service-accounts add-iam-policy-binding \
843958366025-compute@developer.gserviceaccount.com   \
--member="serviceAccount:cloud-build-go@tdsschedules.iam.gserviceaccount.com" \
--role="roles/iam.serviceAccountUser"
```

permissions

```shell
gcloud projects add-iam-policy-binding tdsschedules \
--member="serviceAccount:cloud-build-go@tdsschedules.iam.gserviceaccount.com" \
--role="roles/cloudbuild.builds.builder"

gcloud projects add-iam-policy-binding tdsschedules \
--member="serviceAccount:cloud-build-go@tdsschedules.iam.gserviceaccount.com" \
--role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding tdsschedules \
--member="serviceAccount:cloud-build-go@tdsschedules.iam.gserviceaccount.com" \
--role="roles/storage.objectCreator"

gcloud projects add-iam-policy-binding tdsschedules \
--member="serviceAccount:cloud-build-go@tdsschedules.iam.gserviceaccount.com" \
--role="roles/run.developer"
```
