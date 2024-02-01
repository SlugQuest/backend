# Backend API Interactions

## Getting All User Tasks (GET)

- **Endpoint**: `/main/blah/tasks`
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

- **Endpoint**: `/main/blah/task/:id`
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

- **Endpoint**: `/main/blah/tasks`
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

- **Endpoint**: `/main/blah/tasks/:id`
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

- **Endpoint**: `/main/blah/tasks/:id`
- **Description**: Delete a task by ID.
  - **Request Method**: DELETE
  - **URL Parameters**: 
    - `id` (integer): The ID of the task to be deleted.
  - **Response**:
    - **Status Code**: 200 OK

## User login (GET)
- **Endpoint**: `/login`
- **Description**: Redirects to Auth0's universal login page, then to the `/main/blah/tasks` endpoint after successful login.
  - **Request Method**: GET
- **Body**: None, this redirects to the Auth0 Universal Login page.

## User logout (GET)
- **Endpoint**: `/logout`
- **Description**: Should return back to the host url (i.e `localhost:8080` on manual run)
  - **Request Method**: GET
- **Body**: None

<!-- ## TBA: User info (GET)
In progress, not entirely setup
- **Endpoint**: `/user`
- **Description**: Returns an html page with user information
  - **Request Method**: GET
- **Must be done AFTER a successful login since it depends on user cookies** -->
