FROM --platform=linux/amd64 debian:stable-slim

# Update package repositories and install necessary packages
RUN apt-get update && apt-get install -y ca-certificates

# Copy the ecom1 executable to /usr/bin/
ADD PriceTracker /usr/bin/PriceTracker

# Set the entry point command
CMD ["PriceTracker"]


