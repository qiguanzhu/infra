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
 @Time    : 2024/7/13 -- 13:55
 @Author  : bishop ❤️ MONEY
 @Description:  手机号相关处理函数
*/

package xutils

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

/*
ParsePhone
通用的手机号解析，以及检查
手机号格式为：国家区号-手机号
国家区号和手机号使用"-"进行分割
例如： 86-15000000000
没有-号的，按照86进行处理
*/
func ParsePhone(phone string, default86 bool) (area string, number string) {
	area = ""
	if default86 {
		area = "86"
	}

	idx := strings.Index(phone, "-")
	if idx == -1 {
		number = phone
	} else {
		area = phone[:idx]
		number = phone[idx+1:]
	}

	return
}

// WorldPhone 获取带区号的手机号
func WorldPhone(area, number string) string {
	if len(area) == 0 {
		area = "86"
	}
	return area + "-" + number
}

// WorldPhoneFmt
// 输入的是一个可能不带区号，也可能带的号
// 返回一个规范的，肯定带着区号的表示
func WorldPhoneFmt(phone string) string {
	area, number := ParsePhone(phone, true)

	return WorldPhone(area, number)
}

// GetStandardPhone 返回标准手机号，支持加密
func GetStandardPhone(phone string, star bool) string {
	if phone == "" {
		return phone
	}
	phone = WorldPhoneFmt(phone)
	if star && len(phone) > 8 {
		phone = phone[:len(phone)-4-4] + "****" + phone[len(phone)-4:]
	}
	return phone
}

// 解析手机号成area和phone两部分
// 1、规则所有+会被清掉，所有空格会被清除
// 2、如果有-，则第一个-之前的部分为area，剩下的部分将所有-去掉作为phone
// e.g. 86-13111111111 => 86  13111111111
// +86-13111111111 => 86  13111111111
// +86-1311 1111 111 => 86  13111111111
// +86-1311-1111-111 => 86  13111111111
func ParseFullPhone(ctx context.Context, fullphone string, isDefault86 bool) (area, phone string) {
	area = ""
	if isDefault86 {
		area = "86"
	}
	fullphone = strings.TrimSpace(fullphone)
	fullphone = strings.ReplaceAll(fullphone, " ", "")
	fullphone = strings.ReplaceAll(fullphone, "+", "")
	if len(fullphone) == 0 {
		return area, ""
	}
	if strings.Index(fullphone, "-") == -1 {
		return area, fullphone
	}
	split := strings.SplitN(fullphone, "-", 2)
	if len(split) == 2 {
		phone = strings.ReplaceAll(split[1], "-", "")
		return split[0], phone
	}
	return area, ""
}

// 根据手机号得到用于拨号的手机号格式。
// 参考ParsePhone函数格式得到area和phone
// 如果area为空则默认为86
// 得到的结果是+area-phone，e.g. +86-13111111111
// 国家冠字使用+
// 说明：国家冠子是从中国往国外拨打电话时候的前缀统一都是00或者+
// fullphone 手机号，stateCrest国家冠子 00或者+，exclude86 中国的手机号是否排除掉，因为国内打国内不需要国家冠子和86
func FormatPhoneForCall(ctx context.Context, fullphone string, stateCrest string, exclude86 bool) string {
	area, phone := ParseFullPhone(ctx, fullphone, true)
	if phone == "" {
		return ""
	}
	// 如果没有area认为是国内86的号
	if area == "" {
		area = "86"
	}
	if area == "86" && exclude86 {
		stateCrest = ""
		area = ""
	}
	area = strings.TrimSpace(area)
	phone = strings.TrimSpace(phone)
	// 手机号前面带0的需要全部去掉
	phone = strings.Trim(phone, "0")
	// 拨打的时候中间的-会替换掉
	return fmt.Sprintf("%s%s%s", stateCrest, area, phone)
}

type PhoneRegExp struct {
	AreaNumber string `json:"AreaNumber"`
	RegexpAll  string `json:"regexp_all"`
	Region     string `json:"region"`
	RegionCode string `json:"region_code"`
}

// PhoneVerifyReq 手机号验证请求
type PhoneVerifyReq struct {
	Phone      string `json:"phone"`       // 手机号
	RegionCode string `json:"region_code"` // 国家或地区缩写
}

// PhoneVerifyRes 手机号验证结果
type PhoneVerifyRes struct {
	Ok           bool   `json:"ok"`
	RegularPhone string `json:"regular_phone"` // 规则化的手机号,如 86-18812341919,852-18812341919
	AreaNumber   string `json:"area_number"`   // 区号
	Region       string `json:"region"`        // 国家或地区
	RegionCode   string `json:"region_code"`   // 国家或地区缩写
}

