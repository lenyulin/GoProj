package utils

import (
	"crypto/sha256"
	"encoding/binary"
)

// HashWithSalt 对ID进行哈希加盐处理
// id: 原始业务ID
// salt: 盐值，用于增加哈希的随机性
// 返回: 哈希后的整数结果
func HashWithSalt(id string, salt string) uint64 {
	// 将ID和盐值组合
	data := []byte(id + salt)

	// 使用SHA-256进行哈希计算
	hash := sha256.Sum256(data)

	// 将哈希结果的前8字节转换为uint64
	// 这样可以得到一个范围在0到2^64-1之间的整数
	result := binary.BigEndian.Uint64(hash[:8])

	return result
}

// GetPartition 计算ID对应的分区
// id: 原始业务ID
// salt: 盐值
// partitionCount: 总分区数
// 返回: 分区索引(0 <= partition < partitionCount)
func GetPartition(id string, partitionCount int) int {
	if partitionCount <= 0 {
		return 0
	}

	hashValue := HashWithSalt(id, salt)

	// 对分区数取模，得到最终分区索引
	return int(hashValue % uint64(partitionCount))
}

const salt = "my-secret-salt-123"
