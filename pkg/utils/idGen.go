package utils

import "github.com/bwmarrin/snowflake"

func GenerateId() (int64, error) {
	node, err := snowflake.NewNode(0)
	if err != nil {
		return -1, err
	}
	id := node.Generate().Int64()
	return id, nil
}
