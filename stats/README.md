# Stats service

## Zones of responsibility

- Calculates statistics for posts and comments
- Offers quick access to the latest statistics for posts
- Calculates historical data over time periods
- Allows admins to view aggregate activity stats

## Boundaries

- **Input boundaries:** takes aggregate statistics fetch requests from the gateway, handles requests from the admin dashboard
- **Output boundaries:** returns request data back to the gateway
- **Interface boundaries:** communicates to the statistics database, consumes events from topics (produced by posts service)
