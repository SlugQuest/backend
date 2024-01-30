CREATE TABLE UserTable (
    UserID VARCHAR(255) PRIMARY KEY NOT NULL-- Assuming Auth0 provides a string-based user ID
);

CREATE TABLE TaskTable (
    TaskID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID VARCHAR(255) NOT NULL,
    Category VARCHAR(255) NOT NULL,
    TaskName VARCHAR(255) NOT NULL,
    Description TEXT NOT NULL,
    StartTime DATETIME,
    EndTime DATETIME,
    IsCompleted BOOLEAN NOT NULL,
    IsRecurring BOOLEAN NOT NULL,
    IsAllDay BOOLEAN NOT NULL,
    FOREIGN KEY (UserID) REFERENCES UserTable(UserID)
);

CREATE TABLE RecurrencePatterns (
    TaskID INT NOT NULL,
    RecurringType VARCHAR(15) CHECK(RecurringType IN ('daily','weekly','monthly')),
    DayOfWeek INT check(DayOfWeek >= 0 and DayOfWeek <= 7),
    DayOfMonth INT check(DayOfMonth >= 0 and DayOfMonth <= 31),
    PRIMARY KEY (TaskID, RecurringType, DayOfWeek, DayOfMonth),
    FOREIGN KEY (TaskID) REFERENCES TaskTable(TaskID)
);
