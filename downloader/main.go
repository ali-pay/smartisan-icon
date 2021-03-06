package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/tidwall/gjson"
)

//{"body":{"app_icon":[]},"head":{"code":0,"message":""}}
//{"body":{"app_icon":{"com.tencent.weishi":[{"logo":"http://icon.smartisan.com/drawable/com.tencent.weishi/com.tencent.weishi.png","md5":"http://icon.smartisan.com/info/com.tencent.weishi/com.tencent.weishi.xml"},{"logo":"http://icon.smartisan.com/drawable/com.tencent.weishi/ic_launcher.png","md5":"http://icon.smartisan.com/info/com.tencent.weishi/ic_launcher.xml"}]}},"head":{"code":0,"message":""}}
//{"body":{"app_icon":{"com.tencent.weishi":[{"logo":"http://icon.smartisan.com/drawable/com.tencent.weishi/com.tencent.weishi.png","md5":"http://icon.smartisan.com/info/com.tencent.weishi/com.tencent.weishi.xml"},{"logo":"http://icon.smartisan.com/drawable/com.tencent.weishi/ic_launcher.png","md5":"http://icon.smartisan.com/info/com.tencent.weishi/ic_launcher.xml"}],"com.tencent.news":[{"logo":"http://icon.smartisan.com/drawable/com.tencent.news/z_com.tencent.news.png","md5":"http://icon.smartisan.com/info/com.tencent.news/z_com.tencent.news.xml"},{"logo":"http://icon.smartisan.com/drawable/com.tencent.news/com.tencent.news.png","md5":"http://icon.smartisan.com/info/com.tencent.news/com.tencent.news.xml"}]}},"head":{"code":0,"message":""}}
//{"body":{"app_icon":{"com.tencent.weishi":[{"logo":"http://icon.smartisan.com/drawable/com.tencent.weishi/com.tencent.weishi.png","md5":"http://icon.smartisan.com/info/com.tencent.weishi/com.tencent.weishi.xml"},{"logo":"http://icon.smartisan.com/drawable/com.tencent.weishi/ic_launcher.png","md5":"http://icon.smartisan.com/info/com.tencent.weishi/ic_launcher.xml"}],"me.ele":[{"logo":"http://icon.smartisan.com/drawable/me.ele/icon.png","md5":"http://icon.smartisan.com/info/me.ele/icon.xml"}]}},"head":{"code":0,"message":""}}

func main() {
	//dir := `C:\Users\ali-pay\Downloads\Compressed\Snowboard-IconPack-for-Smartisan-OS-master\IconBundles`
	//fileName(dir, "file_name.json")
	//fileUrl("file_name.json", "file_url.json")
	//download("file_url.json")

	//xmlFiles := []string{
	//	"assets/com.sorcerer.sorcery.iconpack-4.6.7994-7994-22320-appfilter.xml",
	//	"assets/me.morirain.dev.iconpack.pure-7.91-1920091604-150901-appfilter.xml",
	//}
	//xmlName(xmlFiles, "xml_name.json")
	//xmlUrl("xml_name.json", "xml_url.json")
	//download("xml_url.json")

	//readOfficialJson()
	//xmlUrl("official_name.json", "official_url.json")
	//download("official_url.json")

	//uniq("icon_name.json")
	// xmlUrl("icon_name.json", "icon_url.json")
	download("icon_url.json")
}

//??????name.json?????????????????????
func uniq(nameJson string) {
	//????????????
	file, err := ioutil.ReadFile(nameJson)
	if err != nil {
		panic(err)
	}
	names := make([]string, 0)
	json.Unmarshal(file, &names)
	//??????
	list := make(map[string]bool)
	for _, v := range names {
		list[v] = false
	}
	//????????????
	data := make([]string, 0, len(list))
	for key := range list {
		data = append(data, key)
	}
	//??????
	sort.StringSlice(data).Sort()
	//????????????
	j, _ := json.MarshalIndent(data, "", "\t")
	ioutil.WriteFile(nameJson, j, os.ModeAppend)
}

