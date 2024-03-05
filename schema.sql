CREATE TABLE IF NOT EXISTS UserTable (
    UserID VARCHAR(255) PRIMARY KEY NOT NULL,
    Username VARCHAR(32) NOT NULL,
    Picture TEXT,
    Points INTEGER NOT NULL,
    BossId INTEGER NOT NULL,
    SocialCode VARCHAR(8) UNIQUE NOT NULL, -- created in AddUser()
    FOREIGN KEY (BossId) REFERENCES BossTable(BossID)
);

CREATE TABLE TaskTable (
    TaskID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID VARCHAR(255) NOT NULL,
    Category INTEGER NOT NULL,
    TaskName VARCHAR(255) NOT NULL,
    Description TEXT NOT NULL,
    StartTime DATETIME, -- optional
    EndTime DATETIME, --optional 
    Status VARCHAR(15) CHECK(Status IN ('completed','failed', 'todo')),
    IsRecurring BOOLEAN NOT NULL,
    IsAllDay BOOLEAN NOT NULL,
    Difficulty VARCHAR(15) CHECK(Difficulty IN ('easy','medium', 'hard')), -- Backend does conversion of easy/medium/hard to points
    CronExpression VARCHAR(255) NOT NULL,
    TeamID INTEGER,
    FOREIGN KEY (UserID) REFERENCES UserTable(UserID),
    FOREIGN KEY (Category) REFERENCES Category(CatID),
    FOREIGN KEY (TeamID) REFERENCES TeamTable(TeamID) 
);

CREATE TABLE RecurringLog (
	LogId INTEGER PRIMARY KEY AUTOINCREMENT,
	TaskID INTEGER NOT NULL,
    Status VARCHAR(15) CHECK(Status IN ('completed','failed', 'todo')),
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (TaskID) REFERENCES TaskTable(TaskID)
);

CREATE TABLE Friends (
    userA VARCHAR(255),
    userB VARCHAR(255),
    CONSTRAINT diff_users CHECK (userA <> userB),
    CONSTRAINT no_dup CHECK (userA < userB), -- sort before inserting
    UNIQUE(userA, userB),
    PRIMARY KEY (userA, userB),
    FOREIGN KEY (userA) REFERENCES UserTable(UserID),
    FOREIGN KEY (userB) REFERENCES UserTable(UserID)
);

CREATE TABLE BossTable (
    BossID INTEGER PRIMARY KEY AUTOINCREMENT,
    BossName VARCHAR(255) NOT NULL,
    HEALTH INTEGER NOT NULL,
    BossImage VARCHAR(255) NOT NULL--FileName
);

CREATE TABLE Category (
	CatID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID VARCHAR(255) NOT NULL,
    Name VARCHAR(255) NOT NULL,
	Color INT NOT NULL, -- hexcode? Ask frontend,
    FOREIGN KEY (UserID) REFERENCES UserTable(UserID)
);

CREATE TABLE TrophyTable (
    TrophyID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID VARCHAR(255),
    TrophyName VARCHAR (255),
    FOREIGN KEY (UserID) REFERENCES UserTable(UserID)
);

CREATE TABLE TeamTable (
    TeamID INTEGER PRIMARY KEY AUTOINCREMENT,
    TeamName VARCHAR(255) NOT NULL
);

CREATE TABLE TeamMembers (
    TeamID INTEGER NOT NULL,
    UserID VARCHAR(255) NOT NULL,
    PRIMARY KEY (TeamID, UserID),
    FOREIGN KEY (TeamID) REFERENCES TeamTable(TeamID),
    FOREIGN KEY (UserID) REFERENCES UserTable(UserID)
);
