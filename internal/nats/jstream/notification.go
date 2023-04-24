package clientjstream

import (
	"git.selly.red/Selly-Modules/natsio"
	jsmodel "git.selly.red/Selly-Modules/natsio/js/model"
	jssubject "git.selly.red/Selly-Modules/natsio/js/subject"
	"git.selly.red/Selly-Server/affiliate/external/utils/pjson"
)

// PushNotifications ...
func (s ClientJestStreamPull) PushNotifications(payload []jsmodel.PushNotification) (bool, error) {
	var bytesData = pjson.ToBytes(payload)
	return s.publishWithJetStream(natsio.StreamNameSelly, jssubject.Selly.PushNotification, bytesData)
}
