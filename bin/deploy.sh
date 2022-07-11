#!/bin/sh

PROJECT=unipos-dev-0
SERVICE=logger
IMAGE=gcr.io/${PROJECT}/${SERVICE}
LOCATION=asia-northeast1

env_vars="GOOGLE_CLOUD_PROJECT=${PROJECT},"
env_vars="${env_vars}UNIPOS_PROJECT=${PROJECT},"
env_vars="${env_vars}UNIPOS_SERVICE=${SERVICE},"
env_vars="${env_vars}LOCATION=${LOCATION},"


gcloud builds submit --project=$PROJECT && \
gcloud run deploy $SERVICE \
  --image $IMAGE \
  --set-env-vars="$env_vars" \
  --project=$PROJECT \
  --platform managed \
  --allow-unauthenticated \
  --region=asia-northeast1


