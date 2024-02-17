
# Backend API Interactions

## Protected endpoints

All the endpoints below are protected by authenticating the user's session cookies before allowing requests or redirects to happen.

Authentication happens at separate endpoints ([see section on Authentication endpoints](#authentication-endpoints)).

In your requests, make sure to include credentials with the `credentials: 'include'` parameter so these are passed on appropriately:
```js
const response = await fetch(backend/protected, {
  method: 'GET',
  credentials: 'include',
});
```

### Getting All User Tasks (GET)

- **Endpoint**: `/api/v1/tasks`
- **Description**: Get all tasks for a specific user.
  - **Request Method**: GET
  - **Response**:
    - **Status Code**: 200 OK
    - **Body**: JSON Sample Response
      ```json
      {
         "task": {
            "TaskID":         1,
            "UserID":         "testUserId",
            "Category":       "yo",
            "TaskName":       "New Task",
            "Description":    "Description of the new task",
            "StartTime":      "2024-01-01T08:00:00Z",
            "EndTime":        "2024-01-01T17:00:00Z",
            "Status":         "completed",
            "IsRecurring":    false,
            "IsAllDay":       false,
            "Difficulty":     "easy",
            "CronExpression": "" //for now, recurring functions are not supported
        }
      }
      ```

### Get Task by ID (GET)

- **Endpoint**: `/api/v1/task/:id`
  - **NOTE**: Replace `:id` with the actual TaskID.
- **Description**: Get a task by ID.
  - **Request Method**: GET
  - **Parameters**: 
    - `id` (integer): The ID of the task.
  - **Response**:
    - **Status Code**: 200 OK
    - **Body**: JSON
      ```json
      {
         "task": {
            "TaskID":         1,
            "UserID":         "testUserId",
            "Category":       "yo",
            "TaskName":       "New Task",
            "Description":    "Description of the new task",
            "StartTime":      "2024-01-01T08:00:00Z",
            "EndTime":        "2024-01-01T17:00:00Z",
            "Status":         "completed",
            "IsRecurring":    false,
            "IsAllDay":       false,
            "Difficulty":     "easy",
            "CronExpression": "" //for now, recurring functions are not supported
        }
      }
      ```

### Get User Tasks within a Time Range (GET)

- **Endpoint**: `/api/v1/userTasks/:id/:start/:end`
  - **NOTE**: Replace `:id`, `:start`, and `:end` with the actual user ID, start time, and end time respectively.
  - Ex.: `/api/v1/userTasks/123/2024-02-09T00:00:00Z/2024-02-10T00:00:00Z`
- **Description**: Get tasks for the user within a specified time range.
  - **Request Method**: GET
  - **URL Parameters**:
    - `id` (string): The ID of the authenticated user.
    - `start` The start time of the time range.
    - `end` The end time of the time range.
  - **Response**:
    - **Status Code**: 200 OK
    - **Body**: JSON Sample Response Body
      ```json
      {
        "list": [
          {
            "TaskID": 123,
            "UserID": "user123",
            "Category": "Party time",
            "TaskName": "Party",
            "StartTime": "2024-02-09T08:00:00Z",
            "EndTime": "2024-02-09T17:00:00Z",
            "Status": "todo",
            "IsRecurring": false,
            "IsAllDay": false
          },
          // ... more task previews
        ]
      }
      ```


### Get User Points (GET)

- **Endpoint**: `/api/v1/userPoints`
- **Description**: Retrieve the points associated with the authenticated user.
  - **Request Method**: GET
  - **Response**:
    - **Status Code**: 200 OK
    - **Body**: JSON Sample Response Body
      ```json
      {
        "points": 42
      }
      ```

  ### Get User Information (GET)
- **Endpoint**: `/api/v1/user`
- **Description**: Get public user information that can be displayed to user
  - **Request Method**: GET
  - **Response**:
    - **Status Code**: 200 OK
    - **Body**: JSON Sample Response Body
      ```json
      {
        "picture": "https://imagehosterlol/lol.png",
        "points": 42,
        "username": "sluggo"
      }
      ```

### Get Current Boss Health (GET)

- **Endpoint**: `/api/v1/getBossHealth`
- **Description**: Get the current health of the boss associated with the authenticated user.
  - **Request Method**: GET
  - **Response**:
    - **Status Code**: 200 OK
    - **Body**: JSON Sample Response Body
      ```json
      {
        "curr_boss_health": 30
      }
      ```

### Create Category (PUT)

- **Endpoint**: `/api/v1/makeCat`
- **Description**: Create a new category for the authenticated user.
  - **Request Method**: PUT
  - **Request Headers**:
    - `Content-Type: application/json`
  - **Request Body**: JSON Sample Request Body
    ```json
    {
      "UserID": "user123",
      "Name": "Personal",
      "Color": 128
    }
    ```
  - **Response**:
    - **Status Code**: 200 OK
    - **Body**: JSON Sample Response Body
      ```json
      {
        "message": "Success",
        "catID": 789
      }
      ```

### Edit Task (PUT)

- **Endpoint**: `/api/v1/task/:id`
  - **NOTE**: Replace `:id` with the actual TaskID.
- **Description**: Edit an existing task.
  - **Request Method**: PUT
  - **URL Parameters**: 
    - `id` (integer): The ID of the task to be edited.
  - **Body**: JSON
    ```json
      {
         "task": {
            "TaskID":         1,
            "UserID":         "testUserId",
            "Category":       "yo",
            "TaskName":       "New Task",
            "Description":    "Description of the new task",
            "StartTime":      "2024-01-01T08:00:00Z",
            "EndTime":        "2024-01-01T17:00:00Z",
            "Status":         "completed",
            "IsRecurring":    false,
            "IsAllDay":       false,
            "Difficulty":     "easy",
            "CronExpression": "" //for now, recurring functions are not supported
        }
      }
    ```
  - **Response**:
    - **Status Code**: 200 OK


### Pass/Fail Task (POST)

- **Endpoint**: `/api/v1/passtask/:id` and `/api/v1/failtask/:id`
  - **NOTE**: Replace `:id` with the actual TaskID.
  - Ex.: `/api/v1/passtask/123` or `/api/v1/failtask/456`
- **Description**: Mark a task as completed (pass) or uncompleted (fail).
  - **Request Method**: POST
  - **URL Parameters**:
    - `id` (integer): The ID of the task to be marked as completed or uncompleted.
  - **Response**:
    - **Status Code**: 200 OK
    - **Body**: JSON Sample Response Body
      ```json
      {
        "message": "Success"
      }
      ```

### Create Task (POST)

- **Endpoint**: `/api/v1/task`
- **Description**: Create a new task.
  - **Request Method**: POST
  - **Body**: JSON Sample Request Body
    ```json
    {
         "task": {
            "UserID":         "testUserId",
            "Category":       "yo",
            "TaskName":       "New Task",
            "Description":    "Description of the new task",
            "StartTime":      "2024-01-01T08:00:00Z",
            "EndTime":        "2024-01-01T17:00:00Z",
            "Status":         "completed",
            "IsRecurring":    false,
            "IsAllDay":       false,
            "Difficulty":     "easy",
            "CronExpression": "" //for now, recurring functions are not supported
        }
      }
    ```
  - **Response**:
    - **Status Code**: 200 OK
    - **Body**: JSON Sample Response Body
    ```json
    {
      "message": "Success",
      "taskID": 503
    }
    ```

### Delete Task (DELETE)

- **Endpoint**: `/api/v1/task/:id`
  - **NOTE**: Replace `:id` with the actual TaskID.
- **Description**: Delete a task by ID.
  - **Request Method**: DELETE
  - **URL Parameters**: 
    - `id` (integer): The ID of the task to be deleted.
  - **Response**:
    - **Status Code**: 200 OK


## Authentication endpoints

### User login
- **Endpoint**: `/login`
- **Description**: Go **directly** to `{backend_url}/login` to access this endpoint as it loads request headers to send to Auth0. **Do not send a GET** to `/login` or these headers get lost.
- **Response**:
  - **Status Code**: 307
    - Redirects to Auth0 and comes back with another redirect to `backend/callback` (to confirm logged in token)
    - You do not need to route to `/callback`

### User logout
- **Endpoint**: `/logout`
- **Description**: Should return back to the host url (i.e `localhost:8080` on manual run). Go **directly** to `{backend_url}/logout` so the backend performs this logic.

### User signup
- **Endpoint**: `/signup`
- **Description**: Go **directly** to `{backend_url}/signup` to access this endpoint as it loads request headers to send to Auth0. **Do not send a GET** to `/signup` or these headers get lost.
- **Response**:
  - **Status Code**: 307
  - Similar behavior to `/logout` but also adds the newly registered user to the Auth0 database
