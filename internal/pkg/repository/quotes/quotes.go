package quotes

import (
	"context"
	"errors"
	"math/rand"

	"github.com/rs/zerolog"

	"github.com/lisgo88/faraway-test/internal/pkg/repository"
)

type Quotes struct {
	Data []string

	logger zerolog.Logger
}

func New(_ context.Context, log zerolog.Logger) repository.Quotes {
	return &Quotes{
		Data: []string{
			"I didn't think there was anything in the universe more important than homework.",
			"Music; worship and prayer to God.",
			"My life, my strength and my time are my greatest riches.",
			"You must always remember... Whatever their bodies do affects their souls.",
			"Take care of your words and the words will take care of you.",
			"The only way to get what you want is to make them more afraid of you than they are of each other.",
			"Oh darling, your only too wild, to those whom are to tame, don't let opinions change you.",
			"If I disagree with you sometimes, it's because I have a mind of my own.",
			"Life is an affair of mystery; shared with companions of music, dance and poetry.",
			"How dare a person tell a woman, how to dress, how to talk, how to behave! Any being who does that, is no human.",
			"Through synergy of intellect, artistry and grace came into existence the blessing of a dancer.",
			"If you do not have control over your mouth, you will not have control over your future.",
			"In a world of words, anything is possible...",
			"Life is just a slide. Back and forth between loving and leaving, remembering and forgetting, holding on and letting go.",
		},
		logger: log,
	}
}

func (q *Quotes) GetQuote() (string, error) {
	randElem := rand.Intn(len(q.Data))

	if randElem > len(q.Data) {
		return "", errors.New("out of range")
	}

	return q.Data[randElem], nil
}
