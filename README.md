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

# Complexity

Time: O(N + E) where N = users, E = referrals
Space: O(N + E) for storing the graph and results
