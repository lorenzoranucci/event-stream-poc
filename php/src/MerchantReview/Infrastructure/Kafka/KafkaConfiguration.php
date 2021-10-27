<?php

declare(strict_types=1);

namespace App\MerchantReview\Infrastructure\Kafka;

use RdKafka\Conf;

final class KafkaConfiguration
{
    private Conf $config;

    public function __construct(
        private string $brokerList,
    ) {
        $this->config = new Conf();
        $this->config->set(
            'metadata.broker.list',
            $this->brokerList,
        );
    }
}
