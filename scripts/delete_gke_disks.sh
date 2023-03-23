#!/usr/bin/env bash

echo "###-----------------------------------------------------------------###"
echo "###            Delete GCP disks in ${GKE_PROJECT}           ###"
echo "### Only disks not attached to VMs and k8s clusters will be deleted ###"
echo "###-----------------------------------------------------------------###"

DISKS=$(gcloud compute disks list --filter="zone:( ${GKE_ZONE} )" --filter="name~'.*.'" | grep -v "NAME" | cut -f1 -d ' ')

declare -i DELETED_DISKS=0
for DISK in $DISKS
  do
    RESULT=$( { gcloud compute disks delete "${DISK}" --zone="${GKE_ZONE}" --project "${GKE_PROJECT}" --quiet; } 2>&1 )
    echo "${RESULT}"
    NUMBER_DEL="$(echo "${RESULT}" | grep -o -i Deleted | wc -l)"
    DELETED_DISKS+=${NUMBER_DEL}
  done

echo "Deleted ${DELETED_DISKS} disks"

exit 0