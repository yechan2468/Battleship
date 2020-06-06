// Code for 'Battleship' game client
// Written by Yechan Lee, 2020.04. -
// Game version 0.0.0 (on development)
// networks, UI/UX, design, etc. will be updated

package main

import (
	"bufio" // get input
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"os"      //get input
	"strconv" // used in getArrangement()
	"time"
)

// GLOBAL VARIABLES & CONSTANTS
var (
	nickname string
	language string = "en"

	// decoderBuffer bytes.Buffer
	// decoder       *gob.Decoder = gob.NewDecoder(&decoderBuffer)
)

type (
	// have to make sure the struct names and elements start with capitalized alphabet
	// unless they will not be exported
	Message struct {
		Header MessageHeader
		Body   MessageBody
	}
	MessageHeader struct {
		MessageType string
		Nickname    string
		Time        string
	}
	MessageBody struct {
		Content string
	}
)

const (
	gameVersion string = "0.0.0"

	// UI
	boardSize int = 10

	sunkTile       int = -2
	hitTile        int = -1
	oceanTile      int = 0 // my territory: empty field / enemy territory: hit missed
	hiddenTile     int = 1 // my territory: hit missed / enemy territory: empry field
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
	cruiserUnicode    string = "\u2726" // diamond symbol
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

	// make gameboard
	var myBoard [boardSize][boardSize]int
	var enemyBoard [boardSize][boardSize]int
	myBoard = clearBoard(myBoard, oceanTile)
	enemyBoard = clearBoard(enemyBoard, hiddenTile)

	nickname = setNickname()
	// if user enters [q], quit the game
	if nickname == "q" || nickname == "Q" {
		fmt.Println("Program terminated. ")
		return
	}

	// arrange the fleets
	printScript("\nWelcome, commander "+nickname+"! It's time for you to arrange our fleets and ready for battle. \n", "충성! 함선을 배치하고 전투를 준비하십시오.")
	myBoard = getArrangement(myBoard)

	// read continuously while program is executing
	var isReady bool = false
	var enemyReady bool = false
	var myTurn bool = false
	go readServer(connection, &myBoard, &enemyBoard, &isReady, &enemyReady, &myTurn)

	for true {
		// write 'i am ready. are you ready?'
		// if response i am ready, game start
		fmt.Print("Making match queue, please wait for another player")
		for i := 0; i < 60; i++ { // wait for 60sec
			time.Sleep(1 * time.Second)
			writeServer("/r", connection, &myTurn) // myTurn value is changed here
			isReady = true
			if isReady && enemyReady {
				break
			} else {
				fmt.Print(".")
				if i >= 59 {
					os.Exit(1)
				}
				continue
			}
		}
		fmt.Print("\n")

		// make match and start game
		fmt.Println("Found Match!")
		fmt.Print("\n")

		// keyboard shortcuts notice
		writeServer("/h", connection, &myTurn)
		fmt.Print("You can get help on keyboard shortcuts by typing \"/h\".\n")

		// turn notice
		// if myTurn {
		// 	fmt.Println("It is your turn! Attack the enemy's territory using command \"/a\".")
		// } else {
		// 	fmt.Println("It's enemy's turn! Please wait for seconds.")
		// }

		// in-game
		var destroyedCheck [5]bool // checks if carrier, battleship, cruiser, submarine, destroyer is destroyed. if destroyed, the element in bool array is "true".
		for {                      // continuously ready input for chatting or command
			var message string
			in := bufio.NewReader(os.Stdin)
			message, err := in.ReadString('\n')
			if err != nil {
				fmt.Println("Failed to read your command. Please try again.")
				continue
			}
			writeServer(message, connection, &myTurn)

			// If my fleets are destroyed, notice to enemy
			result := isDestoyed(&myBoard)
			for idx, count := range result {
				if count == 0 {
					if destroyedCheck[idx] == false { // if not noticed so far
						tmp := strconv.Itoa(idx)
						writeServer("/"+tmp, connection, &myTurn)
						destroyedCheck[idx] = true // noticed
					}
				}
			}

			// quit game
			winner := isDefeat(myBoard, enemyBoard) // checks winner, at every end of each player's turn. return value -1: no winner yet; 0: I win; 1: Enemy win
			if winner >= 0 {
				declareWinner(winner)
			}
			time.Sleep(10 * time.Second)
		}
	}
}
func isDestoyed(myBoard *[boardSize][boardSize]int) [5]int {
	var (
		carrierCount    int = 0
		battleshipCount int = 0
		cruiserCount    int = 0
		submarineCount  int = 0
		destoryerCount  int = 0
	)

	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			switch (*myBoard)[i][j] {
			case carrierTile:
				carrierCount++
			case battleshipTile:
				battleshipCount++
			case cruiserTile:
				cruiserCount++
			case submarineTile:
				submarineCount++
			case destroyerTile:
				destoryerCount++
			}
		}
	}

	var result [5]int = [5]int{carrierCount, battleshipCount, cruiserCount, submarineCount, destoryerCount}
	return result
}

