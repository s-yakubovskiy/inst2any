# META API HELP

## Useful links for META

1. https://business.facebook.com/
1. https://developers.facebook.com/tools/explorer/
1. https://developers.facebook.com/docs/instagram-api/getting-started
1. https://adsmanager.facebook.com/adsmanager
1. https://business.facebook.com/latest/home?asset_id=101617782984221&business_id=905555497204855

## Pre

**get user token with graph api of instagram**

**exchange user token for long-lasting one (60 days)**

```bash
curl -i -X GET "https://graph.facebook.com/v17.0/oauth/access_token?grant_type=fb_exchange_token&client_id=975208607010307&client_secret=<secret>&fb_exchange_token=<ur-token>"
```

## Useful commands for API

**get all accounts**:

```bash
curl -s "https://graph.facebook.com/v17.0/me/accounts?access_token=$META_INSTAGRAM_TOKEN" | jq '.'
```

**get all media list from instagram business account id: 17841401562555719**:

```bash
curl -s "https://graph.facebook.com/v17.0/17841401562555719?access_token=$meta_instagram_token&fields=media" | jq '.'
```

**get meta from ID of media (and type)**:

```bash
curl -s "https://graph.facebook.com/v17.0/17980121978264246?fields=media_url,caption,id,media_type,permalink&access_token=$META_INSTAGRAM_TOKEN" | jq '.media_url' | xcp
```

**download with wget to debug**:

```bash
wget "https://scontent.cdninstagram.com/o1/v/t16/f1/m82/A740731EF34B5035A5A0170F77C94EAB_video_dashinit.mp4?<credsinfo>" --output-document reels.mp4
```
