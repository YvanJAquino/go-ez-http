# Created with Strapit
FROM    golang:1.18-buster as builder
WORKDIR /app
COPY    . ./
RUN     go build -o service

FROM    debian:buster-slim
WORKDIR /app
ENV		PORT=8010
RUN     set -x && \
		apt-get update && \
		DEBIAN_FRONTEND=noninteractive apt-get install -y \
			ca-certificates && \
			rm -rf /var/lib/apt/lists/* \
		mkdir static
COPY    --from=builder /app/service /app/service
COPY	./static static

CMD     ["/app/service"]
