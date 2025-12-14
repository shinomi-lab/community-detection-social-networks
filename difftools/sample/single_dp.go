package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// `json:"..."` ... JSONのキー(...の文字列)とConfig構造体のフィールドを対応付ける構文
// Configが以下の定義の場合のJSONの例：
//
//	{
//	   "id": 1,
//	   "adj_file": "test.txt",
//	   "capacity": 100,
//	   "use_user": true
//	}
type Config struct {
	//Id       int    `json:"adjFilePath"`
	//Capacity int    `json:"capacity"`
	AdjFilePath string `json:"adjFilePath"`
	UseUser     bool   `json:"use_user"`
	UseInfl     bool   `json:"use_infl"`
	SType       int    `json:"S_f_type"`
}

func compute_wapper(config Config) {
	// ここにconfigを使って処理をするコードを書く
	fmt.Printf("--- 処理開始 ---\n")
	// fmt.Printf("AdjFilePath: %s\n", config.AdjFilePath)
	// fmt.Printf("UseUser:     %t\n", config.UseUser)
	// fmt.Printf("UseInfl:     %t\n", config.UseInfl)
	// fmt.Printf("S_f_type:    %d\n", config.SType)

	compute_maximization_DP(
		adj,
		interest_list,
		assum_list,
		user_weight,
		capacity,
		use_kaiki,
		config.UseUser,
		config.UseInfl,
		use_follower,
		1,
		config.SType,
		false)

	fmt.Printf("--- 処理完了 ---\n")
	//fmt.Println(config.AdjFilePath, config.UseUser, config.UseInfl, config.SType)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("JSONファイルを指定してください")
		return
	}
	configFilePath := os.Args[1]
	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("ファイル %s が存在しません", configFilePath)
		} else {
			fmt.Println("ファイル %s を開けません: %v", configFilePath, err)
		}
		return
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("JSONデータを読み込めません", err)
		return
	}

	var config Config
	json.Unmarshal(jsonData, &config)
	compute_wapper(config)
}
