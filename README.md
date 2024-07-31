# Trivia

A live buzzer to play trivia on your local network

## Setup

### Backend

-   Install [Go version 1.22 or higher](https://go.dev/doc/install)
-   Create a .env file with the path go-backend/.env and add a password; PASSWORD="yourpasswordhere"

### Frontend

-   Install [Node.js](https://nodejs.org/en/download/)
-   Edit the file nextjs-frontend/ip.ts to include your device's IP address on your local network

## How to run

-   Run the Makefile using the _make_ command

### Players

-   Visit _HostIP_:3000 while on the same network
-   Enter your name and advance to the next page
-   You will now see a buzzer, clicking it will notify the host

### Host

-   Visit _HostIP_:3000/host to view the players who have buzzed in sorted in chronological order, as well as the time they buzzed in down to the millisecond
-   This page will also display the players ranked by score, score seniority is used as a tie-breaker when scores are equal
-   Go to _HostIP_:3000/control to access controls
    -   Enter the password you put in your .env file
    -   Enter a player's name and the number of points you'd like to give them
        -   A negative value will subtract from their score while a positive value will add. Negative scores are possible
        -   A positive value will append the current question number to the questions the player has gotten correct
        -   A zero or negative value will append a negative question number to the questions the player has gotten correct
    -   Clear will clear the current buzzed in list but will not affect anything else
    -   Next will increment the question number and clear the buzzed in list

### Additional Pages

#### Leaderboard

-   Used to only view the players ranked by score

#### Stats

-   View players ranked by score as well as the questions each player has received points for

#### Buzzed-In

-   See the players as they buzz in for the current question
