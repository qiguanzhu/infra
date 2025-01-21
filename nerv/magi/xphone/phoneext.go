/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/7/13 -- 14:32
 @Author  : bishop ❤️ MONEY
 @Description: phone_ext.go
*/

package xutils

var PhoneRegExpMap = map[string]*PhoneRegExp{
	"CN": &PhoneRegExp{
		AreaNumber: "86",
		RegexpAll:  "^((00|\\+){0,1}86-){0,1}((13[0-9])|(14[579])|(15[0-9])|(16[567])|(18[0-9])|(17[0-8])|(19[13589]))\\d{8}$",
		Region:     "中国大陆",
		RegionCode: "CN",
	},
	"TEST": &PhoneRegExp{
		AreaNumber: "86",
		RegexpAll:  "^((00|\\+)?86-)?721\\d{8}$",
		Region:     "中国大陆",
		RegionCode: "CN",
	},
	"HK": &PhoneRegExp{
		AreaNumber: "852",
		RegexpAll:  "^((00|\\+){0,1}852-)(5|6|8|9)\\d{7}$",
		Region:     "中国香港",
		RegionCode: "HK",
	},
	"MO": &PhoneRegExp{
		AreaNumber: "853",
		RegexpAll:  "^((00|\\+){0,1}853-6)\\d{7}$",
		Region:     "中国澳门",
		RegionCode: "MO",
	},
	"TW": &PhoneRegExp{
		AreaNumber: "886",
		RegexpAll:  "^((00|\\+){0,1}886-09)\\d{8}$",
		Region:     "中国台湾",
		RegionCode: "TW",
	},
	"UAE": &PhoneRegExp{
		AreaNumber: "971",
		RegexpAll:  "^(00|\\+){0,1}971-0?5(0|2|5|6)\\d{7}$",
		Region:     "阿联酋",
		RegionCode: "UAE",
	},
	"AU": &PhoneRegExp{
		AreaNumber: "61",
		RegexpAll:  "^(00|\\+){0,1}61-0?[45]\\d{8}$",
		Region:     "澳大利亚",
		RegionCode: "AU",
	},
	"TL": &PhoneRegExp{
		AreaNumber: "670",
		RegexpAll:  "^(00|\\+){0,1}670-7\\d{7}$",
		Region:     "东帝汶",
		RegionCode: "TL",
	},
	"PH": &PhoneRegExp{
		AreaNumber: "63",
		RegexpAll:  "^(00|\\+){0,1}63-0?9\\d{9}$",
		Region:     "菲律宾",
		RegionCode: "PH",
	},
	"KR": &PhoneRegExp{
		AreaNumber: "82",
		RegexpAll:  "^(00|\\+){0,1}82-0?1\\d{9}$",
		Region:     "韩国",
		RegionCode: "KR",
	},
	"CA": &PhoneRegExp{
		AreaNumber: "1",
		RegexpAll:  "^(00|\\+){0,1}1-[1-9]\\d{9}$",
		Region:     "加拿大",
		RegionCode: "CA",
	},
	"KH": &PhoneRegExp{
		AreaNumber: "855",
		RegexpAll:  "^(00|\\+){0,1}855-0?85\\d{6}$",
		Region:     "柬埔寨",
		RegionCode: "KH",
	},
	"LA": &PhoneRegExp{
		AreaNumber: "856",
		RegexpAll:  "^(00|\\+){0,1}856-0?20\\d{7,8}$",
		Region:     "老挝",
		RegionCode: "LA",
	},
	"MY": &PhoneRegExp{
		AreaNumber: "60",
		RegexpAll:  "^(00|\\+){0,1}60-0?1(([02346789]\\d{7})|(1\\d{8}))$",
		Region:     "马来西亚",
		RegionCode: "MY",
	},
	"US": &PhoneRegExp{
		AreaNumber: "1",
		RegexpAll:  "^(00|\\+){0,1}1-[1-9]\\d{9}$",
		Region:     "美国",
		RegionCode: "US",
	},
	"MM": &PhoneRegExp{
		AreaNumber: "95",
		RegexpAll:  "^(00|\\+){0,1}95-0?[1-9]\\d{9}$",
		Region:     "缅甸",
		RegionCode: "MM",
	},
	"JP": &PhoneRegExp{
		AreaNumber: "81",
		RegexpAll:  "^(00|\\+){0,1}81-0?[1-9]\\d{9}$",
		Region:     "日本",
		RegionCode: "JP",
	},
	"TH": &PhoneRegExp{
		AreaNumber: "66",
		RegexpAll:  "^(00|\\+){0,1}66-0?[1-9]\\d{8}$",
		Region:     "泰国",
		RegionCode: "TH",
	},
	"BN": &PhoneRegExp{
		AreaNumber: "673",
		RegexpAll:  "^(00|\\+){0,1}673-((22[89])|(71\\d)|(72[0-3])|(8[1236789]\\d))\\d{4}$",
		Region:     "文莱",
		RegionCode: "BN",
	},
	"ES": &PhoneRegExp{
		AreaNumber: "34",
		RegexpAll:  "^(00|\\+){0,1}34-[6-7]\\d{8}$",
		Region:     "西班牙",
		RegionCode: "ES",
	},
	"SG": &PhoneRegExp{
		AreaNumber: "65",
		RegexpAll:  "^(00|\\+){0,1}65-[89]\\d{7}$",
		Region:     "新加坡",
		RegionCode: "SG",
	},
	"NZ": &PhoneRegExp{
		AreaNumber: "64",
		RegexpAll:  "^(00|\\+){0,1}64-0?2[012579]\\d{7,8}$",
		Region:     "新西兰",
		RegionCode: "NZ",
	},
	"ID": &PhoneRegExp{
		AreaNumber: "62",
		RegexpAll:  "^(00|\\+){0,1}62-0?8\\d{8,10}$",
		Region:     "印度尼西亚",
		RegionCode: "ID",
	},
	"UK": &PhoneRegExp{
		AreaNumber: "44",
		RegexpAll:  "^(00|\\+){0,1}44-0?7\\d{9}$",
		Region:     "英国",
		RegionCode: "UK",
	},
	"VN": &PhoneRegExp{
		AreaNumber: "84",
		RegexpAll:  "^(00|\\+){0,1}84-0?[1-9]\\d{9}$",
		Region:     "越南",
		RegionCode: "VN",
	},
}
