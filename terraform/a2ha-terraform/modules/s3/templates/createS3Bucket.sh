#!/usr/bin/env bash

echo $1
echo $2
if sudo hab pkg exec core/aws-cli aws s3 ls s3://$1; then
  echo "Bucket already exists"
else
 sudo hab pkg exec core/aws-cli aws s3 mb s3://$1 --region $2
  echo "Bucket is created"
fi
