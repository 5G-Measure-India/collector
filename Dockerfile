FROM ubuntu:focal AS build

ARG GO_VERSION=1.23.3

RUN \
  apt-get update && \
  apt-get install -y --no-install-recommends \
  ca-certificates make wget

RUN \
  suffix=""; \
  if [ "$TARGETPLATFORM" = "linux/amd64" ] || [ "$(uname -m)" = "x86_64" ]; then \
  suffix="amd64"; \
  elif [ "$TARGETPLATFORM" = "linux/arm64" ] || [ "$(uname -m)" = "aarch64" ]; then \
  suffix="arm64"; \
  else echo "Unsupported platform"; exit 1; \
  fi && \
  wget -qO- https://go.dev/dl/go${GO_VERSION}.linux-${suffix}.tar.gz \
  | tar xz -C /usr/local

ENV PATH=/usr/local/go/bin:${PATH}

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .

RUN make release

FROM mobileinsight-core:dev

WORKDIR /app

RUN \
  apt-get update && \
  apt-get install -y --no-install-recommends adb && \
  rm -rf /var/lib/apt/lists/*

COPY --from=build /app/collector .
COPY mi.py .

ENTRYPOINT [ "./collector" ]
