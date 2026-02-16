package rabbitmq

import (
	"testing"

	"github.com/quickybrains/bankapp/internal/ports/log/zerolog"
	"github.com/quickybrains/bankapp/internal/ports/notifier"
	"github.com/stretchr/testify/require"
)

func Test1(t *testing.T) {
	testConf := &notifier.Config{}
	logConf := &zerolog.Config{
		Filename: "./testlog",
	}
	_, err := NewNotifier(testConf, zerolog.New(logConf).With("component", "notifier_test"))
	require.NoError(t, err)
}
