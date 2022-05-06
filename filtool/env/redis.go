package env

import (
	"bytes"
	"context"
	"encoding/json"
	"filtool/util"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
)

func RedisGetSectorInfo(ctx context.Context, minerId int, sn int) {
	config, err := util.GetConfig()
	if err != nil {
		log.Fatalln("read config failed: ", err)
		return
	}
	log.Println("octopus: init redis client: ", config.Redis.Conn, config.Redis.Password)
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    strings.Split(config.Redis.Conn, ","),
		Password: config.Redis.Password,
		PoolSize: 1,
	})

	key := fmt.Sprintf("%d:%s:%v", minerId, "sectorinfo", sn)
	exists, err := redisClient.Exists(ctx, key).Result()
	log.Println("fetching sector ", sn)
	if err != nil {
		log.Fatalln(err)
		return
	}
	if exists == 1 {
		value, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			log.Println(err)
			return
		}

		var sectorInfo SectorInfo
		err = sectorInfo.UnmarshalCBOR(bytes.NewBuffer([]byte(value)))
		if err != nil {
			log.Println(err)
			return
		}
		sijson, err := json.MarshalIndent(sectorInfo, "", "\t")
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(sijson))
	} else {
		log.Fatalln("not exists.")
	}
}

func RedisClean(ctx context.Context, minerId int, min int, max int) {
	config, err := util.GetConfig()
	if err != nil {
		log.Fatalln("read config failed: ", err)
		return
	}
	log.Println("octopus: init redis client: ", config.Redis.Conn, config.Redis.Password)
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    strings.Split(config.Redis.Conn, ","),
		Password: config.Redis.Password,
		PoolSize: 1,
	})

	var errSn []int
	var errSnMsg []string
	var clrCount int
	var skipCount int
	var logCount int
	for sn := min; sn < max; sn++ {
		key := fmt.Sprintf("%d:%s:%v", minerId, "sectorinfo", sn)
		exists, err := redisClient.Exists(ctx, key).Result()
		log.Println("processing sector ", sn)
		if err != nil {
			errSn = append(errSn, sn)
			errSnMsg = append(errSnMsg, err.Error())
			log.Println(errSnMsg)
			continue
		}
		if exists == 1 {
			value, err := redisClient.Get(ctx, key).Result()
			if err != nil {
				errSn = append(errSn, sn)
				errSnMsg = append(errSnMsg, err.Error())
				log.Println(errSnMsg)
				continue
			}

			var sectorInfo SectorInfo
			err = sectorInfo.UnmarshalCBOR(bytes.NewBuffer([]byte(value)))
			if err != nil {
				errSn = append(errSn, sn)
				errSnMsg = append(errSnMsg, err.Error())
				log.Println(errSnMsg)
				continue
			}
			if sectorInfo.State == "Proving" {
				entries := sectorInfo.Log
				sectorInfo.Log = []Log{}

				buf := new(bytes.Buffer)
				err := sectorInfo.MarshalCBOR(buf)
				if err == nil {
					_, err := redisClient.Set(ctx, key, buf.String(), 0).Result()
					if err != nil {
						errSn = append(errSn, sn)
						errSnMsg = append(errSnMsg, err.Error())
						log.Println(errSnMsg)
						continue
					}
				}

				log.Println("clear ", len(entries), " logs")
				clrCount++
				logCount += len(entries)
			} else {
				log.Println("skip: ", sectorInfo.State)
				skipCount++
			}
		} else {
			log.Println("not exists.")
		}
	}
	log.Println("Summary:")
	log.Println("clear sectors: ", clrCount)
	log.Println("clear logs: ", logCount)
	log.Println("skip sectors: ", skipCount)
	log.Println("error sectors: ")
	for i, sn := range errSn {
		log.Println("sn ", sn, ": ", errSnMsg[i])
	}
}
