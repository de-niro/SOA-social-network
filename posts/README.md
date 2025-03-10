# Posts serivce

## Zones of responsibility

- Handles posts, comments and reactions operations
- The main focus is on the quick access and modification of individual posts/comments
- Doesn't handle aggregate operations, which are not crucial

## Boundaries

- **Input boundaries:** takes requests from the gateway
- **Output boundaries:** returns request data back to the gateway
- **Interface boundaries:** communicates to the posts database, sends events into topics (consumed by stats service)