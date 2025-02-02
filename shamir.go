package main

import (
	"encoding/base32"
	"fmt"
	"hash/crc32"
	"strings"

	"github.com/hashicorp/vault/shamir"
)

var baseEncoder *base32.Encoding = base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ34679+").WithPadding(base32.NoPadding)
var crcTable *crc32.Table = crc32.MakeTable(crc32.Castagnoli)

func ShamirSplit(shares int, threshold int, secret string, groupSize int) ([]string, error) {
	crc32val := crc32.Checksum([]byte(secret), crcTable)

	secretBytes := []byte{byte(crc32val >> 8), byte(crc32val)}
	secretBytes = append(secretBytes, []byte(secret)...)

	keys, err := shamir.Split(secretBytes, shares, threshold)
	if err != nil {
		return nil, err
	}

	shards := make([]string, shares)
	for i, shard := range keys {
		encoded := baseEncoder.EncodeToString(shard)
		finalStr := ""
		if groupSize == 0 {
			finalStr = encoded
		} else {
			for i := 0; i < len(encoded); i += groupSize {
				if i+groupSize < len(encoded) {
					finalStr += encoded[i:i+groupSize] + " "
				} else {
					finalStr += encoded[i:]
				}
			}
		}
		shards[i] = finalStr
	}

	return shards, nil
}

func ShamirCombine(shards []string) (string, error) {
	shardBytes := make([][]byte, len(shards))
	for i, shard := range shards {
		shard = strings.ToUpper(shard)
		shard = strings.ReplaceAll(shard, " ", "")
		shard = strings.ReplaceAll(shard, "0", "O") // replace 0 to O (common mistake)
		shard = strings.ReplaceAll(shard, "1", "I") // replace 1 to I (common mistake)
		shard = strings.ReplaceAll(shard, "2", "Z") // replace 2 to Z (common mistake)
		shard = strings.ReplaceAll(shard, "5", "S") // replace 5 to S (common mistake)
		shard = strings.ReplaceAll(shard, "8", "B") // replace 8 to B (common mistake)
		decodedShard, err := baseEncoder.DecodeString(shard)
		if err != nil {
			return "", fmt.Errorf("error in shard %d: %v", i+1, err)
		}
		shardBytes[i] = []byte(decodedShard)
	}

	decodedSecret, err := shamir.Combine(shardBytes)
	if err != nil {
		return "", err
	}

	if len(decodedSecret) < 3 {
		return "", fmt.Errorf("decoding failed, shards too short")
	}

	crc32val := crc32.Checksum(decodedSecret[2:], crcTable)
	if byte(crc32val>>8) != decodedSecret[0] || byte(crc32val) != decodedSecret[1] {
		return "", fmt.Errorf("checksum mismatch, invalid or not enough shards?")
	}

	return string(decodedSecret[2:]), nil
}
