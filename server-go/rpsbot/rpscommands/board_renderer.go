package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/app"
	"bytes"
	"github.com/prizarena/rock-paper-scissors/server-go/rpstrans"
	"fmt"
		"strconv"
	"net/url"
	"github.com/strongo/bots-api-telegram"
	"github.com/prizarena/turn-based"
	"github.com/prizarena/prizarena-public/pamodels"
	"time"
)

func renderRpsBoardMessage(t strongo.SingleLocaleTranslator, tournament *pamodels.Tournament, board turnbased.Board) (m bots.MessageFromBot, err error) {
	if tournament != nil {
		if board.TournamentID != "" && tournament.ID != "" && board.TournamentID != tournament.ShortTournamentID() {
			panic(fmt.Sprintf("board.TournamentID != tournament.ShortTournamentID() => '%v' != '%v'", board.TournamentID, tournament.ShortTournamentID()))
		}
		if board.TournamentID == "" && tournament.ID != "" {
			board.TournamentID = tournament.ShortTournamentID()
		}
	}
	m.IsEdit = true
	m.Format = bots.MessageFormatHTML
	s := new(bytes.Buffer)
	if tournament.TournamentEntity == nil {
		s.WriteString(t.Translate(rpstrans.GameCardTitle))
	} else {
		fmt.Fprintf(s, "⚔ <b>Tournament</b>: %v", tournament.Name)
	}

	s.WriteString("\n" + t.Translate(rpstrans.RulesShort)+ "\n")

	if board.BoardEntity == nil {
		board.BoardEntity = &turnbased.BoardEntity{Round: 1}
	}
	userMoves := board.UserMoves.Strings()
	movesCount := len(userMoves)
	switch {
	case movesCount == 0: // not started
		if len(board.UserNames) == 1 {
			fmt.Fprintf(s, "<b>%v</b> invited to play.\n\n", board.UserNames[0])
		}
		s.WriteString(t.Translate(rpstrans.AskToMakeMove))
	case movesCount == 1 || userMoves[0] == "" || userMoves[1] == "": // one player made move, waiting for second player
		s.WriteString("\n\n")
		firstPlayerName := board.UserNames[0]
		if firstPlayerName == "" {
			firstPlayerName = "Player #1"
		}
		s.WriteString(t.Translate(rpstrans.FirstMoveDoneAwaitingSecond, firstPlayerName))
	case movesCount == 2: // game completed
		s.WriteString("\n\n")
		s.WriteString(fmt.Sprintf("Game completed, moves: %v v %v", userMoves[0], userMoves[1]))
		s.WriteString("\n\n")
		fmt.Fprintf(s, "<b>%v</b>", t.Translate(rpstrans.AskToMakeMove))
	default:
		panic(fmt.Sprintf("len(game.UserMoves) > 2: %v", movesCount))
	}

	if board.ID != "" && board.UserMoves != "" {
		s.WriteString(fmt.Sprintf("\n\nLast update: <i>%v</i>\n\n", time.Now().Format("2006-01-02 15:04:05.000")))
		s.WriteString(`<b>Sponsor</b>: <a href="https://t.me/DebtusBot?start=ref-playRockPaperScissorBot">@DebtusBot</a> - track your debts & assets`)
	}

	m.Text = s.String()
	if board.Lang == "" {
		board.Lang = "en-US"
	}
	callbackPrefix := fmt.Sprintf("bet?l=%v&r=%v&", board.Lang, strconv.Itoa(board.Round))
	if board.ID != "" {
		callbackPrefix += "b=" + url.QueryEscape(board.ID) + "&"
	} else if board.TournamentID != "" {
		callbackPrefix += "t=" + url.QueryEscape(board.TournamentID) + "&"
	}

	inlineMove := func(text, id string) tgbotapi.InlineKeyboardButton {
		return tgbotapi.InlineKeyboardButton{Text: t.Translate(text), CallbackData: callbackPrefix + "m=" + t.Translate(id)}
	}

	// inlineQuery := ""
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			inlineMove(rpstrans.Option1emoji, rpstrans.Option1code),
			inlineMove(rpstrans.Option2emoji, rpstrans.Option2code),
			inlineMove(rpstrans.Option3emoji, rpstrans.Option3code),
		},
		// []tgbotapi.InlineKeyboardButton{
		// 	{
		// 		Text:              whc.Translate(rpstrans.ChallengeFriendCommandText),
		// 		SwitchInlineQuery: &inlineQuery,
		// 	},
		// },
	)

	if board.ID != "" && board.BoardEntity != nil && board.LeftTournament.IsZero() {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
			{
				Text:         "🚫 Leave tournament",
				CallbackData: fmt.Sprintf("%v?board=%v", turnbased.LeaveTournamentCommandCode, board.ID),
			},
		})
	}

	m.Keyboard = keyboard
	return
}
