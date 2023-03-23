#!/usr/bin/env bash
set -e

echo "--------------------------------------------------------------------"
echo "###                 Create GKE k8s cluster                        ###"
echo "--------------------------------------------------------------------"

printenv

if [ -n "${SERVICE_ACCOUNT_JSON}" ]; then
  echo "Authenticating with service account JSON file"
  gcloud auth activate-service-account --key-file="${SERVICE_ACCOUNT_JSON}"
fi

gcloud container clusters create "${GKE_CLUSTER}" --zone "${GKE_ZONE}" \
      --node-locations "${GKE_ZONE}" --num-nodes "${NUM_NODES:-5}" --enable-autoscaling \
      --machine-type "${MACHINE_TYPE}" \
      --disk-size 50Gi \
      --min-nodes 1 --max-nodes 5 --project "${GKE_PROJECT}" \
      --enable-master-authorized-networks \
      --master-authorized-networks "${ZSCALER_CIDR1}","${ZSCALER_CIDR2}","${WHITELIST_CIDR}"
      # add your NAT CIDR to whitelist local or CI/CD NAT IP. Set WHITELIST_CIDR in CI/CD to add CIDR to the list automatically.
gcloud container clusters get-credentials "${GKE_CLUSTER}" --zone "${GKE_ZONE}" --project "${GKE_PROJECT}"