FROM --platform=${BUILDPLATFORM} golang:1.17-alpine AS build
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
RUN apk add --no-cache git tzdata zip ca-certificates

WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-X main.Version=$VERSION" -o /out/app .

# https://medium.com/@mhcbinder/using-local-time-in-a-golang-docker-container-built-from-scratch-2900af02fbaf
#WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
#RUN zip -q -r -0 /zoneinfo.zip .

FROM scratch as final
WORKDIR /app/
COPY --from=build /out/app .
# the timezone data:
#ENV ZONEINFO /zoneinfo.zip
#COPY --from=build /zoneinfo.zip /
# the tls certificates:
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ 
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Paris

ENTRYPOINT ["/app/app"]
