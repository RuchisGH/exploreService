# ExploreService

ExploreService is a gRPC-based microservice designed to handle decisions (e.g., LIKE, PASS) between users on a platform. It supports key functionalities such as recording user decisions, listing users who liked a particular user, counting likes, and identifying mutual decisions.

---

## Features

- **PutDecision**: Record a user's decision (LIKE or PASS) for another user.
- **ListLikedYou**: Retrieve a list of users who liked a given user.
- **ListNewLikedYou**: Identify users who liked a given user but have not received a like in return.
- **CountLikedYou**: Count the number of likes a user has received.

---

## Requirements

- **Go**: Version 1.23.2 or later
- **MySQL**: Used as the database
- **Docker**: To containerize the application
- **Postman**: For testing gRPC APIs (optional)

---

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/exploreService.git
   cd exploreService

## Running the Application
   **With Docker**
   Build the Docker image: docker-compose build
   Run the container:docker-compose up
   
   **Running Locally**
   Start MySQL locally and ensure the credentials match the .env file.
   Run the server: go run cmd/server/main.go
   
##  Testing
   Use **Postman** or any other gRPC client to test the API methods.
   Verify database changes via SQL queries using a MySQL client.
   gRPC Requests (Postman)
   **PutDecision**:
   {
    "actor_user_id": "user1",
    "liked_recipient": true,
    "recipient_user_id": "user2"
   }
   **ListLikedYou**:
   {
      "recipient_user_id": "user2"
   }
   **CountLikedYou**:
   Copy code
   {
    "recipient_user_id": "user2"
   }

##  Scaling the Application
To handle large-scale user data:

**Database Indexing** : Ensure indexes on frequently queried fields (user_id, target_id, decision).
**Caching** : Use Redis or Memcached to cache frequent queries like CountLikedYou.
**Load Balancing** : Use a load balancer to distribute traffic across multiple instances.
**Asynchronous Processing** : Use message queues for processing decisions asynchronously.
