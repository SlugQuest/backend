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

## Create Task (POST)

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

## Edit Task (PUT)

- **Endpoint**: `/api/v1/task/:id`
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

## Delete Task (DELETE)

- **Endpoint**: `/api/v1/task/:id`
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
