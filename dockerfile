FROM golang:1.20

# Create folder to store my project 
RUN mkdir /app

# Define the desired project address inside the container
WORKDIR /home/app 

# Copy the project from my pc
COPY . .

# Expose the port
EXPOSE 3003

# Run commands inside the container
RUN go get
CMD ["go", "run", "main.go"]
