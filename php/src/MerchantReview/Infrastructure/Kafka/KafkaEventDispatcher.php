<?php

declare(strict_types=1);

namespace App\MerchantReview\Infrastructure\Kafka;

use App\MerchantReview\Application\EventBus;
use App\MerchantReview\Application\IntegrationEvent;
use RdKafka\Producer;

final class KafkaEventDispatcher implements EventBus
{
    public function __construct(
        private KafkaConfiguration $kafkaConfig,
        private KafkaPartitioner $kafkaPartitioner,
    ) {}

    public function dispatch(IntegrationEvent $event): void
    {
        $producer = new Producer($this->kafkaConfig);

        $topic = $producer->newTopic($event->getType());

        $topic->produce(
            RD_KAFKA_PARTITION_UA,
            $this->kafkaPartitioner->getPartitionID(),
            $event->getMessage(),
        );

        $producer->poll(0);
    }
}
