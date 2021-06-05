package consumer

import (
	"context"

	"github.com/pursuit/portal/internal/proto/out/event"
	"github.com/pursuit/portal/internal/service/mutation"

	"github.com/golang/protobuf/proto"

	"github.com/Shopify/sarama"
)

type FreeCoinRegisterConsumer struct {
	Ready chan bool

	MutationSvc mutation.Service
}

func (this *FreeCoinRegisterConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(this.Ready)
	return nil
}

func (this *FreeCoinRegisterConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (this *FreeCoinRegisterConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var protoData pursuit_event_proto.UserCreated
		if err := proto.Unmarshal(message.Value, &protoData); err != nil {
			session.MarkMessage(message, "")
			return err
		}

		if err := this.MutationSvc.Create(context.Background(), int(protoData.Id), 1, "free_coins", 10); err != nil {
			if err.Status != 503 {
				session.MarkMessage(message, "")
			}

			return err
		}
	}

	return nil
}
