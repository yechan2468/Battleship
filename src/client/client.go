// Code for 'Battleship' game client
// Written by Yechan Lee, 2020.04. -
// Game version 0.0.0 (on development)
// networks, UI/UX, design, etc. will be updated

package main

import (
	// "encoding/gob"
	"fmt"
	"net"
	"strconv" // used in getArrangement()
	"time"
)

// GLOBAL VARIABLES & CONSTANTS
var (
	nickname string
	language string = "en"
)

const (
	gameVersion string = "0.0.0"

	// UI
	boardSize int = 10

	sunkTile       int = -2
	hitTile        int = -1
	oceanTile      int = 0 // my territory: empty field / enemy territory: hit missed
	hiddenTile     int = 1 // use only in enemy territory
	destroyerTile  int = 2
	submarineTile  int = 3
	cruiserTile    int = 4
	battleshipTile int = 5
	carrierTile    int = 6
	lightLine      int = 7
	line           int = 8
	doubleLine     int = 9

	sunkUnicode       string = "\u2620" // scull
	hitUnicode        string = "\u2739"
	oceanUnicode      string = "\u2591" // light shade
	hiddenUnicode     string = "\u2592" // medium shade
	destroyerUnicode  string = "\u2638" // rudder sign
	submarineUnicode  string = "\u269B" // nuclear symbol
	cruiserUnicode    string = "\u2693" // anchor symbol
	battleshipUnicode string = "\u274B" // mugunghwa symbol
	carrierUnicode    string = "\u272F" // star symbol
	lightLineUnicode                    // scriptstring = "\u2500"
	lineUnicode       string = "\u2501"
	doubleLineUnicode string = "\u2550"

	// scriptReadingSpeed int = 0 // greater is slower
)

// MAIN FUNCTION
func main() {
	// code for checking
	// checkTile()

	// selectLanguage()

	mainScene()
	fmt.Println("Trying to connect to server...")
	connection := connectToServer()
	defer connection.Close()

	// read & write buffer
	readChannel := make(chan string)
	writeChannel := make(chan string)

	// read and write while program is executing
	go readServer(connection, readChannel)
	go writeServer(connection, writeChannel)

	nickname := getNickname()
	// if user enters [q], quit the game
	if nickname == "q" || nickname == "Q" {
		fmt.Println("Program terminated. ")
		return
	}

	for true {
		// make match and start game

		// make gameboard
		var myBoard [boardSize][boardSize]int
		myBoard = clearBoard(myBoard, oceanTile)

		// arrange the fleets
		printScript("\nWelcome, commander "+nickname+"! It's time for you to arrange our fleets and ready to battle. \n", "충성! 함선을 배치하고 전투를 준비하십시오.")
		myBoard = getArrangement(myBoard)

		// in-game
		var enemyBoard [boardSize][boardSize]int
		enemyBoard = clearBoard(enemyBoard, hiddenTile)

		for i := 0; i < 20; i++ {
			var s string
			fmt.Scan(&s)
			writeChannel <- s

			receiver := <-readChannel
			fmt.Println(receiver)
		}

		// quit game
		break
	}
}

// FUNCTIONS
// * means still on development

// functions related to network
func connectToServer() (connection net.Conn) { // *connect to relay server
	connection, err := net.Dial("tcp" /*"121.159.177.222:8200"*/, "127.0.0.1:8200")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully connected to server.")
	return
}
func readServer(connection net.Conn, readChannel chan string) { // * listen to relay server
	// nickname, currentUser, enemyReady, enemyAttack, isAttackSucceed, enemySurrender, enemyQuitGame
	data := make([]byte, 4096)

	for {
		n, err := connection.Read(data)

		// error handling
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("server: " + string(data[:n]))

		// decode

		// distinguish header and body

		// return
		readChannel <- string(data[:n])
		time.Sleep(500 * time.Millisecond)
	}
}
func writeServer(connection net.Conn, writeChannel chan string) { // *
	// nickname, ready, attack, isEnemyAttackSucceed, surrender, quitGame
	for {
		message := <-writeChannel
		n, err := connection.Write([]byte(message))
		fmt.Println(nickname + "sent the message to the server - " + message)
		if err != nil {
			fmt.Print(err)
			return
		}
		_ = n
		time.Sleep(500 * time.Millisecond)
	}
}
func getCurrentUser() (currentUser int) { // * from relay server
	currentUser = 0
	return
}
func getNickname() (nickname string) { // * get nickname and save it into global variable 'nickname'
	printScript("Enter your nickname to continue, Enter [q] to quit. \n", "닉네임을 입력하세요. [q] 버튼을 누르면 종료됩니다.\n")
	fmt.Print(">> ")
	fmt.Scan(&nickname)
	fmt.Print(nickname, "\n")
	return
}
func disconnectToServer() { // * disconnect to relay server
	// client.Close()
}

