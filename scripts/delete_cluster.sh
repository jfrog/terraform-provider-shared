#!/usr/bin/env bash
set -e

if [[ -z ${GKE_CLUSTER} || -z ${GKE_ZONE} || -z ${GKE_PROJECT} ]]; then
  echo "GKE_CLUSTER, GKE_ZONE, GKE_PROJECT, must be supplied to delete GKE cluster"
  exit 1
fi

if [[ -n ${SERVICE_ACCOUNT_JSON} ]]; then
  echo "Authenticating with service account JSON file"
  gcloud auth activate-service-account --key-file="${SERVICE_ACCOUNT_JSON}"
fi

echo "Deleting GKE cluster ${GKE_CLUSTER} using default authentication"
gcloud container clusters delete "${GKE_CLUSTER}" --zone "${GKE_ZONE}" --project "${GKE_PROJECT}" --quiet
echo "GKE cluster ${GKE_CLUSTER} was successfully deleted"