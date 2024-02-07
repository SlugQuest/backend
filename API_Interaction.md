# Backend API Interactions

## Getting All User Tasks (GET)

- **Endpoint**: `/api/v1/tasks`
- **Description**: Get all tasks for a specific user.
  - **Request Method**: GET
  - **Response**:
    - **Status Code**: 200 OK
    - **Body**: JSON Sample Response
      ```json
      {
        "TaskID": 1,
        "UserID": "user123",
        "Category": "Personal",
        "TaskName": "Go to the Gym",
        "Description": "Exercise for an hour.",
        "StartTime": "2024-01-02T18:00:00Z",
        "EndTime": "2024-01-02T19:00:00Z",
        "IsCompleted": false,
        "IsRecurring": false,
        "IsAllDay": false
      }
      ```

## Get Task by ID (GET)

- **Endpoint**: `/api/v1/task/:id`
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
          "TaskID": 1,
          "UserID": "user123",
          "Category": "Work",
          "TaskName": "Complete Project",
          "Description": "Finish the project by the deadline.",
          "StartTime": "2024-01-01T08:00:00Z",
          "EndTime": "2024-01-01T17:00:00Z",
          "IsCompleted": false,
          "IsRecurring": false,
          "IsAllDay": false
        }
      }
      ```

## Create Task (POST)

- **Endpoint**: `/api/v1/tasks`
- **Description**: Create a new task.
  - **Request Method**: POST
  - **Body**: JSON Sample Request Body
    ```json
    {
      "task": {
        "UserID": "user123",
        "Category": "Work",
        "TaskName": "Complete Project",
        "Description": "Finish the project by the deadline.",
        "StartTime": "2024-01-01T08:00:00Z",
        "EndTime": "2024-01-01T17:00:00Z",
        "IsCompleted": false,
        "IsRecurring": false,
        "IsAllDay": false
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

## Edit Task (PUT)

- **Endpoint**: `/api/v1/tasks/:id`
- **Description**: Edit an existing task.
  - **Request Method**: PUT
  - **URL Parameters**: 
    - `id` (integer): The ID of the task to be edited.
  - **Body**: JSON
    ```json
    {
      "TaskID": 1,
      "UserID": "user123",
      "Category": "Personal",
      "TaskName": "Go to the Gym",
      "Description": "Exercise for an hour.",
      "StartTime": "2024-01-02T18:00:00Z",
      "EndTime": "2024-01-02T19:00:00Z",
      "IsCompleted": true,
      "IsRecurring": false,
      "IsAllDay": false
    }
    ```
  - **Response**:
    - **Status Code**: 200 OK

## Delete Task (DELETE)

- **Endpoint**: `/api/v1/tasks/:id`
- **Description**: Delete a task by ID.
  - **Request Method**: DELETE
  - **URL Parameters**: 
    - `id` (integer): The ID of the task to be deleted.
  - **Response**:
    - **Status Code**: 200 OK

## User login
- **Endpoint**: `/login`
- **Description**: Go **directly** to {backend_url}/login to access this endpoint as it loads request headers to send to Auth0. Do not send a GET to /login or these headers get lost.
- **Response**:
  - **Status Code**: 307
    - Redirects to Auth0 and comes back with another redirect to /callback (to confirm logged in token)

## User logout
- **Endpoint**: `/logout`
- **Description**: Should return back to the host url (i.e `localhost:8080` on manual run). Go **directly** to {backend_url}/logout so the backend performs this logic.