// RegexpPhoneVerify
/**
 * @Description:手机号正则验证
 * 	  支持验证的国家和地区有:{中国大陆-CN,中国香港-HK,中国澳门-MO,中国台湾-TW,阿联酋-UAE,澳大利亚-AU,东帝汶-TL,菲律宾-PH,韩国-KR
 * 	  加拿大-CA,柬埔寨-KH,老挝-LA,马来西亚-MY,美国-US,缅甸-MM,日本-JP,泰国-TH,文莱-BN,西班牙-ES,新加坡-SG,新西兰-NZ,印度尼西亚-ID
 * 	  英国-UK,越南-VN},不支持的一律返回false。
 *    仅中国大陆手机号可不带区号。其余必须为+86-phone，0086-phone，86-phone这样的
 * @param req{Phone:手机号,RegionCode:国家或地区缩写,如"CN","HK"} 注:RegionCode为空则遍历所有Region匹配校验
 * @return PhoneVerifyRes{Ok:正确/错误,RegularPhone:规则化的手机号,AreaNumber:区号,Region:手机号所属国家或地区,RegionCode:国家或地区缩写}
 */
func RegexpPhoneVerify(req PhoneVerifyReq) PhoneVerifyRes {
	var info PhoneVerifyRes
	info.Ok = false
	if req.Phone == "" {
		return info
	}
	if req.RegionCode != "" { // 按照国家缩写去匹配正则
		phoneRegValue, ok := PhoneRegExpMap[req.RegionCode]
		if !ok { // key(国家缩写)不存在
			return info
		}
		info.AreaNumber = phoneRegValue.AreaNumber
		info.Region = phoneRegValue.Region
		info.RegionCode = phoneRegValue.RegionCode
		info.Ok, _ = regexp.MatchString(phoneRegValue.RegexpAll, req.Phone)
		if info.Ok {
			info.RegularPhone = RegularPhone(req.Phone)
		}
		return info
	}
	// 没填写国家缩写,遍历匹配支持国家的正则校验
	for _, phoneRegValue := range PhoneRegExpMap {
		info.Ok, _ = regexp.MatchString(phoneRegValue.RegexpAll, req.Phone)
		if info.Ok {
			info.AreaNumber = phoneRegValue.AreaNumber
			info.Region = phoneRegValue.Region
			info.RegionCode = phoneRegValue.RegionCode
			info.RegularPhone = RegularPhone(req.Phone)
			return info
		}
	}
	return info
}

// RegularPhone
/**
 * @Description: 规则化手机号
 * @param phone 必须是已通过正则校验的手机号,如+86-18812341919，0086-18812341919，86-18812341919，18812341919(仅限中国大陆可不带区号)
 * @return string 返回规则化的手机号，如86-18812341919，852-18812341919
 */
func RegularPhone(phone string) string {
	if strings.Contains(phone, "-") { // 存在"-",说明是86-18812341919，0086-18812341919，+86-18812341919三种情况
		splitPhone := strings.Split(phone, "-")
		if len(splitPhone) < 2 { // 安全防护
			return phone
		}
		if strings.Contains(splitPhone[0], "+") {
			splitPhone[0] = strings.Replace(splitPhone[0], "+", "", 1)
		}
		if strings.Contains(splitPhone[0], "00") {
			splitPhone[0] = strings.Replace(splitPhone[0], "00", "", 1)
		}
		return splitPhone[0] + "-" + splitPhone[1]
	} else { // 不存在"-"说明是如188****1902的中国大陆账号
		return "86-" + phone
	}
}

// 手机号码部分前置0去掉
func TrimLeft0PhoneNum(phone string) (area string, number string) {
	area = ""
	idx := strings.Index(phone, "-")
	if idx == -1 {
		number = phone
	} else {
		area = phone[:idx]
		number = phone[idx+1:]
	}
	number = strings.TrimLeft(number, "0")
	return
}

// 手机号码部分前置0去掉
func SplitTrimLeft0PhoneNum(phone string, isDefault86 bool) (area string, number string) {
	area = ""
	if isDefault86 {
		area = "86"
	}
	idx := strings.Index(phone, "-")
	if idx == -1 {
		number = phone
	} else {
		area = phone[:idx]
		number = phone[idx+1:]
	}
	number = strings.TrimLeft(number, "0")
	return
}

// 手机号码部分前置0去掉
func TrimLeft0PhoneNumFormat(phone string, isDefault86 bool) string {
	area, num := SplitTrimLeft0PhoneNum(phone, isDefault86)
	if len(area) == 0 {
		return fmt.Sprintf("%s", num)
	}
	return fmt.Sprintf("%s-%s", area, num)
}

// 手机号脱敏
func AnonymityPhone(phone string, withoutArea bool) string {
	area, num := ParsePhone(phone, true)
	display := "****"
	if len(num) > 6 {
		num = num[:3] + display + num[7:]
	}
	if withoutArea {
		return num
	}
	return fmt.Sprintf("%s-%s", area, num)
}
