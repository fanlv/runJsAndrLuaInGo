package main

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
	"math"
	"strconv"
	"strings"
	"time"
)

//function add(a, b)
//return a+b
//end
var luaCode = `
function setKeyFunc(data)
  local json = require("json")
  local obj = {"a",1,"b",2,"c",3}
  local jsonStr = json.encode(obj)
  local jsonObj = json.decode(jsonStr)

  obj = covertVersionToInt(data["version"])
  return nil

end
`

//result = {}
//print(data["userId"])
//authChatIDsForLua("6598372565515763972",{})
//if (data["userId"] == nil) then
//result["content"] = "请传userId"
//else
//result["content"] = "userID :"..data["userId"]
//isIn = authChatIDsForLua(data["userId"],{"6599854847647485952","oc", "swift"})
//result["isIn"] = isIn
//end
//result["descripetion"] = "测试判断用户id是否在某个群中"
//config = getConfigForLua(data)
//return config

const (
	GetUserInfoError = iota + 10000
	ConvertStringError
	GetAppInfoListError
	GetAppVersionListError
	GetProtoFileError
	GetAppKeysError
	GetConfigError
)


func Double(L *lua.LState) int {
	lv := L.ToInt(1)             /* get argument */
	L.Push(lua.LNumber(lv * 2)) /* push result */
	return 1                     /* number of results */
}

func covertVersionToInt(L *lua.LState) int {

	start0 := time.Now()


	versionStr := L.ToString(1) //1.10.1-alpha
	if len(versionStr) == 0 {
		logrus.Errorf("[covertVersionToInt] error  versionStr is nil ,versionStr : %s", versionStr)
		L.Push(lua.LNumber(0))
		return 1
	}
	numVersionArray := strings.Split(versionStr, "-")
	numVersion := numVersionArray[0]
	versionArray := strings.Split(numVersion, ".")
	versionNum := 0
	versionLen := len(versionArray)
	for i := 0; i < versionLen; i++ {
		vNum, err1 := strconv.Atoi(versionArray[versionLen-1-i])
		if err1 != nil {
			logrus.Errorf("[covertVersionToInt] covert version : %s to int fail. %d", versionStr, i)
			L.Push(lua.LNumber(0))
			return 1
		}
		versionNum += vNum * int(math.Pow10(3*i))
	}
	L.Push(lua.LNumber(versionNum))
	tmp1 := time.Since(start0).Nanoseconds()
	fmt.Println(tmp1)
	return 1
}



func StrPad(input string, padLength int, padString string, padType string) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] + input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
	}

	return output
}

func getConfigForLua(L *lua.LState) int {

	dic := make(map[string]interface{})
	//dic["userId"] = 6561193189653348615
	dic["config2"] = "ad"
	dic["config1"] = "1.0"


	table := L.NewTable()
	for k, v := range dic {
		vNumber, ok := v.(float64)
		if ok {
			L.SetTable(table, lua.LString(k), lua.LNumber(vNumber))
		} else {
			vStr, ok := v.(string)
			if !ok {
				vStr = fmt.Sprint(v)
			}
			L.SetTable(table, lua.LString(k), lua.LString(vStr))
		}
	}

	L.Push(table) /* push result */
	return 1                     /* number of results */
}

func authChatIDsForLua(L *lua.LState) int {
	userIdStr := L.ToString(1)     /* get argument */
	chatIDLuaTable := L.ToTable(2) /* get argument */
	logrus.Debugf("userId : %d , chatId : %s ", userIdStr, chatIDLuaTable)
	chatIDsInGo := gluamapper.ToGoValue(chatIDLuaTable, gluamapper.Option{NameFunc: getNameFunc})
	chatIDArray, ok := chatIDsInGo.([]interface{})
	userId, err1 := strconv.ParseInt(userIdStr, 10, 64)
	if err1 != nil || !ok {
		L.Push(lua.LBool(false))
		logrus.Errorf("[checkUserInGroup] covert userId: %s  to int64 fail", userIdStr)
		return 1
	}
	chatIDs := make(map[int64]bool, 1)
	for _, v := range chatIDArray {
		vStr, ok := v.(string)
		if !ok {
			vStr = fmt.Sprint(v)
		}
		chatId, err2 := strconv.ParseInt(vStr, 10, 64)
		if err2 == nil {
			chatIDs[chatId] = true
		} else {
			logrus.Errorf("[checkUserInGroup] covert chatId : %s to int64 fail", userIdStr)
		}
	}
	fmt.Print(userId,chatIDs)
	//res, _ := clients.AuthChatIDs(userId, chatIDs)
	L.Push(lua.LBool(false)) /* push result */
	return 1               /* number of results */
}



