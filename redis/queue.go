package redis

import "kamva.ir/libraries/contracts"

type Queue struct {
	redis     *Redis
	queueName string
}

func (q Queue) Push(pushable contracts.Pushable) int {
	defer q.redis.connection.Close()
	return q.redis.LPush("mail_queue", pushable.Serialize())
}

func (q Queue) Has(pushable contracts.Pushable) bool {
	defer q.redis.connection.Close()
	values := q.redis.LRange("mail_queue", 0, -1)

	for _, value := range values {
		switch value := value.(type) {
		case string:
			if value == pushable.Serialize() {
				return true
			}
		case []byte:
			if string(value) == pushable.Serialize() {
				return true
			}
		}
	}

	return false
}

func NewQueue(redis *Redis, queueName string) *Queue {
	return &Queue{
		redis: redis,
		queueName: queueName,
	}
}
