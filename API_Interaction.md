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
  - **Response**:
    - **Status Code**: 200 OK

## Edit Task (PUT)

- **Endpoint**: `/main/blah/tasks/:id`
- **Description**: Edit an existing task.
  - **Request Method**: PUT
  - **URL Parameters**: 
    - `id` (integer): The ID of the task to be edited.
  - **Body**: JSON
    ```json
    {
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
    - **Body**: JSON

## Delete Task (DELETE)

- **Endpoint**: `/main/blah/tasks/:id`
- **Description**: Delete a task by ID.
  - **Request Method**: DELETE
  - **URL Parameters**: 
    - `id` (integer): The ID of the task to be deleted.
  - **Response**:
    - **Status Code**: 200 OK
