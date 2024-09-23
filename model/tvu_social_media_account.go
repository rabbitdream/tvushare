// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package model

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

const SocialAccountTableName string = "tvu_social_media_account"

// TvuSocialMediaAccount is the golang structure for table tvu_social_media_account.
type TvuSocialMediaAccount struct {
	Id                   uint   `orm:"id,primary"              json:"id"`                      //
	TvuAccount           string `orm:"tvu_account"             json:"tvu_account"`             // tvu account
	SocialMediaUserId    string `orm:"social_media_user_id"    json:"social_media_user_id"`    // social media user id
	SocialMediaUserName  string `orm:"social_media_user_name"  json:"social_media_user_name"`  // social media user name
	SocialMediaAvatarUrl string `orm:"social_media_avatar_url" json:"social_media_avatar_url"` // social media avatar url
	SocialMediaPlatform  string `orm:"social_media_platform"   json:"social_media_platform"`   // social media platform
	SocialMediaTokenId   int64  `orm:"social_media_token_id"   json:"social_media_token_id"`   // social media oauth2.0 token id
	CreatedTime          int    `orm:"created_time"            json:"created_time"`            // created_time
	UpdatedTime          int    `orm:"updated_time"            json:"updated_time"`            // updated_time
}

func SaveTvuSocialMediaAccount(tvuSocialMediaAccount *TvuSocialMediaAccount, tx *gdb.TX) (int64, error) {
	tsma := &TvuSocialMediaAccount{}
	err := tx.Model(SocialAccountTableName).Where(g.Map{"tvu_account": tvuSocialMediaAccount.TvuAccount, "social_media_platform": tvuSocialMediaAccount.SocialMediaPlatform, "social_media_user_id": tvuSocialMediaAccount.SocialMediaUserId}).Scan(tsma)
	now := int(time.Now().Unix())
	tvuSocialMediaAccount.UpdatedTime = now
	if err != nil || tsma == nil {
		tvuSocialMediaAccount.CreatedTime = now
		res, err := tx.Model(SocialAccountTableName).Data(tvuSocialMediaAccount).Insert()
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	tvuSocialMediaAccount.CreatedTime = tsma.CreatedTime
	tvuSocialMediaAccount.Id = tsma.Id
	_, err = tx.Model(SocialAccountTableName).Data(tvuSocialMediaAccount).Save()
	return int64(tvuSocialMediaAccount.Id), err
}

func FindAuthListByTvuAccount(ctx context.Context, tvuAccount string) ([]TvuSocialMediaAccount, error) {
	var list []TvuSocialMediaAccount
	now := time.Now().Unix()
	err := g.Model("tvu_social_media_account tsma").InnerJoin("oauth_token ot", "tsma.social_media_token_id = ot.id").Fields("tsma.*").Where("tsma.tvu_account = ? AND (ot.updated_time + ot.expires_in > ? OR ot.expires_in = 0)", tvuAccount, now+600).Scan(&list)
	//g.DB().GetAll(ctx, "SELECT tsma.* FROM tvu_social_media_account AS tsma INNER JOIN oauth2_token AS ot ON tsma.social_media_token_id = ot.id WHERE tsma.tvu_account = ? AND (ot.updated_time + ot.expires_in > ? OR ot.expires_in = 0)", tvuAccount, now)
	return list, err
}
