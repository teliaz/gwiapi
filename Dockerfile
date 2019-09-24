FROM golang:alpine 

# Add Maintainer info
LABEL maintainer="Elias Krontiris <teliaz@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go get github.com/teliaz/goapi


RUN go build -o main . 
EXPOSE 8080
CMD ["/app/main"]


# Set the current working directory inside the container 
# RUN mkdir /app
# WORKDIR /app


# # Copy the source from the current directory to the working Directory inside the container 

# COPY . .
# # Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


# # Start a new stage from scratch
# FROM alpine:latest
# RUN apk --no-cache add ca-certificates

# WORKDIR /root/


# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/main .
# COPY --from=builder /app/.env .

# # Expose port 8080 to the outside world
# EXPOSE 8080


# #Command to run the executable
# CMD ["./main"]

