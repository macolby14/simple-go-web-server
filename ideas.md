## Auth Flow
- Check if we have a user session in memory (already logged in)
- If not:
    - Oauth2 redirect to google authenticates and gets user info
    - Check if user info is in our database, if not... create a new user
    - Assign a session to the user with user info
- create a session with our user info
