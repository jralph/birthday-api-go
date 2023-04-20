FROM golang:1.20 as build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/

# Test
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go test -cover -race ./...

# Build
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod CGO_ENABLED=0 GOOS=linux go build -o /birthdays-api ./cmd/birthdaysApi/main.go

EXPOSE 8080

# Run
CMD ["/birthdays-api"]

# Create a basic user image sollely responsible for making a user
FROM alpine:3.15 as user

# Here we set a default (but configurable) user id and group id to run as
# This should ideally be randomised for each project created to limit any spread of breaking out of a container
# Its good practice to never reuse id's if possible as then you can be sure the same id is not in use anywhere else
ARG uid=10652
ARG gid=10652

RUN echo "scratchuser:x:${uid}:${gid}::/home/scratchuser:/bin/sh" > /scratchpasswd

# Create our app image based off a scratch image
# Our app requires minimal 3rd party libraries or system calls so we can use a scratch image
# This gives us some strong security and a very light weight footprint due the the lack of any other package being
# included in the final docker image
# You can't be vulnerable to a package exploit if you don't have any packages!
FROM scratch as app

# Pull in assets from the other images to create a minimal image
COPY --from=user /scratchpasswd /etc/passwd
COPY --from=build /birthdays-api /birthdays-api

EXPOSE 8080

# Run
CMD ["/birthdays-api"]