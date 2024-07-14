# Price Tracker : https://pricetracker-psqtz2bfpa-uc.a.run.app/

Welcome to the Price Tracker repository! This project is designed to help users monitor and track Amazon product prices over time. Built using Go for the backend and React for the frontend, this web application scrapes Amazon product pages to gather price data and presents trends to the users.

## Objective

The primary objective of this project is to provide users with a tool to track Amazon product prices and identify pricing differences over time. By scraping Amazon product pages, the application collects price data and visualizes trends, enabling users to make informed purchasing decisions.

## Features

- **User Authentication**: Secure login and registration system.
- **Product Listing**: Browse and search through a wide range of tracked products.
- **Price Tracking**: Scrape and track Amazon product prices over time using go Concurrency.
- **Price Trends**: Visualize price trends and historical data.
- **Responsive Design**: Optimized for both desktop and mobile devices.
- **Continuous Deployment**: GitHub Actions for CI/CD in GCP.

## Technologies Used

- **Backend**: Go
- **Frontend**: Js
- **Database**: PostgreSQL
- **API**: RESTful APIs

## Installation

### Prerequisites

- Go (1.21.5 or later)
- Node.js (16.x or later)
- PostgreSQL

### Quick StartUp

1. Clone the repository:

- bash : git clone https://github.com/stonoy/PriceTracker.git

2. Navigate to the backend directory and install dependencies

- cd root
- go mod tidy

3. Set up the database by applying migrations
- goose postgres <database-connection-string> up

4. Navigate to the frontend directory and install dependencies

- cd client
- npm install

5. Start the frontend

- npm run dev

6. Build the frontend and copy the dist to root directory

- npm run build

7. Build and start the server

- go build -o PriceTracker && ./PriceTracker

# Feel free to customize it further according to your project's specifics and requirements.