//????????????????????????icon
func download(urlJson string) {
	file, err := ioutil.ReadFile(urlJson)
	if err != nil {
		panic(err)
	}
	urls := make([]string, 0)
	json.Unmarshal(file, &urls)

	//???????????????????????????????????????
	suffix := path.Ext(urlJson)
	dir := strings.TrimSuffix(urlJson, suffix)
	err = os.Mkdir(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for i, url := range urls {
		name := strings.ReplaceAll(url, "http://icon.smartisan.com/drawable/", "")
		name = strings.ReplaceAll(name, "/icon_provided_by_smartisan", "")
		//name = strings.ReplaceAll(name, "/launcher_icon", "")
		//name = strings.ReplaceAll(name, "/ic_launcher", "")
		//name = strings.ReplaceAll(name, "/icon", "")
		//name = strings.ReplaceAll(name, "/logo", "")
		name = strings.ReplaceAll(name, "/", ".")
		err := downloadFile(url, dir+"/"+name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("over: %d, download: %s\n", len(urls)-i, url)
	}
}

//??????xml????????????????????????????????????
//files: xml???????????????
//output: ?????????????????????
func xmlName(files []string, output string) {
	//??????map?????????????????????????????????
	apps := make(map[string]bool)
	subs := []string{"item", "calendar"}
	for _, file := range files {
		doc := etree.NewDocument()
		if err := doc.ReadFromFile(file); err != nil {
			panic(err)
		}
		root := doc.SelectElement("resources")
		for _, sub := range subs {
			for _, item := range root.SelectElements(sub) {
				for _, attr := range item.Attr {
					if attr.Key == "component" {
						s := strings.Index(attr.Value, "{")
						e := strings.Index(attr.Value, "/")
						if s < e {
							key := strings.TrimSpace(attr.Value[s+1 : e])
							apps[key] = false
						}
					}
				}
			}
		}
		fmt.Printf("?????????%s, ?????????%d\n", file, len(apps))
	}
	//?????????????????????
	data := make([]string, 0)
	for app := range apps {
		data = append(data, app)
	}
	j, _ := json.MarshalIndent(data, "", "\t")
	ioutil.WriteFile(output, j, os.ModeAppend)
}

type Pgk struct {
	Package string `json:"package"`
}

//??????????????????????????????????????????????????????
//input: ?????????????????????
//output: ?????????????????????
func xmlUrl(input, output string) {
	//??????url.json?????????icon??????
	data := make([]string, 0)

	//??????name.json??????
	file, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}
	pkg := strings.ReplaceAll(string(file), "\t\"", "\t{\"package\":\"")
	pkg = strings.ReplaceAll(pkg, "\",", "\"},")
	pkg = strings.ReplaceAll(pkg, "\n]", "}\n]")
	pkgs := make([]Pgk, 0)
	json.Unmarshal([]byte(pkg), &pkgs)

	//????????????10000???
	for start := 0; start < len(pkgs); start += 9999 {
		end := start + 9999
		if end >= len(pkgs) {
			end = len(pkgs) - 1
		}
		body, _ := json.Marshal(pkgs[start:end])
		req, _ := http.NewRequest(
			"POST",
			"http://setting.smartisan.com/app/icon/",
			bytes.NewBuffer(body))
		req.Header.Add("content-type", "application/json")
		defer req.Body.Close()
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		result, _ := ioutil.ReadAll(resp.Body)
		apps := gjson.GetBytes(result, "body.app_icon").Map()
		for _, app := range apps {
			for _, logo := range app.Get("#.logo").Array() {
				data = append(data, logo.String())
			}
		}
	}

	//??????url.json??????
	j, _ := json.MarshalIndent(data, "", "\t")
	ioutil.WriteFile(output, j, os.ModeAppend)
}

//??????Snowboard-IconPack-for-Smartisan-OS?????????icon????????????????????????????????????
//dir: icon???????????????
//output: ?????????????????????
func fileName(dir, output string) {
	file, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	data := make([]string, 0)
	for _, v := range file {
		name := v.Name()
		data = append(data, name[:len(name)-4]) //??????.png??????
	}
	j, _ := json.MarshalIndent(data, "", "\t")
	ioutil.WriteFile(output, j, os.ModeAppend)
}

//??????????????????????????????????????????????????????
//input: ?????????????????????
//output: ?????????????????????
func fileUrl(input, output string) {
	file, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}
	data := make([]string, 0)
	json.Unmarshal(file, &data)
	for i, v := range data {
		//???????????????https://github.com/Sunbelife/get_smartisan_icon_pack/blob/master/download_icon.py
		//???????????????"http://icon.smartisan.com/drawable/" + bundle_id + "/icon_provided_by_smartisan.png"
		data[i] = "http://icon.smartisan.com/drawable/" + v + "/icon_provided_by_smartisan.png"
	}
	j, _ := json.MarshalIndent(data, "", "\t")
	ioutil.WriteFile(output, j, os.ModeAppend)
}

//???????????????app????????????????????????????????????
func readOfficialJson() {
	file, err := ioutil.ReadFile("./assets/apps_category.json")
	if err != nil {
		panic(err)
	}
	names := map[string]interface{}{}
	err = json.Unmarshal(file, &names)
	if err != nil {
		panic(err)
	}
	data := make([]string, 0)
	for name := range names {
		data = append(data, name)
	}
	j, _ := json.MarshalIndent(data, "", "\t")
	ioutil.WriteFile("official_name.json", j, os.ModeAppend)
}
