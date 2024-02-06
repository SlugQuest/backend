CREATE TABLE IF NOT EXISTS UserTable (
    UserID VARCHAR(255) PRIMARY KEY NOT NULL,
    Points INTEGER NOT NULL,
    BossId INTEGER NOT NULL,
    FOREIGN KEY (BossId) REFERENCES BossTable(BossID)
);

CREATE TABLE TaskTable (
    TaskID INT PRIMARY KEY AUTOINCREMENT,
    UserID VARCHAR(255) NOT NULL,
    Category INT NOT NULL,
    TaskName VARCHAR(255) NOT NULL,
    Description TEXT NOT NULL,
    StartTime DATETIME, -- optional
    EndTime DATETIME, --optional 
    Status VARCHAR(15) CHECK(Status IN ('completed','failed')),
    IsRecurring BOOLEAN NOT NULL,
    IsAllDay BOOLEAN NOT NULL,
    Points INTEGER NOT NULL, -- Backend does conversion of easy/medium/hard to points
    CronExpression VARCHAR(255) NOT NULL,
    FOREIGN KEY (UserID) REFERENCES UserTable(UserID),
    FOREIGN KEY (Category) REFERENCES Category(CatID)
);


CREATE TABLE BossTable (
    BossID INTEGER PRIMARY KEY AUTOINCREMENT,
    BossName VARCHAR(255) NOT NULL,
    HEALTH INTEGER NOT NULL,
    BossImage VARCHAR(255) NOT NULL--FileName
);

CREATE TABLE Category (
	CatID INTEGER PRIMARY KEY AUTOINCREMENT,
	Color INT NOT NULL -- hexcode? Ask frontend	
);

CREATE TABLE TrophyTable (
    TrophyID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID VARCHAR(255),
    TrophyName VARCHAR (255),
    FOREIGN KEY (UserID) REFERENCES UserTable(UserID)
);

