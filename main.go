package main

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/expectedsh/go-sonic/sonic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Mongo2Sonic()
	// Insert2Sonic()
	SonicQuery()
	// GetRecipeById()
	// GetMaterialByUrl("https://www.haodou.com/recipe/160294")

}

// GetRecipeList 获取菜单列表
func GetRecipeList() {
	file, err := os.OpenFile("abc.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	offset := 0

	for i := 1; i <= 250; i++ {
		currentTime := strconv.FormatInt(time.Now().Unix(), 10)
		offsetStr := strconv.Itoa(offset)
		secret := "Sigere127bb33345b5e3c9b50ed0be4e35da8_HOP_.actionapi.www.recipe.category_HOP_.current_time" + currentTime + "_HOP_.secret_id5722f877e4b0d4512e3fd872_HOP_.version1.0.0adcode100000appid100frommvuehduid0last%7B%22offset%22%3A" + offsetStr + "%2C%22limit%22%3A40%7DmoduleId5d35709cfd96c61a103a13c2numbers%5B%5Duid0uuid0vc177vn1.0.01bc0d50feafb484b863d4100a561a9cf"
		offset = (i - 1) * 40
		has := md5.Sum([]byte(secret))
		sign := fmt.Sprintf("%x", has)
		// fmt.Println(sign)

		data := "numbers=%5B%5D&moduleId=5d35709cSonicQueryfd96c61a103a13c2&_HOP_=%7B%22version%22%3A%221.0.0%22%2C%22action%22%3A%22api.www.recipe.category%22%2C%22secret_id%22%3A%225722f877e4b0d4512e3fd872%22%2C%22current_time%22%3A" + currentTime + "%2C%22sign%22%3A%22" + sign + "%22%7D&from=mvue&adcode=100000&appid=100&Siger=e127bb33345b5e3c9b50ed0be4e35da8&uuid=0&uid=0&hduid=0&vc=177&vn=1.0.0&last=%7B%22offset%22%3A" + offsetStr + "%2C%22limit%22%3A40%7D"

		url := "https://vhop.haodou.com/hop/router/rest.json"
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
		req.Header.Set("Accept", "application/json, text/plain")
		req.Header.Set("Accept-Language", "Accept-Language")
		req.Header.Set("Host", "vhop.haodou.com")
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			panic(err) // > div:nth-child(" + strconv.Itoa(j) + ") > a > div > p.name
		}
		defer resp.Body.Close()
		if le, err := io.Copy(file, resp.Body); err != nil {
			panic(err)
		} else {
			fmt.Println(le)
		}
		file.WriteString("\n")

		time.Sleep(time.Second)

	}
}

func GetRecipeById() {
	var wg sync.WaitGroup

	// 连接mongodb
	url := "mongodb://mws_mongo:mws_mongo@127.0.0.1:27017/mws"
	clientOptions := options.Client().ApplyURI(url)

	conte, _ := context.WithTimeout(context.Background(), time.Second*5)

	client, err := mongo.Connect(conte, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	coll := client.Database("mws").Collection("recipe_list")

	// 读取文件记录
	file, err := os.Open("abc.record")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	// 全局go程池
	goPoll := make(chan bool, 40)

	count := 0

	for {
		buf, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			break
		}

		a := make(map[string]interface{})
		if err := json.Unmarshal(buf, &a); err != nil {
			panic(err)
		}

		b := a["data"].(map[string]interface{})

		c := b["dataset"].([]interface{})

		for _, v := range c {

			if count >= 1713 {
				goPoll <- true
				wg.Add(1)
				go func(record interface{}) {
					defer func() {
						_ = <-goPoll
					}()
					defer wg.Done()

					// 插入mongodb里面的数据 （json)
					m := make(map[string]interface{})

					d := record.(map[string]interface{})
					// title := d["title"].(string)
					id := d["id"].(float64)

					url := "https://www.haodou.com/recipe/" + strconv.FormatFloat(id, 'f', -1, 64)
					fmt.Println(url)

					client := &http.Client{}
					req, err := http.NewRequest(http.MethodGet, url, nil)
					if err != nil {
						fmt.Println(err)
						return
					}

					req.Header.Set("Accept", "text/html, application/xhtml+xml")
					req.Header.Set("Accept-Language", "zh,zh-CN")
					req.Header.Set("Host", "www.haodou.com")
					req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36")

					resp, err := client.Do(req)
					if err != nil {
						fmt.Println(err)
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode != 200 {
						fmt.Println("err")
						return
					}

					Doc, err := goquery.NewDocumentFromReader(resp.Body)
					if err != nil {
						fmt.Println(err)
						return
					}

					// 主料
					MaterialData := make([]string, 0)
					Doc.Find("#__layout > div > div > article > div > div.recipe-left > div:nth-child(1) > div.practice > div.ingredient > div.paixu > div.ingredients ").Each(func(i int, selection *goquery.Selection) {
						MaterialData = append(MaterialData, selection.Find("a > div > p.name").Text())
					})
					m["material"] = MaterialData

					// 辅料
					SecMaterialData := make([]string, 0)
					Doc.Find("#__layout > div > div > article > div > div.recipe-left > div:nth-child(1) > div.practice > div.accessories > div.paixu > div.condiment").Each(func(i int, selection *goquery.Selection) {
						SecMaterialData = append(SecMaterialData, selection.Find("div.condiment-weight").Text())
					})
					m["sec_material"] = SecMaterialData

					// 步骤
					stepsData := make([]string, 0)
					Doc.Find("#__layout > div > div > article > div > div.recipe-left > div:nth-child(1) > div.practice > div.practices > div.pai > div").Each(func(i int, selection *goquery.Selection) {
						stepsData = append(stepsData, selection.Find("div > div").Text())
					})
					m["steps"] = stepsData

					m["title"] = d["title"].(string)

					// 插入mongodb
					if _, err := coll.InsertOne(context.TODO(), m); err != nil {
						fmt.Println(err)
					}

					count++
					fmt.Println(count)

					// time.Sleep(time.Second)

				}(v)
			} else {
				count++
			}

		}

	}
	wg.Wait()

}

func GetMaterialByUrl(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("err")
	}

	Doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	MaterialData := make([]string, 0)

	Doc.Find("#__layout > div > div > article > div > div.recipe-left > div:nth-child(1) > div.practice > div.practices > div.pai > div").Each(func(i int, selection *goquery.Selection) {
		MaterialData = append(MaterialData, selection.Find("div > div").Text())
		fmt.Println(selection.Find("div > div").Text())
	})

}

