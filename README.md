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
    where ```YOUR_CLIENT_SECRET``` is your Strava client secret. Alternatively, you can use a tool such as Postman instead.

5. Save the ```refresh_token``` value from the response. This is the token you will use to retrieve new access tokens for reading and writing Strava data.

### Retrieving an access token

If you have just completed the steps in the previous section, the access token is available from the ```access_token``` field in the response. If your access token has expired, run the command

```
curl -X POST https://www.strava.com/oauth/token \
  -d client_id=YOUR_CLIENT_ID \
  -d client_secret=YOUR_CLIENT_SECRET \
  -d grant_type=refresh_token \
  -d refresh_token=YOUR_REFRESH_TOKEN
```

where ```YOUR_REFRESH_TOKEN``` is the refresh token from the previous section. The new access token will be the value of ```access_token``` in the response.

### Creating a webhook subscription

1. In a terminal, run the command
    ```
    curl -X POST https://www.strava.com/api/v3/push_subscriptions \
       -F client_id=YOUR_CLIENT_ID \
       -F client_secret=YOUR_CLIENT_SECRET \
       -F callback_url=YOUR_CALLBACK_URL \
       -F verify_token=YOUR_VERIFY_TOKEN
    ```
    where ```YOUR_CALLBACK_URL``` is the URL of your webhook endpoint and ```YOUR_VERIFY_TOKEN``` is a secret string of your choice.
