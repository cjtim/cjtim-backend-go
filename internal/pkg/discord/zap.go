package discord

import "runtime"

type DiscordZapAddSync struct{}

func (l *DiscordZapAddSync) Write(p []byte) (n int, err error) {
	client, err := newClient()
	if err != nil {
		return 1, err
	}
	pc, _, _, _ := runtime.Caller(5)
	err = client.SendMsg(ErrorsChannel, runtime.FuncForPC(pc).Name(), string(p))
	defer client.Disconnect()
	return 0, err
}

// Close implements io.Closer, and closes the current logfile.
func (l *DiscordZapAddSync) Close() error {
	return nil
}