// InsertRecipe2Mongo 将菜单列表插入mongodb中
func InsertRecipe2Mongo() {

	file, err := os.Open("abc.record")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	url := "mongodb://mws_mongo:mws_mongo@127.0.0.1:27017/mws"
	// Set client options
	clientOptions := options.Client().ApplyURI(url)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	coll := client.Database("mws").Collection("test")

	count := 0

	for {
		buf, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		a := make(map[string]interface{})
		if err := json.Unmarshal(buf, &a); err != nil {
			panic(err)
		}

		b := a["data"].(map[string]interface{})

		c := b["dataset"].([]interface{})

		count += len(c)

		coll.InsertMany(context.TODO(), c)
	}
	fmt.Println(count)
}

// Insert2Sonic 在sonic中创建索引
func Insert2Sonic() {

	file, err := os.Open("abc.record")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	ingester, err := sonic.NewIngester("127.0.0.1", 27016, "SecretPassword")
	if err != nil {
		panic(err)
	}

	count := 0

	for {
		buf, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		a := make(map[string]interface{})
		if err := json.Unmarshal(buf, &a); err != nil {
			panic(err)
		}

		b := a["data"].(map[string]interface{})

		c := b["dataset"].([]interface{})

		record := make([]sonic.IngestBulkRecord, 0)

		for _, v := range c {
			d := v.(map[string]interface{})
			title := d["title"].(string)
			material := d["material"].(string)

			title = base64.StdEncoding.EncodeToString([]byte(title))
			material = base64.StdEncoding.EncodeToString([]byte(material))

			record = append(record, sonic.IngestBulkRecord{Object: title, Text: material})
		}

		err1 := ingester.BulkPush("mws", "test", 1, record)
		fmt.Println(err1)

		count += len(c)

	}

	fmt.Println(count)

}

// SonicQuery 查询sonic
func SonicQuery() {

	// ingester, err := sonic.NewIngester("127.0.0.1", 27016, "SecretPassword")
	// if err != nil {
	// 	panic(err)
	// }

	// ingester.BulkPush("mws", "test", 3, []sonic.IngestBulkRecord{
	// 	{"xxxx", s},查询sonic
	search, err := sonic.NewSearch("127.0.0.1", 27016, "SecretPassword")
	if err != nil {
		panic(err)
	}

	results, _ := search.Query("mws", "recipe_list", "五花", 10000, 0)

	for _, v := range results {
		if strings.Contains(v, "5ead7b93258d72477c4c3718") {
			fmt.Println("get")
		}
	}
}

// Mongo2Sonic mongodb 中的数据转化为sonic的索引
func Mongo2Sonic() {

	ingester, err := sonic.NewIngester("127.0.0.1", 27016, "SecretPassword")
	if err != nil {
		panic(err)
	}
	defer ingester.Quit()

	url := "mongodb://mws_mongo:mws_mongo@127.0.0.1:27017/mws"
	// Set client options
	clientOptions := options.Client().ApplyURI(url)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	cursor, err := client.Database("mws").Collection("recipe_list").Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
		return
	}

	count := 0

	for cursor.Next(context.TODO()) {
		m := make(map[string]interface{})
		bson.Unmarshal(cursor.Current, &m)

		// if v, ok := m["_id"].(string); ok {
		// 	fmt.Println(v)
		id := m["_id"].(primitive.ObjectID).String()
		object := id

		text := ""

		meterial := m["material"].(primitive.A)
		for _, v := range meterial {
			s := v.(string)
			text = text + s + "、"
			ingester.BulkPush("mws", "recipe_list", 1, []sonic.IngestBulkRecord{
				{Object: object, Text: text},
			})
		}

		count++
	}

	fmt.Println(count)

}
