package pw

import "fmt"

func identiconOf(fullName, mainPass string) (string, error) {
	leftArm := []rune{'╔', '╚', '╰', '═'}
	rightArm := []rune{'╗', '╝', '╯', '═'}
	body := []rune{'█', '░', '▒', '▓', '☺', '☻'}
	accessory := []rune{
		'◈', '◎', '◐', '◑', '◒', '◓', '☀', '☁', '☂', '☃', '☄', '★', '☆', '☎', '☏', '⎈', '⌂', '☘',
		'☢', '☣', '☕', '⌚', '⌛', '⏰', '⚡', '⛄', '⛅', '☔', '♔', '♕', '♖', '♗', '♘', '♙',
		'♚', '♛', '♜', '♝', '♞', '♟', '♨', '♩', '♪', '♫', '⚐', '⚑', '⚔', '⚖', '⚙', '⚠', '⌘', '⏎',
		'✄', '✆', '✈', '✉', '✌',
	}

	seed, err := hmacSha256([]byte(mainPass), []byte(fullName))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%c%c%c%c",
		leftArm[int(seed[0])%len(leftArm)],
		body[int(seed[1])%len(body)],
		rightArm[int(seed[2])%len(rightArm)],
		accessory[int(seed[3])%len(accessory)],
	), nil
}
