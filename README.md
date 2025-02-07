# Upvote System

## Project Architecture

The **Upvote System** is a simple web application built using the **Gin Framework** in Go. It allows users to create posts, upvote or downvote posts, and retrieve posts with a caching mechanism powered by **Redis**.

### Components:

1. **Controllers**:
   - Handle HTTP requests and responses.
   - Manage the main logic for creating posts, voting, and fetching posts.
   
2. **Services**:
   - The **PostService** handles business logic related to voting on posts (upvote and downvote).
   - The **CacheService** manages caching mechanisms, storing and retrieving posts based on vote thresholds using Redis.
   
3. **Models**:
   - Define data structures for posts and votes.
   - Include validation and serialization methods for these models.
   
4. **Database**:
   - **Redis** is used as the data store, allowing fast access to posts and votes, and also used for caching the most popular posts based on a minimum vote threshold.
## Routes & API Endpoints

The API provides several routes for interacting with posts and votes.

### 1. **Create Post**
   - **Endpoint**: `POST /posts`
   - **Description**: Allows the creation of a new post with a title and content.
   - **Request Body**:
     ```json
     {
       "title": "Post Title",
       "content": "Post content",
       "by_user": 1
     }
     ```
   - **Response**:
     ```json
     {
       "id": "UUID",
       "title": "Post Title",
       "content": "Post content",
       "by_user": 1,
       "votes_count": 0,
       "created_date": "2025-02-07T00:00:00Z"
     }
     ```

### 2. **Get Posts**
   - **Endpoint**: `GET /posts`
   - **Description**: Retrieves all posts, potentially from cache if they have a sufficient number of votes.
   - **Response**: A list of posts in JSON format.
### 3. **Get Post by ID**
   - **Endpoint**: `GET /post/:post_id`
   - **Description**: Fetches a single post by its unique ID.
   - **Response**:
     ```json
     {
       "id": "UUID",
       "title": "Post Title",
       "content": "Post content",
       "by_user": 1,
       "votes_count": 5,
       "created_date": "2025-02-07T00:00:00Z"
     }
     ```

### 4. **Upvote Post**
   - **Endpoint**: `POST /post/upvote`
   - **Description**: Allows a user to upvote a post, increasing the vote count.
   - **Request Body**:
     ```json
     {
       "by_user": 1,
       "post_id": "UUID"
     }
     ```
   - **Response**:
     ```json
     {
       "id": "UUID",
       "title": "Post Title",
       "content": "Post content",
       "by_user": 1,
       "votes_count": 6,
       "created_date": "2025-02-07T00:00:00Z"
     }
     ```
### 5. **Downvote Post**
   - **Endpoint**: `POST /post/downvote`
   - **Description**: Allows a user to downvote a post, decreasing the vote count.
   - **Request Body**:
     ```json
     {
       "by_user": 1,
       "post_id": "UUID"
     }
     ```
   - **Response**:
     ```json
     {
       "id": "UUID",
       "title": "Post Title",
       "content": "Post content",
       "by_user": 1,
       "votes_count": 4,
       "created_date": "2025-02-07T00:00:00Z"
     }
     ```

## Caching System and Redis Usage
This project demonstrates the implementation of a caching system using Redis, aimed at reducing the load on the database by storing frequently accessed data in memory. Redis is used here to cache posts based on the minimum votes, and it ensures fast retrieval of posts without querying the database repeatedly. 

The cache is invalidated automatically whenever a post's vote count exceeds the cached threshold, ensuring that the cache is up-to-date.
## Redis Integration and Caching Details

### Cache Storage:
- **Key Format**: 
  The cache key format for posts is `top_posts_<min_vote>`. For instance, if a post has 5 votes, the cache key will be `top_posts_5`.

- **Cache Duration**: 
  The cache is set to expire after 280 seconds (4 minutes and 40 seconds) to ensure fresh data is fetched periodically.

- **Cache Invalidating**:
  Whenever a new post's votes exceed the minimum threshold stored in the cache, the associated cache for lower thresholds is invalidated. This ensures that only relevant, up-to-date posts are retrieved from the cache.

### Redis Commands Used:
- **ZAdd**: To add the minimum vote values to the sorted set `cached_min_likes`.
- **ZRangeByScore**: Used to fetch keys that should be invalidated based on the vote count threshold.
- **HSet and HGetAll**: Used to store and retrieve individual post and vote data in Redis hashes.

This caching approach significantly boosts performance by reducing unnecessary database queries, especially when users request posts with specific vote thresholds.

## Project Features

### 1. Caching with Redis
This project demonstrates the use of Redis as a caching layer to store and retrieve top posts based on their vote count. It provides fast access to frequently requested data and optimizes response time for clients.

#### How Caching Works:
- When a post's vote count exceeds a certain threshold, the system caches the list of top posts with a minimum vote count in Redis.
- If a new vote pushes a post's vote count above this threshold, the cache is updated accordingly.
- The cache also gets invalidated when a post with a lower vote count is added or when votes drop below the minimum threshold.

### 2. Redis Operations
- **Setting cache**: Posts are cached with a key based on their vote count to optimize for quick access.
- **Cache invalidation**: When a post's vote count is updated, the cache is invalidated to ensure that the data is up-to-date.
- **Cache retrieval**: Cached posts are fetched from Redis using a key based on the minimum vote count.

This approach helps minimize redundant database queries and improves the overall performance of the system.

### 3. Technologies Used

This project utilizes several technologies and libraries that showcase my skills:

- **Go**: The main programming language, used to write the entire backend logic.
- **Gin**: A fast and efficient web framework used for routing and handling HTTP requests.
- **Redis**: Used for caching and fast data retrieval. It helps improve performance by reducing database load.
- **Validator**: For input validation, ensuring the API accepts correct data types and values.
- **JSON**: Used for serializing and deserializing data, particularly for caching and responses.

### 4. Conclusion

This project demonstrates my ability to design and implement an efficient backend system that leverages caching mechanisms using Redis. By building the system with **Go**, I have been able to focus on performance and scalability, particularly in scenarios involving frequent data retrieval and updates. The use of Redis as a cache helps in optimizing query performance and reducing load on the main database, especially for popular posts.

The clean architecture and structured approach show my proficiency in designing maintainable and efficient systems. Additionally, the implementation of voting mechanisms demonstrates my understanding of business logic and my ability to create complex workflows.

Feel free to contribute or explore the repository to learn more about Redis integration and Go programming.
