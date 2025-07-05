#! /usr/bin/env bash
declare PROJECT_NAME
declare WORK_DIR="$(dirname "$(readlink -f $0)")"
declare FOLDER_NAME="${WORK_DIR##*/}"
if [ -z "$1" ]; then
    PROJECT_NAME="${FOLDER_NAME}"
else
    PROJECT_NAME="${1}"
fi

echo "deploying to $PROJECT_NAME"
gcloud run deploy $PROJECT_NAME --source . --region us-west1 --allow-unauthenticated --platform managed
echo "deployed to $PROJECT_NAME"

exit 0
