CREATE TABLE UserTable (
    UserID VARCHAR(255) PRIMARY KEY -- Assuming Auth0 provides a string-based user ID
);

CREATE TABLE TaskTable (
    TaskID INT PRIMARY KEY,
    UserID VARCHAR(255),
    Category VARCHAR(255),
    TaskName VARCHAR(255) NOT NULL,
    Description TEXT,
    StartTime DATETIME,
    EndTime DATETIME,
    IsCompleted BOOLEAN,
    IsRecurring BOOLEAN,
    IsAllDay BOOLEAN,
    FOREIGN KEY (UserID) REFERENCES UserTable(UserID)
);

CREATE TABLE RecurrencePatterns (
    TaskID INT,
    RecurringType VARCHAR(15) CHECK(RecurringType IN ('daily','weekly','monthly')),
    DayOfWeek INT check(DayOfWeek >= 0 and DayOfWeek <= 7),
    DayOfMonth INT check(DayOfMonth >= 0 and DayOfMonth <= 31),
    PRIMARY KEY (TaskID, RecurringType, DayOfWeek, DayOfMonth),
    FOREIGN KEY (TaskID) REFERENCES TaskTable(TaskID)
);
