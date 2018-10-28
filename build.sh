#! /bin/bash

docker build -t kklipsch/reagled --build-arg TEST_FLAG=$1 --build-arg REAGLE_LOCAL_LOCATION=$REAGLE_LOCAL_LOCATION --build-arg REAGLE_LOCAL_USER=$REAGLE_LOCAL_USER --build-arg REAGLE_LOCAL_PASSWORD=$REAGLE_LOCAL_PASSWORD --no-cache .
