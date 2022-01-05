package Megrez

import (
	"Megrez/util/errorinfo"
	"context"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
	"strings"
)

func Profile(data *dto.WSATMessageData, statu *Statu, ctx context.Context, api openapi.OpenAPI) error {
	if statu.Code == 0 {
		if strings.Contains(data.Content, "/个人资料") {
			profile, err := Db.GetProfile(data.Author.ID)
			if err != nil {
				return err
			}
			_, err = api.PostMessage(ctx, data.ChannelID, errorinfo.SendProfile(data.ID, profile.Uid, profile.Icon, profile.Success, profile.Sum))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
