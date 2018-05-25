package scheduler

import (
	parserModel "github.com/schwarzeni/go-get-v2/parser/model"
)

// 任务等待队列
type VideoQueue struct {
	VideoLists []parserModel.Video
}

func (q VideoQueue) IsEmpty() bool {
	return len(q.VideoLists) == 0
}

func (q *VideoQueue) Pop() parserModel.Video {
	if q.IsEmpty() {
		return nil
	}
	v := q.VideoLists[0]
	q.VideoLists = q.VideoLists[1:]
	return v
}

func (q *VideoQueue) Push(v parserModel.Video) {
	q.VideoLists = append(q.VideoLists, v)
}
