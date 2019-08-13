FROM scratch

# Common configuration
ENV GRPC_SERVER_PORT 8000
ENV GRPC_INFO_PORT 8080
ENV GRPC_LOGGER_LEVEL 0

# Exposing ports
EXPOSE $GRPC_SERVER_PORT
EXPOSE $GRPC_INFO_PORT

# Copy dependecies
COPY certs /etc/ssl/certs/
COPY bin/linux-amd64/grpc /

CMD ["/grpc", "serve"]