// FUNCTIONS
// * means still on development

// functions related to network
func connectToServer() (connection net.Conn) { // connect to relay server
	connection, err := net.Dial("tcp", "121.159.177.222:8200")
	//connection, err := net.Dial("tcp", "127.0.0.1:8200")

	if err != nil {
		fmt.Println("Failed to connect to server: ", err)
		fmt.Println("Please make sure that the server is available now. If so, please try again.")
		fmt.Println("Battleship client shut down.")
		time.Sleep(5 * time.Second)
		return
	}
	fmt.Println("Successfully connected to server.")
	return
}
func disconnectToServer(connection net.Conn) {
	connection.Close()
}
func readServer(connection net.Conn, pMyBoard *[boardSize][boardSize]int, pEnemyBoard *[boardSize][boardSize]int, pIsReady *bool, pEnemyReady *bool, pMyTurn *bool) { // listen to relay server
	// received data is stored in it
	data := make([]byte, 256)
	// decoded message is stored in it
	var message Message

	for {
		// decoder
		var decoderBuffer bytes.Buffer
		var decoder *gob.Decoder = gob.NewDecoder(&decoderBuffer)

		// read message
		n, err := connection.Read(data)
		if err != nil {
			fmt.Println("Failed to read message from server: ", err)
			fmt.Println("Server connection lost.")
			time.Sleep(5 * time.Second)
			os.Exit(1)
		}

		// decode message
		if n > 0 {
			message = Message{}
			decoderBuffer.Write(data[:n])
			err = decoder.Decode(&message)
			if err != nil {
				fmt.Println("Failed to decode message from client:", err)
				continue
			}
		} else {
			fmt.Println("message ignored because size of message <= 0.")
		}
		decoderBuffer.Reset()

		// ignore if the message is sent by myself
		if message.Header.Nickname == nickname {
			continue
		}

		// handle request
		// for debugging
		// fmt.Print("Received: ", message, "\n")
		if message.Header.MessageType == "chat" {
			fmt.Print("[", message.Header.Time, "] ", message.Header.Nickname, ": ", message.Body.Content, "\n")
		} else if message.Header.MessageType == "command" {
			// var replyMessage string
			switch message.Body.Content[0] {
			case 'a': // if enemy attacked my territory
				// get underattack coordinate
				var row int = int(message.Body.Content[2] - 'A')   // row in integer value
				var col int = int((message.Body.Content[3]) - '1') // column in integer value

				isAttackSucceed := "0"
				var tile int
				printScript("The enemy shot my territory "+message.Body.Content[2:4]+"; \n", "")
				if checkAttackSucceed(row, col, *(pMyBoard)) {
					tile = hitTile // my territory, hit
					isAttackSucceed = "1"
					updateBoard(pMyBoard, row, col, tile)
					showBoard(*pMyBoard, *pEnemyBoard)
					printScript("Hit!\n", "")
				} else {
					tile = hiddenTile // my territory, missed
					updateBoard(pMyBoard, row, col, tile)
					showBoard(*pMyBoard, *pEnemyBoard)
					printScript("Missed!\n", "")
				}
				// writes "/? 1rc" when attack succeeded, else "/? 0rc"
				// rc stands for integer value of row and column
				writeServer("/? "+isAttackSucceed+strconv.Itoa(row)+strconv.Itoa(col), connection, pMyTurn)

				*pMyTurn = true // bool myTurn is set false by default
				fmt.Println("")
				fmt.Println("[Your Turn]")

			case '?': // is enemy attack succeed - attack response handling
				// return attack coordinate, attack succeeded yes/no
				isAttackSucceed := int(message.Body.Content[2] - '0') // 0 or 1 (int)
				var row int = int(message.Body.Content[3] - '0')      // row in integer value
				var col int = int((message.Body.Content[4]) - '0')    // column in integer value
				// if attack succeeded
				if isAttackSucceed == 1 {
					updateBoard(pEnemyBoard, row, col, hitTile) // enemy territory, hit
					showBoard(*pMyBoard, *pEnemyBoard)
					printScript("HIT!\n", "명중!")
				} else {
					updateBoard(pEnemyBoard, row, col, oceanTile) // enemy territory, missed
					showBoard(*pMyBoard, *pEnemyBoard)
					printScript("missed.\n", "타격 실패.")
				}
				*pMyTurn = false
				fmt.Println("")
				fmt.Println("[Enemy's Turn]")

			case 'r': // if enemy ready
				// check whether i am ready; if ready, start game
				*pEnemyReady = true
				if *pIsReady && *pEnemyReady {
					turnInt := rand.Int() % 2 // turnInt[0,1] is used to determine who will start the first turn.
					if turnInt == 1 {         // if random number is 1, I start first and enemy start later
						*pMyTurn = true
					} else {
						*pMyTurn = false
					}
					writeServer("/R "+strconv.Itoa(turnInt), connection, pMyTurn)
				}
			case 'R': // I ready and wait->Enemy ready-> confirm sign; decide first turn
				*pEnemyReady = true
				if message.Body.Content[2] == '0' { // if enemy starts later
					*pMyTurn = true // I start the first turn
				} else {
					*pMyTurn = false
				}
			case '!': // if enemy quit game
				fmt.Println("Enemy Declared Surrender!")
				declareWinner(0)
				time.Sleep(10 * time.Second)
				os.Exit(1)
			// below: enemy fleet sunk notice
			case '0':
				fmt.Println("Enemy Carrier Sunk!")
			case '1':
				fmt.Println("Enemy Battleship Sunk!")
			case '2':
				fmt.Println("Enemy Cruiser Sunk!")
			case '3':
				fmt.Println("Enemy Submarine Sunk!")
			case '4':
				fmt.Println("Enemy Destroyer Sunk!")
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}
func writeServer(message string, connection net.Conn, pMyTurn *bool) { // * out of range occurs on user input
	// ready, attack, isEnemyAttackSucceed, surrender, quitGame

	var (
		encoderBuffer bytes.Buffer
		encoder       *gob.Encoder = gob.NewEncoder(&encoderBuffer)
		msgType, ctt  string       // ctt: content
	)

	if message[len(message)-1] == '\n' { // to strip '\n'
		message = message[0 : len(message)-1]
	}

	if len(message) == 0 { // ignore empty messages
		return
	}

	// distinguish type of message (command, chatting, gameBoard, etc.)
	if message[0] == '/' { // if message is command
		msgType = "command"

		// command type: 'a' for attack, 'q' for quit or surrender, 'r' for ready, 'R' for confirm ready
		if message[1] == 'h' { // help
			fmt.Println("-------------------------------------------------------")
			fmt.Println("[CHATS]")
			fmt.Println("  If you enter what you want to say, you can chat with")
			fmt.Println("  another player.    e.g. Hello!")
			fmt.Println("[COMMANDS]")
			fmt.Println("  /a [coordinate] : attack.	e.g. /a A1")
			fmt.Println("  /q : quit game")
			fmt.Println("-------------------------------------------------------")
		} else if message[1] == 'q' { // quit game\
			writeServer("/!", connection, pMyTurn)
			disconnectToServer(connection)
			os.Exit(1)
		} else if message[1] == 'a' { // attack
			if *pMyTurn {
				fmt.Println("[command] ATTACK " + message[3:])
			} else {
				fmt.Println("It's not your turn. Please wait for another player to attack.")
				return
			}

		} else if message[1] == 'r' || message[1] == 'R' || message[1] == '?' || message[1] == '!' || message[1] == '0' || message[1] == '1' || message[1] == '2' || message[1] == '3' || message[1] == '4' {
			/* Do nothing */
		} else {
			fmt.Println("Failed to find command: ", message)
			return
		}

		ctt = message[1:] // command & command parameter
	} else {
		msgType = "chat"
		ctt = message
		time := (strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()) + ":" + strconv.Itoa(time.Now().Second()))
		fmt.Println("[" + time + "] " + nickname + ": " + ctt)
	}

	// make message struct form
	msg := Message{
		Header: MessageHeader{
			MessageType: msgType,
			Nickname:    nickname,
			Time:        (strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()) + ":" + strconv.Itoa(time.Now().Second())),
		},
		Body: MessageBody{
			Content: ctt,
		},
	}

	// for debugging
	//fmt.Print("Sent: ", msg, "\n")

	// encode message
	err := encoder.Encode(msg)
	if err != nil {
		fmt.Println("Failed to encode message: ", err)
	}

	// write message
	data := make([]byte, 0, 256)                  // bytes slice that consists length 0 and capacity 256
	data = append(data, encoderBuffer.Bytes()...) // to make data size consistent
	_, err = connection.Write(data)
	if err != nil {
		fmt.Println("Failed to write message: ", err)
	}

	encoderBuffer.Reset()
}
func getCurrentUser() (currentUser int) { // *
	currentUser = 0
	return
}
func setNickname() (nickname string) { // get nickname and save it into global variable 'nickname'
	printScript("Enter your nickname to continue, Enter [q] to quit. \n", "닉네임을 입력하세요. [q] 버튼을 누르면 종료됩니다.\n")
	fmt.Print(">> ")

	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	if err != nil {
		nickname = "nickname"
	}
	nickname = line
	nickname = nickname[0 : len(nickname)-1]

	fmt.Print(nickname, "\n")
	return
}

// functions related to gameboard
func checkAttackSucceed(row int, col int, myBoard [boardSize][boardSize]int) bool {
	if myBoard[row][col] == oceanTile || myBoard[row][col] == hiddenTile {
		return false
	} else {
		return true
	}
}
func clearBoard(board [boardSize][boardSize]int, tileNum int) [boardSize][boardSize]int { // clear the whole board tile to designated tile
	for row := 0; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {
			board[row][col] = tileNum
		}
	}
	return board
}
func showBoard(myBoard [boardSize][boardSize]int, enemyBoard [boardSize][boardSize]int) { // boards consist of integers -> convert them to unicode -> print
	println("======================================================================")
	switch language {
	case "kr":
		fmt.Print("          << 아군 영해 >>        \t           << 적군 영해 >>       \n    ")
	case "en":
		fmt.Print("         << My Territory >>      \t         << Enemy Territory >>      \n    ")
	}
	for i := 0; i < boardSize; i++ {
		fmt.Print(i+1, " ")
	}
	fmt.Print("\t    ")
	for i := 0; i < boardSize; i++ {
		fmt.Print(i+1, " ")
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
	println("======================================================================")
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
			col--

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
func updateBoard(board *[boardSize][boardSize]int, row int, col int, tileNum int) {
	// read board changes from server and update
	(*board)[row][col] = tileNum
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
		return 1
	} else if enemyDamagedCount >= 20 {
		return 0
	} else {
		return -1
	}
}

// functions related to UI
func mainScene() { // main scene UI
	printScript("Welcome to 'Battleship' game!\n", "'Battleship' 게임의 메인 화면입니다.\n")
	fmt.Println("                                                                                                    ")
	fmt.Println("                                                                      `                             ")
	fmt.Println("                                                                     ./                             ")
	fmt.Println("                                                          `-:-       .s                             ")
	fmt.Println("                                                          dMMM.     -mM+                            ")
	fmt.Println("                                                         `NMMMh     oyoo ````                       ")
	fmt.Println("                                                         -MMMMM:   -y/:dyNMMy                       ")
	fmt.Println("                `.....`   ```````                 -dy----oMMMMMm-mNNNNNMMMMMN:                      ")
	fmt.Println("                `ohMyo/o/ `+yNhs++o            +oohMMNNMMMMMMMMMMMMMMMMMMMMMMs  `-..`               ")
	fmt.Println("             `---ydyy-.----smdd/-----....----+dMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMy--+MMMN::-----..      ")
	fmt.Println("             :hNMmhhNMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMd.      ")
	fmt.Println("      smmmmmmmmMMNmmNMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMd`       ")
	fmt.Println("      -NMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMm`        ")
	fmt.Println("       .hMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM/         ")
	fmt.Println("                                                                                                    ")
	fmt.Println("                                                                                                    ")
	fmt.Println("====================================================================================================")
	fmt.Println("                                                                                                    ")
	fmt.Println("    MMMMMMm      MMMM    MMMMMMMMM  MMMMMMMMM MMM      MMMMMMM   hNMMm    MMM  MMM  mMMM  MMMMMN ")
	fmt.Println("    MMMM MMMN    MMMM    MMMMMMMMM  MMMMMMMMM MMM      MMMMMMM  MMMMMMMN  MMM  MMM  mMMM  MMM  MMM")
	fmt.Println("    MMMM  MMM   MMMMMM      MMM        MMM    MMM      MMM     hMMMN      MMM  MMM  mMMM  MMM  MMN ")
	fmt.Println("    MMMMdMMN    MMMMMM      MMM        MMM    MMM      MMMMMM   MMMMd     MMMMMMMM  mMMM  MMMNdMM")
	fmt.Println("    MMMMMMMd   NMM  MMN     MMM        MMM    MMM      MMMMMM    dMMMMN   MMMMMMMM  mMMM  MMM ")
	fmt.Println("    MMMM  MMm  MMMMMMMM     MMM        MMM    MMM      MMM          MMMM  MMM  MMM  mMMM  MMM")
	fmt.Println("    MMMM  MMM MMMM  MMMM    MMM        MMM    MMM      MMM     dMMN NMMM  MMM  MMM  mMMM  MMM")
	fmt.Println("    MMMM  MMM MMMM  MMMM    MMM        MMM    MMMMMMMM MMMhMMM dMMN+NMMM  MMM  MMM  mMMM  MMM")
	fmt.Println("    MMMMNMMM  NMMM  MMMN    MMM        MMM    MMMMMMMM MMMMMMM  NMMMMMM   MMM  MMM  mMMM  MMM")
	fmt.Println("                                                                                                    ")
	fmt.Println("====================================================================================================")
	//currentUser := getCurrentUser()
	//fmt.Println("Number of current user: ", currentUser)
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
func declareWinner(winner int) {
	fmt.Println("Game Over!")
	// who is winner?
	switch winner {
	case 0:
		fmt.Println("===========================")
		fmt.Println("          VICTORY          ")
		fmt.Println("===========================")
	case 1:
		fmt.Println("===========================")
		fmt.Println("          DEFEAT           ")
		fmt.Println("===========================")
	}
}
