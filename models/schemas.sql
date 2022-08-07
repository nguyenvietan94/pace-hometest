DROP TABLE IF EXISTS merchants;
DROP TABLE IF EXISTS members;

CREATE TABLE merchants (
    merchantID SERIAL,
    name VARCHAR(255) NOT NULL,
    age INT,
    location TEXT,
    PRIMARY KEY (merchantID)
);

CREATE TABLE members (
    memberID SERIAL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(320) NOT NULL,
    merchantID INT,
    PRIMARY KEY (memberID),
    FOREIGN KEY (merchantID) REFERENCES merchants(merchantID) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX email_idx ON members(email);
