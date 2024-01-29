# -- Build --
FROM docker.io/library/golang:1.21.5-alpine3.19 AS builder

# Install upx
WORKDIR /source
RUN apk --no-cache add git upx

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN go build -o /source/bin/artifact

# -- Run --
FROM docker.io/library/alpine:3.18 AS runner

# Add non-root user
RUN addgroup -g 1000 codelabs
RUN adduser -D -h /opt/codelabs -s /bin/nologin -G codelabs -u 1000 codelabs
USER codelabs

# Copy binary and run
WORKDIR /opt/codelabs
COPY --from=builder /source/bin/artifact /source/bin/artifact

# Create default folder to save archives
RUN mkdir -p /opt/codelabs/archives
RUN mkdir -p /opt/codelabs/archives/submissions
RUN mkdir -p /opt/codelabs/archives/tests
ENV ARCHIVES_VOLUME_PATH /opt/codelabs/archives

# Copy the templates folder
RUN mkdir -p /opt/codelabs/templates
COPY files/templates /opt/codelabs/templates
ENV TEMPLATES_PATH /opt/codelabs/templates

# Run
EXPOSE 8080
ENTRYPOINT ["/source/bin/artifact"]