# User management service

## Zones of responsibility

- Handles user sessions (login/logout)
- Handles registration, deletion and other operations on users
- Performs email verification
- Handles account incidents

## Boundaries

- **Input boundaries:** takes user operation requests from gateway
- **Output boundaries:** returns request results back to the gateway, sends email to clients
- **Interface boundaries:** communicates to the users database, synchronizes user creation/deletion events with posts service