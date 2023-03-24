#!/usr/bin/env bash
set -e

echo "--------------------------------------------------------------------"
echo "###                 Create GKE k8s cluster                        ###"
echo "--------------------------------------------------------------------"

if [[ -z ${GKE_CLUSTER} || -z ${GKE_ZONE} || -z ${GKE_PROJECT} || -z ${MACHINE_TYPE} || -z ${ZSCALER_CIDR1} || -z ${ZSCALER_CIDR2} ]]; then
  echo "GKE_CLUSTER, GKE_ZONE, GKE_PROJECT, MACHINE_TYPE, ZSCALER_CIDR1 and ZSCALER_CIDR2 must be supplied to create GKE cluster"
  exit 1
fi

if [[ -n ${SERVICE_ACCOUNT_JSON} ]]; then
  echo "Authenticating with service account JSON file"
  gcloud auth activate-service-account --key-file="${SERVICE_ACCOUNT_JSON}"
fi

echo "Creating GKE cluster ${GKE_CLUSTER} using default authentication"
gcloud container clusters create "${GKE_CLUSTER}" --zone "${GKE_ZONE}" \
      --node-locations "${GKE_ZONE}" --num-nodes "${NUM_NODES:-5}" --enable-autoscaling \
      --machine-type "${MACHINE_TYPE}" \
      --disk-size 50Gi \
      --min-nodes 1 --max-nodes 5 --project "${GKE_PROJECT}" \
      --enable-master-authorized-networks \
      --master-authorized-networks "${ZSCALER_CIDR1}","${ZSCALER_CIDR2}","${WHITELIST_CIDR}"
      # add your NAT CIDR to whitelist local or CI/CD NAT IP. Set WHITELIST_CIDR in CI/CD to add CIDR to the list automatically.
gcloud container clusters get-credentials "${GKE_CLUSTER}" --zone "${GKE_ZONE}" --project "${GKE_PROJECT}"