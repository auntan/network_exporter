FROM golang:1.17 AS build

WORKDIR /src/
COPY . /src/
RUN make

FROM scratch
COPY --from=build /src/configs/config.yaml /bin/configs/config.yaml
COPY --from=build /src/bin/network_exporter /bin/network_exporter
WORKDIR /bin/
ENTRYPOINT ["/bin/network_exporter"]
