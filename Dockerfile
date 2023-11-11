#!/bin/bash

FROM gcr.io/distroless/static-debian11
 
WORKDIR /app

COPY "bin/coop_case" /app/coop_case

ENTRYPOINT ["/app/coop_case"]