func main() {

	dic := make(map[string]interface{})
	//dic["userId"] = 6561193189653348615
	dic["os"] = "ad"
	dic["version"] = "1.10.1-alpha"
	//dic["version"] = ""

	dic["json"]= "{\"Fid\":3,\"Key\":\"hello\",\"Value\":\"world\",\"AdministratorID\":1,\"CreateTime\":1539509193,\"Status\":0}"


	byteArray, _ := json.Marshal(dic)

	data := make(map[string]interface{})
	err := json.Unmarshal(byteArray, &data)

	fmt.Println(err)

	start0 := time.Now()
	count := 10000
	for i := 0; i < count; i++ {
		LuaTest(dic)
		//fmt.Println(dic)
	}
	tmp1 := time.Since(start0).Nanoseconds() / 1000 / 1000

	start1 := time.Now()
	for i := 0; i < count; i++ {
		//JsTest(dic)
		fmt.Println(dic)
	}
	tmp2 := time.Since(start1).Nanoseconds() / 1000 / 1000
	fmt.Printf("\r\nLuaTest : %d,JsTest : %d", tmp1, tmp2)

}

func LuaTest(dic map[string]interface{}) {
	L := lua.NewState()
	luajson.Preload(L)
	defer L.Close()
	L.SetGlobal("getConfigForLua", L.NewFunction(getConfigForLua)) /* Original lua_setglobal uses stack... */
	L.SetGlobal("authChatIDsForLua", L.NewFunction(authChatIDsForLua)) /* Original lua_setglobal uses stack... */
	L.SetGlobal("covertVersionToInt", L.NewFunction(covertVersionToInt)) /* Original lua_setglobal uses stack... */

	if err := L.DoString(luaCode); err != nil {
		panic(err)
	}
	table := L.NewTable()
	for k, v := range dic {
		vNumber, ok := v.(float64)
		if ok {
			L.SetTable(table, lua.LString(k), lua.LNumber(vNumber))
		} else {
			vStr, ok := v.(string)
			if !ok {
				vStr = fmt.Sprint(v)
			}
			L.SetTable(table, lua.LString(k), lua.LString(vStr))
		}
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("setKeyFunc"),
		NRet:    1,
		Protect: true,
	}, table); err != nil {
		panic(err)
	}
	ret := L.Get(-1) // returned value
	L.Pop(1)         // remove received value
	obj := gluamapper.ToGoValue(ret, gluamapper.Option{NameFunc: getNameFunc})
	obj1,err := convertMapStringInterface(obj)
	fmt.Println(obj1,err)
}

func convertMapStringInterface(obj interface{})(interface{},error){

	objDic, ok := obj.(map[interface{}]interface{})
	if ok {
		result := make(map[string]interface{})
		for k, v := range objDic {
			kStr, ok := k.(string)
			if !ok {
				kStr = fmt.Sprint(k)
			}
			result[kStr] = v
		}
		return result, nil
	}

	objArray, ok := obj.([]interface{})
	result := make([]interface{},0,0)
	for _, v := range objArray {
		objResult,err := convertMapStringInterface(v)
		if err != nil {
			return nil,err
		}
		result = append(result,objResult)
	}
	return result, nil

	return obj, nil
}

func getNameFunc(s string) string {
	return s
}

func JsTest(dic map[string]interface{}) {
	vm := otto.New()
	v, err := vm.Run(`
function testFun(tab) {
	result = {}
	result["key"] = "test"
	result["key1"] = "val2"
 	if(tab["user"]=="test"){
       result["title"]="good"
    }
    if(tab["os"]=="ios"){
        result["url"]="http://www.google.com"
	}else{
        result["url"]="http://www.baidu.com"
    }
	return result
}
`)
	if err == nil {
		fmt.Println(v)
	}
	jsa, err := vm.ToValue(dic)
	if err != nil {
		panic(err)
	}
	result, err := vm.Call("testFun", nil, jsa)
	//fmt.Printf("result: %s", result.String())

	//if result.IsString() {
	//	fmt.Println("result: ", result.String())
	//} else if result.IsNumber() {
	//	fmt.Println("result: ", result.String())
	//} else if result.IsObject() {
	//	fmt.Println("result: ", result.Object())
	//}

	tmpR, err := result.Export()
	fmt.Println("object: ", tmpR)
	//resultDic := make(map[string]interface{})
	//resultDic["data"] = tmpR
	//ret, err := json.Marshal(resultDic)
	//if err == nil {
	//	fmt.Println("msg: ", string(ret))
	//}

	//obj, _ := vm.Object(`obj`)
	//fmt.Println("object:", obj)
	//fmt.Println("object:", obj.Value())

}



