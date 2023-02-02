# Code Things
This social media application built using the Spin SDK is intended to be an in-depth reference or guide for Spin developers. It focuses on more common real-world use cases like CRUD APIs, Token auth, etc.

## SDK Requirements
The following SDKs are required to build this application.
- [Spin SDK](https://developer.fermyon.com/spin/install)
- [rustup](https://rustup.rs)
- [Node.js](https://nodejs.org)

## External Resources
The following resources are required to run this application. For local development, these resources can be started via docker.
- [PostreSQL](https://www.postgresql.org)

## Auth0 Setup

To complete this setup for Fermyon Cloud, you must have run `spin deploy` at least once to capture the app's URL. For example, the application's URL in the following example is `https://code-things-xxx.fermyon.app`:
```
% spin deploy
Uploading code-things version 0.1.0+rcf68d278...
Deploying...
Waiting for application to become ready.......... ready
Available Routes:
  web: https://code-things-xxx.fermyon.app (wildcard)
  profile-api: https://code-things-xxx.fermyon.app/api/profile (wildcard)
```

1. Sign up for Auth0 account (free)
2. Create a "Single Page Web" application
    a. Configure callback URLs: `http://127.0.0.1:3000, https://code-things-xxx.fermyon.app`
    b. Configure logout URLs: `http://127.0.0.1:3000, https://code-things-xxx.fermyon.app`
    c. Allowed web origins: `http://127.0.0.1:3000, https://code-things-xxx.fermyon.app`
    d. Add GitHub Connection
3. Create API
    a. Name: 'Code Things API'
    b. Identifier: `https://code-things-xxx.fermyon.app/`
    c. Signing Algorithm: `RS256`
4. Add the Auth0 configuration to Vue.js:
    a. Create a file at `./web/.env.local` (this is gitignored)
    b. Add domain: `VITE_AUTH0_DOMAIN = "dev-xxx.us.auth0.com"`
    c. Add client id: `VITE_AUTH0_CLIENT_ID = "xxx"`
    c. Add audience: `VITE_AUTH0_AUDIENCE = "https://code-things.fermyon.app/api"`
