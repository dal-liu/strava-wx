## Setup

### Authorizing the application

These steps assume that you have set up a Strava API application already. If you have not, follow the instructions at https://developers.strava.com/docs/getting-started/.

1. In a browser, navigate to 
    ```
    https://www.strava.com/oauth/authorize?client_id=YOUR_CLIENT_ID&redirect_uri=http://localhost&response_type=code&scope=activity:read_all,activity:write
    ```
    where ```YOUR_CLIENT_ID``` is your Strava client ID.

2. Click "Authorize" to grant Strava read and write access to your app.

3. You will be redirected to an error page with a URL that looks like
    ```
    http://localhost/?state=&code=YOUR_CODE&scope=read,activity:read_all,activity:write
    ```
    Copy ```YOUR_CODE``` from the URL.

4. In a terminal, run the command
    ```
    curl -X POST https://www.strava.com/oauth/token \
      -d client_id=YOUR_CLIENT_ID \
      -d client_secret=YOUR_CLIENT_SECRET \
      -d code=YOUR_CODE \
      -d grant_type=authorization_code
    ```
    where ```YOUR_CLIENT_SECRET``` is your Strava client secret.

5. Copy the ```refresh_token``` value from the response. This is the token you will use to retrieve new access tokens for reading and writing Strava data.
