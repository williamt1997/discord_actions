ENV DOCKERVERSION=18.03.1-ce
RUN curl -fsSLO https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKERVERSION}.tgz \
  && tar xzvf docker-${DOCKERVERSION}.tgz --strip 1 \
                 -C /usr/local/bin docker/docker \
  && rm docker-${DOCKERVERSION}.tgz

# Use image
FROM golang:1-alpine AS build
# Create working directory 
WORKDIR /app
# Copy files to working directory
COPY go.mod go.sum ./
# Download packages
RUN go mod download
# Copy rest of files
COPY . .
# Create executable
RUN go build -o main .

# Use image
FROM alpine:3.20
# Create working directory
WORKDIR /app
# Copy built app
COPY --from=build /app/main /app/main
# Allow connections on 8080
EXPOSE 8080
# Run executable 
CMD [ "./main" ]
