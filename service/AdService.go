package service

import (
	"math/rand"
	"time"
)

//1=首页广告，2=开奖历史广告，3=抽奖详情激励广告
func GetAd(adType int) string {
	var adList []string
	switch adType {
		case 1:
			adList = []string{
				"adunit-edd1b0f37c9d930a",
				"adunit-2bb29f8ddb4eb600",
				"adunit-23299c2f1f1974a0",
				"adunit-d28443e86d48b8d2",
				"adunit-73cb298a0c339231",
				"adunit-888bca1b93a8f629",
				"adunit-12967438ca0af99c",
				"adunit-bef365c3caa91342",
			}
			break
		case 2:
			adList = []string{
				"adunit-d956f27342f353dc",
				"adunit-0b6f59837aab5d11",
				"adunit-90b3aa89e61d2bbe",
				"adunit-416adf39532efc67",
				"adunit-505791df34be0462",
				"adunit-0b9ffb8e6c96f174",
				"adunit-0c47d181e1f30302",
				"adunit-e25ed81e65c49347",
			}
			break
		case 3:
			adList = []string{
				"adunit-a3f14ff6cd7ca3d7",
				"adunit-ef76c0045ef6891a",
				"adunit-58efe24b715a45bf",
				"adunit-1c356cdd4b5f567a",
				"adunit-89669742521f73d6",
				"adunit-3bf601568a0de494",
				"adunit-d0e23026ea75bc23",
				"adunit-8b42b7319d507a66",
			}
			break
		case 4:
			adList = []string{
				"adunit-bec8cc6a4abbcbfd",
				"adunit-a0068ea8277c26b0",
				"adunit-3d97e32f8de4dd08",
				"adunit-3c37fce804325aed",
				"adunit-993f26f47a7e6644",
				"adunit-0a893b23673d132b",
				"adunit-f11b0acff4cefdbf",
				"adunit-3926f2288b4c25d5",
			}
			break
		case 5:
			adList = []string{
				"adunit-aed442de4ea3f8dc",
				"adunit-fb267a11ddcb8050",
				"adunit-0e3cc57bd58dcaef",
				"adunit-a69e292150f0d0ca",
				"adunit-08366384ec43d4c9",
				"adunit-63fe4678f9189b5a",
			}
			break
		default:
			adList = []string{
				"adunit-a3f14ff6cd7ca3d7",
				"adunit-ef76c0045ef6891a",
				"adunit-58efe24b715a45bf",
				"adunit-1c356cdd4b5f567a",
				"adunit-89669742521f73d6",
				"adunit-3bf601568a0de494",
				"adunit-d0e23026ea75bc23",
				"adunit-8b42b7319d507a66",
			}
			break
	}

	if len(adList) <= 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())
	adIndex := rand.Intn(len(adList))
	if adIndex < 0 {
		adIndex = 0
	}else if adIndex > len(adList) {
		adIndex = len(adList) - 1
	}

	return adList[adIndex]
}