// in-game network
func attack(x string, y string) {

}
func checkAttackSucceed(myBoard [boardSize][boardSize]int) bool { // *
	// read
	return false
}

// functions related to gameboard
func clearBoard(board [boardSize][boardSize]int, tileNum int) [boardSize][boardSize]int { // clear the whole board tile to designated tile
	for row := 0; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {
			board[row][col] = tileNum
		}
	}
	return board
}
func showBoard(myBoard [boardSize][boardSize]int, enemyBoard [boardSize][boardSize]int) { // boards consist of integers -> convert them to unicode -> print
	switch language {
	case "kr":
		fmt.Print("          << 아군 영해 >>        \t           << 적군 영해 >>       \n    ")
	case "en":
		fmt.Print("         << My Territory >>      \t         << Enemy Territory >>      \n    ")
	}

	for i := 0; i < boardSize; i++ {
		fmt.Print(i+1, "  ")
	}
	fmt.Print("\t    ")
	for i := 0; i < boardSize; i++ {
		fmt.Print(i+1, "  ")
	}
	fmt.Print("\n")

	fmt.Print("")
	for row := 0; row < boardSize; row++ {
		// print my board
		fmt.Printf(" %c  ", 'A'+row)
		for col := 0; col < boardSize; col++ {
			var tileUnicode string = convertToUnicode((myBoard)[row][col])
			fmt.Print(tileUnicode, " ")
		}
		fmt.Print("\t")
		fmt.Printf(" %c  ", 'A'+row)
		// print enemy board
		for col := 0; col < boardSize; col++ {
			var tileUnicode string = convertToUnicode((enemyBoard)[row][col])
			fmt.Print(tileUnicode, " ")
		}
		fmt.Print("\n")
	}
}
func convertToUnicode(tileNum int) (unicode string) { // converts tile number -> corresponding unicode string
	switch tileNum {
	case -2:
		unicode = sunkUnicode
	case -1:
		unicode = hitUnicode
	case 0:
		unicode = oceanUnicode
	case 1:
		unicode = hiddenUnicode
	case 2:
		unicode = destroyerUnicode
	case 3:
		unicode = submarineUnicode
	case 4:
		unicode = cruiserUnicode
	case 5:
		unicode = battleshipUnicode
	case 6:
		unicode = carrierUnicode
	case 7:
		unicode = lightLineUnicode
	case 8:
		unicode = lineUnicode
	case 9:
		unicode = doubleLineUnicode
	}
	return
}
func isArrangeInputValid(row int, col int, dir string, shipSize int, myBoard [boardSize][boardSize]int) (validity bool) {
	validity = false
	switch dir[0] {
	case 'v', 'V':
		if (0 <= row) && (row <= boardSize-shipSize+1) {
			if (0 <= col) && (col <= boardSize) {
				for i := row; i < (row + shipSize); i++ {
					if myBoard[i][col] == oceanTile {
						validity = true
					} else {
						validity = false
						fmt.Print("Already occupied.\n")
						return
					}
				}
				return
			}
		}
		fmt.Print("The fleets must not escape our territory. \n")
		return

	case 'h', 'H':
		if (0 <= row) && (row < boardSize) {
			if (0 <= col) && (col < boardSize-shipSize+1) {
				for i := col; i < (col + shipSize); i++ {
					if myBoard[row][i] == oceanTile {
						validity = true
					} else {
						validity = false
						fmt.Print("Already occupied.\n")
						return
					}
				}
				return
			}
		}
		fmt.Print("The fleets must not escape our territory. \n")
		return
	default:
		fmt.Print("the second parameter must be 'v'('V') or 'h'('H'). \n")
		return
	}
}
func arrangeShip(row int, col int, dir string, tileNum int, myBoard [boardSize][boardSize]int) [boardSize][boardSize]int {
	switch dir[0] {
	case 'v', 'V':
		for i := row; i < row+tileNum; i++ {
			myBoard[i][col] = tileNum
		}
	case 'h', 'H':
		for i := col; i < col+tileNum; i++ {
			myBoard[row][i] = tileNum
		}
	default:
		fmt.Print("Err in func arrangeShip()")
	}
	return myBoard
}
func getArrangement(myBoard [boardSize][boardSize]int) [boardSize][boardSize]int {
	/*
		There are 5 kinds of ships: Carrier, Battleship, Cruiser, Submarine, and destroyer.
		and they respectively have 6*1, 5*1, 4*1, 3*1, 2*1 size.
		The ships are arranged in 10*10 grid gameboard, either horizontally or vertically, and they cannot be overlap.

		In the gameboard(10*10 2-d int array), each integer numbers stands for something like below:
			-1: destroyed ship
			0: empty field (ocean)
			1: hidden field
			2: destroyer, 3: submarine, 4: cruiser, 5: battleship, 6: carrier
		Each tile numbers can be converted to unicode string value, using convertToUnicode().
	*/
	var enemyBoard [boardSize][boardSize]int
	enemyBoard = clearBoard(enemyBoard, hiddenTile)

	showBoard(myBoard, enemyBoard)
	printLine(doubleLine, 40)
	printScript(
		" How to arrange your fleets:\n  1. Select starting point of the fleet.\n  2. Select whether you arrange it horizontally[h] or vertically[v].\n Examples: if you want to put carrier at A1 - A6, type in [A1 h]\n", "함선 배치 방법:\n  1. 뱃머리의 위치를 입력하십시오.\n  2. 가로[h]/세로[v]를 입력하십시오.\n 예시: 항공모함을 A1 ~ A6에 배치하려 한다면, [A1 h]을 입력하십시오.\n")
	printLine(doubleLine, 40)

	var pos, dir string
	var shipList [5]string = [5]string{"destroyer", "submarine", "cruiser", "battleship", "carrier"}
	var krShipList [5]string = [5]string{"구축함", "잠수함", "순양함", "전함", "항공모함"}

ARRANGE:
	for true {
		for i := 0; i < 5; { // arrange ships by orders
			printScript("Arrange your "+shipList[4-i]+": (length: "+(string('0'+6-i))+")\n>> ", "Arrange your "+krShipList[4-i]+": \n>> ")
			fmt.Scan(&pos, &dir)
			row := int(pos[0]) - 'A'        // row index starting from 0
			col, _ := strconv.Atoi(pos[1:]) // column index starting from 0
			col -= 1

			if isArrangeInputValid(row, col, dir, 6-i, myBoard) { // 6-i refers to each length of the ship
				fmt.Print(pos, " ", dir, "\n")
				myBoard = arrangeShip(row, col, dir, 6-i, myBoard) // 6-i refers to tile number of the ship
				showBoard(myBoard, enemyBoard)
				i++
			} else {
				continue
			}
		}

		printScript("Are you sure with the arrangement? confirm: [y] / undo: [n]\n>> ", "배치가 완료되었습니까? 확인: [y] / 취소: [n]")
		var tmp string
		fmt.Scan(&tmp)

		switch tmp[0] { // user selects y/n
		case 'y', 'Y':
			goto END_ARRANGE
		case 'n', 'N':
			myBoard = clearBoard(myBoard, oceanTile)
			showBoard(myBoard, enemyBoard)
			goto ARRANGE
		default:
			fmt.Print("Invalid Input. \n>> ")
			continue
		}
	}
END_ARRANGE:
	printScript("All right, arrangement all done. \n", "알겠습니다. 전투 준비를 완료했습니다!")
	return myBoard
}
func updateBoard(myBoard [boardSize][boardSize]int, enemyBoard [boardSize][boardSize]int) ([boardSize][boardSize]int, [boardSize][boardSize]int) { // *
	// read board changes from server
	// update
	return myBoard, enemyBoard
}
func isDefeat(myBoard [boardSize][boardSize]int, enemyBoard [boardSize][boardSize]int) int { // checks winner, at every end of each player's turn. return value -1: no winner yet; 0: I win; 1: Enemy win
	damagedCount := 0
	enemyDamagedCount := 0

	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			if myBoard[i][j] == hitTile {
				damagedCount++
			}
			if enemyBoard[i][j] == hitTile {
				enemyDamagedCount++
			}
		}
	}

	if damagedCount >= 20 {
		return 0
	} else if enemyDamagedCount >= 20 {
		return 1
	} else {
		return -1
	}
}

// functions related to UI
func mainScene() { // main scene UI
	printScript("Welcome to 'Battleship' game. This is main scene.\n", "'Battleship' 게임의 메인 화면입니다.\n")
	currentUser := getCurrentUser()
	fmt.Println("Number of current user: ", currentUser)
}
func selectLanguage() { // needs more development
	fmt.Print("Select language: [en/kr]\n")
	fmt.Print(">> ")
	fmt.Scan(&language)
	if !(language == "kr" || language == "en") {
		language = "en"
	}
	fmt.Print("\n")
}
func printScript(enScript string, krScript string) {
	var script string
	if language == "kr" {
		script = krScript
	} else {
		script = enScript
	}

	for i := 0; i < len(script); i++ {
		fmt.Printf("%c", script[i])
		if script[i] == '\n' {
			time.Sleep(0 * time.Second)
		}
		time.Sleep(0 * time.Millisecond)
	}
}
func printLine(lineNum int, length int) {
	for i := 0; i < length; i++ {
		fmt.Print(convertToUnicode(lineNum))
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Print("\n")
}
func checkTile() { // just for checking, have no effect on operation
	for i := -2; i < 10; i++ {
		fmt.Print(i, ": ")
		printLine(i, 1)
	}
}