# Running the User-actions API

A simple API server that provides endpoints for user and action analytics.

## Building and Running with Docker

`docker compose up --build`

## Running with Docker

`docker compose up`

## Example Requests

```bash
curl http://localhost:3000/users/1

curl http://localhost:3000/users/1/actions/count

curl http://localhost:3000/actions/REFER_USER/next

curl http://localhost:3000/referral-indices
```

## Referral Index API Endpoint

- The /referral-indices endpoint calculates the total number of users directly or indirectly referred by each user.

**Time Complexity**: O(n) or -> O(A + U)

- A = number of actions
- U = number of users

# Response Time

api-1 | 17:17:05 | 200 | 745.776µs | 172.22.0.1 | GET | /referral-indices | -
api-1 | 17:17:12 | 200 | 9.585µs | 172.22.0.1 | GET | /referral-indices | -

- 745.776µs is less than 1 millisecond, its 3/4 of a millisecond
- As for the second request here, it used the cache, thats why its about 80x times faster

`Note that: If we add a real db connection, this time would be expected to increase significantly`
