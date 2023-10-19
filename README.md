Certainly! Below is a simple README.md file based on the provided Go code. Feel free to customize it further to meet your specific needs.

```markdown
# Simple Receipt Processing API

This is a simple Go application that provides an API for processing receipts and calculating points based on certain rules.

# Kindly note
The getPoints function gives incorrect points and could have been rectified. But I kept a time limit due to which I forced myself to submit this codebase. Besides the Node.js project contains accurate points calculation logic

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Endpoints](#endpoints)
- [Rules](#rules)
- [License](#license)

## Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/dl/) (tested with Go version 1.16)
- [Git](https://git-scm.com/)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/chetan-plrch/simple-receipt-app-go.git
   ```

2. Navigate to the project directory:

   ```bash
   cd simple-receipt-app-go
   ```

3. Run the application:

   ```bash
   go run main.go
   ```

## Usage

The application will start the server on http://localhost:3000.

## Endpoints

- **POST /receipts/process**: Store a new receipt.

  Example:

  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"retailer": "Example Retailer", "purchaseDate": "2023-10-18", "purchaseTime": "15:30", "items": [{"shortDescription": "Item 1", "price": "20.99"}], "total": "20.99"}' http://localhost:3000/receipts/process
  ```

- **GET /receipts/{id}/points**: Calculate points for a receipt by ID.

  Example:

  ```bash
  curl http://localhost:3000/receipts/{your-receipt-id}/points
  ```

## Rules

The points calculation is based on the following rules:

1. One point for every alphanumeric character in the retailer name.
2. 50 points if the total is a round dollar amount with no cents.
3. 25 points if the total is a multiple of 0.25.
4. 5 points for every two items on the receipt.
5. If the trimmed length of the item description is a multiple of 3, calculate points.
6. 6 points if the day in the purchase date is odd.
7. 10 points if the time of purchase is after 2:00 pm and before 4:00 pm.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```
