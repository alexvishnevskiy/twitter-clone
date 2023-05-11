CREATE DATABASE IF NOT EXISTS twitter;
USE twitter;

# TODO: think how to make authorization without storing password
CREATE TABLE IF NOT EXISTS User (
    user_id INT NOT NULL AUTO_INCREMENT,
    nickname VARCHAR(15) NOT NULL,
    first_name VARCHAR(10) NOT NULL,
    last_name VARCHAR(15) NOT NULL,
    email VARCHAR(20) NOT NULL,
    password VARCHAR(20) NOT NULL,
    PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS Tweets (
    tweet_id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    retweet_id INT NULL,
    content VARCHAR(500) NOT NULL,
    media_url VARCHAR(50) NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (tweet_id),
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Likes (
    user_id INT NOT NULL,
    tweet_id INT NOT NULL,
    PRIMARY KEY (user_id, tweet_id),
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    FOREIGN KEY (tweet_id) REFERENCES Tweets(tweet_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Followers (
    user_id INT NOT NULL,
    following_id INT NOT NULL,
    PRIMARY KEY (user_id, following_id),
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES User(user_id) ON DELETE CASCADE
);
