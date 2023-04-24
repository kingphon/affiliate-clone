package clientjstream

import "git.selly.red/Selly-Modules/natsio"

// ClientJestStreamPull ...
type ClientJestStreamPull struct{}

// publishWithJetStream ...
func (ClientJestStreamPull) publishWithJetStream(streamName, subject string, data []byte) (isPublished bool, err error) {
	if err = natsio.GetJetStream().Publish(streamName, subject, data); err != nil {
		return
	}
	isPublished = true
	return
}
