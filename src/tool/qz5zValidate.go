package tool

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type loginData struct {
	Head struct {
		MsgType string `json:"msgType"`
	} `json:"head"`
	Account     string `json:"account"`
	Psw         string `json:"psw"`
	CheckSum    string `json:"checkSum"`
	DeviceToken string `json:"deviceToken"`
}

type loginRespond struct {
	Head struct {
		MsgType   string `json:"msgType"`
		SessionID string `json:"sessionID"`
	} `json:"head"`
	WorkID string `json:"workID"`
	Result struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"result"`
}

type item struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type infoRespond struct {
	Items [87]item `json:"items"`
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func getCheckSum(account string, psw string) string {
	h := md5.New()
	h.Write([]byte("login" + account + psw))
	cipherStr := h.Sum(nil)
	original := hex.EncodeToString(cipherStr)
	a := strings.ToUpper(original)
	b := a[8:32]
	return b
}

//{"head":{"msgType":"login"},"account":"黄泽源","psw":"1min=59s","checkSum":"4EF043E634A1D001F9D40E57","deviceToken":"123456"}
//`{"head":{"msgType":"login"},"account":"黄泽源","psw":"1min=59s","checkSum":"34A1D001F9D40E574179EACC","deviceToken":"123456"}`

//Validate 验证是否为五中学生并返回ID、sessionID、姓名、年段与班级
func Validate(username string, password string) (bool, string, string, string, string) {
	data := loginData{}
	data.Head.MsgType = "login"
	data.Account = username
	data.Psw = password
	data.CheckSum = getCheckSum(data.Account, data.Psw)
	data.DeviceToken = "123456"
	b, err := json.Marshal(data)
	checkError(err)

	resp, err := http.Post("http://www.qz5z.com//?action=stu", "application/x-www-form-urlencoded", strings.NewReader(url.QueryEscape(string(b))))
	checkError(err)

	defer resp.Body.Close()
	var lr loginRespond
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	err = json.Unmarshal(body, &lr)
	checkError(err)

	var success bool
	var uID, name, class string
	if lr.Result.Code == "0" {
		success = true

		resp, err = http.Post("http://www.qz5z.com//?action=stu", "application/x-www-form-urlencoded", strings.NewReader(`{"head":{"msgType":"getSchoolInfo","sessionID":"`+lr.Head.SessionID+`"}}`))
		checkError(err)

		defer resp.Body.Close()
		var ir infoRespond
		body, err = ioutil.ReadAll(resp.Body)
		checkError(err)
		err = json.Unmarshal(body, &ir)
		checkError(err)
		uID = ir.Items[0].Value
		name = ir.Items[4].Value
		class = "" + ir.Items[25].Value + ir.Items[26].Value
	} else if lr.Result.Code == "1" {
		success = false
		uID = ""
		name = ""
		class = ""
	}

	return success, uID, lr.Head.SessionID, name, class
}

/* {"head":{"msgType":"getSchoolInfo","sessionID":"7612"},"result":{"code":"0"},"items":[{"name":"id","value":"11960"},{"name":"个人
标识码(学籍号)","value":""},{"name":"学校名称","value":""},{"name":"学校标识码","value":""},{"name":"姓名","value":"黄泽源"},{"name":"性别","value":"男"},{"name":"出生年月","value":"20010326"},{"name":"出生地","value":"福建省泉州市鲤城区开元街道县后社区"},{"name":"籍贯","value":"泉州"},{"name":"民族","value":"汉族"},{"name":"国籍/地区","value":"中国"},{"name":"身份证件类型","value":"居民身份证"},{"name":"身份证件号","value":"350502200103260515"},{"name":"港澳台侨外","value":"否"},{"name":"政治面貌","value":"共青团
员"},{"name":"健康状况","value":"健康或良好"},{"name":"照片","value":"/upload/jsda_img/20160903102511901.JPG"},{"name":"姓名拼音","value":""},{"name":"曾用名","value":""},{"name":"身份证件有效期","value":""},{"name":"户口所在地","value":"福建省泉州市丰泽区东海街道云谷社区"},{"name":"户口性质","value":"非农业户口"},{"name":"特长","value":""},{"name":"学籍辅号","value":""},{"name":"班内学
号","value":"8"},{"name":"年级","value":"高中2016级"},{"name":"班级","value":"10班"},{"name":"入学年月","value":""},{"name":"入学
方式","value":""},{"name":"就读方式","value":"走读"},{"name":"住宿类型","value":"走读"},{"name":"学生来源","value":""},{"name":"学籍变动日期","value":""},{"name":"学籍变动情况","value":""},{"name":"现住址","value":"福建省泉州市丰泽区东海街道云谷社区盛世天骄12#201"},{"name":"通信地址","value":"福建省泉州市丰泽区东海街道云谷社区盛世天骄12#201"},{"name":"家庭地址","value":"福建省泉州市丰泽
区东海街道云谷社区盛世天骄12#201"},{"name":"联系电话","value":"22799365"},{"name":"邮政编码","value":"362000"},{"name":"电子信箱","value":""},{"name":"主页地址","value":""},{"name":"是否独生子女","value":"是"},{"name":"是否受过学前教育","value":"是"},{"name":"是否留守儿童","value":"非留守儿童"},{"name":"是否进城务工人员随迁子女","value":"否"},{"name":"是否孤儿","value":"否"},{"name":"是
否烈士或优抚子女","value":"否"},{"name":"随班就读","value":"非随班就读"},{"name":"残疾人类型","value":"无残疾"},{"name":"是否由政
府购买学位","value":""},{"name":"是否需要申请资助","value":"否"},{"name":"是否享受一补","value":"否"},{"name":"上下学距离","value":""},{"name":"上下学交通方式","value":""},{"name":"是否需要乘坐校车","value":""},{"name":"家庭成员或监护人姓名","value":"黄国森"},{"name":"关系","value":"父亲"},{"name":"关系说明","value":""},{"name":"民族","value":"汉族"},{"name":"工作单位","value":"泉州市培
元中学教务处"},{"name":"现住址","value":"福建省泉州市丰泽区东海街道云谷社区盛世天骄12#201"},{"name":"户口所在地","value":"福建省泉州市丰泽区东海街道云谷社区"},{"name":"联系电话","value":"18959893365"},{"name":"是否监护人","value":"是"},{"name":"身份证件类型","value":""},{"name":"身份证件号","value":""},{"name":"职务","value":"学籍管理员"},{"name":"家庭成员或监护人姓名","value":"李秀梅"},{"name":"关系","value":"母亲"},{"name":"关系说明","value":""},{"name":"民族","value":"汉族"},{"name":"工作单位","value":"泉州市纪
委驻泉州市水利局纪检组"},{"name":"现住址","value":"福建省泉州市丰泽区东海街道云谷社区盛世天骄12#201"},{"name":"户口所在地","value":"福建省泉州市丰泽区东海街道云谷社区"},{"name":"联系电话","value":"15392234465"},{"name":"是否监护人","value":"是"},{"name":"身份
证件类型","value":""},{"name":"身份证件号","value":""},{"name":"职务","value":"副组长"},{"name":"是否部队子女","value":"否"},{"name":"毕业学校","value":"泉州实验中学"},{"name":"报名号","value":"3683050767"},{"name":"校服（男装）规格","value":""},{"name":"校服
（女装）规格","value":""},{"name":"奥赛学科第一志愿","value":""},{"name":"奥赛学科第二志愿","value":""},{"name":"是否志愿者","value":"是志愿者"}]} */
