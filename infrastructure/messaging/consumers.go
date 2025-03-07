package messaging

import "log"

func (mq *MessageQueue) startConsumers() {
	go mq.startQueueConsumer(CreateUserQueueName)
}

func (mq *MessageQueue) startQueueConsumer(queueName string) {
	msgs, err := mq.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register consumer for %s: %s", queueName, err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received message from %s: %s", queueName, d.Body)
			err := mq.router.RouteMessage(queueName, d.Body)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages in %s. To exit press CTRL+C", queueName)
}
