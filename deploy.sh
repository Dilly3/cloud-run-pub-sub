#! /usr/bin/env bash
declare PROJECT_NAME
if [ -z "$1" ]; then
    echo "enter project name:"
    read PROJECT_NAME
else
    PROJECT_NAME="${1}"
fi
echo "deploying to $PROJECT_NAME"
gcloud run deploy $PROJECT_NAME --source . --region us-west1 --allow-unauthenticated
echo "deployed to $PROJECT_NAME"

exit 0