#! /bin/bash

docker run -e REAGLE_LOCAL_LOCATION=$REAGLE_LOCAL_LOCATION -e REAGLE_LOCAL_USER=$REAGLE_LOCAL_USER -e REAGLE_LOCAL_PASSWORD=$REAGLE_LOCAL_PASSWORD -p 9000:9000 kklipsch/reagled