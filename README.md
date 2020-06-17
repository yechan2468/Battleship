# Battleship

## 0. About Game
- Turn-based strategy game.
- It was made based on the boardgame, 'Battleship'.
- It was made using the Go language and TCP protocol.

## 1. How to Play Game
- 1) Turn on the server. (/build/YOUR_OS/server)
- 2) Turn on the client program(s) and input your nickname.
- 3) Arrange your Ships using commands (described in 2-1 below)
- 4) If two players finished arranging their ships, the match will be made.
- 5) Each players attack another player by turns.
- 6) If one of the player's ship is destroyed completely, it will be announced to both players.
- 7) If all of a player's ship is destroyed, game over!

## 2. User Commands and Other Settings
### 1) client
#### arrangement stage
`[row][column] [direction]`
- Arrange ship. `[row]` value is presented as alphabet, and `[column]` value as an integer. `[direction]` value is either `v`(vertical) or `h`(horizontal)  e.g. `A1 h`, `F3 V`, `d4v`
#### in-game stage
`[anything you want to say]`
- You can chat with another player.
`/a [row][column]`
- Attack enemy territory. The coordinate is given as `[row]`(alphabet), `[column]`(integer). e.g. `/a a1`, `/a f3`
`/q`
- Quit game. If you were playing game with another player, the game is considered a defeat.

### 2) server
`/user`
- Show how many users are connected to the server, connection information, and the ip addresses of the users.
`/time`
- Show the current time in Hour:Minute:Second'Milli Micro Nanoseconds.
`/refresh`
- Refresh server.
`/exit`
- Shut down server.

### 3) other settings
- You can change the size of the gameboard. In the source code(/src/client.go), you can find `const ( ... boardSize int = 10 ... )`. If you replace 10 with other integer numbers, you can decrease or increase the board width. *Caution*: the longest ship size is 6(carrier or mothership), so if the size of the gameboard is smaller than 6 the game will not operate properly.
- You can also change the UI (unicode) of the gameboard. In the same location `const (...)`, there are unicodes tables like `sunkUnicode    string = "\u2620" / hitUnicode string = "\u2739" / ...`. You can change unicode value. It does not affect on gaming, or network. It just changes the way the gameboard is printed on your screen, so don't worry!

## 3. System Requirements
- OS: Windows(64bit, 32bit), Linux(64bit, 32bit)
- The network has to be available

## 4. Bug Reports
- Please report bugs, or improvements to yechan24680@gmail.com
- Thank you for playing!